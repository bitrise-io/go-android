package avdmanager

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-android/sdk"
	"github.com/bitrise-io/go-android/sdkcomponent"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Model ...
type Model struct {
	legacy bool
	binPth string
}

// New ...
func New(sdk sdk.AndroidSdkInterface) (*Model, error) {
	cmdlineTools, err := sdk.CmdlineToolsPath()
	if err != nil {
		return nil, err
	}

	avdmanagerPath := filepath.Join(cmdlineTools, "avdmanager")
	if exists, err := pathutil.IsPathExists(avdmanagerPath); err != nil {
		return nil, err
	} else if exists {
		return &Model{
			binPth: avdmanagerPath,
		}, nil
	}

	legacyAvdmanagerPath := filepath.Join(cmdlineTools, "android")
	if exists, err := pathutil.IsPathExists(legacyAvdmanagerPath); err != nil {
		return nil, err
	} else if exists {
		return &Model{
			legacy: true,
			binPth: avdmanagerPath,
		}, nil
	}

	return nil, fmt.Errorf("no avdmanager found at %s", avdmanagerPath)
}

// CreateAVDCommand ...
func (model Model) CreateAVDCommand(name string, systemImage sdkcomponent.SystemImage, options ...string) *command.Model {
	args := []string{"--verbose", "create", "avd", "--force", "--name", name, "--abi", systemImage.ABI}

	if model.legacy {
		args = append(args, "--target", systemImage.Platform)
	} else {
		args = append(args, "--package", systemImage.GetSDKStylePath())
	}

	if systemImage.Tag != "" && systemImage.Tag != "default" {
		args = append(args, "--tag", systemImage.Tag)
	}

	args = append(args, options...)
	return command.New(model.binPth, args...)
}
