package dns_certs

import (
	"fmt"
	"strings"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
)

// Add all domains to the Google Cloud project DNS service, setting the single correct A record for each.
func updateDNS(domains []string, svc *dns.Service, setup *projConfig) (msgs []string, err error) {

	for _, d := range domains {

		mz := &dns.ManagedZone{
			DnsName:     d + ".",
			Name:        strings.Replace(d, ".", "-", -1),
			Description: "zone",
		}

		zoneExists := false
		managedZone, err := svc.ManagedZones.Get(setup.projID, strings.Replace(d, ".", "-", -1)).Do()
		if err != nil {
			if googleapi.IsNotModified(err) {
				zoneExists = true
				managedZone = mz
			} else if err.(*googleapi.Error).Code != 404 {
				err = fmt.Errorf("could not check for zone %s; %s", d, err)
				return
			}
		} else if managedZone != nil { // Double checking; shouldn't be nil at this point.
			zoneExists = true
		}

		if !zoneExists {
			msgs = append(msgs, "Creating zone for domain "+d)
			managedZone, err = svc.ManagedZones.Create(setup.projID, mz).Do()
			if err != nil {
				err = fmt.Errorf("could not create zone for %s; %s", d, err)
				return
			}
		}

		aRecordsList, err := svc.ResourceRecordSets.List(setup.projID, managedZone.Name).Name(managedZone.DnsName).Type("A").Do()
		if err != nil {
			err = fmt.Errorf("could not check zone %s for A records; %s", d, err)
			return
		}

		// Check if the current value of the A record is correct.
		if aRecordsList != nil {
			if sets := aRecordsList.Rrsets; len(aRecordsList.Rrsets) > 0 && len(sets) > 0 {

				if data := sets[0].Rrdatas; len(data) == 1 && data[0] == setup.loadBalancerIP {
					continue
				}

				// Delete the current A record values.
				deletion := &dns.Change{Deletions: sets}
				if _, err := svc.Changes.Create(setup.projID, managedZone.Name, deletion).Do(); err != nil {
					err = fmt.Errorf("could not delete old A records for domain %s; %s", d, err)
					return
				}

			}
		}

		addition := &dns.Change{
			Additions: []*dns.ResourceRecordSet{
				{
					Name:    managedZone.DnsName,
					Type:    "A",
					Rrdatas: []string{setup.loadBalancerIP},
					Ttl:     3600, // 30 minutes
				},
			},
		}

		if _, err = svc.Changes.Create(setup.projID, managedZone.Name, addition).Do(); err != nil {
			if !googleapi.IsNotModified(err) {
				err = fmt.Errorf("could not create A record for %s; %s", d, err)
				return
			}
		}

		msgs = append(msgs, "Updated zone for domain "+d)

	}

	return

}
