package keystore

import (
	"crypto/x509"
	"errors"
	"fmt"
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

type Decoder interface {
	Decode(data []byte, password, alias, keyPassword string) (privateKey interface{}, certificate *x509.Certificate, err error)
	IsInvalidCredentialsError(err error) bool
}

type Parser struct {
	decoders []Decoder
}

func NewParser(decoders []Decoder) Parser {
	return Parser{decoders: decoders}
}

func NewDefaultParser() Parser {
	return NewParser([]Decoder{PKCS12KeystoreDecoder{}, JKSKeystoreDecoder{}})
}

func (p Parser) Parse(data []byte, password, alias, keyPassword string) (*KeyStore, error) {
	var cert *x509.Certificate
	var decodeErrs []error

	for _, decoder := range p.decoders {
		_, c, err := decoder.Decode(data, password, alias, keyPassword)
		if err != nil {
			if decoder.IsInvalidCredentialsError(err) {
				return nil, err
			}

			decodeErrs = append(decodeErrs, err)
		} else {
			cert = c
			break
		}
	}

	if cert == nil && len(decodeErrs) > 0 {
		decodeErrsStr := ""
		for i, decodeErr := range decodeErrs {
			decodeErrsStr += "- " + decodeErr.Error()
			if i < len(decodeErrs)-1 {
				decodeErrsStr += "\n"
			}
		}
		return nil, fmt.Errorf("failed to decode keystore:\n%s", decodeErrsStr)
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
