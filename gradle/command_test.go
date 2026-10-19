package gradle

import (
	"bytes"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

func TestCommand_Run_failureReasonFromBothStreams(t *testing.T) {
	cmd := NewCommand(newFactory(), "sh", []string{"-c", failingGradleRun}, command.Opts{})

	err := cmd.Run()

	require.Error(t, err)
	require.Contains(t, err.Error(), "command failed with exit status 1")
	require.Contains(t, err.Error(), "Execution failed for task ':app:compileDebugKotlin'.")
	require.Contains(t, err.Error(), "/bitrise/src/app/Foo.kt:4:27: Unresolved reference 'bar'.")
}

func TestCommand_Run_keepsTheCallersOutput(t *testing.T) {
	var out, errOut bytes.Buffer

	cmd := NewCommand(newFactory(), "sh", []string{"-c", failingGradleRun}, command.Opts{Stdout: &out, Stderr: &errOut})

	require.Error(t, cmd.Run())
	require.Contains(t, out.String(), "> Task :app:compileDebugKotlin FAILED")
	require.Contains(t, errOut.String(), "* What went wrong:")
}

func TestCommand_Run_oneWriterForBothStreams(t *testing.T) {
	var out bytes.Buffer

	cmd := NewCommand(newFactory(), "sh", []string{"-c", failingGradleRun}, command.Opts{Stdout: &out, Stderr: &out})

	err := cmd.Run()

	require.Error(t, err)
	require.Contains(t, err.Error(), "/bitrise/src/app/Foo.kt:4:27: Unresolved reference 'bar'.")
	require.Contains(t, out.String(), "> Task :app:compileDebugKotlin FAILED")
	require.Contains(t, out.String(), "* What went wrong:")
}

func TestCommand_Run_noReasonInTheOutput(t *testing.T) {
	cmd := NewCommand(newFactory(), "sh", []string{"-c", "exit 3"}, command.Opts{})

	err := cmd.Run()

	require.Error(t, err)
	require.Contains(t, err.Error(), "command failed with exit status 3")
	require.Contains(t, err.Error(), "check the command's output for details")
}

func TestCommand_Run_successfulRun(t *testing.T) {
	cmd := NewCommand(newFactory(), "sh", []string{"-c", "exit 0"}, command.Opts{})

	require.NoError(t, cmd.Run())
}

func TestCommand_Run_commandCouldNotStart(t *testing.T) {
	cmd := NewCommand(newFactory(), "./no-such-gradlew", nil, command.Opts{})

	require.ErrorContains(t, cmd.Run(), "no-such-gradlew")
}

func newFactory() command.Factory {
	return command.NewFactory(env.NewRepository())
}

const failingGradleRun = `
printf '> Task :app:compileDebugKotlin FAILED\n'
printf "e: file:///bitrise/src/app/Foo.kt:4:27 Unresolved reference 'bar'.\n" >&2
printf 'FAILURE: Build failed with an exception.\n' >&2
printf '* What went wrong:\n' >&2
printf "Execution failed for task ':app:compileDebugKotlin'.\n" >&2
printf '* Try:\n' >&2
exit 1
`
