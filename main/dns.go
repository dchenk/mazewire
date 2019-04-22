package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/env"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/users"
	"github.com/dchenk/mazewire/pkg/util"
	"golang.org/x/oauth2"
	"google.golang.org/api/dns/v1"
)

// list all zones for the Google project
// API params: (none)
func dnsListZones(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	if !users.IsSuper(u.Id) {
		return errLowPrivileges()
	}

	ctx, service, err := setupDnsApi(r)
	if err != nil {
		return errProcessing()
	}

	body := make(RespDnsListZones, 0, 1)

	call := service.ManagedZones.List(env.GCP_PROJECT) // returns *ManagedZonesListCall
	if err := call.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for _, zone := range page.ManagedZones {
			body = append(body, &DnsManagedZone{
				CreationTime:  zone.CreationTime,
				DnsName:       zone.DnsName,
				Id:            zone.Id,
				Name:          zone.Name,
				NameServerSet: zone.NameServerSet,
				NameServers:   zone.NameServers,
			})
		}
		return nil
	}); err != nil {
		log.Err(r, "could not list all DNS zones", err)
		return errProcessing()
	}

	return &APIResponse{Body: body}
}

type RespDnsListZones []*DnsManagedZone

// A mirror of the important fields of google.golang.org/api/dns.ManagedZone
type DnsManagedZone struct {
	CreationTime  string   `msgp:"created"`
	DnsName       string   `msgp:"dns_name"`
	Id            uint64   `msgp:"id"`
	Name          string   `msgp:"name"`
	NameServerSet string   `msgp:"ns_set"` // maybe this is not necessary?
	NameServers   []string `msgp:"name_servers"`
}

// list all records for a zone
func dnsListZoneRecords(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	var reqData ReqDnsListZoneRecords
	err := msgp.Decode(r.Body, &reqData)
	if err != nil {
		log.Err(r, errDecodingMsgp, err)
		return errProcessing()
	}

	// reqData.Site defaults to the current host
	if reqData.Site == 0 {
		reqData.Site = s.Id
	}

	role := u.Role
	if reqData.Site != s.Id {
		role, err = users.SiteRole(u.Id, reqData.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
	}

	if !users.RoleAtLeast(role, users.RoleAdmin) {
		if u.Id == 0 {
			return errMustLogin()
		}
		return errLowPrivileges()
	}

	if util.IsAnyStringBlank(reqData.Zone) {
		return APIResponseErr("The 'zone' parameter is missing")
	}

	ctx, service, err := setupDnsApi(r)
	if err != nil {
		return errProcessing()
	}

	body := make(RespDnsListZoneRecords, 0, 1)

	call := service.ResourceRecordSets.List(env.GCP_PROJECT, reqData.Zone)
	if err := call.Pages(ctx, func(setsList *dns.ResourceRecordSetsListResponse) error {
		for i := range setsList.Rrsets {
			body = append(body, &DnsResourceRecordSet{
				Name:    setsList.Rrsets[i].Name,
				Rrdatas: setsList.Rrsets[i].Rrdatas,
				Ttl:     setsList.Rrsets[i].Ttl,
				Type:    setsList.Rrsets[i].Type,
			})
		}
		return nil
	}); err != nil {
		log.Err(r, "could not list all DNS records for zone "+reqData.Zone, err)
		return errProcessing()
	}

	return &APIResponse{Body: body}
}

type ReqDnsListZoneRecords struct {
	Zone string `json:"zone"` // the managed zone name or id -- not DNS name
	Site int64  `json:"site"` // the site ID; defaults to ID of current host
}

// A mirror of the important fields of google.golang.org/api/dns.ResourceRecordSet
type DnsResourceRecordSet struct {
	Name    string   `msgp:"name"`    // For example, www.example.com.
	Rrdatas []string `msgp:"rrdatas"` // As defined in RFC 1035 (section 5) and RFC 1034 (section 3.6.1).
	Ttl     int64    `msgp:"ttl"`     // Number of seconds that this ResourceRecordSet can be cached by resolvers.
	Type    string   `msgp:"type"`    // The identifier of a supported record type, for example, A, AAAA, MX, TXT.
}

type RespDnsListZoneRecords []*DnsResourceRecordSet

// create a zone
func dnsCreateZone(r *http.Request, _ *data.Site, u *data.User) *APIResponse {
	var reqData ReqDnsCreateZone
	if err := msgp.Decode(r.Body, &reqData); err != nil {
		log.Err(r, errDecodingMsgp, err)
		return errProcessing()
	}

	// TODO: check if user has the right priveleges

	reqData.Domain = strings.TrimSpace(reqData.Domain)
	if reqData.Domain == "" {
		return APIResponseErr("The 'domain' parameter is missing")
	}

	// TODO: validate the domain format by regexp

	// the name of the zone in Google's database
	name := strings.Replace(reqData.Domain, ".", "-", -1)

	// Add a dot to the end of the domain (the way Google wants it).
	if reqData.Domain[len(reqData.Domain)-1:] == "." {
		name = name[:len(name)-2] // Remove the last dash that replaced the dot.
	} else {
		reqData.Domain += "."
	}

	rb := &dns.ManagedZone{
		DnsName:     reqData.Domain,
		Name:        name,
		Description: name,
	}

	ctx, service, err := setupDnsApi(r)
	if err != nil {
		return errProcessing()
	}

	body := new(RespDnsCreateZone)

	managedZone, err := service.ManagedZones.Create(env.GCP_PROJECT, rb).Context(ctx).Do()
	if err != nil {
		log.Err(r, "could not set up DNS API service", err)
		return errProcessing()
	}

	if managedZone.DnsName != reqData.Domain {
		log.Err(r, "could not zone for domain "+reqData.Domain, err)
		return errProcessing()
	}

	return &APIResponse{Body: body}
}

type ReqDnsCreateZone struct {
	Domain string `json:"domain"` // the domain name for the new zone
}

type RespDnsCreateZone struct { // TODO
}

// set up context and service for Google DNS API
func setupDnsApi(r *http.Request) (context.Context, *dns.Service, error) {

	ctx := context.Background()

	service, err := dns.New(oauth2.NewClient(ctx, gcpDefaultCreds.TokenSource))
	if err != nil {
		log.Err(r, "could not set up DNS API service", err)
		return ctx, nil, err
	}

	return ctx, service, nil

}
