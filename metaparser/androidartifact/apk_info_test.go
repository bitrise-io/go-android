package androidartifact

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-android/v2/metaparser/github"
	"github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/stretchr/testify/mock"
)

const testArtifactAndroidManifest string = `<manifest xmlns:android="http://schemas.android.com/apk/res/android" android:versionCode="1" android:versionName="1.0" package="com.example.birmachera.myapplication">
	<uses-sdk android:minSdkVersion="17" android:targetSdkVersion="28"></uses-sdk>
	<uses-permission android:name="android.permission.INTERNET"></uses-permission>
	<application android:theme="null" android:label="My Application" android:icon="res/mipmap-xxxhdpi-v4/ic_launcher.png" android:debuggable="true" android:allowBackup="true" android:supportsRtl="true" android:roundIcon="res/mipmap-xxxhdpi-v4/ic_launcher_round.png" android:appComponentFactory="android.support.v4.app.CoreComponentFactory">
		<activity android:name="com.example.birmachera.myapplication.MainActivity">
			<intent-filter>
				<action android:name="android.intent.action.MAIN"></action>
				<category android:name="android.intent.category.LAUNCHER"></category>
			</intent-filter>
		</activity>
	</application>
</manifest>`

// fakeAaptDumpBadgingOutput is what `aapt dump badging` would print for the
// "apk with unicode app name" fixture; used to exercise the aapt fallback path
// without needing an installed Android SDK.
const fakeAaptDumpBadgingOutput = `package: name='com.example.myapplicationUnicode' versionCode='1' versionName='1.0'
sdkVersion:'33'
application-label:'My Application Unicode 🤪'
application: label='My Application Unicode 🤪' icon='res/mipmap-mdpi-v4/ic_launcher.png'`

// fakeSDKLocator satisfies SDKLocator without requiring an installed Android SDK.
type fakeSDKLocator struct {
	path string
	err  error
}

func (f fakeSDKLocator) LatestBuildToolPath(_ string) (string, error) {
	return f.path, f.err
}

func TestGetAPKInfoWithFallback(t *testing.T) {
	tLogger := &testLogger{}

	aaptCmd := mocks.NewCommand(t)
	aaptCmd.On("RunAndReturnTrimmedCombinedOutput").Return(fakeAaptDumpBadgingOutput, nil)

	cmdFactory := mocks.NewFactory(t)
	cmdFactory.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(aaptCmd)

	sdkModel := fakeSDKLocator{path: "aapt"}

	tests := []struct {
		name    string
		apkPath string
		want    Info
		wantErr bool
	}{
		{
			name:    "valid apk",
			apkPath: "apks/app-debug.apk",
			want: Info{
				AppName:           "My Application",
				PackageName:       "com.example.birmachera.myapplication",
				VersionCode:       "1",
				VersionName:       "1.0",
				MinSDKVersion:     "17",
				RawPackageContent: testArtifactAndroidManifest,
			},
		},
		{
			name:    "apk with unicode app name",
			apkPath: "apks/app-debug-uncode-name.apk",
			want: Info{
				AppName:       "My Application Unicode 🤪",
				PackageName:   "com.example.myapplicationUnicode",
				VersionCode:   "1",
				VersionName:   "1.0",
				MinSDKVersion: "33",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Faster than cloning repo
			contents, err := github.FetchFile(
				"bitrise-io",
				"sample-artifacts",
				tt.apkPath,
				"master",
				os.Getenv("GITHUB_TOKEN"),
			)
			if err != nil {
				t.Fatalf("failed to fetch test artifact: %s", err)
			}

			tmpDir := t.TempDir()
			apkPath := filepath.Join(tmpDir, "test.apk")
			if err := os.WriteFile(apkPath, contents, os.ModePerm); err != nil {
				t.Fatalf("failed to write file, error: %s", err)
			}

			got, err := GetAPKInfoWithFallback(tLogger, cmdFactory, sdkModel, apkPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPKInfoWithFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.AppName != tt.want.AppName {
				t.Fatalf("GetAPKInfo() app name = %+v, want %+v", got.AppName, tt.want.AppName)
			}
			if got.PackageName != tt.want.PackageName {
				t.Fatalf("GetAPKInfo() package name = %+v, want %+v", got.PackageName, tt.want.PackageName)
			}
			if got.VersionCode != tt.want.VersionCode {
				t.Fatalf("GetAPKInfo() version code = %+v, want %+v", got.VersionCode, tt.want.VersionCode)
			}
			if got.VersionName != tt.want.VersionName {
				t.Fatalf("GetAPKInfo() version name = %+v, want %+v", got.VersionName, tt.want.VersionName)
			}
			if got.MinSDKVersion != tt.want.MinSDKVersion {
				t.Fatalf("GetAPKInfo() min sdk version = %+v, want %+v", got.MinSDKVersion, tt.want.MinSDKVersion)
			}
			if tt.want.RawPackageContent != "" && got.RawPackageContent != tt.want.RawPackageContent {
				t.Fatalf("GetAPKInfo() raw package content = %+v, want %+v", got.RawPackageContent, tt.want.RawPackageContent)
			}
		})
	}
}
