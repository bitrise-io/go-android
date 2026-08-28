package gradle

import (
	"path/filepath"

	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-io/go-utils/v2/ziputil"
)

// Artifact ...
type Artifact struct {
	Path string
	Name string
}

// Export ...
func (artifact Artifact) Export(destination string) error {
	return fileutil.NewFileManager().CopyFile(artifact.Path, filepath.Join(destination, artifact.Name), &fileutil.CopyOptions{Overwrite: true})
}

// ExportZIP ...
func (artifact Artifact) ExportZIP(destination string) error {
	zipManager := ziputil.NewZipManager(pathutil.NewPathChecker())
	return zipManager.ZipDir(artifact.Path, filepath.Join(destination, artifact.Name), true)
}
