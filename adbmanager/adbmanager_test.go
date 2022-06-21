package adbmanager

import (
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
)

func TestModel_InstallAPKCmd(t *testing.T) {
	// Given
	mockAPKPath := "/path/to/apk"

	// When
	testCommand := mockModel().InstallAPKCmd(mockAPKPath, &command.Opts{})

	// Then
	actualCMDArgs := testCommand.PrintableCommandArgs()
	expectedCommandArgs := ` "install" "` + mockAPKPath + `"`
	require.Equal(t, expectedCommandArgs, actualCMDArgs)
}

func TestModel_RunInstrumentedTestsCmd_WithAdditionalTestingOptions(t *testing.T) {
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
	actualCMDArgs := testCommand.PrintableCommandArgs()
	expectedCommandArgs := ` "shell" "am" "instrument" "-w"`
	expectedCommandArgs += ` "` + mockPackageName + `/` + mockTestRunnerClass + `"`
	require.Equal(t, expectedCommandArgs, actualCMDArgs)
}

func TestModel_RunInstrumentedTestsCmd_WithoutAdditionalTestingOptions(t *testing.T) {
	// Given
	mockPackageName := "com.package.name"
	mockTestRunnerClass := "mock.testrunner.class"
	mockTestingOpt1 := "opt1"
	mockTestingOpt2 := "opt2"
	mockAdditionalTestingOptions := []string{
		mockTestingOpt1, mockTestingOpt2,
	}

	// When
	testCommand := mockModel().RunInstrumentedTestsCmd(
		mockPackageName,
		mockTestRunnerClass,
		mockAdditionalTestingOptions,
		&command.Opts{},
	)

	// Then
	actualCMDArgs := testCommand.PrintableCommandArgs()
	expectedCommandArgs := ` "shell" "am" "instrument" "-w"`
	expectedCommandArgs += ` "-e" "` + mockTestingOpt1 + `" "` + mockTestingOpt2 + `"`
	expectedCommandArgs += ` "` + mockPackageName + `/` + mockTestRunnerClass + `"`
	require.Equal(t, expectedCommandArgs, actualCMDArgs)
}

// Helpers

func mockModel() Model {
	return Model{
		binPth:     "",
		cmdFactory: command.NewFactory(env.NewRepository()),
	}
}
