package keystore

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name               string
		pth                string
		password           string
		privateKeyAlias    string
		privateKeyPassword string
		want               *KeyStore
	}{
		{
			name:               "PKCS12 keystore test",
			pth:                filepath.Join("testdata", "pkcs12_type_keystore.jks"),
			password:           "storepass",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypass",
			want: &KeyStore{
				Details{
					FirstAndLastName:   "Bitrise Bot",
					OrganizationalUnit: "",
					Organization:       "",
					CityOrLocality:     "",
					StateOrProvince:    "",
					CountryCode:        "",
					ValidFrom:          "2023-09-11 13:18:53 +0000 UTC",
					ValidUntil:         "2048-09-04 13:18:53 +0000 UTC",
				},
			},
		},
		{
			name:               "JKS keystore test",
			pth:                filepath.Join("testdata", "jks_type_keystore.keystore"),
			password:           "keystore",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "keystore",
			want: &KeyStore{
				Details{
					FirstAndLastName:   "asdf",
					OrganizationalUnit: "asdf",
					Organization:       "asdf",
					CityOrLocality:     "asdf",
					StateOrProvince:    "asdf",
					CountryCode:        "as",
					ValidFrom:          "2016-07-14 10:10:41 +0000 UTC",
					ValidUntil:         "2043-11-30 10:10:41 +0000 UTC",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.pth, tt.password, tt.privateKeyAlias, tt.privateKeyPassword)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestOpenErrors(t *testing.T) {
	tests := []struct {
		name               string
		pth                string
		password           string
		privateKeyAlias    string
		privateKeyPassword string
		wantError          string
	}{
		{
			name:               "PKCS12 keystore test - incorrect password",
			pth:                filepath.Join("testdata", "pkcs12_type_keystore.jks"),
			password:           "incorrect-password",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypass",
			wantError:          IncorrectKeystorePasswordError.Error(),
		},
		{
			name:               "PKCS12 keystore test - incorrect alias",
			pth:                filepath.Join("testdata", "pkcs12_type_keystore.jks"),
			password:           "storepass",
			privateKeyAlias:    "incorrect-alias",
			privateKeyPassword: "keypass",
			wantError:          IncorrectAliasError.Error(),
		},
		{
			name:               "PKCS12 keystore test - incorrect key password",
			pth:                filepath.Join("testdata", "pkcs12_type_keystore.jks"),
			password:           "storepass",
			privateKeyAlias:    "key0",
			privateKeyPassword: "incorrect-keypassword",
			wantError:          IncorrectKeyPasswordError.Error(),
		},
		{
			name:               "JKS keystore test - incorrect password",
			pth:                filepath.Join("testdata", "jks_type_keystore.keystore"),
			password:           "incorrect-password",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "keystore",
			wantError:          IncorrectKeystorePasswordError.Error(),
		},
		{
			name:               "JKS keystore test - incorrect alias",
			pth:                filepath.Join("testdata", "jks_type_keystore.keystore"),
			password:           "keystore",
			privateKeyAlias:    "incorrect-alias",
			privateKeyPassword: "keystore",
			wantError:          IncorrectAliasError.Error(),
		},
		{
			name:               "JKS keystore test - incorrect key password",
			pth:                filepath.Join("testdata", "jks_type_keystore.keystore"),
			password:           "keystore",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "incorrect-keypassword",
			wantError:          IncorrectKeyPasswordError.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.pth, tt.password, tt.privateKeyAlias, tt.privateKeyPassword)
			require.EqualError(t, err, tt.wantError)
			require.Nil(t, got)
		})
	}
}
