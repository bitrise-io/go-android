package sdk

import (
	"testing"

	"github.com/bitrise-tools/go-android/sdkcomponent"
	"github.com/stretchr/testify/require"
)

/*
// GetSystemImage ...
func (sdk *Model) GetSystemImage(platform string, abi string, componentType string) (sdkcomponent.SystemImage, bool) {
	systemImages := []sdkcomponent.SystemImage{}

	for _, systemImage := range sdk.systemImages {
		if systemImage.Platform == platform && systemImage.ABI == abi {
			systemImages = append(systemImages, systemImage)
		}
	}

	if componentType == "" {
		componentType = "default"
	}

	for _, systemImage := range systemImages {
		if systemImage.Type == componentType {
			return systemImage, true
		}
	}

	if componentType == "" {
		componentType = "google_apis"
	}

	for _, systemImage := range systemImages {
		if systemImage.Type == componentType {
			return systemImage, true
		}
	}

	return sdkcomponent.SystemImage{}, false
}
*/

func TestGetSystemImage(t *testing.T) {
	sdk := Model{
		systemImages: []sdkcomponent.SystemImage{
			sdkcomponent.SystemImage{
				Platform: "android-25",
				Type:     "google_apis",
				ABI:      "x86",
			},
		},
	}

	t.Log("SystemImage exist")
	{
		systemImage, found := sdk.GetSystemImage("android-25", "x86", "google_apis")
		require.Equal(t, true, found)
		require.Equal(t, sdkcomponent.SystemImage{
			Platform: "android-25",
			Type:     "google_apis",
			ABI:      "x86",
		}, systemImage)
	}

	t.Log("SystemImage exist - type not specified")
	{
		systemImage, found := sdk.GetSystemImage("android-25", "x86", "")
		require.Equal(t, true, found)
		require.Equal(t, sdkcomponent.SystemImage{
			Platform: "android-25",
			Type:     "google_apis",
			ABI:      "x86",
		}, systemImage)
	}

	t.Log("SystemImage NOT exist")
	{
		systemImage, found := sdk.GetSystemImage("android-23", "x86", "google_apis")
		require.Equal(t, false, found)
		require.Equal(t, sdkcomponent.SystemImage{}, systemImage)
	}
}

func TestGetPlatform(t *testing.T) {
	sdk := Model{
		platforms: []sdkcomponent.Platform{
			sdkcomponent.Platform{
				Version: "android-25",
			},
		},
	}

	t.Log("Platform exist")
	{
		platform, found := sdk.GetPlatform("android-25")
		require.Equal(t, true, found)
		require.Equal(t, sdkcomponent.Platform{
			Version: "android-25",
		}, platform)
	}

	t.Log("Platform NOT exist")
	{
		platform, found := sdk.GetPlatform("android-23")
		require.Equal(t, false, found)
		require.Equal(t, sdkcomponent.Platform{}, platform)
	}
}

func TestGetBuildTool(t *testing.T) {
	sdk := Model{
		buildTools: []sdkcomponent.BuildTool{
			sdkcomponent.BuildTool{
				Version: "19.0.1",
			},
		},
	}

	t.Log("BuildTool exist")
	{
		buildTool, found := sdk.GetBuildTool("19.0.1")
		require.Equal(t, true, found)
		require.Equal(t, sdkcomponent.BuildTool{
			Version: "19.0.1",
		}, buildTool)
	}

	t.Log("BuildTool NOT exist")
	{
		buildTool, found := sdk.GetBuildTool("25.0.2")
		require.Equal(t, false, found)
		require.Equal(t, sdkcomponent.BuildTool{}, buildTool)
	}
}

func TestAdd(t *testing.T) {
	sdk := Model{}

	t.Log("addBuildTool")
	{
		sdk.addBuildTool(sdkcomponent.BuildTool{
			Version: "19.0.1",
		})
		require.Equal(t, 1, len(sdk.buildTools))
		require.Equal(t, sdkcomponent.BuildTool{
			Version: "19.0.1",
		}, sdk.buildTools[0])

		sdk.addBuildTool(sdkcomponent.BuildTool{
			Version: "25.0.2",
		})
		require.Equal(t, 2, len(sdk.buildTools))
		require.Equal(t, sdkcomponent.BuildTool{
			Version: "19.0.1",
		}, sdk.buildTools[0])
		require.Equal(t, sdkcomponent.BuildTool{
			Version: "25.0.2",
		}, sdk.buildTools[1])
	}

	t.Log("addPlatform")
	{
		sdk.addPlatform(sdkcomponent.Platform{
			Version: "android-25",
		})
		require.Equal(t, 1, len(sdk.platforms))
		require.Equal(t, sdkcomponent.Platform{
			Version: "android-25",
		}, sdk.platforms[0])

		sdk.addPlatform(sdkcomponent.Platform{
			Version: "android-23",
		})
		require.Equal(t, 2, len(sdk.platforms))
		require.Equal(t, sdkcomponent.Platform{
			Version: "android-25",
		}, sdk.platforms[0])
		require.Equal(t, sdkcomponent.Platform{
			Version: "android-23",
		}, sdk.platforms[1])
	}

	t.Log("addSystemImage")
	{
		sdk.addSystemImage(sdkcomponent.SystemImage{
			Platform: "android-25",
			ABI:      "armeabi-v7a",
		})
		require.Equal(t, 1, len(sdk.systemImages))
		require.Equal(t, sdkcomponent.SystemImage{
			Platform: "android-25",
			ABI:      "armeabi-v7a",
		}, sdk.systemImages[0])

		sdk.addSystemImage(sdkcomponent.SystemImage{
			Platform: "android-23",
			ABI:      "x86",
			Type:     "google_apis",
		})
		require.Equal(t, 2, len(sdk.systemImages))
		require.Equal(t, sdkcomponent.SystemImage{
			Platform: "android-25",
			ABI:      "armeabi-v7a",
		}, sdk.systemImages[0])
		require.Equal(t, sdkcomponent.SystemImage{
			Platform: "android-23",
			ABI:      "x86",
			Type:     "google_apis",
		}, sdk.systemImages[1])
	}
}

func TestParseSDKOut(t *testing.T) {
	sdk, err := parseSDKOut(sdkOut)
	require.NoError(t, err)

	require.Equal(t, 2, len(sdk.buildTools))
	require.Equal(t, sdkcomponent.BuildTool{
		Version:      "19.1.0",
		SDKStylePath: "build-tools;19.1.0",
	}, sdk.buildTools[0])
	require.Equal(t, sdkcomponent.BuildTool{
		Version:      "19.1.0",
		SDKStylePath: "build-tools;19.1.0",
	}, sdk.buildTools[1])

	require.Equal(t, 2, len(sdk.platforms))
	require.Equal(t, sdkcomponent.Platform{
		Version:      "android-25",
		SDKStylePath: "platforms;android-25",
	}, sdk.platforms[0])
	require.Equal(t, sdkcomponent.Platform{
		Version:      "android-10",
		SDKStylePath: "platforms;android-10",
	}, sdk.platforms[1])

	require.Equal(t, 5, len(sdk.systemImages))
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:     "android-25",
		Type:         "google_apis",
		ABI:          "x86_64",
		SDKStylePath: "system-images;android-25;google_apis;x86_64",
	}, sdk.systemImages[0])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:     "android-10",
		Type:         "default",
		ABI:          "armeabi-v7a",
		SDKStylePath: "system-images;android-10;default;armeabi-v7a",
	}, sdk.systemImages[1])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:     "android-10",
		Type:         "google_apis",
		ABI:          "armeabi-v7a",
		SDKStylePath: "system-images;android-10;google_apis;armeabi-v7a",
	}, sdk.systemImages[2])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:     "android-21",
		Type:         "android-tv",
		ABI:          "armeabi-v7a",
		SDKStylePath: "system-images;android-21;android-tv;armeabi-v7a",
	}, sdk.systemImages[3])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:     "android-23",
		Type:         "android-wear",
		ABI:          "armeabi-v7a",
		SDKStylePath: "system-images;android-23;android-wear;armeabi-v7a",
	}, sdk.systemImages[4])
}

func TestParseLegacySDKOut(t *testing.T) {
	sdk, err := parseLegacySDKOut(legacySDKOut)
	require.NoError(t, err)

	require.Equal(t, 1, len(sdk.buildTools))
	require.Equal(t, sdkcomponent.BuildTool{
		Version:            "25.0.2",
		LegacySDKStylePath: "build-tools-25.0.2",
	}, sdk.buildTools[0])

	require.Equal(t, 1, len(sdk.platforms))
	require.Equal(t, sdkcomponent.Platform{
		Version:            "android-25",
		LegacySDKStylePath: "android-25",
	}, sdk.platforms[0])

	require.Equal(t, 9, len(sdk.systemImages))
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "android-tv",
		ABI:                "x86",
		LegacySDKStylePath: "sys-img-x86-android-tv-25",
	}, sdk.systemImages[0])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "android-wear",
		ABI:                "armeabi-v7a",
		LegacySDKStylePath: "sys-img-armeabi-v7a-android-wear-25",
	}, sdk.systemImages[1])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "android-wear",
		ABI:                "x86",
		LegacySDKStylePath: "sys-img-x86-android-wear-25",
	}, sdk.systemImages[2])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "google_apis",
		ABI:                "arm64-v8a",
		LegacySDKStylePath: "sys-img-arm64-v8a-google_apis-25",
	}, sdk.systemImages[3])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "google_apis",
		ABI:                "armeabi-v7a",
		LegacySDKStylePath: "sys-img-armeabi-v7a-google_apis-25",
	}, sdk.systemImages[4])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "google_apis",
		ABI:                "x86_64",
		LegacySDKStylePath: "sys-img-x86_64-google_apis-25",
	}, sdk.systemImages[5])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-25",
		Type:               "google_apis",
		ABI:                "x86",
		LegacySDKStylePath: "sys-img-x86-google_apis-25",
	}, sdk.systemImages[6])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-24",
		Type:               "android",
		ABI:                "arm64-v8a",
		LegacySDKStylePath: "sys-img-arm64-v8a-android-24",
	}, sdk.systemImages[7])
	require.Equal(t, sdkcomponent.SystemImage{
		Platform:           "android-24",
		Type:               "android",
		ABI:                "armeabi-v7a",
		LegacySDKStylePath: "sys-img-armeabi-v7a-android-24",
	}, sdk.systemImages[8])
}

const sdkOut = `Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/19.1.0/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/20.0.0/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/21.0.1/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/22.0.1/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/23.0.2/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/23.0.3/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/24.0.2/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/25.0.2/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/emulator/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/android/m2repository/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/android/support/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/auto/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/google_play_services/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/m2repository/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/market_apk_expansion/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/market_licensing/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/play_billing/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/simulators/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/webdriver/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/intel/Hardware_Accelerated_Execution_Manager/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/m2repository/com/android/support/constraint/constraint-layout-solver/1.0.0-beta3/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/m2repository/com/android/support/constraint/constraint-layout/1.0.0-beta3/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/patcher/v1/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/patcher/v4/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platform-tools/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platforms/android-23/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platforms/android-24/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platforms/android-25/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/sources/android-23/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/sources/android-25/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/system-images/android-25/google_apis/x86_64/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/tools/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/19.1.0/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/20.0.0/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/21.0.1/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/22.0.1/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/23.0.2/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/23.0.3/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/24.0.2/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/build-tools/25.0.2/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/emulator/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/android/m2repository/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/android/support/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/auto/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/google_play_services/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/m2repository/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/market_apk_expansion/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/market_licensing/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/play_billing/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/simulators/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/google/webdriver/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/intel/Hardware_Accelerated_Execution_Manager/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/m2repository/com/android/support/constraint/constraint-layout-solver/1.0.0-beta3/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/extras/m2repository/com/android/support/constraint/constraint-layout/1.0.0-beta3/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/patcher/v1/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/patcher/v4/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platform-tools/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platforms/android-23/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platforms/android-24/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/platforms/android-25/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/sources/android-23/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/sources/android-25/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/system-images/android-25/google_apis/x86_64/package.xml
Info: Parsing /Users/godrei/Library/Android/sdk/tools/package.xml
Installed packages:
--------------------------------------
build-tools;19.1.0
    Description:        Android SDK Build-Tools 19.1
    Version:            19.1.0
    Installed Location: /Users/godrei/Library/Android/sdk/build-tools/19.1.0

emulator
    Description:        Android Emulator
    Version:            25.3.0
    Installed Location: /Users/godrei/Library/Android/sdk/emulator

extras;android;m2repository
    Description:        Android Support Repository
    Version:            44.0.0
    Installed Location: /Users/godrei/Library/Android/sdk/extras/android/m2repository

extras;android;support
    Description:        Android Support Library, rev 23.2.1
    Version:            23.2.1
    Installed Location: /Users/godrei/Library/Android/sdk/extras/android/support

extras;google;google_play_services
    Description:        Google Play services
    Version:            39
    Installed Location: /Users/godrei/Library/Android/sdk/extras/google/google_play_services

extras;google;m2repository
    Description:        Google Repository
    Version:            44
    Installed Location: /Users/godrei/Library/Android/sdk/extras/google/m2repository

extras;intel;Hardware_Accelerated_Execution_Manager
    Description:        Intel x86 Emulator Accelerator (HAXM installer)
    Version:            6.0.5
    Installed Location: /Users/godrei/Library/Android/sdk/extras/intel/Hardware_Accelerated_Execution_Manager

extras;m2repository;com;android;support;constraint;constraint-layout-solver;1.0.0-beta3
    Description:        Solver for ConstraintLayout 1.0.0-beta3
    Version:            1
    Installed Location: /Users/godrei/Library/Android/sdk/extras/m2repository/com/android/support/constraint/constraint-layout-solver/1.0.0-beta3

extras;m2repository;com;android;support;constraint;constraint-layout;1.0.0-beta3
    Description:        ConstraintLayout for Android 1.0.0-beta3
    Version:            1
    Installed Location: /Users/godrei/Library/Android/sdk/extras/m2repository/com/android/support/constraint/constraint-layout/1.0.0-beta3

patcher;v1
    Description:        SDK Patch Applier v1
    Version:            1
    Installed Location: /Users/godrei/Library/Android/sdk/patcher/v1

platform-tools
    Description:        Android SDK Platform-Tools 25.0.3
    Version:            25.0.3
    Installed Location: /Users/godrei/Library/Android/sdk/platform-tools

platforms;android-25
    Description:        Android SDK Platform 25
    Version:            3
    Installed Location: /Users/godrei/Library/Android/sdk/platforms/android-25

sources;android-23
    Description:        Sources for Android 23
    Version:            1
    Installed Location: /Users/godrei/Library/Android/sdk/sources/android-23

system-images;android-25;google_apis;x86_64
    Description:        Google APIs Intel x86 Atom_64 System Image
    Version:            4
    Installed Location: /Users/godrei/Library/Android/sdk/system-images/android-25/google_apis/x86_64

tools
    Description:        Android SDK Tools
    Version:            25.3.1
    Installed Location: /Users/godrei/Library/Android/sdk/tools

Available Packages:
--------------------------------------
add-ons;addon-google_apis-google-15
    Description:        Google APIs
    Version:            3

build-tools;19.1.0
    Description:        Android SDK Build-Tools 19.1
    Version:            19.1.0

cmake;3.6.3155560
    Description:        CMake 3.6.3155560
    Version:            3.6.3155560

docs
    Description:        Documentation for Android SDK
    Version:            1

emulator
    Description:        Android Emulator
    Version:            25.3.1
    Dependencies:
        tools Revision 25.3

extras;google;google_play_services
    Description:        Google Play services
    Version:            39
    Dependencies:
        patcher;v4

extras;google;m2repository
    Description:        Google Repository
    Version:            44
    Dependencies:
        patcher;v4

extras;intel;Hardware_Accelerated_Execution_Manager
    Description:        Intel x86 Emulator Accelerator (HAXM installer)
    Version:            6.0.5

extras;m2repository;com;android;support;constraint;constraint-layout-solver;1.0.0
    Description:        Solver for ConstraintLayout 1.0.0
    Version:            1

extras;m2repository;com;android;support;constraint;constraint-layout-solver;1.0.0-alpha2
    Description:        com.android.support.constraint:constraint-layout-solver:1.0.0-alpha2
    Version:            1

lldb;2.0
    Description:        LLDB 2.0
    Version:            2.0.2558144

ndk-bundle
    Description:        NDK
    Version:            14.0.3770861

patcher;v4
    Description:        SDK Patch Applier v4
    Version:            1

platform-tools
    Description:        Android SDK Platform-Tools
    Version:            25.0.4

platforms;android-10
    Description:        Android SDK Platform 10
    Version:            2

sources;android-15
    Description:        Sources for Android 15
    Version:            2

system-images;android-10;default;armeabi-v7a
    Description:        ARM EABI v7a System Image
    Version:            4       4

system-images;android-10;google_apis;armeabi-v7a
    Description:        Google APIs ARM EABI v7a System Image
    Version:            5

system-images;android-21;android-tv;armeabi-v7a
    Description:        Android TV ARM EABI v7a System Image
    Version:            3

system-images;android-23;android-wear;armeabi-v7a
    Description:        Android Wear ARM EABI v7a System Image
    Version:            6
    Dependencies:
        patcher;v4

tools
    Description:        Android SDK Tools
    Version:            25.3.1
    Dependencies:
        emulator
        platform-tools Revision 20

Available Updates:
--------------------------------------
emulator
    Local Version:  25.3.0
    Remote Version: 25.3.1
extras;android;m2repository
    Local Version:  44.0.0
    Remote Version: 45.0.0
platform-tools
    Local Version:  25.0.3
    Remote Version: 25.0.4
done`

const legacySDKOut = `Refresh Sources:
  Fetching https://dl.google.com/android/repository/addons_list-2.xml
  Validate XML
  Parse XML
  Fetched Add-ons List successfully
  Refresh Sources
  Fetching URL: https://dl.google.com/android/repository/repository-11.xml
  Validate XML: https://dl.google.com/android/repository/repository-11.xml
  Parse XML:    https://dl.google.com/android/repository/repository-11.xml
  Fetching URL: https://dl.google.com/android/repository/addon.xml
  Validate XML: https://dl.google.com/android/repository/addon.xml
  Parse XML:    https://dl.google.com/android/repository/addon.xml
  Fetching URL: https://dl.google.com/android/repository/glass/addon.xml
  Validate XML: https://dl.google.com/android/repository/glass/addon.xml
  Parse XML:    https://dl.google.com/android/repository/glass/addon.xml
  Fetching URL: https://dl.google.com/android/repository/extras/intel/addon.xml
  Validate XML: https://dl.google.com/android/repository/extras/intel/addon.xml
  Parse XML:    https://dl.google.com/android/repository/extras/intel/addon.xml
  Fetching URL: https://dl.google.com/android/repository/sys-img/android/sys-img.xml
  Validate XML: https://dl.google.com/android/repository/sys-img/android/sys-img.xml
  Parse XML:    https://dl.google.com/android/repository/sys-img/android/sys-img.xml
  Fetching URL: https://dl.google.com/android/repository/sys-img/android-wear/sys-img.xml
  Validate XML: https://dl.google.com/android/repository/sys-img/android-wear/sys-img.xml
  Parse XML:    https://dl.google.com/android/repository/sys-img/android-wear/sys-img.xml
  Fetching URL: https://dl.google.com/android/repository/sys-img/android-tv/sys-img.xml
  Validate XML: https://dl.google.com/android/repository/sys-img/android-tv/sys-img.xml
  Parse XML:    https://dl.google.com/android/repository/sys-img/android-tv/sys-img.xml
  Fetching URL: https://dl.google.com/android/repository/sys-img/google_apis/sys-img.xml
  Validate XML: https://dl.google.com/android/repository/sys-img/google_apis/sys-img.xml
  Parse XML:    https://dl.google.com/android/repository/sys-img/google_apis/sys-img.xml
Packages available for installation or update: 174
----------
id: 1 or "tools"
     Type: Tool
     Desc: Android SDK Tools, revision 25.2.5
----------
id: 2 or "platform-tools"
     Type: PlatformTool
     Desc: Android SDK Platform-tools, revision 25.0.3
----------
id: 3 or "build-tools-25.0.2"
     Type: BuildTool
     Desc: Android SDK Build-tools, revision 25.0.2
----------
id: 32 or "doc-24"
     Type: Doc
     Desc: Documentation for Android SDK, API 24, revision 1
----------
id: 33 or "android-25"
     Type: Platform
     Desc: Android SDK Platform 25
           Revision 3
----------
id: 57 or "sys-img-x86-android-tv-25"
     Type: SystemImage
     Desc: Android TV Intel x86 Atom System Image
           Revision 3
           Requires SDK Platform Android API 25
----------
id: 58 or "sys-img-armeabi-v7a-android-wear-25"
     Type: SystemImage
     Desc: Android Wear ARM EABI v7a System Image
           Revision 3
           Requires SDK Platform Android API 25
----------
id: 59 or "sys-img-x86-android-wear-25"
     Type: SystemImage
     Desc: Android Wear Intel x86 Atom System Image
           Revision 3
           Requires SDK Platform Android API 25
----------
id: 60 or "sys-img-arm64-v8a-google_apis-25"
     Type: SystemImage
     Desc: Google APIs ARM 64 v8a System Image
           Revision 4
           Requires SDK Platform Android API 25
----------
id: 61 or "sys-img-armeabi-v7a-google_apis-25"
     Type: SystemImage
     Desc: Google APIs ARM EABI v7a System Image
           Revision 4
           Requires SDK Platform Android API 25
----------
id: 62 or "sys-img-x86_64-google_apis-25"
     Type: SystemImage
     Desc: Google APIs Intel x86 Atom_64 System Image
           Revision 4
           Requires SDK Platform Android API 25
----------
id: 63 or "sys-img-x86-google_apis-25"
     Type: SystemImage
     Desc: Google APIs Intel x86 Atom System Image
           Revision 4
           Requires SDK Platform Android API 25
----------
id: 65 or "sys-img-arm64-v8a-android-24"
     Type: SystemImage
     Desc: ARM 64 v8a System Image
           Revision 7
           Requires SDK Platform Android API 24
----------
id: 66 or "sys-img-armeabi-v7a-android-24"
     Type: SystemImage
     Desc: ARM EABI v7a System Image
           Revision 7
           Requires SDK Platform Android API 24
----------
id: 126 or "addon-google_apis-google-24"
     Type: Addon
     Desc: Google APIs, Android API 24, revision 1
           By Google Inc.
           Android + Google APIs
           Requires SDK Platform Android API 24
----------
id: 150 or "source-25"
     Type: Source
     Desc: Sources for Android SDK, API 25, revision 1
----------
id: 162 or "extra-android-m2repository"
     Type: Extra
     Desc: Android Support Repository, revision 45
           By Android
           Local Maven repository for Support Libraries
           Install path: extras/android/m2repository
----------
id: 167 or "extra-google-google_play_services_froyo"
     Type: Extra
     Desc: Google Play services for Froyo, revision 12 (Obsolete)
           By Google Inc.
           Google Play services client library and sample code
           Install path: extras/google/google_play_services_froyo
----------
id: 168 or "extra-google-google_play_services"
     Type: Extra
     Desc: Google Play services, revision 39
           By Google Inc.
           Google Play services Javadocs and sample code
           Install path: extras/google/google_play_services
----------
id: 169 or "extra-google-m2repository"
     Type: Extra
     Desc: Google Repository, revision 44
           By Google Inc.
           Local Maven repository for Support Libraries
           Install path: extras/google/m2repository`
