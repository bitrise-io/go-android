package sdk

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestLatestBuildToolsDir(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	buildToolsVersions := []string{"25.0.2", "25.0.3", "22.0.4"}
	for _, buildToolsVersion := range buildToolsVersions {
		buildToolsVersionPth := filepath.Join(tmpDir, "build-tools", buildToolsVersion)
		require.NoError(t, os.MkdirAll(buildToolsVersionPth, 0700))
	}

	sdk, err := New(tmpDir)
	require.NoError(t, err)

	latestBuildToolsDir, err := sdk.LatestBuildToolsDir()
	require.NoError(t, err)
	require.Equal(t, true, strings.Contains(latestBuildToolsDir, filepath.Join("build-tools", "25.0.3")), latestBuildToolsDir)
}

func TestNoBuildToolsDir(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	sdk, err := New(tmpDir)
	require.NoError(t, err)

	_, err = sdk.LatestBuildToolsDir()
	require.EqualError(t, err, "failed to find latest build-tools dir")
}

func TestLatestBuildToolPath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	buildToolsVersions := []string{"25.0.2", "25.0.3", "22.0.4"}
	for _, buildToolsVersion := range buildToolsVersions {
		buildToolsVersionPth := filepath.Join(tmpDir, "build-tools", buildToolsVersion)
		require.NoError(t, os.MkdirAll(buildToolsVersionPth, 0700))
	}

	latestBuildToolsVersions := filepath.Join(tmpDir, "build-tools", "25.0.3")
	zipalignPth := filepath.Join(latestBuildToolsVersions, "zipalign")
	require.NoError(t, fileutil.WriteStringToFile(zipalignPth, ""))

	sdk, err := New(tmpDir)
	require.NoError(t, err)

	t.Log("zipalign - exist")
	{
		pth, err := sdk.LatestBuildToolPath("zipalign")
		require.NoError(t, err)
		require.Equal(t, true, strings.Contains(pth, filepath.Join("build-tools", "25.0.3", "zipalign")), pth)
	}

	t.Log("aapt - NOT exist")
	{
		pth, err := sdk.LatestBuildToolPath("aapt")
		require.Equal(t, true, strings.Contains(err.Error(), "tool (aapt) not found at:"))
		require.Equal(t, "", pth)
	}

}

func TestNewDefaultModel(t *testing.T) {
	androidHome, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	if androidHome, err = filepath.EvalSymlinks(androidHome); err != nil {
		t.Fatalf("failed to eval symlink: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(androidHome); err != nil {
			t.Errorf("failed to remove temp dir: %v", err)
		}
	}()

	sdkRoot, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	if sdkRoot, err = filepath.EvalSymlinks(sdkRoot); err != nil {
		t.Fatalf("failed to eval symlink: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(sdkRoot); err != nil {
			t.Errorf("failed to remove temp dir: %v", err)
		}
	}()

	tests := []struct {
		name    string
		envs    Environment
		want    *Model
		wantErr bool
	}{
		{
			name: "ANDROID_HOME set",
			envs: Environment{
				AndroidHome:    androidHome,
				AndroidSDKRoot: "",
			},
			want: &Model{
				androidHome: androidHome,
			},
		},
		{
			name: "ANDROID_HOME, ANDROID_SDK_ROOT set",
			envs: Environment{
				AndroidHome:    androidHome,
				AndroidSDKRoot: sdkRoot,
			},
			want: &Model{
				androidHome: androidHome,
			},
		},
		{
			name: "ANDROID_SDK_ROOT set",
			envs: Environment{
				AndroidHome:    "",
				AndroidSDKRoot: sdkRoot,
			},
			want: &Model{
				androidHome: sdkRoot,
			},
		},
		{
			name:    "neither ANDROID_HOME, ANDROID_SDK_ROOT set",
			envs:    Environment{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ANDROID_HOME set, non-existent invalid path",
			envs: Environment{
				AndroidHome: filepath.Join(os.TempDir(), "_invalid", "_android_sdk_path"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDefaultModel(tt.envs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_CmdlineToolsPath(t *testing.T) {
	tests := []struct {
		name      string
		SDKlayout []string
		wantPath  string
		wantErr   bool
	}{
		{
			name: "Tools",
			SDKlayout: []string{
				"tools/bin",
			},
			wantPath: "tools/bin",
		},
		{
			name: "Command-line tools latest",
			SDKlayout: []string{
				"tools/bin",
				"cmdline-tools/latest/bin",
				"cmdline-tools/1.0/bin",
			},
			wantPath: "cmdline-tools/latest/bin",
		},
		{
			name: "Command-line tools fixed version",
			SDKlayout: []string{
				"cmdline-tools/1.0/bin",
			},
			wantPath: "cmdline-tools/1.0/bin",
		},
		{
			name: "Command-line tools fixed version",
			SDKlayout: []string{
				"cmdline-tools/1.0/bin",
			},
			wantPath: "cmdline-tools/1.0/bin",
		},
		{
			name:      "No valid path found",
			SDKlayout: []string{},
			wantPath:  "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SDKRoot, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatalf("failed to create temp dir")
			}

			for _, path := range tt.SDKlayout {
				if err := os.MkdirAll(filepath.Join(SDKRoot, path), 0700); err != nil {
					t.Fatalf("failed  to create SDK layout: %v", err)
				}
			}

			model := &Model{
				androidHome: SDKRoot,
			}
			want := ""
			if !tt.wantErr {
				want = filepath.Join(SDKRoot, tt.wantPath)
			}

			got, err := model.CmdlineToolsPath()
			if err != nil {
				t.Log(err.Error())
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Model.CmdlineToolsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != want {
				t.Errorf("Model.CmdlineToolsPath() = %v, want %v", got, want)
			}
		})
	}
}
