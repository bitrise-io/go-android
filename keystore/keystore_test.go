package keystore

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name               string
		pth                string
		password           string
		privateKeyAlias    string
		privateKeyPassword string
		want               *CertificateInformation
		wantError          string
	}{
		{
			name:               "PKCS12 keystore test",
			pth:                filepath.Join("testdata", "keystore.pkcs12"),
			password:           "storepass",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypass",
			want: &CertificateInformation{
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
		{
			name:               "JKS keystore test",
			pth:                filepath.Join("testdata", "keystore.jks"),
			password:           "keystore",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "keystore",
			want: &CertificateInformation{
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
		{
			name:               "PKCS12 Keystore with upper case letters in the alias",
			pth:                filepath.Join("testdata", "upper_case_alias_keystore.pkcs12"),
			password:           "keystore",
			privateKeyAlias:    "MyKey",
			privateKeyPassword: "keystore",
			want: &CertificateInformation{
				Organization: "Bitrise",
				ValidFrom:    "2024-01-31 14:08:42 +0000 UTC",
				ValidUntil:   "2049-01-24 14:08:42 +0000 UTC",
			},
		},
		{
			name:               "JKS Keystore with upper case letters in the alias",
			pth:                filepath.Join("testdata", "upper_case_alias_keystore.jks"),
			password:           "keystore",
			privateKeyAlias:    "Alias0",
			privateKeyPassword: "keystore",
			want: &CertificateInformation{
				FirstAndLastName:   "Unknown",
				OrganizationalUnit: "Unknown",
				Organization:       "Bitrise",
				CityOrLocality:     "Unknown",
				StateOrProvince:    "Unknown",
				CountryCode:        "Unknown",
				ValidFrom:          "2024-01-31 14:34:34 +0000 UTC",
				ValidUntil:         "2051-06-18 14:34:34 +0000 UTC",
			},
		},
		{
			name:               "Keystore with multiple keys - key0",
			pth:                filepath.Join("testdata", "multiple_keys_keystore.pkcs12"),
			password:           "storepass",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypass0",
			want: &CertificateInformation{
				FirstAndLastName:   "",
				OrganizationalUnit: "",
				Organization:       "Bitrise",
				CityOrLocality:     "",
				StateOrProvince:    "",
				CountryCode:        "",
				ValidFrom:          "2024-11-18 14:41:36 +0000 UTC",
				ValidUntil:         "2025-11-18 14:41:36 +0000 UTC",
			},
		},
		{
			name:               "Keystore with multiple keys - key1",
			pth:                filepath.Join("testdata", "multiple_keys_keystore.pkcs12"),
			password:           "storepass",
			privateKeyAlias:    "key1",
			privateKeyPassword: "keypass1",
			want: &CertificateInformation{
				FirstAndLastName:   "",
				OrganizationalUnit: "",
				Organization:       "Bitrise",
				CityOrLocality:     "",
				StateOrProvince:    "",
				CountryCode:        "",
				ValidFrom:          "2024-11-18 14:43:38 +0000 UTC",
				ValidUntil:         "2025-11-18 14:43:38 +0000 UTC",
			},
		},
		{
			name:               "Invalid file",
			pth:                filepath.Join("testdata", "empty_file"),
			password:           "keystore",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "keystore",
			wantError:          "failed to decode keystore:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.pth)
			require.NoError(t, err)
			defer func() {
				err := f.Close()
				require.NoError(t, err)
			}()

			b, err := io.ReadAll(f)
			require.NoError(t, err)

			parser := NewDefaultReader()
			got, err := parser.ReadCertificateInformation(b, tt.password, tt.privateKeyAlias, tt.privateKeyPassword)
			if tt.wantError != "" {
				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), tt.wantError))
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIncorrectKeystoreCredentials(t *testing.T) {
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
			pth:                filepath.Join("testdata", "keystore.pkcs12"),
			password:           "incorrect-password",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypass",
			wantError:          IncorrectKeystorePasswordError.Error(),
		},
		{
			name:               "PKCS12 keystore test - incorrect alias",
			pth:                filepath.Join("testdata", "keystore.pkcs12"),
			password:           "storepass",
			privateKeyAlias:    "incorrect-alias",
			privateKeyPassword: "keypass",
			wantError:          IncorrectAliasError.Error(),
		},
		{
			name:               "PKCS12 keystore test - incorrect key password",
			pth:                filepath.Join("testdata", "keystore.pkcs12"),
			password:           "storepass",
			privateKeyAlias:    "key0",
			privateKeyPassword: "incorrect-keypassword",
			wantError:          IncorrectKeyPasswordError.Error(),
		},
		{
			name:               "JKS keystore test - incorrect password",
			pth:                filepath.Join("testdata", "keystore.jks"),
			password:           "incorrect-password",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "keystore",
			wantError:          IncorrectKeystorePasswordError.Error(),
		},
		{
			name:               "JKS keystore test - incorrect alias",
			pth:                filepath.Join("testdata", "keystore.jks"),
			password:           "keystore",
			privateKeyAlias:    "incorrect-alias",
			privateKeyPassword: "keystore",
			wantError:          IncorrectAliasError.Error(),
		},
		{
			name:               "JKS keystore test - incorrect key password",
			pth:                filepath.Join("testdata", "keystore.jks"),
			password:           "keystore",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "incorrect-keypassword",
			wantError:          IncorrectKeyPasswordError.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.pth)
			require.NoError(t, err)
			defer func() {
				err := f.Close()
				require.NoError(t, err)
			}()

			b, err := io.ReadAll(f)
			require.NoError(t, err)

			parser := NewDefaultReader()
			got, err := parser.ReadCertificateInformation(b, tt.password, tt.privateKeyAlias, tt.privateKeyPassword)
			require.EqualError(t, err, tt.wantError)
			require.Nil(t, got)
		})
	}
}

func TestIsInvalidCredentialsError(t *testing.T) {
	tests := []struct {
		name               string
		decoder            Decoder
		pth                string
		password           string
		privateKeyAlias    string
		privateKeyPassword string
		wantError          string
	}{
		{
			name:               "PKCS12 keystore, JKS decoder",
			decoder:            JKSKeystoreDecoder{},
			pth:                filepath.Join("testdata", "keystore.pkcs12"),
			password:           "storepass",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypass",
			wantError:          IncorrectKeystorePasswordError.Error(),
		},
		{
			name:               "JKS keystore, PKCS12 decoder",
			decoder:            PKCS12KeystoreDecoder{},
			pth:                filepath.Join("testdata", "keystore.jks"),
			password:           "keystore",
			privateKeyAlias:    "mykey",
			privateKeyPassword: "keystore",
			wantError:          IncorrectAliasError.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.pth)
			require.NoError(t, err)
			defer func() {
				err := f.Close()
				require.NoError(t, err)
			}()

			b, err := io.ReadAll(f)
			require.NoError(t, err)

			_, _, err = tt.decoder.Decode(b, tt.password, tt.privateKeyAlias, tt.privateKeyPassword)
			require.Error(t, err)

			wrongType := tt.decoder.IsInvalidCredentialsError(err)
			require.False(t, wrongType)
		})
	}
}
