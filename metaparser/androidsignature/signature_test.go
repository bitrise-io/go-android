package androidsignature

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/go-utils/command/git"
)

func Test_ReadAABSignature(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir, error: %s", err)
		}
	}()

	gitCommand, err := git.New(tmpDir)
	require.NoError(t, err)
	err = gitCommand.Clone("https://github.com/bitrise-io/sample-artifacts.git").Run()
	require.NoError(t, err)

	tests := []struct {
		name          string
		apkPath       string
		wantSignature string
	}{
		{
			name:          "Reads AAB signature",
			apkPath:       path.Join(tmpDir, "app-bitrise-signed.aab"),
			wantSignature: "CN=Bitrise",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, gotError := ReadAABSignature(tt.apkPath)
			require.NoError(t, gotError)
			require.Equal(t, tt.wantSignature, gotSignature)
		})
	}
}

func Test_ReadASignature(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir, error: %s", err)
		}
	}()

	gitCommand, err := git.New(tmpDir)
	require.NoError(t, err)
	err = gitCommand.Clone("https://github.com/bitrise-io/sample-artifacts.git").Run()
	require.NoError(t, err)

	tests := []struct {
		name          string
		apkPath       string
		idsigPath     string
		wantSignature string
		wantError     string
	}{
		{
			name:      "Unsigned APK",
			apkPath:   path.Join(tmpDir, "apks", "app-release-unsigned.apk"),
			wantError: NotVerifiedError.Error(),
		},
		{
			name:          "Debug signed APK",
			apkPath:       path.Join(tmpDir, "apks", "app-debug.apk"),
			wantSignature: "C=US, O=Android, CN=Android Debug",
		},
		{
			name:          "Reads V1 signature",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v1-signature.apk"),
			wantSignature: "O=Bitrise",
		},
		{
			name:          "Reads V2 signature",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v2-signature.apk"),
			wantSignature: "O=Bitrise",
		},
		{
			name:          "Reads V2+v3 signature",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v2-v3-signature.apk"),
			wantSignature: "O=Bitrise",
		},
		{
			name:          "Reads V2+v4 signature",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v2-v4-signature.apk"),
			wantSignature: "O=Bitrise",
		},
		{
			name:          "Reads V2+v4 signature with .idsig file",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v2-v4-signature.apk"),
			idsigPath:     path.Join(tmpDir, "apks", "app-release-v2-v4-signature.apk.idsig"),
			wantSignature: "O=Bitrise",
		},
		{
			name:          "Reads V3+v4 signature",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v3-v4-signature.apk"),
			wantSignature: "O=Bitrise",
		},
		{
			name:          "Reads V3+v4 signature with .idsig file",
			apkPath:       path.Join(tmpDir, "apks", "app-release-v3-v4-signature.apk"),
			idsigPath:     path.Join(tmpDir, "apks", "app-release-v3-v4-signature.apk.idsig"),
			wantSignature: "O=Bitrise",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, gotError := ReadAPKSignature(tt.apkPath)
			if tt.wantError != "" {
				require.EqualError(t, gotError, tt.wantError)
			} else {
				require.NoError(t, gotError)
			}
			require.Equal(t, tt.wantSignature, gotSignature)
		})
	}
}
