package keystore

import (
	"crypto/x509"
	"fmt"

	"github.com/bitrise-io/go-pkcs12"
)

type PKCS12KeystoreDecoder struct {
}

func (d PKCS12KeystoreDecoder) Decode(data []byte, password, alias, keyPassword string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	key, cert, err := pkcs12.DecodeKeystore(data, password, alias, keyPassword)
	if err != nil {
		switch err {
		case pkcs12.IncorrectKeystorePasswordError:
			return nil, nil, IncorrectKeystorePasswordError
		case pkcs12.IncorrectAliasError:
			return nil, nil, IncorrectAliasError
		case pkcs12.IncorrectKeyPasswordError:
			return nil, nil, IncorrectKeyPasswordError
		default:
			return nil, nil, err
		}
	}

	if cert == nil {
		return nil, nil, fmt.Errorf("certificate not found")
	}

	return key, cert, nil
}

func (d PKCS12KeystoreDecoder) IsInvalidCredentialsError(err error) bool {
	return err == IncorrectKeystorePasswordError ||
		err == IncorrectAliasError ||
		err == IncorrectKeyPasswordError
}
