package dns_certs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/providers/dns/googlecloud"
	"google.golang.org/api/dns/v1"
)

// Create an SSL certificate for the domains.
func getCert(domains []string, svc *dns.Service, config *projConfig) (cr *acme.CertificateResource, msgs []string, err error) {

	dnsProvider, err := googlecloud.NewDNSProviderService(config.projID, svc)
	if err != nil {
		err = fmt.Errorf("could not create a DNSProvider object to get a certificate; %s", err)
		return
	}

	acctFileExists := true
	rdr, err := gcsBucket.Object(sslUserFile).NewReader(context.Background())
	if err != nil {
		if err != storage.ErrObjectNotExist {
			err = fmt.Errorf("could not begin reading user info file; %v", err)
			return
		}
		acctFileExists = false
	}
	acctFile, err := ioutil.ReadAll(rdr) // Where the single account record is (for now).
	if err != nil {
		err = fmt.Errorf("could not read user info file; %v", err)
		return
	}

	acct := &Account{Email: sslUserEmail} // Set email here in case we'll be creating a new account.
	if acctFileExists {
		if err = json.Unmarshal(acctFile, acct); err != nil {
			err = fmt.Errorf("could not decode account JSON file; %s", err)
			return
		}
		userPrivKey, err := loadPrivateKey()
		if err != nil {
			if err != storage.ErrObjectNotExist {
				err = fmt.Errorf("could not read user key file; %s", err)
				return
			}
			// The error is simply saying that no private key exists, so create one.
			if err = acct.genPrivateKey(); err != nil {
				err = fmt.Errorf("could not generate user private key; %s", err)
				return
			}
		} else {
			acct.key = userPrivKey
		}
	} else if err = acct.genPrivateKey(); err != nil {
		err = fmt.Errorf("could not generate user private key; %s", err)
		return
	}

	client, err := acme.NewClient("https://acme-v01.api.letsencrypt.org/directory", acct, acme.EC256)
	if err != nil {
		return
	}

	// New users need to register.
	if !acctFileExists {

		msgs = append(msgs, "Registering a new user")
		acct.Registration, err = client.Register()
		if err != nil {
			err = fmt.Errorf("could not register a new user; %s", err)
			return
		}

		// The client has a URL to the current Let's Encrypt Subscriber Agreement. The user needs to agree to it.
		if err = client.AgreeToTOS(); err != nil {
			err = fmt.Errorf("could not agree to the terms; %s", err)
			return
		}

		if err = acct.save(); err != nil {
			err = fmt.Errorf("could not save new user info; %s", err)
			return
		}

	}

	if len(domains) == 0 {
		err = errors.New("the list of domains is empty")
		return
	}

	// Only use the DNS challenge.
	client.ExcludeChallenges([]acme.Challenge{acme.HTTP01, acme.TLSSNI01})
	client.SetChallengeProvider(acme.DNS01, dnsProvider)

	// Complete the challenge for each domain and obtain the certificate.
	cert, failures := client.ObtainCertificate(domains, false, nil, false)
	if len(failures) > 0 {
		for f, err := range failures {
			msgs = append(msgs, fmt.Sprintf("Error getting cert: %q => %v", f, err))
		}
		err = errors.New("could not obtain a certificate")
		return
	}

	cr = &cert
	return

}

func certName(projID, timestamp string) string {
	return projID + "-cert-" + timestamp + ".pem"
}

func certAuthorityName(projID, timestamp string) string {
	return projID + "-cert-authority-" + timestamp + ".pem"
}

func certKeyName(projID, timestamp string) string {
	return projID + "-cert-key-" + timestamp + ".pem"
}
