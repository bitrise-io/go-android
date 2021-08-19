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
	binPth     string
	cmdFactory command.Factory
}

// New ...
func New(sdk sdk.AndroidSdkInterface, cmdFactory command.Factory) (*Model, error) {
	binPth := filepath.Join(sdk.GetAndroidHome(), "platform-tools", "adb")
	if exist, err := pathutil.IsPathExists(binPth); err != nil {
		return nil, fmt.Errorf("failed to check if adb exist, error: %s", err)
	} else if !exist {
		return nil, fmt.Errorf("adb not exist at: %s", binPth)
	}

	return &Model{
		binPth:     binPth,
		cmdFactory: cmdFactory,
	}, nil
}

// DevicesCmd ...
func (model Model) DevicesCmd() *command.Command {
	cmd := model.cmdFactory.Create(model.binPth, []string{"devices"}, nil)
	return &cmd
}

// IsDeviceBooted ...
func (model Model) IsDeviceBooted(serial string) (bool, error) {
	devBootCmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "getprop dev.bootcomplete"}, nil)
	devBootOut, err := devBootCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return false, err
	}

	sysBootCmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "getprop sys.boot_completed"}, nil)
	sysBootOut, err := sysBootCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return false, err
	}

	bootAnimCmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "getprop init.svc.bootanim"}, nil)
	bootAnimOut, err := bootAnimCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return false, err
	}

	return devBootOut == "1" && sysBootOut == "1" && bootAnimOut == "stopped", nil
}

// UnlockDevice ...
func (model Model) UnlockDevice(serial string) error {
	keyEvent82Cmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "input", "82", "&"}, nil)
	if err := keyEvent82Cmd.Run(); err != nil {
		return err
	}

	keyEvent1Cmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "input", "1", "&"}, nil)
	return keyEvent1Cmd.Run()
}
