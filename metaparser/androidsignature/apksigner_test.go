package androidsignature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// The apksigner outputs below were captured from `apksigner verify --print-certs -v` of build-tools 35.0.0 and 37.0.0.
// Build-tools 36.0.0 and 36.1.0 print the same as 35.0.0.

const apkSignerOutputSigned35 = `Verifies
Verified using v1 scheme (JAR signing): true
Verified using v2 scheme (APK Signature Scheme v2): true
Verified using v3 scheme (APK Signature Scheme v3): true
Verified using v3.1 scheme (APK Signature Scheme v3.1): false
Verified using v4 scheme (APK Signature Scheme v4): false
Verified for SourceStamp: false
Number of signers: 1
Signer #1 certificate DN: CN=k1, O=Test
Signer #1 certificate SHA-256 digest: 3959a0bc7d2d3143da2565ef1aa7315499ef94ab26da14b7feb259f708b5a93f
Signer #1 certificate SHA-1 digest: 7e79d463ccbcece9b4cec1c732085254e950bf6f
Signer #1 certificate MD5 digest: dd1829b733e0ffc685ec539105f858b1
`

const apkSignerOutputSigned37 = `Verifies
Verified using v1 scheme (JAR signing): true
Verified using v2 scheme (APK Signature Scheme v2): true
Verified using v3 scheme (APK Signature Scheme v3): true
Verified using v3.1 scheme (APK Signature Scheme v3.1): false
Verified using v3.2 scheme (APK Signature Scheme v3.2): false
Verified using v4 scheme (APK Signature Scheme v4): false
Verified for SourceStamp: false
Number of signers: 1
V3.0 Signer: certificate DN: CN=k1, O=Test
V3.0 Signer: certificate SHA-256 digest: 3959a0bc7d2d3143da2565ef1aa7315499ef94ab26da14b7feb259f708b5a93f
V3.0 Signer: certificate SHA-1 digest: 7e79d463ccbcece9b4cec1c732085254e950bf6f
V3.0 Signer: certificate MD5 digest: dd1829b733e0ffc685ec539105f858b1
`

// An APK rotated from k1 to k2 with `apksigner rotate` and signed with `--lineage --rotation-min-sdk-version 33`.
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

// An unsigned APK, or one without a v1 signature while its minSdkVersion is below 24.
const apkSignerOutputNotVerified = `DOES NOT VERIFY
ERROR: Missing META-INF/MANIFEST.MF
`

// An APK with only a v1 signature while targeting SDK 35 (WARNING lines shortened).
const apkSignerOutputNotVerifiedWithWarnings = `DOES NOT VERIFY
ERROR: Target SDK version 35 requires a minimum of signature scheme v2; the APK is not signed with this or a later signature scheme
WARNING: META-INF/com/android/build/gradle/app-metadata.properties not protected by signature. Unauthorized modifications to this JAR entry will not be detected. Delete or move the entry outside of META-INF/.
WARNING: META-INF/version-control-info.textproto not protected by signature. Unauthorized modifications to this JAR entry will not be detected. Delete or move the entry outside of META-INF/.
`

// A truncated APK.
const apkSignerOutputMalformedAPK = `Exception in thread "main" com.android.apksig.apk.ApkFormatException: Malformed APK: not a ZIP archive
	at com.android.apksig.ApkVerifier.verify(ApkVerifier.java:191)
	at com.android.apksig.ApkVerifier.verify(ApkVerifier.java:164)
	at com.android.apksigner.ApkSignerTool.verify(ApkSignerTool.java:602)
	at com.android.apksigner.ApkSignerTool.main(ApkSignerTool.java:97)
Caused by: com.android.apksig.zip.ZipFormatException: ZIP End of Central Directory record not found
	at com.android.apksig.apk.ApkUtilsLite.findZipSections(ApkUtilsLite.java:49)
	at com.android.apksig.apk.ApkUtils.findZipSections(ApkUtils.java:60)
	at com.android.apksig.ApkVerifier.verify(ApkVerifier.java:189)
	... 3 more
`

const jarSignerOutputSigned = `s = signature was verified
m = entry is listed in manifest
k = at least one certificate was found in keystore

- Signed by "CN=k1, O=Test"
    Digest algorithm: SHA-256
    Signature algorithm: SHA256withRSA, 2048-bit key

jar verified.
`

const jarSignerOutputUnsigned = `s = signature was verified
m = entry is listed in manifest
k = at least one certificate was found in keystore

jar is unsigned.
`

// writeFakeTool writes an executable at path that prints output and exits with exitCode.
// It returns whether the executable has been run.
func writeFakeTool(t *testing.T, path, output string, exitCode int) (ran func() bool) {
	t.Helper()

	outputPath := path + ".output"
	require.NoError(t, os.WriteFile(outputPath, []byte(output), 0o600))

	ranPath := path + ".ran"
	script := fmt.Sprintf("#!/bin/sh\ntouch %q\ncat %q\nexit %d\n", ranPath, outputPath, exitCode)
	require.NoError(t, os.WriteFile(path, []byte(script), 0o700))

	return func() bool {
		_, err := os.Stat(ranPath)
		return err == nil
	}
}

// fakeAndroidSDK points ANDROID_HOME at an SDK whose only apksigner (in build-tools/<buildToolsVersion>)
// prints output and exits with exitCode, so that ReadAPKSignature runs it through the real command runner.
func fakeAndroidSDK(t *testing.T, buildToolsVersion, output string, exitCode int) {
	t.Helper()

	sdkRoot := t.TempDir()
	buildToolsDir := filepath.Join(sdkRoot, "build-tools", buildToolsVersion)
	require.NoError(t, os.MkdirAll(buildToolsDir, 0o700))
	writeFakeTool(t, filepath.Join(buildToolsDir, "apksigner"), output, exitCode)

	t.Setenv("ANDROID_HOME", sdkRoot)
}

// fakeJarsigner puts a jarsigner that prints output first on the PATH and returns whether it has been run.
func fakeJarsigner(t *testing.T, output string) (ran func() bool) {
	t.Helper()

	dir := t.TempDir()
	ran = writeFakeTool(t, filepath.Join(dir, "jarsigner"), output, 0)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	return ran
}

func TestReadAPKSignature_apkSignerOutput(t *testing.T) {
	tests := []struct {
		name              string
		buildToolsVersion string
		apkSignerOutput   string
		apkSignerExitCode int
		jarSignerOutput   string
		wantSignature     string
		wantErrContains   []string
		wantErrIs         []error
		wantJarSigner     bool
	}{
		{
			name:              "signed APK (build-tools 35)",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputSigned35,
			wantSignature:     "CN=k1, O=Test",
		},
		{
			name:              "signed APK (build-tools 37)",
			buildToolsVersion: "37.0.0",
			apkSignerOutput:   apkSignerOutputSigned37,
			wantSignature:     "CN=k1, O=Test",
		},
		{
			name:              "APK signed with key rotation reports the rotated signer (build-tools 35)",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputKeyRotation35,
			wantSignature:     "CN=k2, O=Test",
		},
		{
			name:              "APK signed with key rotation reports the rotated signer (build-tools 37)",
			buildToolsVersion: "37.0.0",
			apkSignerOutput:   apkSignerOutputKeyRotation37,
			wantSignature:     "CN=k2, O=Test",
		},
		{
			name:              "APK that does not verify falls back to its JAR signature",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputNotVerifiedWithWarnings,
			apkSignerExitCode: 1,
			jarSignerOutput:   jarSignerOutputSigned,
			wantSignature:     "CN=k1, O=Test",
			wantJarSigner:     true,
		},
		{
			name:              "unsigned APK keeps apksigner's reason",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputNotVerified,
			apkSignerExitCode: 1,
			jarSignerOutput:   jarSignerOutputUnsigned,
			wantErrContains:   []string{"no signature found (apksigner: not verified: ERROR: Missing META-INF/MANIFEST.MF)"},
			wantErrIs:         []error{ErrNoSignatureFound, ErrNotVerified},
			wantJarSigner:     true,
		},
		{
			name:              "verified APK without a recognised signer falls back to its JAR signature",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   "Verifies\nNumber of signers: 1\n",
			jarSignerOutput:   jarSignerOutputUnsigned,
			wantErrContains:   []string{"no signature found (apksigner: no signature found)"},
			wantErrIs:         []error{ErrNoSignatureFound},
			wantJarSigner:     true,
		},
		{
			name:              "apksigner crash surfaces the exception and does not fall back",
			buildToolsVersion: "35.0.0",
			apkSignerOutput:   apkSignerOutputMalformedAPK,
			apkSignerExitCode: 1,
			wantErrContains: []string{
				"command failed with exit status 1",
				`"verify" "--print-certs" "-v" "app.apk"`,
				`: Exception in thread "main" com.android.apksig.apk.ApkFormatException: Malformed APK: not a ZIP archive; ` +
					"Caused by: com.android.apksig.zip.ZipFormatException: ZIP End of Central Directory record not found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAndroidSDK(t, tt.buildToolsVersion, tt.apkSignerOutput, tt.apkSignerExitCode)
			jarSignerRan := fakeJarsigner(t, tt.jarSignerOutput)

			gotSignature, gotError := ReadAPKSignature("app.apk")

			require.Equal(t, tt.wantSignature, gotSignature)
			require.Equal(t, tt.wantJarSigner, jarSignerRan())

			if len(tt.wantErrContains) == 0 {
				require.NoError(t, gotError)
				return
			}

			require.Error(t, gotError)
			for _, want := range tt.wantErrContains {
				require.Contains(t, gotError.Error(), want)
			}
			for _, want := range tt.wantErrIs {
				require.ErrorIs(t, gotError, want)
			}
			require.NotContains(t, gotError.Error(), "check the command's output for details")
			require.NotContains(t, gotError.Error(), "\tat ")
		})
	}
}

func TestReadAPKSignature_apkSignerNotFound(t *testing.T) {
	t.Setenv("ANDROID_HOME", t.TempDir())
	jarSignerRan := fakeJarsigner(t, jarSignerOutputSigned)

	gotSignature, gotError := ReadAPKSignature("app.apk")

	require.EqualError(t, gotError, "failed to find latest apksigner binary, error: failed to find latest build-tools dir")
	require.Empty(t, gotSignature)
	require.False(t, jarSignerRan())
}

func TestAPKSignerDNRegex(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		wantDN string
	}{
		{
			name:   "build-tools < 37",
			line:   "Signer #1 certificate DN: C=Aa, ST=Bbbbb, L=Ccccc, O=Ddddd, OU=Eeeee, CN=Fffff",
			wantDN: "C=Aa, ST=Bbbbb, L=Ccccc, O=Ddddd, OU=Eeeee, CN=Fffff",
		},
		{
			name:   "build-tools < 37 with key rotation",
			line:   "Signer (minSdkVersion=33, maxSdkVersion=2147483647) certificate DN: CN=k2, O=Test",
			wantDN: "CN=k2, O=Test",
		},
		{
			// apksigner adds this when the rotation targets a development release.
			name:   "build-tools < 37 with key rotation targeting a dev release",
			line:   "Signer (minSdkVersion=36 (dev release=true), maxSdkVersion=2147483647) certificate DN: CN=k2, O=Test",
			wantDN: "CN=k2, O=Test",
		},
		{
			name:   "build-tools >= 37",
			line:   "V3.0 Signer: certificate DN: CN=k1, O=Test",
			wantDN: "CN=k1, O=Test",
		},
		{
			name:   "build-tools >= 37 with key rotation",
			line:   "V3.1 Signer: (minSdkVersion=33, maxSdkVersion=2147483647) certificate DN: CN=k2, O=Test",
			wantDN: "CN=k2, O=Test",
		},
		{
			name: "second signer",
			line: "Signer #2 certificate DN: CN=k2, O=Test",
		},
		{
			name: "source stamp signer",
			line: "Source Stamp Signer certificate DN: CN=stamp, O=Test",
		},
		{
			name: "not a certificate DN line",
			line: "Signer #1 certificate SHA-256 digest: 3959a0bc7d2d3143da2565ef1aa7315499ef94ab26da14b7feb259f708b5a93f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := "Verifies\nNumber of signers: 1\n" + tt.line + "\n"

			match := apkSignerDNRegex.FindStringSubmatch(output)

			if tt.wantDN == "" {
				require.Nil(t, match)
				return
			}
			require.Len(t, match, 2)
			require.Equal(t, tt.wantDN, match[1])
		})
	}
}

func TestAPKSignerFailureReason(t *testing.T) {
	var manyErrors strings.Builder
	manyErrors.WriteString("DOES NOT VERIFY\n")
	for i := 1; i <= 7; i++ {
		fmt.Fprintf(&manyErrors, "ERROR: JAR signer CERT.RSA: entry %d digest mismatch\n", i)
		fmt.Fprintf(&manyErrors, "WARNING: entry %d not protected by signature\n", i)
	}

	tests := []struct {
		name       string
		output     string
		wantReason string
	}{
		{
			name:       "does not verify",
			output:     apkSignerOutputNotVerified,
			wantReason: "ERROR: Missing META-INF/MANIFEST.MF",
		},
		{
			name:       "does not verify with warnings",
			output:     apkSignerOutputNotVerifiedWithWarnings,
			wantReason: "ERROR: Target SDK version 35 requires a minimum of signature scheme v2; the APK is not signed with this or a later signature scheme",
		},
		{
			name:   "crash",
			output: apkSignerOutputMalformedAPK,
			wantReason: `Exception in thread "main" com.android.apksig.apk.ApkFormatException: Malformed APK: not a ZIP archive; ` +
				"Caused by: com.android.apksig.zip.ZipFormatException: ZIP End of Central Directory record not found",
		},
		{
			name:   "more errors than lines kept",
			output: manyErrors.String(),
			wantReason: "ERROR: JAR signer CERT.RSA: entry 1 digest mismatch; " +
				"ERROR: JAR signer CERT.RSA: entry 2 digest mismatch; " +
				"ERROR: JAR signer CERT.RSA: entry 3 digest mismatch; " +
				"ERROR: JAR signer CERT.RSA: entry 4 digest mismatch; " +
				"ERROR: JAR signer CERT.RSA: entry 5 digest mismatch; " +
				"(2 more lines)",
		},
		{
			name:   "no output",
			output: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantReason, apkSignerFailureReason(tt.output))
		})
	}
}
