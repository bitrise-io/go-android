package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskParseVariants(t *testing.T) {
	task := &Task{
		name: "lint",
		module: Module{
			name: "",
		},
	}

	parsedVariants := task.parseVariants(sampleGradleOutput)

	expected := Variants{
		"Myflavor2Debug",
		"Myflavor2Release",
		"Myflavor2Staging",
		"MyflavorDebug",
		"MyflavorokDebug",
		"MyflavorokRelease",
		"MyflavorokStaging",
		"MyflavorRelease",
		"MyflavorStaging",
		"InvArm7LocalDebug",
		"InvArm7LocalRelease",
		"InvArm7ProdDebug",
		"InvArm7ProdRelease",
		"InvArm7StageDebug",
		"InvArm7StageRelease",
		"InvX86LocalDebug",
		"InvX86LocalRelease",
		"InvX86ProdDebug",
		"InvX86ProdRelease",
		"InvX86StageDebug",
		"InvX86StageRelease",
	}

	require.Equal(t, expected, parsedVariants)
}

const sampleGradleOutput = `
Android tasks
-------------
androidDependencies - Displays the Android dependencies of the project.
signingReport - Displays the signing info for each variant.
sourceSets - Prints out all the source sets defined in this project.

Build tasks
-----------
assemble - Assembles all variants of all applications and secondary packages.
assembleAndroidTest - Assembles all the Test applications.
assembleArm7 - Assembles all Arm7 builds.
assembleDebug - Assembles all Debug builds.
assembleMyflavor - Assembles all Myflavor builds.
assembleMyflavor2 - Assembles all Myflavor2 builds.
assembleMyflavorok - Assembles all Myflavorok builds.
assembleInv - Assembles all Inv builds.
assembleInvArm7Local - Assembles all builds for flavor combination: InvArm7Local
assembleInvArm7Prod - Assembles all builds for flavor combination: InvArm7Prod
assembleInvArm7Stage - Assembles all builds for flavor combination: InvArm7Stage
assembleInvX86Local - Assembles all builds for flavor combination: InvX86Local
assembleInvX86Prod - Assembles all builds for flavor combination: InvX86Prod
assembleInvX86Stage - Assembles all builds for flavor combination: InvX86Stage
assembleLocal - Assembles all Local builds.
assembleProd - Assembles all Prod builds.
assembleRelease - Assembles all Release builds.
assembleStage - Assembles all Stage builds.
assembleStaging - Assembles all Staging builds.
assembleX86 - Assembles all X86 builds.
build - Assembles and tests this project.
buildDependents - Assembles and tests this project and all projects that depend on it.
buildNeeded - Assembles and tests this project and all projects it depends on.
clean - Deletes the build directory.
cleanBuildCache - Deletes the build cache directory.
compileMyflavor2DebugAndroidTestSources
compileMyflavor2DebugSources
compileMyflavor2DebugUnitTestSources
compileMyflavor2ReleaseSources
compileMyflavor2ReleaseUnitTestSources
compileMyflavor2StagingSources
compileMyflavor2StagingUnitTestSources
compileMyflavorDebugAndroidTestSources
compileMyflavorDebugSources
compileMyflavorDebugUnitTestSources
compileMyflavorokDebugAndroidTestSources
compileMyflavorokDebugSources
compileMyflavorokDebugUnitTestSources
compileMyflavorokReleaseSources
compileMyflavorokReleaseUnitTestSources
compileMyflavorokStagingSources
compileMyflavorokStagingUnitTestSources
compileMyflavorReleaseSources
compileMyflavorReleaseUnitTestSources
compileMyflavorStagingSources
compileMyflavorStagingUnitTestSources
compileInvArm7LocalDebugAndroidTestSources
compileInvArm7LocalDebugSources
compileInvArm7LocalDebugUnitTestSources
compileInvArm7LocalReleaseSources
compileInvArm7LocalReleaseUnitTestSources
compileInvArm7ProdDebugAndroidTestSources
compileInvArm7ProdDebugSources
compileInvArm7ProdDebugUnitTestSources
compileInvArm7ProdReleaseSources
compileInvArm7ProdReleaseUnitTestSources
compileInvArm7StageDebugAndroidTestSources
compileInvArm7StageDebugSources
compileInvArm7StageDebugUnitTestSources
compileInvArm7StageReleaseSources
compileInvArm7StageReleaseUnitTestSources
compileInvX86LocalDebugAndroidTestSources
compileInvX86LocalDebugSources
compileInvX86LocalDebugUnitTestSources
compileInvX86LocalReleaseSources
compileInvX86LocalReleaseUnitTestSources
compileInvX86ProdDebugAndroidTestSources
compileInvX86ProdDebugSources
compileInvX86ProdDebugUnitTestSources
compileInvX86ProdReleaseSources
compileInvX86ProdReleaseUnitTestSources
compileInvX86StageDebugAndroidTestSources
compileInvX86StageDebugSources
compileInvX86StageDebugUnitTestSources
compileInvX86StageReleaseSources
compileInvX86StageReleaseUnitTestSources
mockableAndroidJar - Creates a version of android.jar that's suitable for unit tests.

Build Setup tasks
-----------------
init - Initializes a new Gradle build.
wrapper - Generates Gradle wrapper files.

Help tasks
----------
buildEnvironment - Displays all buildscript dependencies declared in root project '_tmp'.
components - Displays the components produced by root project '_tmp'. [incubating]
dependencies - Displays all dependencies declared in root project '_tmp'.
dependencyInsight - Displays the insight into a specific dependency in root project '_tmp'.
dependentComponents - Displays the dependent components of components in root project '_tmp'. [incubating]
help - Displays a help message.
model - Displays the configuration model of root project '_tmp'. [incubating]
projects - Displays the sub-projects of root project '_tmp'.
properties - Displays the properties of root project '_tmp'.
tasks - Displays the tasks runnable from root project '_tmp' (some of the displayed tasks may belong to subprojects).

Install tasks
-------------
installMyflavor2Debug - Installs the DebugMyflavor2 build.
installMyflavor2DebugAndroidTest - Installs the android (on device) tests for the Myflavor2Debug build.
installMyflavorDebug - Installs the DebugMyflavor build.
installMyflavorDebugAndroidTest - Installs the android (on device) tests for the MyflavorDebug build.
installMyflavorokDebug - Installs the DebugMyflavorok build.
installMyflavorokDebugAndroidTest - Installs the android (on device) tests for the MyflavorokDebug build.
installInvArm7LocalDebug - Installs the DebugInvArm7Local build.
installInvArm7LocalDebugAndroidTest - Installs the android (on device) tests for the InvArm7LocalDebug build.
installInvArm7ProdDebug - Installs the DebugInvArm7Prod build.
installInvArm7ProdDebugAndroidTest - Installs the android (on device) tests for the InvArm7ProdDebug build.
installInvArm7StageDebug - Installs the DebugInvArm7Stage build.
installInvArm7StageDebugAndroidTest - Installs the android (on device) tests for the InvArm7StageDebug build.
installInvX86LocalDebug - Installs the DebugInvX86Local build.
installInvX86LocalDebugAndroidTest - Installs the android (on device) tests for the InvX86LocalDebug build.
installInvX86ProdDebug - Installs the DebugInvX86Prod build.
installInvX86ProdDebugAndroidTest - Installs the android (on device) tests for the InvX86ProdDebug build.
installInvX86StageDebug - Installs the DebugInvX86Stage build.
installInvX86StageDebugAndroidTest - Installs the android (on device) tests for the InvX86StageDebug build.
uninstallAll - Uninstall all applications.
uninstallMyflavor2Debug - Uninstalls the DebugMyflavor2 build.
uninstallMyflavor2DebugAndroidTest - Uninstalls the android (on device) tests for the Myflavor2Debug build.
uninstallMyflavor2Release - Uninstalls the ReleaseMyflavor2 build.
uninstallMyflavor2Staging - Uninstalls the StagingMyflavor2 build.
uninstallMyflavorDebug - Uninstalls the DebugMyflavor build.
uninstallMyflavorDebugAndroidTest - Uninstalls the android (on device) tests for the MyflavorDebug build.
uninstallMyflavorokDebug - Uninstalls the DebugMyflavorok build.
uninstallMyflavorokDebugAndroidTest - Uninstalls the android (on device) tests for the MyflavorokDebug build.
uninstallMyflavorokRelease - Uninstalls the ReleaseMyflavorok build.
uninstallMyflavorokStaging - Uninstalls the StagingMyflavorok build.
uninstallMyflavorRelease - Uninstalls the ReleaseMyflavor build.
uninstallMyflavorStaging - Uninstalls the StagingMyflavor build.
uninstallInvArm7LocalDebug - Uninstalls the DebugInvArm7Local build.
uninstallInvArm7LocalDebugAndroidTest - Uninstalls the android (on device) tests for the InvArm7LocalDebug build.
uninstallInvArm7LocalRelease - Uninstalls the ReleaseInvArm7Local build.
uninstallInvArm7ProdDebug - Uninstalls the DebugInvArm7Prod build.
uninstallInvArm7ProdDebugAndroidTest - Uninstalls the android (on device) tests for the InvArm7ProdDebug build.
uninstallInvArm7ProdRelease - Uninstalls the ReleaseInvArm7Prod build.
uninstallInvArm7StageDebug - Uninstalls the DebugInvArm7Stage build.
uninstallInvArm7StageDebugAndroidTest - Uninstalls the android (on device) tests for the InvArm7StageDebug build.
uninstallInvArm7StageRelease - Uninstalls the ReleaseInvArm7Stage build.
uninstallInvX86LocalDebug - Uninstalls the DebugInvX86Local build.
uninstallInvX86LocalDebugAndroidTest - Uninstalls the android (on device) tests for the InvX86LocalDebug build.
uninstallInvX86LocalRelease - Uninstalls the ReleaseInvX86Local build.
uninstallInvX86ProdDebug - Uninstalls the DebugInvX86Prod build.
uninstallInvX86ProdDebugAndroidTest - Uninstalls the android (on device) tests for the InvX86ProdDebug build.
uninstallInvX86ProdRelease - Uninstalls the ReleaseInvX86Prod build.
uninstallInvX86StageDebug - Uninstalls the DebugInvX86Stage build.
uninstallInvX86StageDebugAndroidTest - Uninstalls the android (on device) tests for the InvX86StageDebug build.
uninstallInvX86StageRelease - Uninstalls the ReleaseInvX86Stage build.

Verification tasks
------------------
check - Runs all checks.
connectedAndroidTest - Installs and runs instrumentation tests for all flavors on connected devices.
connectedCheck - Runs all device checks on currently connected devices.
connectedMyflavor2DebugAndroidTest - Installs and runs the tests for Myflavor2Debug on connected devices.
connectedMyflavorDebugAndroidTest - Installs and runs the tests for MyflavorDebug on connected devices.
connectedMyflavorokDebugAndroidTest - Installs and runs the tests for MyflavorokDebug on connected devices.
connectedInvArm7LocalDebugAndroidTest - Installs and runs the tests for invArm7LocalDebug on connected devices.
connectedInvArm7ProdDebugAndroidTest - Installs and runs the tests for invArm7ProdDebug on connected devices.
connectedInvArm7StageDebugAndroidTest - Installs and runs the tests for invArm7StageDebug on connected devices.
connectedInvX86LocalDebugAndroidTest - Installs and runs the tests for invX86LocalDebug on connected devices.
connectedInvX86ProdDebugAndroidTest - Installs and runs the tests for invX86ProdDebug on connected devices.
connectedInvX86StageDebugAndroidTest - Installs and runs the tests for invX86StageDebug on connected devices.
deviceAndroidTest - Installs and runs instrumentation tests using all Device Providers.
deviceCheck - Runs all device checks using Device Providers and Test Servers.
lint - Runs lint on all variants.
lintMyflavor2Debug - Runs lint on the Myflavor2Debug build.
lintMyflavor2Release - Runs lint on the Myflavor2Release build.
lintMyflavor2Staging - Runs lint on the Myflavor2Staging build.
lintMyflavorDebug - Runs lint on the MyflavorDebug build.
lintMyflavorokDebug - Runs lint on the MyflavorokDebug build.
lintMyflavorokRelease - Runs lint on the MyflavorokRelease build.
lintMyflavorokStaging - Runs lint on the MyflavorokStaging build.
lintMyflavorRelease - Runs lint on the MyflavorRelease build.
lintMyflavorStaging - Runs lint on the MyflavorStaging build.
lintInvArm7LocalDebug - Runs lint on the InvArm7LocalDebug build.
lintInvArm7LocalRelease - Runs lint on the InvArm7LocalRelease build.
lintInvArm7ProdDebug - Runs lint on the InvArm7ProdDebug build.
lintInvArm7ProdRelease - Runs lint on the InvArm7ProdRelease build.
lintInvArm7StageDebug - Runs lint on the InvArm7StageDebug build.
lintInvArm7StageRelease - Runs lint on the InvArm7StageRelease build.
lintInvX86LocalDebug - Runs lint on the InvX86LocalDebug build.
lintInvX86LocalRelease - Runs lint on the InvX86LocalRelease build.
lintInvX86ProdDebug - Runs lint on the InvX86ProdDebug build.
lintInvX86ProdRelease - Runs lint on the InvX86ProdRelease build.
lintInvX86StageDebug - Runs lint on the InvX86StageDebug build.
lintInvX86StageRelease - Runs lint on the InvX86StageRelease build.
lintVitalMyflavor2Release - Runs lint on just the fatal issues in the Myflavor2Release build.
lintVitalMyflavor2Staging - Runs lint on just the fatal issues in the Myflavor2Staging build.
lintVitalMyflavorokRelease - Runs lint on just the fatal issues in the MyflavorokRelease build.
lintVitalMyflavorokStaging - Runs lint on just the fatal issues in the MyflavorokStaging build.
lintVitalMyflavorRelease - Runs lint on just the fatal issues in the MyflavorRelease build.
lintVitalMyflavorStaging - Runs lint on just the fatal issues in the MyflavorStaging build.
lintVitalInvArm7LocalRelease - Runs lint on just the fatal issues in the invArm7LocalRelease build.
lintVitalInvArm7ProdRelease - Runs lint on just the fatal issues in the invArm7ProdRelease build.
lintVitalInvArm7StageRelease - Runs lint on just the fatal issues in the invArm7StageRelease build.
lintVitalInvX86LocalRelease - Runs lint on just the fatal issues in the invX86LocalRelease build.
lintVitalInvX86ProdRelease - Runs lint on just the fatal issues in the invX86ProdRelease build.
lintVitalInvX86StageRelease - Runs lint on just the fatal issues in the invX86StageRelease build.
test - Run unit tests for all variants.
testMyflavor2DebugUnitTest - Run unit tests for the Myflavor2Debug build.
testMyflavor2ReleaseUnitTest - Run unit tests for the Myflavor2Release build.
testMyflavor2StagingUnitTest - Run unit tests for the Myflavor2Staging build.
testMyflavorDebugUnitTest - Run unit tests for the MyflavorDebug build.
testMyflavorokDebugUnitTest - Run unit tests for the MyflavorokDebug build.
testMyflavorokReleaseUnitTest - Run unit tests for the MyflavorokRelease build.
testMyflavorokStagingUnitTest - Run unit tests for the MyflavorokStaging build.
testMyflavorReleaseUnitTest - Run unit tests for the MyflavorRelease build.
testMyflavorStagingUnitTest - Run unit tests for the MyflavorStaging build.
testInvArm7LocalDebugUnitTest - Run unit tests for the invArm7LocalDebug build.
testInvArm7LocalReleaseUnitTest - Run unit tests for the invArm7LocalRelease build.
testInvArm7ProdDebugUnitTest - Run unit tests for the invArm7ProdDebug build.
testInvArm7ProdReleaseUnitTest - Run unit tests for the invArm7ProdRelease build.
testInvArm7StageDebugUnitTest - Run unit tests for the invArm7StageDebug build.
testInvArm7StageReleaseUnitTest - Run unit tests for the invArm7StageRelease build.
testInvX86LocalDebugUnitTest - Run unit tests for the invX86LocalDebug build.
testInvX86LocalReleaseUnitTest - Run unit tests for the invX86LocalRelease build.
testInvX86ProdDebugUnitTest - Run unit tests for the invX86ProdDebug build.
testInvX86ProdReleaseUnitTest - Run unit tests for the invX86ProdRelease build.
testInvX86StageDebugUnitTest - Run unit tests for the invX86StageDebug build.
testInvX86StageReleaseUnitTest - Run unit tests for the invX86StageRelease build.`
