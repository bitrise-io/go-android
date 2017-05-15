package emulatormanager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestEmulatorBinPth(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	emulatorDir := filepath.Join(tmpDir, "emulator")
	require.NoError(t, os.MkdirAll(emulatorDir, 0700))

	emulatorPth := filepath.Join(emulatorDir, "emulator")
	require.NoError(t, fileutil.WriteStringToFile(emulatorPth, ""))

	emulator64Pth := filepath.Join(emulatorDir, "emulator64-arm")
	require.NoError(t, fileutil.WriteStringToFile(emulator64Pth, ""))

	t.Log("prefer emulator64-arm over emulator")
	{
		pth, err := emulatorBinPth(tmpDir, false)
		require.NoError(t, err)
		require.Equal(t, true, strings.Contains(pth, filepath.Join("emulator", "emulator64-arm")), pth)
	}

	t.Log("fallback to emulator")
	{
		require.NoError(t, os.RemoveAll(emulator64Pth))
		pth, err := emulatorBinPth(tmpDir, false)
		require.NoError(t, err)
		require.Equal(t, true, strings.Contains(pth, filepath.Join("emulator", "emulator")), pth)
	}

	t.Log("fail if no emulator bin found")
	{
		require.NoError(t, os.RemoveAll(emulatorPth))
		pth, err := emulatorBinPth(tmpDir, false)
		require.EqualError(t, err, "no emulator binary found in $ANDROID_HOME/emulator")
		require.Equal(t, "", pth)
	}
}

func TestLegacyEmulatorBinPth(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	emulatorDir := filepath.Join(tmpDir, "tools")
	require.NoError(t, os.MkdirAll(emulatorDir, 0700))

	emulatorPth := filepath.Join(emulatorDir, "emulator")
	require.NoError(t, fileutil.WriteStringToFile(emulatorPth, ""))

	emulator64Pth := filepath.Join(emulatorDir, "emulator64-arm")
	require.NoError(t, fileutil.WriteStringToFile(emulator64Pth, ""))

	t.Log("prefer emulator64-arm over emulator")
	{
		pth, err := emulatorBinPth(tmpDir, true)
		require.NoError(t, err)
		require.Equal(t, true, strings.Contains(pth, filepath.Join("tools", "emulator64-arm")), pth)
	}

	t.Log("fallback to emulator")
	{
		require.NoError(t, os.RemoveAll(emulator64Pth))
		pth, err := emulatorBinPth(tmpDir, true)
		require.NoError(t, err)
		require.Equal(t, true, strings.Contains(pth, filepath.Join("tools", "emulator")), pth)
	}

	t.Log("fail if no emulator bin found")
	{
		require.NoError(t, os.RemoveAll(emulatorPth))
		pth, err := emulatorBinPth(tmpDir, true)
		require.EqualError(t, err, "no emulator binary found in $ANDROID_HOME/tools")
		require.Equal(t, "", pth)
	}
}

func TestLib64QTLibEnv(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	lib64QTDir := filepath.Join(tmpDir, "emulator", "lib64", "qt", "lib")
	require.NoError(t, os.MkdirAll(lib64QTDir, 0700))

	t.Log("lib qt env on linux")
	{
		env, err := lib64QTLibEnv(tmpDir, "linux")
		require.NoError(t, err)
		require.Equal(t, true, strings.HasPrefix(env, "LD_LIBRARY_PATH="), env)
		require.Equal(t, true, strings.HasSuffix(env, "emulator/lib64/qt/lib"), env)
	}

	t.Log("lib qt env on osx")
	{
		env, err := lib64QTLibEnv(tmpDir, "darwin")
		require.NoError(t, err)
		require.Equal(t, true, strings.HasPrefix(env, "DYLD_LIBRARY_PATH="), env)
		require.Equal(t, true, strings.HasSuffix(env, "emulator/lib64/qt/lib"), env)
	}

	t.Log("lib qt missing")
	{
		require.NoError(t, os.RemoveAll(lib64QTDir))

		env, err := lib64QTLibEnv(tmpDir, "linux")
		require.Error(t, err)
		require.Equal(t, true, strings.HasPrefix(err.Error(), "qt lib does not exist at:"))
		require.Equal(t, "", env)

		env, err = lib64QTLibEnv(tmpDir, "darwin")
		require.Error(t, err)
		require.Equal(t, true, strings.HasPrefix(err.Error(), "qt lib does not exist at:"))
		require.Equal(t, "", env)
	}

	t.Log("unspported os")
	{
		env, err := lib64QTLibEnv(tmpDir, "windows")
		require.Error(t, err)
		require.EqualError(t, err, "unsupported os windows")
		require.Equal(t, "", env)
	}
}
