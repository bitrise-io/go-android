package adbmanager

import (
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

func Test_GivenAPKPath_WhenCreateInstallAPKCmd_ThenCreatesExpectedCommand(t *testing.T) {
	// Given
	mockAPKPath := "/path/to/apk"

	// When
	testCommand := mockModel().InstallAPKCmd(
		mockAPKPath,
		&command.Opts{},
	)

	// Then
	actualArgs := testCommand.PrintableCommandArgs()
	expectedArgs := `adb "install" "/path/to/apk"`
	require.Equal(t, expectedArgs, actualArgs)
}

func Test_GivenNoAdditionalTestingOptions_WhenCreateRunInstrumentedTestsCmd_ThenCreatesExpectedCommand(t *testing.T) {
	// Given
	mockPackageName := "com.package.name"
	mockTestRunnerClass := "mock.testrunner.class"

	// When
	testCommand := mockModel().RunInstrumentedTestsCmd(
		mockPackageName,
		mockTestRunnerClass,
		nil,
		&command.Opts{},
	)

	// Then
	actualArgs := testCommand.PrintableCommandArgs()
	expectedArgs := `adb "shell" "am" "instrument" "-w" "com.package.name/mock.testrunner.class"`
	require.Equal(t, expectedArgs, actualArgs)
}

func Test_GivenAdditionalTestingOptions_WhenCreateRunInstrumentedTestsCmd_ThenCreatesExpectedCommand(t *testing.T) {
	// Given
	mockPackageName := "com.package.name"
	mockTestRunnerClass := "mock.testrunner.class"
	mockAdditionalTestingOptions := []string{"opt1", "opt2"}

	// When
	testCommand := mockModel().RunInstrumentedTestsCmd(
		mockPackageName,
		mockTestRunnerClass,
		mockAdditionalTestingOptions,
		&command.Opts{},
	)

	// Then
	actualArgs := testCommand.PrintableCommandArgs()
	expectedArgs := `adb "shell" "am" "instrument" "-w" "-e" "opt1" "opt2" "com.package.name/mock.testrunner.class"`
	require.Equal(t, expectedArgs, actualArgs)
}

func Test_GivenShellCommand_WhenCreateWaitForDeviceCmd_ThenCreatesExpectedCommand(t *testing.T) {
	// Given
	serial := "emulator-5554"
	shellCmd := "getprop sys.boot_completed"

	// When
	testCommand := mockModel().WaitForDeviceThenShellCmd(
		serial,
		&command.Opts{},
		shellCmd,
	)

	// Then
	actualArgs := testCommand.PrintableCommandArgs()
	expectedArgs := `adb "-s" "emulator-5554" "wait-for-device" "shell" "getprop sys.boot_completed"`
	require.Equal(t, expectedArgs, actualArgs)
}

// Helpers

func mockModel() Model {
	return Model{
		binPth:     "adb",
		cmdFactory: command.NewFactory(env.NewRepository()),
	}
}
