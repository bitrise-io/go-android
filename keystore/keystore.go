package keystore

import (
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-pkcs12"
	"github.com/lwithers/minijks/jks"
)

type KeyStore struct {
	KeyStoreDetails
}

func Open(pth string, password string, privateKeyAlias, privateKeyPassword string) (*KeyStore, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	b, err := os.ReadFile(pth)
	if err != nil {
		return nil, err
	}

	// a keystore is either a PKCS12 type keystore
	keystore, err := parsePKCS12Keystore(b, password, privateKeyAlias, privateKeyPassword)
	if err != nil && !strings.Contains(err.Error(), "pkcs12: error reading P12 data: asn1: structure error: length too large") {
		return nil, err
	}
	// or a JKS type keystore (might pkcs8)
	if keystore == nil {
		keystore, err = parseJKSKeystore(b, password, privateKeyAlias, privateKeyPassword)
		if err != nil {
			return nil, err
		}
	}

	return keystore, nil
}

func parseJKSKeystore(content []byte, password, privateKeyAlias, privateKeyPassword string) (*KeyStore, error) {
	ks, err := jks.Parse(content, &jks.Options{
		Password:         password,
		SkipVerifyDigest: false,
		KeyPasswords:     map[string]string{privateKeyAlias: privateKeyPassword},
	})
	if err != nil {
		return nil, err
	}

	certificate := ks.Keypairs[0].CertChain[0]
	details := parseCertificate(certificate.Cert)
	return &KeyStore{details}, nil
}

type KeyStoreDetails struct {
	FirstAndLastName   string
	OrganizationalUnit string
	Organization       string
	CityOrLocality     string
	StateOrProvince    string
	CountryCode        string
	ValidFrom          string
	ValidUntil         string
}

func parsePKCS12Keystore(content []byte, password, privateKeyAlias, privateKeyPassword string) (*KeyStore, error) {
	_, certificate, err := pkcs12.DecodeKeystore(content, password, privateKeyAlias, privateKeyPassword)
	if err != nil {
		return nil, err
	}

	details := parseCertificate(certificate)
	return &KeyStore{details}, nil
}

func parseCertificate(certificate *x509.Certificate) KeyStoreDetails {
	issuer := certificate.Issuer
	details := KeyStoreDetails{}

	details.FirstAndLastName = issuer.CommonName

	if len(issuer.OrganizationalUnit) > 0 {
		details.OrganizationalUnit = issuer.OrganizationalUnit[0]
	}

	if len(issuer.Organization) > 0 {
		details.Organization = issuer.Organization[0]
	}

	if len(issuer.Locality) > 0 {
		details.CityOrLocality = issuer.Locality[0]
	}

	if len(issuer.Province) > 0 {
		details.StateOrProvince = issuer.Province[0]
	}

	if len(issuer.Country) > 0 {
		details.CountryCode = issuer.Country[0]
	}

	details.ValidFrom = certificate.NotBefore.String()

	details.ValidUntil = certificate.NotAfter.String()

	return details
}
