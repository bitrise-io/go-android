package parser

import "github.com/bitrise-io/go-android/v2/metaparser/androidartifact"

type ArtifactMetadata struct {
	AppInfo        androidartifact.Info
	FileSizeBytes  int64
	Module         string
	ProductFlavour string
	BuildType      string
	SignedBy       string
	Warnings       []string
}
