package androidsignature

import (
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/kr/pretty"

	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/log"
)

func Test_ReadAPKSignature(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("setup: failed to create temp dir, error: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Warnf("failed to remove temp dir, error: %s", err)
		}
	}()

	gitCommand, err := git.New(tmpDir)
	if err != nil {
		t.Fatalf("setup: failed to create git project, error: %s", err)
	}
	if err := gitCommand.Clone("https://github.com/bitrise-io/sample-artifacts.git").Run(); err != nil {
		t.Fatalf("setup: failed to clone test artifact repo, error: %s", err)
	}

	tests := []struct {
		name            string
		path            string
		wantSignature   string
		wantErrorsCount int
	}{
		{
			name:            "Return signature by jarsigner",
			path:            path.Join(tmpDir, "apks", "app-release-v1-signature.apk"),
			wantSignature:   "O=Bitrise",
			wantErrorsCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, gotErrors := ReadAPKSignature(tt.path)
			if diffs := pretty.Diff(gotSignature, tt.wantSignature); len(diffs) > 0 {
				t.Errorf(
					"\nReadAPKSignature()\n - got:\t\t%+v\n - want:\t%+v\n diff:\n\t%s",
					gotSignature,
					tt.wantSignature,
					strings.Join(diffs, "\n"),
				)
			}
			if len(gotErrors) != tt.wantErrorsCount {
				t.Errorf(
					"\nReadAPKSignature()\n - got_array:\t\t%+v\n - got:\t\t%+v\n - want:\t%+v",
					errors.Join(gotErrors...),
					len(gotErrors),
					tt.wantErrorsCount,
				)
			}
		})
	}
}
