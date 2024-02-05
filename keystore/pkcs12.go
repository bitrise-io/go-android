package keystore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"github.com/bitrise-io/go-pkcs12"
)

type PKCS12KeystoreDecoder struct {
}

func (d PKCS12KeystoreDecoder) Decode(data []byte, password, alias, keyPassword string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	alias = strings.ToLower(alias)
	key, cert, err := pkcs12.DecodeKeystore(data, password, alias, keyPassword)
	if err != nil {
		return nil, nil, keystoreErrorFromPKCS12Error(err)
	}

	if cert == nil {
		return nil, nil, fmt.Errorf("certificate not found")
	}

	return key, cert, nil
}

func (d PKCS12KeystoreDecoder) IsInvalidCredentialsError(err error) bool {
	return errors.Is(err, IncorrectKeystorePasswordError) ||
		errors.Is(err, IncorrectAliasError) ||
		errors.Is(err, IncorrectKeyPasswordError)
}

func keystoreErrorFromPKCS12Error(err error) error {
	switch {
	case errors.Is(err, pkcs12.IncorrectKeystorePasswordError):
		return IncorrectKeystorePasswordError
	case errors.Is(err, pkcs12.IncorrectAliasError):
		return IncorrectAliasError
	case errors.Is(err, pkcs12.IncorrectKeyPasswordError):
		return IncorrectKeyPasswordError
	default:
		return err
	}
}
