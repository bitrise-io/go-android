package adbmanager

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-android/sdk"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Model ...
type Model struct {
	binPth string
}

// New ...
func New(sdk sdk.AndroidSdkInterface) (*Model, error) {
	binPth := filepath.Join(sdk.GetAndroidHome(), "platform-tools", "adb")
	if exist, err := pathutil.IsPathExists(binPth); err != nil {
		return nil, fmt.Errorf("failed to check if adb exist, error: %s", err)
	} else if !exist {
		return nil, fmt.Errorf("adb not exist at: %s", binPth)
	}

	return &Model{
		binPth: binPth,
	}, nil
}

// UnlockDevice ...
func (model Model) UnlockDevice(serial string) error {
	keyEvent82Cmd := command.New(model.binPth, "-s", serial, "shell", "input", "82", "&")
	if err := keyEvent82Cmd.Run(); err != nil {
		return err
	}

	keyEvent1Cmd := command.New(model.binPth, "-s", serial, "shell", "input", "1", "&")
	return keyEvent1Cmd.Run()
}
