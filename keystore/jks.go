package keystore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"github.com/lwithers/minijks/jks"
)

type JKSKeystoreDecoder struct {
}

func (d JKSKeystoreDecoder) Decode(data []byte, password, alias, keyPassword string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	alias = strings.ToLower(alias)
	ks, err := jks.Parse(data, &jks.Options{
		Password:         password,
		SkipVerifyDigest: false,
		KeyPasswords:     map[string]string{alias: keyPassword},
	})
	if err != nil {
		if err.Error() == "digest mismatch" {
			return nil, nil, ErrIncorrectKeystorePassword
		}
		return nil, nil, err
	}

	var keypair *jks.Keypair
	for _, kp := range ks.Keypairs {
		if kp.Alias == alias {
			keypair = kp
			break
		}
	}
	if keypair == nil {
		return nil, nil, ErrIncorrectAlias
	}
	if keypair.PrivKeyErr != nil {
		if keypair.PrivKeyErr.Error() == "invalid password" {
			return nil, nil, ErrIncorrectKeyPassword
		}
		return nil, nil, fmt.Errorf("failed to decrypt key: %s", keypair.PrivKeyErr)
	}

	var cert *x509.Certificate
	if len(keypair.CertChain) > 0 {
		cert = keypair.CertChain[0].Cert
	}
	if cert == nil {
		return nil, nil, fmt.Errorf("certificate not found")
	}

	return keypair.PrivateKey, cert, nil
}

func (d JKSKeystoreDecoder) IsInvalidCredentialsError(err error) bool {
	return errors.Is(err, ErrIncorrectKeystorePassword) ||
		errors.Is(err, ErrIncorrectAlias) ||
		errors.Is(err, ErrIncorrectKeyPassword)
}
