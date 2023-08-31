package keystore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: finalise tests
func TestOpen(t *testing.T) {
	tests := []struct {
		name               string
		pth                string
		password           string
		privateKeyAlias    string
		privateKeyPassword string
		want               KeyStore
	}{
		{
			name:               "PKCS12 keystore test",
			pth:                "/Users/godrei/Development/keystores/keystore_2.jks",
			password:           "password",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keypassword",
			want:               KeyStore{},
		},
		{
			name:               "JKS keystore test",
			pth:                "/Users/godrei/Downloads/ANDROID_SIGN_DEFAULT_ALIAS_KEYSTORE.keystore",
			password:           "keystore",
			privateKeyAlias:    "key0",
			privateKeyPassword: "keystore",
			want:               KeyStore{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.pth, tt.password, tt.privateKeyAlias, tt.privateKeyPassword)
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
