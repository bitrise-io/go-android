package keystore

import (
	"crypto/x509"
	"errors"
	"fmt"
)

type CertificateInformation struct {
	FirstAndLastName   string
	OrganizationalUnit string
	Organization       string
	CityOrLocality     string
	StateOrProvince    string
	CountryCode        string
	ValidFrom          string
	ValidUntil         string
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

type Reader struct {
	decoders []Decoder
}

func NewReader(decoders []Decoder) Reader {
	return Reader{decoders: decoders}
}

func NewDefaultReader() Reader {
	return NewReader([]Decoder{PKCS12KeystoreDecoder{}, JKSKeystoreDecoder{}})
}

func (p Reader) ReadCertificateInformation(data []byte, password, alias, keyPassword string) (*CertificateInformation, error) {
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

	certInfo := parseCertificate(cert)

	return &certInfo, nil

}

func parseCertificate(certificate *x509.Certificate) CertificateInformation {
	certInfo := CertificateInformation{}

	issuer := certificate.Issuer
	certInfo.FirstAndLastName = issuer.CommonName

	if len(issuer.OrganizationalUnit) > 0 {
		certInfo.OrganizationalUnit = issuer.OrganizationalUnit[0]
	}

	if len(issuer.Organization) > 0 {
		certInfo.Organization = issuer.Organization[0]
	}

	if len(issuer.Locality) > 0 {
		certInfo.CityOrLocality = issuer.Locality[0]
	}

	if len(issuer.Province) > 0 {
		certInfo.StateOrProvince = issuer.Province[0]
	}

	if len(issuer.Country) > 0 {
		certInfo.CountryCode = issuer.Country[0]
	}

	certInfo.ValidFrom = certificate.NotBefore.String()
	certInfo.ValidUntil = certificate.NotAfter.String()

	return certInfo
}
