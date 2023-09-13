package keystore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-pkcs12"
	"github.com/lwithers/minijks/jks"
)

type Details struct {
	FirstAndLastName   string
	OrganizationalUnit string
	Organization       string
	CityOrLocality     string
	StateOrProvince    string
	CountryCode        string
	ValidFrom          string
	ValidUntil         string
}

type KeyStore struct {
	Details
}

var (
	IncorrectAliasError            = errors.New("incorrect key alias")
	IncorrectKeystorePasswordError = errors.New("incorrect keystore password")
	IncorrectKeyPasswordError      = errors.New("incorrect key password")
)

func Open(pth string, password string, alias, keyPassword string) (*KeyStore, error) {
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
	keystore, err := parsePKCS12Keystore(b, password, alias, keyPassword)
	if err != nil && !strings.Contains(err.Error(), "pkcs12: error reading P12 data: asn1: structure error: length too large") {
		return nil, err
	}
	// or a JKS type keystore
	if keystore == nil {
		keystore, err = parseJKSKeystore(b, password, alias, keyPassword)
		if err != nil {
			return nil, err
		}
	}

	return keystore, nil
}

func parseJKSKeystore(content []byte, password, alias, keyPassword string) (*KeyStore, error) {
	ks, err := jks.Parse(content, &jks.Options{
		Password:         password,
		SkipVerifyDigest: false,
		KeyPasswords:     map[string]string{alias: keyPassword},
	})
	if err != nil {
		if err.Error() == "digest mismatch" {
			return nil, IncorrectKeystorePasswordError
		}
		return nil, err
	}

	var keypair *jks.Keypair
	for _, kp := range ks.Keypairs {
		if kp.Alias == alias {
			keypair = kp
			break
		}
	}
	if keypair == nil {
		return nil, IncorrectAliasError
	}
	if keypair.PrivKeyErr != nil {
		if keypair.PrivKeyErr.Error() == "invalid password" {
			return nil, IncorrectKeyPasswordError
		}
		return nil, fmt.Errorf("failed to decrypt key: %s", keypair.PrivKeyErr)
	}

	var cert *x509.Certificate
	if len(keypair.CertChain) > 0 {
		cert = keypair.CertChain[0].Cert
	}
	if cert == nil {
		return nil, fmt.Errorf("certificate not found")
	}

	details := parseCertificate(cert)
	return &KeyStore{details}, nil
}

func parsePKCS12Keystore(content []byte, password, alias, keyPassword string) (*KeyStore, error) {
	_, cert, err := pkcs12.DecodeKeystore(content, password, alias, keyPassword)
	if err != nil {
		switch err {
		case pkcs12.IncorrectKeystorePasswordError:
			return nil, IncorrectKeystorePasswordError
		case pkcs12.IncorrectAliasError:
			return nil, IncorrectAliasError
		case pkcs12.IncorrectKeyPasswordError:
			return nil, IncorrectKeyPasswordError
		default:
			return nil, err
		}
	}

	if cert == nil {
		return nil, fmt.Errorf("certificate not found")
	}

	details := parseCertificate(cert)
	return &KeyStore{details}, nil
}

func parseCertificate(certificate *x509.Certificate) Details {
	issuer := certificate.Issuer
	details := Details{}

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
