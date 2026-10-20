package androidsignature

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-android/v2/sdk"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

var cmdFactory = command.NewFactory(env.NewRepository())

const (
	unsignedJarSignatureMessage       = "jar is unsigned"
	validJarSignatureMessage          = "jar verified"
	validV2PlusSignatureMessage       = "Verifies"
	notVerifiedV2PlusSignatureMessage = "DOES NOT VERIFY"

	maxFailureReasonLines = 5
)

var (
	ErrNotVerified      = errors.New("not verified")
	ErrNoSignatureFound = errors.New("no signature found")
)

// apksigner's signer label depends on the build-tools version, the number of signers and key rotation,
// e.g. "Signer #1", "V2 Signer #1:" or "V3.1 Signer: (minSdkVersion=33, maxSdkVersion=2147483647)".
var apkSignerDNRegex = regexp.MustCompile(`(?m)^(?:V\d+(?:\.\d+)?(?: \w+)* )?Signer(?: #1)?:?(?: \(minSdkVersion=.*?\))? certificate DN: (.*)`)

// Read ...
//
// Deprecated: Read is deprecated. Use ReadAABSignature or ReadAPKSignature method instead.
func Read(path string) (string, error) {
	return ReadAABSignature(path)
}

// ReadAABSignature returns the signature of the provided AAB file.
// If the signature can't be read (unsigned, unexpected certificate printing format, ...), it returns a ErrNoSignatureFound.
// If the signature is not verified, it returns a ErrNotVerified.
func ReadAABSignature(path string) (string, error) {
	return getJarSignature(path)
}

// ReadAPKSignature returns the signature of the provided APK file.
// If the signature can't be read (unsigned, unexpected certificate printing format, ...), it returns a ErrNoSignatureFound.
// If the signature is not verified, it returns a ErrNotVerified.
func ReadAPKSignature(apkPath string) (string, error) {
	idSigPath := apkPath + ".idsig"
	if _, err := os.Stat(idSigPath); err == nil {
		signature, err := getV4Signature(apkPath, idSigPath)
		if err != nil && !errors.Is(err, ErrNotVerified) && !errors.Is(err, ErrNoSignatureFound) {
			return "", err
		}
		if signature != "" {
			return signature, nil
		}
	}

	signature, apkSignerErr := getV23Signature(apkPath)
	if apkSignerErr != nil && !errors.Is(apkSignerErr, ErrNotVerified) && !errors.Is(apkSignerErr, ErrNoSignatureFound) {
		return "", apkSignerErr
	}
	if signature != "" {
		return signature, nil
	}

	signature, err := getJarSignature(apkPath)
	if err != nil && apkSignerErr != nil {
		// Only apksigner's verdict tells an unsigned APK apart from one signed without v1 below minSdkVersion 24.
		return "", fmt.Errorf("%w (apksigner: %w)", err, apkSignerErr)
	}

	return signature, err
}

func getV4Signature(apkPath string, idsigPath string) (string, error) {
	if _, err := os.Stat(idsigPath); err != nil {
		return "", fmt.Errorf("failed to check if detached signature file (.idsig) exist: %s", err)
	}

	pathParams := []string{"-v4-signature-file", idsigPath, apkPath}
	return getV2PlusSignature(pathParams)
}

func getV23Signature(path string) (string, error) {
	pathParams := []string{path}
	return getV2PlusSignature(pathParams)
}

func getV2PlusSignature(pathParams []string) (string, error) {
	sdkModel, err := sdk.NewDefaultModel(sdk.Environment{
		AndroidHome:    os.Getenv("ANDROID_HOME"),
		AndroidSDKRoot: os.Getenv("ANDROID_SDK_ROOT"),
	}, pathutil.NewPathChecker())
	if err != nil {
		return "", fmt.Errorf("failed to create sdk model, error: %s", err)
	}

	apkSignerPath, err := sdkModel.LatestBuildToolPath("apksigner")
	if err != nil {
		return "", fmt.Errorf("failed to find latest apksigner binary, error: %s", err)
	}

	params := append([]string{"verify", "--print-certs", "-v"}, pathParams...)
	opts := &command.Opts{ErrorFinder: errorFinder}
	apkSignerOutput, err := cmdFactory.Create(apkSignerPath, params, opts).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && strings.Contains(apkSignerOutput, notVerifiedV2PlusSignatureMessage) {
			if reason := failureReason(apkSignerOutput); reason != "" {
				return "", fmt.Errorf("%w: %s", ErrNotVerified, reason)
			}
			return "", ErrNotVerified
		}
		// apksigner also exits with 1 when it crashes, e.g. on a malformed APK.
		return "", err
	}

	if !strings.Contains(apkSignerOutput, validV2PlusSignatureMessage) {
		return "", ErrNotVerified
	}

	if match := apkSignerDNRegex.FindStringSubmatch(apkSignerOutput); match != nil {
		return match[1], nil
	}

	return "", ErrNoSignatureFound
}

// Without an ErrorFinder the command error only says "check the command's output for details".
func errorFinder(output string) []string {
	if reason := failureReason(output); reason != "" {
		return []string{reason}
	}

	return nil
}

func failureReason(output string) string {
	var lines []string

	for _, line := range strings.Split(output, "\n") {
		// Stack frames are indented.
		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" || line == notVerifiedV2PlusSignatureMessage || strings.HasPrefix(line, "WARNING:") {
			continue
		}

		lines = append(lines, line)
	}

	if len(lines) > maxFailureReasonLines {
		omitted := len(lines) - maxFailureReasonLines
		lines = append(lines[:maxFailureReasonLines], fmt.Sprintf("(%d more lines)", omitted))
	}

	return strings.Join(lines, "; ")
}

func getJarSignature(path string) (string, error) {
	params := []string{"-verify", "-certs", "-verbose", path}
	opts := &command.Opts{ErrorFinder: errorFinder}
	output, err := cmdFactory.Create("jarsigner", params, opts).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", err
	}

	if strings.Contains(output, unsignedJarSignatureMessage) {
		return "", ErrNoSignatureFound
	}

	if !strings.Contains(output, validJarSignatureMessage) {
		return "", ErrNotVerified
	}

	var signature string

	// The signature details appear in the output in the following format:
	// - Signed by "C=Aa, ST=Bbbbb, L=Ccccc, O=Ddddd, OU=Eeeee, CN=Fffff"
	regex := regexp.MustCompile(`- Signed by ".*"`)
	sig := regex.FindString(output)
	if sig != "" {
		signature = strings.TrimPrefix(sig, "- Signed by \"")
		signature = strings.TrimSuffix(signature, "\"")
		return signature, nil
	}

	return "", ErrNoSignatureFound
}
