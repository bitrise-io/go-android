package metaparser

import (
	"github.com/bitrise-io/go-android/v2/metaparser/androidartifact"
	"github.com/bitrise-io/go-android/v2/metaparser/bundletool"
	"github.com/bitrise-io/go-utils/v2/fileutil"
)

type ArtifactMetadata struct {
	AppInfo        androidartifact.Info `json:"app_info"`
	FileSizeBytes  int64                `json:"file_size_bytes"`
	Module         string               `json:"module"`
	ProductFlavour string               `json:"product_flavour"`
	BuildType      string               `json:"build_type"`
	SignedBy       string               `json:"signed_by"`
	androidartifact.Artifact
}

type MetaParser interface {
	ParseAABData(pth string) (*ArtifactMetadata, error)
	ParseAPKData(pth string) (*ArtifactMetadata, error)
}

type metaparser struct {
	Logger         androidartifact.Logger
	BundletoolPath bundletool.Path
	FileManager    fileutil.FileManager
}

// NewMetaParser ...
func NewMetaParser(logger androidartifact.Logger, bundletoolPath bundletool.Path) MetaParser {
	return &metaparser{
		Logger:         logger,
		BundletoolPath: bundletoolPath,
		FileManager:    fileutil.NewFileManager(),
	}
}
