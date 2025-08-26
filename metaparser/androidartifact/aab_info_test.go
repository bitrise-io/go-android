package androidartifact

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-android/v2/metaparser/bundletool"
	"github.com/bitrise-io/go-android/v2/metaparser/github"
	"github.com/kr/pretty"
)

func Test_GetAABInfo(t *testing.T) {
	contents, err := github.FetchFile(
		"bitrise-io",
		"sample-artifacts",
		path.Join("aab", "app-weak-algorithm-signed.aab"),
		"master",
		os.Getenv("GITHUB_API_TOKEN"),
	)
	if err != nil {
		t.Fatalf("failed to fetch test artifact: %s", err)
	}

	tmpDir := t.TempDir()
	aabPath := filepath.Join(tmpDir, "temp.aab")
	if err := os.WriteFile(aabPath, contents, os.ModePerm); err != nil {
		t.Fatalf("failed to write file, error: %s", err)
	}

	bt, err := bundletool.New("1.15.0")
	if err != nil {
		t.Fatalf("setup: failed to initialize bundletool, error: %s", err)
	}

	got, err := GetAABInfo(bt, aabPath)
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
