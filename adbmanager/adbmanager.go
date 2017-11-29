package adbmanager

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-android/sdk"
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

// Command ...
func (adb Model) Command(args ...string) *command.Model {
	return command.New(adb.binPth, args...)
}

// IsDeviceBooted ...
func (adb Model) IsDeviceBooted(serial string) (bool, error) {
	getpropCmd := adb.shellCommand(serial, "shell", "getprop dev.bootcomplete '0' && getprop sys.boot_completed '0' && getprop init.svc.bootanim 'running'")
	getpropOut, err := getpropCmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return false, err
	}

	return getpropOut == "1\n1\nstopped", nil
}

// UnlockDevice ...
func (adb Model) UnlockDevice(serial string) error {
	if err := adb.shellCommand(serial, "shell", "input", "82", "&").Run(); err != nil {
		return err
	}
	return adb.shellCommand(serial, "shell", "input", "1", "&").Run()
}

// GetDevices ...
func (adb Model) GetDevices() (map[string]string, error) {
	cmd := adb.Command("devices")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return map[string]string{}, fmt.Errorf("command failed, error: %s", err)
	}

	// List of devices attached
	// emulator-5554	device
	deviceListItemPattern := `^(?P<emulator>emulator-\d*)[\s+](?P<state>.*)`
	deviceListItemRegexp := regexp.MustCompile(deviceListItemPattern)

	deviceStateMap := map[string]string{}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		matches := deviceListItemRegexp.FindStringSubmatch(line)
		if len(matches) == 3 {
			serial := matches[1]
			state := matches[2]

			deviceStateMap[serial] = state
		}

	}
	if scanner.Err() != nil {
		return map[string]string{}, fmt.Errorf("scanner failed, error: %s", err)
	}

	return deviceStateMap, nil
}

func (adb Model) shellCommand(serial string, args ...string) *command.Model {
	return command.New(adb.binPth, append([]string{"-s", serial}, args...)...)
}
