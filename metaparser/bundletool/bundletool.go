package bundletool

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/filedownloader"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

var cmdFactory = command.NewFactory(env.NewRepository())

// Path ...
type Path string

// New ...
func New(logger log.Logger, version string) (Path, error) {
	return fetchAny(
		logger,
		"https://github.com/google/bundletool/releases/download/"+version+"/bundletool-all-"+version+".jar",
		"https://github.com/google/bundletool/releases/download/"+version+"/bundletool-all.jar",
	)
}

func fetchAny(logger log.Logger, source string, fallbackSources ...string) (Path, error) {
	tmpPth, err := pathutil.NewPathProvider().CreateTempDir("tool")
	if err != nil {
		return "", err
	}

	downloader := filedownloader.NewDownloader(logger)

	toolPath := filepath.Join(tmpPth, "bundletool-all.jar")
	if err := downloader.DownloadWithFallback(context.Background(), toolPath, source, fallbackSources...); err != nil {
		return "", err
	}

	return Path(toolPath), nil
}

// Command ...
func (p Path) Command(cmd string, args ...string) command.Command {
	return cmdFactory.Create("java", append([]string{"-Djdk.util.zip.disableZip64ExtraFieldValidation=true", "-jar", string(p), cmd}, args...), nil)
}

// Exec ...
func (p Path) Exec(cmd string, args ...string) (string, error) {
	c := p.Command(cmd, args...)

	out, err := c.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("command failed with exit status %d (%s): %s", exitErr.ExitCode(), c.PrintableCommandArgs(), out)
		}
		return "", fmt.Errorf("executing command failed (%s): %w", c.PrintableCommandArgs(), err)
	}

	return out, nil
}
