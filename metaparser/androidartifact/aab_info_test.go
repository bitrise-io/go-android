package androidartifact

import (
	"slices"
	"strings"
	"testing"

	"github.com/bitrise-io/go-android/v2/metaparser/bundletool"
	"github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/mock"
)

func Test_GetAABInfo(t *testing.T) {
	manifestCmd := mocks.NewCommand(t)
	manifestCmd.On("RunAndReturnTrimmedCombinedOutput").Return(testAABArtifactAndroidManifest, nil)

	resourcesCmd := mocks.NewCommand(t)
	resourcesCmd.On("RunAndReturnTrimmedCombinedOutput").Return(`"en" "sample-apps-android-simple"`, nil)

	cmdFactory := mocks.NewFactory(t)
	cmdFactory.On("Create", mock.Anything, mock.MatchedBy(func(args []string) bool {
		return slices.Contains(args, "resources")
	}), mock.Anything).Return(resourcesCmd)
	cmdFactory.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(manifestCmd)

	bt := bundletool.Path("fake-bundletool-all.jar")

	got, err := GetAABInfo(cmdFactory, bt, "fake.aab")
	if err != nil {
		t.Fatalf("GetAABInfo() error = %v", err)
	}

	want := Info{
		AppName:           "sample-apps-android-simple",
		PackageName:       "com.bitrise_io.sample_apps_android_simple",
		VersionCode:       "189",
		VersionName:       "1.0",
		MinSDKVersion:     "15",
		RawPackageContent: testAABArtifactAndroidManifest,
	}
	if diffs := pretty.Diff(got, want); len(diffs) > 0 {
		t.Errorf(
			"\nGetAABInfo()\n - got:\t\t%+v\n - want:\t%+v\n diff:\n\t%s",
			got,
			want,
			strings.Join(diffs, "\n"),
		)
	}
}

// testAABArtifactAndroidManifest is the manifest `bundletool dump manifest` would print
// for a real AAB; used here as canned output for the fake command factory.
const testAABArtifactAndroidManifest string = `<manifest xmlns:android="http://schemas.android.com/apk/res/android" android:versionCode="189" android:versionName="1.0" package="com.bitrise_io.sample_apps_android_simple" platformBuildVersionCode="189" platformBuildVersionName="1.0">

  <uses-sdk android:minSdkVersion="15" android:targetSdkVersion="26"/>

  <application android:allowBackup="true" android:icon="@mipmap/ic_launcher" android:label="@string/app_name" android:supportsRtl="true" android:theme="@style/AppTheme">

    <activity android:label="@string/app_name" android:name="com.bitrise_io.sample_apps_android_simple.MainActivity">

      <intent-filter>

        <action android:name="android.intent.action.MAIN"/>

        <category android:name="android.intent.category.LAUNCHER"/>

      </intent-filter>

    </activity>

    <meta-data android:name="android.support.VERSION" android:value="26.1.0"/>

    <meta-data android:name="android.arch.lifecycle.VERSION" android:value="27.0.0-SNAPSHOT"/>

  </application>

</manifest>`
