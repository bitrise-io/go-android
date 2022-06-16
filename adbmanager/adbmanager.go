package adbmanager

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-android/v2/sdk"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/v2/command"
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

// UnlockDevice ...
func (model Model) UnlockDevice(serial string) error {
	keyEvent82Cmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "input", "82", "&"}, nil)
	if err := keyEvent82Cmd.Run(); err != nil {
		return err
	}

	keyEvent1Cmd := model.cmdFactory.Create(model.binPth, []string{"-s", serial, "shell", "input", "1", "&"}, nil)
	return keyEvent1Cmd.Run()
}

// InstallAPKCmd builds and returns a `Command` for installing APKs on an attached device or emulator.
// The `Command` can than be run by the consumer without needing to know the implementation details.
func (model Model) InstallAPKCmd(pathToAPK string, commandOptions *command.Opts) command.Command {
	cmd := model.cmdFactory.Create(model.binPth, []string{"install", pathToAPK}, commandOptions)
	return cmd
}

// RunInstrumentedTestsCmd builds and returns a `Command` for running instrumented tests on an attached device or emulator.
// The `Command` can than be run by the consumer without needing to know the implementation details.
func (model Model) RunInstrumentedTestsCmd(
	packageName string,
	testRunnerClass string,
	additionalTestingOptions []string,
	commandOptions *command.Opts,
) command.Command {
	args := []string{"shell", "am", "instrument"}
	if len(additionalTestingOptions) > 0 {
		args = append(args, "-e")
		args = append(args, additionalTestingOptions...)
	}
	args = append(args, "-w", packageName+"/"+testRunnerClass)
	cmd := model.cmdFactory.Create(model.binPth, args, commandOptions)
	return cmd
}
