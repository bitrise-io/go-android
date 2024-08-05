package metaparser

import (
	"github.com/bitrise-io/go-android/v2/metaparser/androidartifact"
	"github.com/bitrise-io/go-android/v2/metaparser/androidsignature"
)

// ParseAABData ...
func (m *metaparser) ParseAABData(pth string) (*ArtifactMetadata, error) {
	aabInfo, err := androidartifact.GetAABInfo(m.BundletoolPath, pth)
	if err != nil {
		m.Logger.Warnf("Failed to parse AAB info: %s", err)
		m.Logger.AABParseWarnf("aab-parse", "aabparser package failed to parse AAB, error: %s", err)
		return nil, err
	}

	fileSize, err := m.FileManager.FileSizeInBytes(pth)
	if err != nil {
		m.Logger.Warnf("Failed to get apk size, error: %s", err)
	}

	info := androidartifact.ParseArtifactPath(pth)

	signature, err := androidsignature.Read(pth)
	if err != nil {
		m.Logger.Warnf("Failed to read signature: %s", err)
	}

	return &ArtifactMetadata{
		AppInfo:        aabInfo,
		FileSizeBytes:  fileSize,
		Module:         info.Module,
		ProductFlavour: info.ProductFlavour,
		BuildType:      info.BuildType,
		SignedBy:       signature,
	}, nil
}
