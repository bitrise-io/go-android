package metaparser

import (
	"fmt"

	"github.com/bitrise-io/go-android/v2/metaparser/androidartifact"
	"github.com/bitrise-io/go-android/v2/metaparser/androidsignature"
)

// ParseAPKData ...
func (m *Parser) ParseAPKData(pth string) (*ArtifactMetadata, error) {
	apkInfo, err := androidartifact.GetAPKInfoWithFallback(m.logger, pth)
	if err != nil {
		return nil, err
	}

	fileSize, err := m.fileManager.FileSizeInBytes(pth)
	if err != nil {
		m.logger.Warnf("Failed to get apk size, error: %s", err)
	}

	info := androidartifact.ParseArtifactPath(pth)

	signature, errs := androidsignature.ReadAPKSignature(pth)
	if len(errs) > 0 {
		format := "Failed to read signature of `%s`: \n%s"
		if signature != "" {
			format = "Errors during reading the signature of `%s`: \n%s"
		}
		var output string
		for _, err := range errs {
			output += fmt.Sprintf("- %s\n", err.Error())
		}
		m.logger.Warnf(format, pth, output)
	}

	return &ArtifactMetadata{
		AppInfo:        apkInfo,
		FileSizeBytes:  fileSize,
		Module:         info.Module,
		ProductFlavour: info.ProductFlavour,
		BuildType:      info.BuildType,
		SignedBy:       signature,
	}, nil
}
