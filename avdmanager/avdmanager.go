package avdmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-android/sdk"
	"github.com/bitrise-io/go-android/sdkcomponent"
	"github.com/bitrise-io/go-android/sdkmanager"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
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
		return &Model{binPth: avdmanagerPath}, nil
	}

	if SDKManager, err := sdkmanager.New(sdk); err != nil {
		log.Warnf("%v", err)
	} else if !SDKManager.IsLegacySDK() {
		fmt.Println()
		log.Warnf("Found sdkmanager but no avdmanager, updating SDK Tools...")

		sdkToolComponent := sdkcomponent.SDKTool{}
		updateCmd := SDKManager.InstallCommand(sdkToolComponent)
		updateCmd.SetStderr(os.Stderr).SetStdout(os.Stdout)
		log.Infof("$ %s", updateCmd.PrintableCommandArgs())
		if err := updateCmd.Run(); err != nil {
			log.Errorf("Installing avdmanager failed: %v", err)
		} else {
			if exists, err := pathutil.IsPathExists(avdmanagerPath); err != nil {
				return nil, err
			} else if exists {
				return &Model{binPth: avdmanagerPath}, nil
			}
		}
		log.Printf("Updating SDK tools was unsuccessful, continuing with legacy avdmanager...")
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
