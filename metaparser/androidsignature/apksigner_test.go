package androidsignature

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Outputs of `apksigner verify --print-certs -v` captured from build-tools 35.0.0 and 37.0.0 (36.x prints like 35.0.0).

// Rotated from k1 to k2 with `apksigner rotate` and signed with `--lineage --rotation-min-sdk-version 33`.
const apkSignerOutputKeyRotation35 = `Verifies
Verified using v1 scheme (JAR signing): true
Verified using v2 scheme (APK Signature Scheme v2): true
Verified using v3 scheme (APK Signature Scheme v3): true
Verified using v3.1 scheme (APK Signature Scheme v3.1): true
Verified using v4 scheme (APK Signature Scheme v4): false
Verified for SourceStamp: false
Number of signers: 1
Signer (minSdkVersion=33, maxSdkVersion=2147483647) certificate DN: CN=k2, O=Test
Signer (minSdkVersion=33, maxSdkVersion=2147483647) certificate SHA-256 digest: eb3b2fb1d9d32f27222ae12767db7ca06e589bf896ded795ce659761b93ce383
Signer (minSdkVersion=24, maxSdkVersion=32) certificate DN: CN=k1, O=Test
Signer (minSdkVersion=24, maxSdkVersion=32) certificate SHA-256 digest: 3959a0bc7d2d3143da2565ef1aa7315499ef94ab26da14b7feb259f708b5a93f
`

const apkSignerOutputKeyRotation37 = `Verifies
Verified using v1 scheme (JAR signing): true
Verified using v2 scheme (APK Signature Scheme v2): true
Verified using v3 scheme (APK Signature Scheme v3): true
Verified using v3.1 scheme (APK Signature Scheme v3.1): true
Verified using v3.2 scheme (APK Signature Scheme v3.2): false
Verified using v4 scheme (APK Signature Scheme v4): false
Verified for SourceStamp: false
Number of signers: 1
V3.1 Signer: (minSdkVersion=33, maxSdkVersion=2147483647) certificate DN: CN=k2, O=Test
V3.1 Signer: (minSdkVersion=33, maxSdkVersion=2147483647) certificate SHA-256 digest: eb3b2fb1d9d32f27222ae12767db7ca06e589bf896ded795ce659761b93ce383
V3.0 Signer: (minSdkVersion=24, maxSdkVersion=32) certificate DN: CN=k1, O=Test
V3.0 Signer: (minSdkVersion=24, maxSdkVersion=32) certificate SHA-256 digest: 3959a0bc7d2d3143da2565ef1aa7315499ef94ab26da14b7feb259f708b5a93f
`

// Signed by k1 and k2 (v1 and v2 schemes).
const apkSignerOutputTwoSigners37 = `Verifies
Verified using v1 scheme (JAR signing): true
Verified using v2 scheme (APK Signature Scheme v2): true
Verified using v3 scheme (APK Signature Scheme v3): false
Verified using v3.1 scheme (APK Signature Scheme v3.1): false
Verified using v3.2 scheme (APK Signature Scheme v3.2): false
Verified using v4 scheme (APK Signature Scheme v4): false
Verified for SourceStamp: false
Number of signers: 2
V2 Signer #1: certificate DN: CN=k1, O=Test
V2 Signer #1: certificate SHA-256 digest: 3959a0bc7d2d3143da2565ef1aa7315499ef94ab26da14b7feb259f708b5a93f
V2 Signer #2: certificate DN: CN=k2, O=Test
V2 Signer #2: certificate SHA-256 digest: eb3b2fb1d9d32f27222ae12767db7ca06e589bf896ded795ce659761b93ce383
`

const apkSignerOutputUnsigned = `DOES NOT VERIFY
ERROR: Missing META-INF/MANIFEST.MF
`

const apkSignerOutputMalformedAPK = `Exception in thread "main" com.android.apksig.apk.ApkFormatException: Malformed APK: not a ZIP archive
	at com.android.apksig.ApkVerifier.verify(ApkVerifier.java:191)
	at com.android.apksigner.ApkSignerTool.main(ApkSignerTool.java:97)
Caused by: com.android.apksig.zip.ZipFormatException: ZIP End of Central Directory record not found
	at com.android.apksig.apk.ApkUtilsLite.findZipSections(ApkUtilsLite.java:49)
	... 3 more
`

func TestReadAPKSignature_apkSignerOutput(t *testing.T) {
	tests := []struct {
		name              string
		buildToolsVersion string
		apkSignerOutput   string
		apkSignerExitCode int
		wantSignature     string
		wantErrContains   []string
	}{
		{
			name:              "key rotation (build-tools 35)",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputKeyRotation35,
			wantSignature:     "CN=k2, O=Test",
		},
		{
			name:              "key rotation (build-tools 37)",
			buildToolsVersion: "37.0.0",
			apkSignerOutput:   apkSignerOutputKeyRotation37,
			wantSignature:     "CN=k2, O=Test",
		},
		{
			name:              "two signers (build-tools 37)",
			buildToolsVersion: "37.0.0",
			apkSignerOutput:   apkSignerOutputTwoSigners37,
			wantSignature:     "CN=k1, O=Test",
		},
		{
			name:              "unsigned APK keeps apksigner's reason",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputUnsigned,
			apkSignerExitCode: 1,
			wantErrContains:   []string{"no signature found (apksigner: not verified: ERROR: Missing META-INF/MANIFEST.MF)"},
		},
		{
			name:              "apksigner crash surfaces the exception",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputMalformedAPK,
			apkSignerExitCode: 1,
			wantErrContains: []string{
				"command failed with exit status 1",
				`: Exception in thread "main" com.android.apksig.apk.ApkFormatException: Malformed APK: not a ZIP archive; ` +
					"Caused by: com.android.apksig.zip.ZipFormatException: ZIP End of Central Directory record not found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAndroidSDK(t, tt.buildToolsVersion, tt.apkSignerOutput, tt.apkSignerExitCode)
			fakeJarsigner(t, "jar is unsigned.", 0)

			gotSignature, gotError := ReadAPKSignature("app.apk")

			require.Equal(t, tt.wantSignature, gotSignature)
			if len(tt.wantErrContains) == 0 {
				require.NoError(t, gotError)
				return
			}
			require.Error(t, gotError)
			for _, want := range tt.wantErrContains {
				require.Contains(t, gotError.Error(), want)
			}
			require.NotContains(t, gotError.Error(), "check the command's output for details")
		})
	}
}

func TestReadAABSignature_jarSignerCrash(t *testing.T) {
	fakeJarsigner(t, "jarsigner: java.util.zip.ZipException: zip END header not found\n", 1)

	gotSignature, gotError := ReadAABSignature("app.aab")

	require.Empty(t, gotSignature)
	require.Error(t, gotError)
	require.Contains(t, gotError.Error(), "command failed with exit status 1")
	require.Contains(t, gotError.Error(), ": jarsigner: java.util.zip.ZipException: zip END header not found")
	require.NotContains(t, gotError.Error(), "check the command's output for details")
}

func fakeAndroidSDK(t *testing.T, buildToolsVersion, output string, exitCode int) {
	t.Helper()

	sdkRoot := t.TempDir()
	buildToolsDir := filepath.Join(sdkRoot, "build-tools", buildToolsVersion)
	require.NoError(t, os.MkdirAll(buildToolsDir, 0o700))
	writeFakeTool(t, filepath.Join(buildToolsDir, "apksigner"), output, exitCode)
	t.Setenv("ANDROID_HOME", sdkRoot)
}

func fakeJarsigner(t *testing.T, output string, exitCode int) {
	t.Helper()

	dir := t.TempDir()
	writeFakeTool(t, filepath.Join(dir, "jarsigner"), output, exitCode)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func writeFakeTool(t *testing.T, path, output string, exitCode int) {
	t.Helper()

	outputPath := path + ".output"
	require.NoError(t, os.WriteFile(outputPath, []byte(output), 0o600))
	script := fmt.Sprintf("#!/bin/sh\ncat %q\nexit %d\n", outputPath, exitCode)
	require.NoError(t, os.WriteFile(path, []byte(script), 0o700))
}
