package androidsignature

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-android/v2/sdk"
	"github.com/bitrise-io/go-utils/command"
)

const (
	validJarSignatureMessage    = "jar verified"
	validV2PlusSignatureMessage = "Verifies"
)

var (
	NotVerifiedError = errors.New("not verified")
)

// Read ...
//
// Deprecated: Read is deprecated. Use ReadAABSignature method instead.
func Read(path string) (string, error) {
	return ReadAABSignature(path)
}

// ReadAABSignature ...
func ReadAABSignature(path string) (string, error) {
	return getJarSignature(path)
}

// ReadAPKSignature ...
func ReadAPKSignature(path string) (string, []error) {
	var signature string
	var errs []error
	signature, err := getV4Signature(path)
	if err != nil {
		errs = append(errs, fmt.Errorf("v4 signature: %+v", err))
	}
	if signature != "" {
		return signature, errs
	}

	signature, err = getV23Signature(path)
	if err != nil {
		errs = append(errs, fmt.Errorf("v2/v3 signature: %+v", err))
	}
	if signature != "" {
		return signature, errs
	}

	signature, err = getJarSignature(path)
	if err != nil {
		errs = append(errs, fmt.Errorf("jar (v1) signature: %+v", err))
	}

	return signature, errs
}

func getV4Signature(path string) (string, error) {
	idSigPath := path + ".idsig"

	if _, err := os.Stat(idSigPath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("failed to find the v4 signature file (*.idsig), error: %s", err)
	}

	pathParams := []string{"-v4-signature-file", idSigPath, path}
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
	})
	if err != nil {
		return "", fmt.Errorf("failed to create sdk model, error: %s", err)
	}

	apkSignerPath, err := sdkModel.LatestBuildToolPath("apksigner")
	if err != nil {
		return "", fmt.Errorf("failed to find latest aapt binary, error: %s", err)
	}

	params := append([]string{"verify", "--print-certs", "-v"}, pathParams...)
	apkSignerOutput, err := command.New(apkSignerPath, params...).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1" {
			regex := regexp.MustCompile(`DOES NOT VERIFY`)
			sig := regex.FindString(apkSignerOutput)
			if sig != "" {
				return "", NotVerifiedError
			}
		}
		return "", err
	}

	var signature string

	if strings.Contains(apkSignerOutput, validV2PlusSignatureMessage) {
		// The signature details appear in the output in the following format:
		// Signer #1 certificate DN: C=Aa, ST=Bbbbb, L=Ccccc, O=Ddddd, OU=Eeeee, CN=Fffff
		// Signer #1 certificate SHA-256 digest: <hash>
		// Signer #1 certificate SHA-1 digest: <hash>
		// Signer #1 certificate MD5 digest: <hash>
		regex := regexp.MustCompile("Signer #1 certificate DN: (.*)")
		res := regex.FindAllStringSubmatch(apkSignerOutput, 1)
		if len(res) > 0 && len(res[0]) > 1 {
			return res[0][1], nil
		}
	} else {
		err = NotVerifiedError
	}

	return signature, err
}

func getJarSignature(path string) (string, error) {
	params := []string{"-verify", "-certs", "-verbose", path}
	output, err := command.New("jarsigner", params...).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", err
	}

	var signature string

	if strings.Contains(output, validJarSignatureMessage) {
		// The signature details appear in the output in the following format:
		// - Signed by "C=Aa, ST=Bbbbb, L=Ccccc, O=Ddddd, OU=Eeeee, CN=Fffff"
		regex := regexp.MustCompile(`- Signed by ".*"`)
		sig := regex.FindString(output)
		if sig != "" {
			signature = strings.TrimPrefix(sig, "- Signed by \"")
			signature = strings.TrimSuffix(signature, "\"")
		}
	} else {
		err = NotVerifiedError
	}

	return signature, err
}
