package integration

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-tools/go-android/adbmanager"
	"github.com/bitrise-tools/go-android/avdmanager"
	"github.com/bitrise-tools/go-android/emulatormanager"
	"github.com/bitrise-tools/go-android/sdk"
	"github.com/bitrise-tools/go-android/sdkcomponent"
	"github.com/bitrise-tools/go-android/sdkmanager"
	"github.com/stretchr/testify/require"
)

const (
	testEmulatorName = "test_emu_name"
	testEmulatorSkin = "768x1280"
)

type emulator struct {
	platform string
	tag      string
	abi      string
}

func TestIsLegacyAVDManager(t *testing.T) {
	_, err := avdmanager.IsLegacyAVDManager(os.Getenv("ANDROID_HOME"))
	require.NoError(t, err)
}

func TestIntegrity(t *testing.T) {
	for _, emu := range getEmulatorConfigList() {
		createEmulator(t, emu.platform, emu.tag, emu.abi)
		startEmulator(t)
	}
}

func getEmulatorConfigList() []emulator {
	return []emulator{
		emulator{platform: "android-24", tag: "google_apis", abi: "armeabi-v7a"},
		emulator{platform: "android-24", tag: "google_apis", abi: "arm64-v8a"},
		emulator{platform: "android-25", tag: "android-wear", abi: "armeabi-v7a"},
		emulator{platform: "android-23", tag: "android-tv", abi: "armeabi-v7a"},
		emulator{platform: "android-24", tag: "default", abi: "armeabi-v7a"},
		emulator{platform: "android-24", tag: "default", abi: "arm64-v8a"},
		emulator{platform: "android-17", tag: "default", abi: "mips"},
	}
}

func createEmulator(t *testing.T, platform string, tag string, abi string) {
	log.Printf("\nRunning test: %s - %s - %s", platform, tag, abi)

	log.Printf("\n-Check if platform installed")

	androidSdk, err := sdk.New(os.Getenv("ANDROID_HOME"))
	require.NoError(t, err)

	manager, err := sdkmanager.New(androidSdk)
	require.NoError(t, err)

	platformComponent := sdkcomponent.Platform{
		Version: platform,
	}

	platformInstalled, err := manager.IsInstalled(platformComponent)
	require.NoError(t, err)

	log.Printf("\n-installed: %v", platformInstalled)

	if !platformInstalled {

		log.Printf("\n-Installing: %s", platform)

		installCmd := manager.InstallCommand(platformComponent)
		installCmd.SetStdin(strings.NewReader("y"))

		log.Printf("\n-$ %s", installCmd.PrintableCommandArgs())

		out, err := installCmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		installed, err := manager.IsInstalled(platformComponent)
		require.NoError(t, err)
		require.Equal(t, true, installed)

		log.Printf("\n-Installed")
	}

	log.Printf("\n-Check if system image installed")

	systemImageComponent := sdkcomponent.SystemImage{
		Platform: platform,
		Tag:      tag,
		ABI:      abi,
	}

	log.Printf("\n-Checking path: %s", systemImageComponent.InstallPathInAndroidHome())

	systemImageInstalled, err := manager.IsInstalled(systemImageComponent)
	require.NoError(t, err)

	log.Printf("\n-installed: %v", systemImageInstalled)

	if !systemImageInstalled {
		log.Printf("\n-Installing system image (platform: %s abi: %s tag: %s)", systemImageComponent.Platform, systemImageComponent.ABI, systemImageComponent.Tag)

		installCmd := manager.InstallCommand(systemImageComponent)
		installCmd.SetStdin(strings.NewReader("y"))

		log.Printf("\n-$ %s", installCmd.PrintableCommandArgs())

		out, err := installCmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		installed, err := manager.IsInstalled(systemImageComponent)
		require.NoError(t, err)
		require.Equal(t, true, installed)

		log.Printf("\n-Installed")
	}

	log.Printf("\n-Creating AVD image")

	avdManager, err := avdmanager.New(androidSdk)
	require.NoError(t, err)

	cmd := avdManager.CreateAVDCommand(testEmulatorName, systemImageComponent)
	cmd.SetStdin(strings.NewReader("n"))

	log.Printf("\n-$ %s", cmd.PrintableCommandArgs())

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	require.NoError(t, err, out)
	log.Printf("\n...DONE...\n\n")
}

func startEmulator(t *testing.T) {
	log.Printf("Validate AVD image")

	avdImages, err := listAVDImages()
	require.NoError(t, err)

	if !sliceutil.IsStringInSlice(testEmulatorName, avdImages) {

		if len(avdImages) > 0 {
			log.Printf("Available avd images:")
			for _, avdImage := range avdImages {
				log.Printf("* %s", avdImage)
			}
		}

		os.Exit(1)
	}

	log.Donef("AVD image (%s) exist", testEmulatorName)
	// ---

	androidSdk, err := sdk.New(os.Getenv("ANDROID_HOME"))
	require.NoError(t, err)

	adb, err := adbmanager.New(androidSdk)
	require.NoError(t, err)

	//
	// Print running devices Info
	deviceStateMap, err := runningDeviceInfos(*adb)
	require.NoError(t, err)

	if len(deviceStateMap) > 0 {
		fmt.Println()
		log.Infof("Running devices:")

		for serial, state := range deviceStateMap {
			log.Printf("* %s (%s)", serial, state)
		}
	}
	// ---

	emulator, err := emulatormanager.New(androidSdk)
	require.NoError(t, err)

	//
	// Start AVD image
	fmt.Println()
	log.Infof("Start AVD image")

	options := []string{"-no-boot-anim", "-no-window"}

	startEmulatorCommand := emulator.StartEmulatorCommand(testEmulatorName, testEmulatorSkin, options...)
	startEmulatorCmd := startEmulatorCommand.GetCmd()

	e := make(chan error)

	// Redirect output
	stdoutReader, err := startEmulatorCmd.StdoutPipe()
	require.NoError(t, err)

	outScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for outScanner.Scan() {
			line := outScanner.Text()
			fmt.Println(line)
			if strings.Contains(strings.ToLower(line), "invalid cpu architecture") {
				require.FailNow(t, line)
			}
		}
	}()
	err = outScanner.Err()
	require.NoError(t, err)

	// Redirect error
	stderrReader, err := startEmulatorCmd.StderrPipe()
	require.NoError(t, err)

	errScanner := bufio.NewScanner(stderrReader)
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			log.Warnf(line)
		}
	}()
	err = errScanner.Err()
	require.NoError(t, err)
	// ---

	serial := ""

	go func() {
		// Start emulator
		log.Printf("$ %s", command.PrintableCommandArgs(false, startEmulatorCmd.Args))
		fmt.Println()

		if err := startEmulatorCommand.Run(); err != nil {
			e <- err
			return
		}
	}()

	go func() {
		// Wait until device appears in device list
		for len(serial) == 0 {
			time.Sleep(5 * time.Second)

			log.Printf("> Checking for started device serial...")

			currentDeviceStateMap, err := runningDeviceInfos(*adb)
			if err != nil {
				e <- err
				return
			}

			serial = currentlyStartedDeviceSerial(deviceStateMap, currentDeviceStateMap)
		}

		log.Donef("> Started device serial: %s", serial)

		// Wait until device is booted

		bootInProgress := true
		for bootInProgress {
			time.Sleep(5 * time.Second)

			log.Printf("> Checking if device booted...")

			booted, err := adb.IsDeviceBooted(serial)
			if err != nil {
				e <- err
				return
			}

			bootInProgress = !booted
		}

		err := adb.UnlockDevice(serial)
		require.NoError(t, err)

		log.Donef("> Device booted")

		e <- nil
	}()

	timeout := 600

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		err := startEmulatorCmd.Process.Kill()
		require.NoError(t, err)
		require.FailNow(t, "Boot timed out...", timeout)

	case err := <-e:
		require.NoError(t, err)

	}
	// ---
	os.Setenv("EMULATOR_SERIAL", serial)

	fmt.Println()
	log.Donef("Emulator (%s) booted", serial)
	err = startEmulatorCmd.Process.Kill()
	require.NoError(t, err)
}

func listAVDImages() ([]string, error) {
	homeDir := pathutil.UserHomeDir()
	avdDir := filepath.Join(homeDir, ".android", "avd")

	avdImagePattern := filepath.Join(avdDir, "*.ini")
	avdImages, err := filepath.Glob(avdImagePattern)
	if err != nil {
		return []string{}, fmt.Errorf("glob failed with pattern (%s), error: %s", avdImagePattern, err)
	}

	avdImageNames := []string{}

	for _, avdImage := range avdImages {
		imageName := filepath.Base(avdImage)
		imageName = strings.TrimSuffix(imageName, filepath.Ext(avdImage))
		avdImageNames = append(avdImageNames, imageName)
	}

	return avdImageNames, nil
}

func avdImageDir(name string) string {
	homeDir := pathutil.UserHomeDir()
	return filepath.Join(homeDir, ".android", "avd", name+".avd")
}

func currentlyStartedDeviceSerial(alreadyRunningDeviceInfos, currentlyRunningDeviceInfos map[string]string) string {
	startedSerial := ""

	for serial := range currentlyRunningDeviceInfos {
		_, found := alreadyRunningDeviceInfos[serial]
		if !found {
			startedSerial = serial
			break
		}
	}

	if len(startedSerial) > 0 {
		state := currentlyRunningDeviceInfos[startedSerial]
		if state == "device" {
			return startedSerial
		}
	}

	return ""
}

func runningDeviceInfos(adb adbmanager.Model) (map[string]string, error) {
	cmd := adb.DevicesCmd()
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
