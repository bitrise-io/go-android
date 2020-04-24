package gradle

import (
	"reflect"
	"testing"
)

func TestTask_parseVariants(t *testing.T) {
	tests := []struct {
		name         string
		taskName     string
		gradleOutput string
		want         Variants
	}{
		{
			name:         "No output",
			taskName:     "lint",
			gradleOutput: ``,
			want:         map[string][]string{},
		},
		{
			name:         "One module with one nested submodule",
			taskName:     "lint",
			gradleOutput: gradleOutput,
			want: map[string][]string{
				"app":           {"Debug", "Release"},
				"app:mylibrary": {"Debug", "Release"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				name:    tt.taskName,
				project: Project{},
			}
			if got := task.parseVariants(tt.gradleOutput); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.parseVariants() = %v, want %v", got, tt.want)
			}
		})
	}
}

const gradleOutput = `
------------------------------------------------------------
Tasks runnable from root project
------------------------------------------------------------

Android tasks
-------------
app:androidDependencies - Displays the Android dependencies of the project.
app:mylibrary:androidDependencies - Displays the Android dependencies of the project.
app:signingReport - Displays the signing info for the base and test modules
app:mylibrary:signingReport - Displays the signing info for the base and test modules
app:sourceSets - Prints out all the source sets defined in this project.
app:mylibrary:sourceSets - Prints out all the source sets defined in this project.

Build tasks
-----------
app:assemble - Assemble main outputs for all the variants.
app:mylibrary:assemble - Assemble main outputs for all the variants.
app:assembleAndroidTest - Assembles all the Test applications.
app:mylibrary:assembleAndroidTest - Assembles all the Test applications.
app:build - Assembles and tests this project.
app:mylibrary:build - Assembles and tests this project.
app:buildDependents - Assembles and tests this project and all projects that depend on it.
app:mylibrary:buildDependents - Assembles and tests this project and all projects that depend on it.
app:buildNeeded - Assembles and tests this project and all projects it depends on.
app:mylibrary:buildNeeded - Assembles and tests this project and all projects it depends on.
app:bundle - Assemble bundles for all the variants.
app:clean - Deletes the build directory.
app:mylibrary:clean - Deletes the build directory.
app:cleanBuildCache - Deletes the build cache directory.
app:mylibrary:cleanBuildCache - Deletes the build cache directory.
app:compileDebugAndroidTestSources
app:mylibrary:compileDebugAndroidTestSources
app:compileDebugSources
app:mylibrary:compileDebugSources
app:compileDebugUnitTestSources
app:mylibrary:compileDebugUnitTestSources
app:compileReleaseSources
app:mylibrary:compileReleaseSources
app:compileReleaseUnitTestSources
app:mylibrary:compileReleaseUnitTestSources
app:mylibrary:extractDebugAnnotations - Extracts Android annotations for the debug variant into the archive file
app:mylibrary:extractReleaseAnnotations - Extracts Android annotations for the release variant into the archive file

Build Setup tasks
-----------------
init - Initializes a new Gradle build.
wrapper - Generates Gradle wrapper files.

Cleanup tasks
-------------
app:lintFix - Runs lint on all variants and applies any safe suggestions to the source code.
app:mylibrary:lintFix - Runs lint on all variants and applies any safe suggestions to the source code.

Help tasks
----------
buildEnvironment - Displays all buildscript dependencies declared in root project 'My Application'.
app:buildEnvironment - Displays all buildscript dependencies declared in project ':app'.
app:mylibrary:buildEnvironment - Displays all buildscript dependencies declared in project ':app:mylibrary'.
components - Displays the components produced by root project 'My Application'. [incubating]
app:components - Displays the components produced by project ':app'. [incubating]
app:mylibrary:components - Displays the components produced by project ':app:mylibrary'. [incubating]
dependencies - Displays all dependencies declared in root project 'My Application'.
app:dependencies - Displays all dependencies declared in project ':app'.
app:mylibrary:dependencies - Displays all dependencies declared in project ':app:mylibrary'.
dependencyInsight - Displays the insight into a specific dependency in root project 'My Application'.
app:dependencyInsight - Displays the insight into a specific dependency in project ':app'.
app:mylibrary:dependencyInsight - Displays the insight into a specific dependency in project ':app:mylibrary'.
dependentComponents - Displays the dependent components of components in root project 'My Application'. [incubating]
app:dependentComponents - Displays the dependent components of components in project ':app'. [incubating]
app:mylibrary:dependentComponents - Displays the dependent components of components in project ':app:mylibrary'. [incubating]
help - Displays a help message.
app:help - Displays a help message.
app:mylibrary:help - Displays a help message.
model - Displays the configuration model of root project 'My Application'. [incubating]
app:model - Displays the configuration model of project ':app'. [incubating]
app:mylibrary:model - Displays the configuration model of project ':app:mylibrary'. [incubating]
projects - Displays the sub-projects of root project 'My Application'.
app:projects - Displays the sub-projects of project ':app'.
app:mylibrary:projects - Displays the sub-projects of project ':app:mylibrary'.
properties - Displays the properties of root project 'My Application'.
app:properties - Displays the properties of project ':app'.
app:mylibrary:properties - Displays the properties of project ':app:mylibrary'.
tasks - Displays the tasks runnable from root project 'My Application' (some of the displayed tasks may belong to subprojects).
app:tasks - Displays the tasks runnable from project ':app' (some of the displayed tasks may belong to subprojects).
app:mylibrary:tasks - Displays the tasks runnable from project ':app:mylibrary'.

Install tasks
-------------
app:installDebug - Installs the Debug build.
app:installDebugAndroidTest - Installs the android (on device) tests for the Debug build.
app:mylibrary:installDebugAndroidTest - Installs the android (on device) tests for the Debug build.
app:uninstallAll - Uninstall all applications.
app:mylibrary:uninstallAll - Uninstall all applications.
app:uninstallDebug - Uninstalls the Debug build.
app:uninstallDebugAndroidTest - Uninstalls the android (on device) tests for the Debug build.
app:mylibrary:uninstallDebugAndroidTest - Uninstalls the android (on device) tests for the Debug build.
app:uninstallRelease - Uninstalls the Release build.

Verification tasks
------------------
app:check - Runs all checks.
app:mylibrary:check - Runs all checks.
app:connectedAndroidTest - Installs and runs instrumentation tests for all flavors on connected devices.
app:mylibrary:connectedAndroidTest - Installs and runs instrumentation tests for all flavors on connected devices.
app:connectedCheck - Runs all device checks on currently connected devices.
app:mylibrary:connectedCheck - Runs all device checks on currently connected devices.
app:connectedDebugAndroidTest - Installs and runs the tests for debug on connected devices.
app:mylibrary:connectedDebugAndroidTest - Installs and runs the tests for debug on connected devices.
app:deviceAndroidTest - Installs and runs instrumentation tests using all Device Providers.
app:mylibrary:deviceAndroidTest - Installs and runs instrumentation tests using all Device Providers.
app:deviceCheck - Runs all device checks using Device Providers and Test Servers.
app:mylibrary:deviceCheck - Runs all device checks using Device Providers and Test Servers.
app:lint - Runs lint on all variants.
app:mylibrary:lint - Runs lint on all variants.
app:lintDebug - Runs lint on the Debug build.
app:mylibrary:lintDebug - Runs lint on the Debug build.
app:lintRelease - Runs lint on the Release build.
app:mylibrary:lintRelease - Runs lint on the Release build.
app:lintVitalRelease - Runs lint on just the fatal issues in the release build.
app:test - Run unit tests for all variants.
app:mylibrary:test - Run unit tests for all variants.
app:testDebugUnitTest - Run unit tests for the debug build.
app:mylibrary:testDebugUnitTest - Run unit tests for the debug build.
app:testReleaseUnitTest - Run unit tests for the release build.
app:mylibrary:testReleaseUnitTest - Run unit tests for the release build.

Other tasks
-----------
app:assembleDebug - Assembles main output for variant debug
app:mylibrary:assembleDebug - Assembles main output for variant debug
app:assembleDebugAndroidTest - Assembles main output for variant debugAndroidTest
app:mylibrary:assembleDebugAndroidTest - Assembles main output for variant debugAndroidTest
app:assembleDebugUnitTest - Assembles main output for variant debugUnitTest
app:mylibrary:assembleDebugUnitTest - Assembles main output for variant debugUnitTest
app:assembleRelease - Assembles main output for variant release
app:mylibrary:assembleRelease - Assembles main output for variant release
app:assembleReleaseUnitTest - Assembles main output for variant releaseUnitTest
app:mylibrary:assembleReleaseUnitTest - Assembles main output for variant releaseUnitTest
app:buildDebugPreBundle
app:buildReleasePreBundle
app:bundleDebug - Assembles bundle for variant debug
app:mylibrary:bundleDebugAar - Assembles a bundle containing the library in debug.
app:bundleDebugAndroidTestClasses
app:bundleDebugAndroidTestResources
app:mylibrary:bundleDebugAndroidTestResources
app:bundleDebugClasses
app:bundleDebugResources
app:bundleDebugUnitTestClasses
app:mylibrary:bundleLibCompileDebug
app:mylibrary:bundleLibCompileDebugAndroidTest
app:mylibrary:bundleLibCompileDebugUnitTest
app:mylibrary:bundleLibCompileRelease
app:mylibrary:bundleLibCompileReleaseUnitTest
app:mylibrary:bundleLibResDebug
app:mylibrary:bundleLibResRelease
app:mylibrary:bundleLibRuntimeDebug
app:mylibrary:bundleLibRuntimeRelease
app:bundleRelease - Assembles bundle for variant release
app:mylibrary:bundleReleaseAar - Assembles a bundle containing the library in release.
app:bundleReleaseClasses
app:bundleReleaseResources
app:bundleReleaseUnitTestClasses
app:checkDebugAndroidTestDuplicateClasses
app:mylibrary:checkDebugAndroidTestDuplicateClasses
app:checkDebugDuplicateClasses
app:checkDebugManifest
app:mylibrary:checkDebugManifest
app:checkReleaseDuplicateClasses
app:checkReleaseManifest
app:mylibrary:checkReleaseManifest
clean
app:collectDebugDependencies
app:collectReleaseDependencies
app:compileDebugAidl
app:mylibrary:compileDebugAidl
app:compileDebugAndroidTestAidl
app:mylibrary:compileDebugAndroidTestAidl
app:compileDebugAndroidTestJavaWithJavac
app:mylibrary:compileDebugAndroidTestJavaWithJavac
app:compileDebugAndroidTestKotlin - Compiles the debugAndroidTest kotlin.
app:mylibrary:compileDebugAndroidTestKotlin - Compiles the debugAndroidTest kotlin.
app:compileDebugAndroidTestRenderscript
app:mylibrary:compileDebugAndroidTestRenderscript
app:compileDebugAndroidTestShaders
app:mylibrary:compileDebugAndroidTestShaders
app:compileDebugJavaWithJavac
app:mylibrary:compileDebugJavaWithJavac
app:compileDebugKotlin - Compiles the debug kotlin.
app:mylibrary:compileDebugKotlin - Compiles the debug kotlin.
app:compileDebugRenderscript
app:mylibrary:compileDebugRenderscript
app:compileDebugShaders
app:mylibrary:compileDebugShaders
app:compileDebugUnitTestJavaWithJavac
app:mylibrary:compileDebugUnitTestJavaWithJavac
app:compileDebugUnitTestKotlin - Compiles the debugUnitTest kotlin.
app:mylibrary:compileDebugUnitTestKotlin - Compiles the debugUnitTest kotlin.
app:compileLint
app:mylibrary:compileLint
app:compileReleaseAidl
app:mylibrary:compileReleaseAidl
app:compileReleaseJavaWithJavac
app:mylibrary:compileReleaseJavaWithJavac
app:compileReleaseKotlin - Compiles the release kotlin.
app:mylibrary:compileReleaseKotlin - Compiles the release kotlin.
app:compileReleaseRenderscript
app:mylibrary:compileReleaseRenderscript
app:compileReleaseShaders
app:mylibrary:compileReleaseShaders
app:compileReleaseUnitTestJavaWithJavac
app:mylibrary:compileReleaseUnitTestJavaWithJavac
app:compileReleaseUnitTestKotlin - Compiles the releaseUnitTest kotlin.
app:mylibrary:compileReleaseUnitTestKotlin - Compiles the releaseUnitTest kotlin.
app:configureDebugDependencies
app:configureReleaseDependencies
app:consumeConfigAttr
app:mylibrary:consumeConfigAttr
app:createDebugCompatibleScreenManifests
app:mylibrary:createFullJarDebug
app:mylibrary:createFullJarRelease
app:createMockableJar
app:mylibrary:createMockableJar
app:createReleaseCompatibleScreenManifests
app:mylibrary:dexDebug
app:mylibrary:dexRelease
app:dummydebugUnitTest
app:mylibrary:dummydebugUnitTest
app:dummyreleaseUnitTest
app:mylibrary:dummyreleaseUnitTest
app:extractApksForDebug
app:extractApksForRelease
app:extractProguardFiles
app:mylibrary:extractProguardFiles
app:generateDebugAndroidTestAssets
app:mylibrary:generateDebugAndroidTestAssets
app:generateDebugAndroidTestBuildConfig
app:mylibrary:generateDebugAndroidTestBuildConfig
app:generateDebugAndroidTestResources
app:mylibrary:generateDebugAndroidTestResources
app:generateDebugAndroidTestResValues
app:mylibrary:generateDebugAndroidTestResValues
app:generateDebugAndroidTestSources
app:mylibrary:generateDebugAndroidTestSources
app:generateDebugAssets
app:mylibrary:generateDebugAssets
app:generateDebugBuildConfig
app:mylibrary:generateDebugBuildConfig
app:generateDebugFeatureMetadata
app:generateDebugFeatureTransitiveDeps
app:generateDebugResources
app:mylibrary:generateDebugResources
app:generateDebugResValues
app:mylibrary:generateDebugResValues
app:mylibrary:generateDebugRFile
app:generateDebugSources
app:mylibrary:generateDebugSources
app:generateDebugUnitTestAssets
app:mylibrary:generateDebugUnitTestAssets
app:generateDebugUnitTestResources
app:mylibrary:generateDebugUnitTestResources
app:generateDebugUnitTestSources
app:mylibrary:generateDebugUnitTestSources
app:generateReleaseAssets
app:mylibrary:generateReleaseAssets
app:generateReleaseBuildConfig
app:mylibrary:generateReleaseBuildConfig
app:generateReleaseFeatureMetadata
app:generateReleaseFeatureTransitiveDeps
app:generateReleaseResources
app:mylibrary:generateReleaseResources
app:generateReleaseResValues
app:mylibrary:generateReleaseResValues
app:mylibrary:generateReleaseRFile
app:generateReleaseSources
app:mylibrary:generateReleaseSources
app:generateReleaseUnitTestAssets
app:mylibrary:generateReleaseUnitTestAssets
app:generateReleaseUnitTestResources
app:mylibrary:generateReleaseUnitTestResources
app:generateReleaseUnitTestSources
app:mylibrary:generateReleaseUnitTestSources
app:javaPreCompileDebug
app:mylibrary:javaPreCompileDebug
app:javaPreCompileDebugAndroidTest
app:mylibrary:javaPreCompileDebugAndroidTest
app:javaPreCompileDebugUnitTest
app:mylibrary:javaPreCompileDebugUnitTest
app:javaPreCompileRelease
app:mylibrary:javaPreCompileRelease
app:javaPreCompileReleaseUnitTest
app:mylibrary:javaPreCompileReleaseUnitTest
app:mainApkListPersistenceDebug
app:mainApkListPersistenceDebugAndroidTest
app:mylibrary:mainApkListPersistenceDebugAndroidTest
app:mainApkListPersistenceRelease
app:makeApkFromBundleForDebug
app:makeApkFromBundleForRelease
app:mergeDebugAndroidTestAssets
app:mylibrary:mergeDebugAndroidTestAssets
app:mergeDebugAndroidTestGeneratedProguardFiles
app:mylibrary:mergeDebugAndroidTestGeneratedProguardFiles
app:mergeDebugAndroidTestJavaResource
app:mylibrary:mergeDebugAndroidTestJavaResource
app:mergeDebugAndroidTestJniLibFolders
app:mylibrary:mergeDebugAndroidTestJniLibFolders
app:mergeDebugAndroidTestNativeLibs
app:mylibrary:mergeDebugAndroidTestNativeLibs
app:mergeDebugAndroidTestResources
app:mylibrary:mergeDebugAndroidTestResources
app:mergeDebugAndroidTestShaders
app:mylibrary:mergeDebugAndroidTestShaders
app:mergeDebugAssets
app:mylibrary:mergeDebugAssets
app:mylibrary:mergeDebugConsumerProguardFiles
app:mergeDebugGeneratedProguardFiles
app:mylibrary:mergeDebugGeneratedProguardFiles
app:mergeDebugJavaResource
app:mylibrary:mergeDebugJavaResource
app:mergeDebugJniLibFolders
app:mylibrary:mergeDebugJniLibFolders
app:mergeDebugNativeLibs
app:mylibrary:mergeDebugNativeLibs
app:mergeDebugResources
app:mylibrary:mergeDebugResources
app:mergeDebugShaders
app:mylibrary:mergeDebugShaders
app:mergeDexRelease
app:mergeExtDexDebug
app:mergeExtDexDebugAndroidTest
app:mylibrary:mergeExtDexDebugAndroidTest
app:mergeExtDexRelease
app:mergeLibDexDebug
app:mergeLibDexDebugAndroidTest
app:mylibrary:mergeLibDexDebugAndroidTest
app:mergeProjectDexDebug
app:mergeProjectDexDebugAndroidTest
app:mylibrary:mergeProjectDexDebugAndroidTest
app:mergeReleaseAssets
app:mylibrary:mergeReleaseAssets
app:mylibrary:mergeReleaseConsumerProguardFiles
app:mergeReleaseGeneratedProguardFiles
app:mylibrary:mergeReleaseGeneratedProguardFiles
app:mergeReleaseJavaResource
app:mylibrary:mergeReleaseJavaResource
app:mergeReleaseJniLibFolders
app:mylibrary:mergeReleaseJniLibFolders
app:mergeReleaseNativeLibs
app:mylibrary:mergeReleaseNativeLibs
app:mergeReleaseResources
app:mylibrary:mergeReleaseResources
app:mergeReleaseShaders
app:mylibrary:mergeReleaseShaders
app:packageDebug
app:packageDebugAndroidTest
app:mylibrary:packageDebugAndroidTest
app:mylibrary:packageDebugAssets
app:packageDebugBundle
app:mylibrary:packageDebugRenderscript
app:mylibrary:packageDebugResources
app:packageDebugUniversalApk
app:packageRelease
app:mylibrary:packageReleaseAssets
app:packageReleaseBundle
app:mylibrary:packageReleaseRenderscript
app:mylibrary:packageReleaseResources
app:packageReleaseUniversalApk
app:mylibrary:parseDebugLibraryResources
app:mylibrary:parseReleaseLibraryResources
app:preBuild
app:mylibrary:preBuild
app:preDebugAndroidTestBuild
app:mylibrary:preDebugAndroidTestBuild
app:preDebugBuild
app:mylibrary:preDebugBuild
app:preDebugUnitTestBuild
app:mylibrary:preDebugUnitTestBuild
prepareKotlinBuildScriptModel
app:prepareKotlinBuildScriptModel
app:mylibrary:prepareKotlinBuildScriptModel
app:prepareLintJar
app:mylibrary:prepareLintJar
app:prepareLintJarForPublish
app:mylibrary:prepareLintJarForPublish
app:preReleaseBuild
app:mylibrary:preReleaseBuild
app:preReleaseUnitTestBuild
app:mylibrary:preReleaseUnitTestBuild
app:processDebugAndroidTestJavaRes
app:mylibrary:processDebugAndroidTestJavaRes
app:processDebugAndroidTestManifest
app:mylibrary:processDebugAndroidTestManifest
app:processDebugAndroidTestResources
app:mylibrary:processDebugAndroidTestResources
app:processDebugJavaRes
app:mylibrary:processDebugJavaRes
app:processDebugManifest
app:mylibrary:processDebugManifest
app:processDebugResources
app:processDebugUnitTestJavaRes
app:mylibrary:processDebugUnitTestJavaRes
app:processReleaseJavaRes
app:mylibrary:processReleaseJavaRes
app:processReleaseManifest
app:mylibrary:processReleaseManifest
app:processReleaseResources
app:processReleaseUnitTestJavaRes
app:mylibrary:processReleaseUnitTestJavaRes
app:reportBuildArtifactsDebug
app:mylibrary:reportBuildArtifactsDebug
app:reportBuildArtifactsRelease
app:mylibrary:reportBuildArtifactsRelease
app:reportSourceSetTransformAndroidTest
app:mylibrary:reportSourceSetTransformAndroidTest
app:reportSourceSetTransformAndroidTestDebug
app:mylibrary:reportSourceSetTransformAndroidTestDebug
app:reportSourceSetTransformDebug
app:mylibrary:reportSourceSetTransformDebug
app:reportSourceSetTransformMain
app:mylibrary:reportSourceSetTransformMain
app:reportSourceSetTransformRelease
app:mylibrary:reportSourceSetTransformRelease
app:reportSourceSetTransformTest
app:mylibrary:reportSourceSetTransformTest
app:reportSourceSetTransformTestDebug
app:mylibrary:reportSourceSetTransformTestDebug
app:reportSourceSetTransformTestRelease
app:mylibrary:reportSourceSetTransformTestRelease
app:resolveConfigAttr
app:mylibrary:resolveConfigAttr
app:signDebugBundle
app:signingConfigWriterDebug
app:signingConfigWriterDebugAndroidTest
app:mylibrary:signingConfigWriterDebugAndroidTest
app:signingConfigWriterRelease
app:signReleaseBundle
app:stripDebugDebugSymbols
app:mylibrary:stripDebugDebugSymbols
app:stripReleaseDebugSymbols
app:mylibrary:stripReleaseDebugSymbols
app:mylibrary:transformClassesAndResourcesWithSyncLibJarsForDebug
app:mylibrary:transformClassesAndResourcesWithSyncLibJarsForRelease
app:transformClassesWithDexBuilderForDebug
app:transformClassesWithDexBuilderForDebugAndroidTest
app:mylibrary:transformClassesWithDexBuilderForDebugAndroidTest
app:transformClassesWithDexBuilderForRelease
app:mylibrary:transformNativeLibsWithIntermediateJniLibsForDebug
app:mylibrary:transformNativeLibsWithIntermediateJniLibsForRelease
app:mylibrary:transformNativeLibsWithSyncJniLibsForDebug
app:mylibrary:transformNativeLibsWithSyncJniLibsForRelease
app:validateSigningDebug
app:validateSigningDebugAndroidTest
app:mylibrary:validateSigningDebugAndroidTest
app:mylibrary:verifyReleaseResources
app:writeDebugApplicationId
app:writeDebugModuleMetadata
app:writeReleaseApplicationId
app:writeReleaseModuleMetadata
`
