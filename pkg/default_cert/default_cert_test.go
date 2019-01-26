package cert

import (
	"crypto/tls"
	"testing"
)

func TestDefaultCert(t *testing.T) {
	cert, err := DefaultCert()
	if err != nil {
		t.Fatalf("could not create certificate; %v", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{*cert},
	}
	conf.BuildNameToCertificate()
	if len(conf.NameToCertificate) != 1 {
		t.Errorf("got not exactly one name-to-certificate element; %v", conf.NameToCertificate)
	}
	if len(conf.NameToCertificate) > 0 && conf.NameToCertificate["localhost"] == nil {
		t.Errorf("got a nil NameToCertificate[\"localhost\"]")
	}
}
