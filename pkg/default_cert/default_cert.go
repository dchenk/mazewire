// Package cert helps you generate a self-signed certificate for a TLS server (for host "localhost").
package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"
)

var hosts = []string{"localhost"}

var (
	mu      sync.Mutex
	created *tls.Certificate
)

func DefaultCert() (*tls.Certificate, error) {

	mu.Lock()
	defer mu.Unlock()

	if created != nil {
		return created, nil
	}

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("cert: could not generate private key: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("cert: could not generate serial number: %v", err)
	}

	now := time.Now()

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Mazewire Server"},
		},
		NotBefore:             now,
		NotAfter:              now.Add(2 * 365 * 24 * time.Hour), // Two years
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, fmt.Errorf("cert: could not create certificate: %v", err)
	}

	created = &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privKey,
	}

	return created, nil

}
