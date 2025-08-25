package androidartifact

import (
	"os"
	"path"
	"testing"

	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/log"
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

func TestGetAPKInfoWithFallback(t *testing.T) {
	tLogger := &testLogger{}
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("setup: failed to create temp dir, error: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Warnf("failed to remove temp dir, error: %s", err)
		}
	}()

	gitCommand, err := git.New(tmpDir)
	if err != nil {
		t.Fatalf("setup: failed to create git project, error: %s", err)
	}
	if err := gitCommand.Clone("https://github.com/bitrise-io/sample-artifacts.git").Run(); err != nil {
		t.Fatalf("setup: failed to clone test artifact repo, error: %s", err)
	}

	type args struct {
		logger Logger
		apkPth string
	}
	tests := []struct {
		name    string
		apkPath string
		want    Info
		wantErr bool
	}{
		{
			name:    "valid apk",
			apkPath: path.Join(tmpDir, "apks", "app-debug.apk"),
			want: Info{
				AppName:           "My Application",
				PackageName:       "com.example.birmachera.myapplication",
				VersionCode:       "1",
				VersionName:       "1.0",
				MinSDKVersion:     "17",
				RawPackageContent: testArtifactAndroidManifest,
			},
			wantErr: false,
		},
		{
			name:    "apk with unicode app name",
			apkPath: path.Join(tmpDir, "apks", "app-debug-uncode-name.apk"),
			want: Info{
				AppName:       "My Application Unicode 🤪",
				PackageName:   "com.example.myapplicationUnicode",
				VersionCode:   "1",
				VersionName:   "1.0",
				MinSDKVersion: "33",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAPKInfoWithFallback(tLogger, tt.apkPath)
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
