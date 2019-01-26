package dns_certs

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/iterator"
)

const (
	sslUserEmail   = "" // An email address of the user managing the SSL certificates.
	sslUserFile    = "acct.json"
	sslUserKeyFile = "acct.key"
)

type projConfig struct {
	projID, loadBalancerIP string
}

var (
	gcpDefaultCreds *google.DefaultCredentials
	gcsBucket       *storage.BucketHandle
)

func init() {

	var err error
	gcpDefaultCreds, err = google.FindDefaultCredentials(context.Background(), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		panic(fmt.Sprintf("could not get GCP credentials; %v", err))
	}

	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		panic(fmt.Sprintf("could not create GCS cloud.Client; %v", err))
	}
	gcsBucket = gcsClient.Bucket(bucketName)

}

// DomainsZones retrieves the latest list of domains and zones for the project. Some or all of the domains may not have a certificate.
func DomainsZones(projShortName string) (dz []string, err error) {

	pc, ok := configs[projShortName]
	if !ok {
		err = fmt.Errorf("the given project short name %q is invalid", projShortName)
		return
	}

	service, err := dns.New(oauth2.NewClient(context.Background(), gcpDefaultCreds.TokenSource))
	if err != nil {
		err = fmt.Errorf("could not set up DNS API; %v", err)
		return
	}

	timestamp, domains, err := readDomainsList(pc.projID)
	if err != nil {
		err = fmt.Errorf("could not read domains.txt; %v", err)
		return
	}

	// Check if a certificate exists for these domains.
	gcsBucket.Object()

	for i := range domains {

	}

	return

}

type DomainZoneInfo struct {
	Domain, ARecord string
}

func UpdateDomains(projShortName string) (msgs []string, err error) {

	pc, ok := configs[projShortName]
	if !ok {
		return nil, fmt.Errorf("the given project short name %q is invalid", projShortName)
	}

	service, err := dns.New(oauth2.NewClient(context.Background(), gcpDefaultCreds.TokenSource))
	if err != nil {
		err = fmt.Errorf("could not set up DNS API; %v", err)
		return
	}

	timestamp, domains, err := readDomainsList(pc.projID)
	if err != nil {
		return nil, fmt.Errorf("could not read domains.txt; %v", err)
	}

	msgs, err = updateDNS(domains, service, &pc)
	if err != nil {
		return
	}

	cert, msgs2, err := getCert(domains, service, &pc)
	msgs = append(msgs, msgs2...)
	if err != nil {
		return
	}

	ctx := context.Background() // TODO: make this with cancel?
	errs := make(chan error, 3)

	go func() {
		wrCert := gcsBucket.Object(certName(pc.projID, timestamp)).NewWriter(ctx)
		if _, err = wrCert.Write(cert.Certificate); err != nil {
			errs <- fmt.Errorf("could not write certificate file; %s", err)
			return
		}
		errs <- wrCert.Close()
	}()

	go func() {
		wrCertAuthority := gcsBucket.Object(certAuthorityName(pc.projID, timestamp)).NewWriter(ctx)
		if _, err = wrCertAuthority.Write(cert.IssuerCertificate); err != nil {
			errs <- fmt.Errorf("could not write CA certificate file; %s", err)
			return
		}
		errs <- wrCertAuthority.Close()
	}()

	go func() {
		wrKey := gcsBucket.Object(certKeyName(pc.projID, timestamp)).NewWriter(ctx)
		if _, err = wrKey.Write(cert.PrivateKey); err != nil {
			errs <- fmt.Errorf("could not write private key file; %s", err)
			return
		}
		errs <- wrKey.Close()
	}()

	for i := 0; i < 3; i++ {
		if err = <-errs; err != nil {
			return // Only report the first error.
		}
	}

	return

}

// SetDomainsList saves a new version of the list of domains for a project.
func SetDomainsList(projShortName string, newList []string) error {
	return nil
}

// readDomainsList returns, for the given project, the timestamp (as string) of the latest version of the domains list and
// the cleaned up list of domains (not fully qualified).
func readDomainsList(projID string) (string, []string, error) {

	listPrefix := projID + "-domains-" // The name of the list of domains has this prefix.

	q := &storage.Query{
		Prefix:   listPrefix,
		Versions: false,
	}
	oi := gcsBucket.Objects(context.Background(), q)

	var oal objAttrsList
	oal.oa = make([]*storage.ObjectAttrs, 0, 4)

	for {
		oa, err := oi.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", nil, fmt.Errorf("could not iterate over list of objects; %v", err)
		}
		oal.oa = append(oal.oa, oa)
	}

	if len(oal.oa) == 0 {
		return "", nil, errNoDomains
	}
	sort.Sort(&oal)
	dl := oal.oa[len(oal.oa)-1] // The latest list

	objReader, err := gcsBucket.Object(dl.Name).NewReader(context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("could not begin reading list of domains; %v", err)
	}

	data, err := ioutil.ReadAll(objReader)
	if err != nil {
		return "", nil, fmt.Errorf("could not read list of domains; %v", err)
	}

	timestamp := strings.TrimPrefix(dl.Name, listPrefix)
	return timestamp, removeBlank(strings.Split(string(data), "\n")), nil

}

type objAttrsList struct {
	oa []*storage.ObjectAttrs
}

func (oal *objAttrsList) Len() int { return len(oal.oa) }

func (oal *objAttrsList) Less(i, j int) bool {
	return oal.oa[i].Name < oal.oa[j].Name
}

func (oal *objAttrsList) Swap(i, j int) {
	oal.oa[i], oal.oa[j] = oal.oa[j], oal.oa[i]
}

var errNoDomains = errors.New("no domains list found")

// removeBlank returns the list of given domains (separated by \n) trimmed and pruned of blank lines.
func removeBlank(ds []string) []string {
	for i := 0; i < len(ds); i++ {
		ds[i] = strings.TrimSpace(ds[i])
		if len(ds[i]) == 0 {
			ds = append(ds[:i], ds[i+1:]...)
			i--
		}
	}
	return ds
}
