package errorfinder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatcher_Reasons(t *testing.T) {
	tests := []struct {
		name   string
		stdout string
		stderr string
		want   []string
	}{
		{
			name:   "Kotlin compile error",
			stdout: kotlinCompileStdout,
			stderr: kotlinCompileStderr,
			want: []string{
				"Execution failed for task ':app:compileDebugKotlin'.",
				"  A failure occurred while executing org.jetbrains.kotlin.compilerRunner.GradleCompilerRunnerWithWorkers$GradleKotlinCompilerWorkAction",
				"    /bitrise/src/app/src/main/java/com/example/Foo.kt:4:27: Unresolved reference 'bar'.",
			},
		},
		{
			name:   "Kotlin compile error, path and position format",
			stderr: "e: /bitrise/src/app/src/main/java/com/example/Foo.kt: (4, 27): Unresolved reference: bar",
			want:   []string{"/bitrise/src/app/src/main/java/com/example/Foo.kt:4:27: Unresolved reference: bar"},
		},
		{
			name:   "javac error",
			stdout: "> Task :app:compileDebugJavaWithJavac FAILED",
			stderr: javacCompileStderr,
			want: []string{
				"Execution failed for task ':app:compileDebugJavaWithJavac'.",
				"  /bitrise/src/app/src/main/java/com/example/Bar.java:4: cannot find symbol",
			},
		},
		{
			name:   "javac error, with the compiler output under the cause",
			stdout: "> Task :app:compileDebugJavaWithJavac FAILED",
			stderr: javacCompileOutputStderr,
			want: []string{
				"Execution failed for task ':app:compileDebugJavaWithJavac'.",
				"  Compilation failed; see the compiler output below.",
				"    /bitrise/src/app/src/main/java/com/example/Bar.java:4: error: cannot find symbol",
				"    int x = y;",
				"    symbol:   variable y",
				"    location: class Bar",
				"    1 error",
			},
		},
		{
			name:   "resource linking error",
			stdout: "> Task :app:processDebugResources FAILED",
			stderr: resourceLinkingStderr,
			want: []string{
				"Execution failed for task ':app:processDebugResources'.",
				"  A failure occurred while executing com.android.build.gradle.internal.res.LinkApplicationAndroidResourcesTask$TaskAction",
				"    Android resource linking failed",
				"      com.example.app-main-5:/layout/broken_screen.xml:5: error: resource string/missing (aka com.example:string/missing) not found.",
				"      error: failed linking file resources.",
			},
		},
		{
			name:   "more detail lines than the message holds",
			stdout: "> Task :app:processDebugResources FAILED",
			stderr: manyResourceErrorsStderr,
			want: []string{
				"Execution failed for task ':app:processDebugResources'.",
				"  A failure occurred while executing com.android.build.gradle.internal.res.LinkApplicationAndroidResourcesTask$TaskAction",
				"    Android resource linking failed",
				"      com.example.app-main-5:/layout/a.xml:1: error: resource string/a not found.",
				"      com.example.app-main-5:/layout/b.xml:2: error: resource string/b not found.",
				"      com.example.app-main-5:/layout/c.xml:3: error: resource string/c not found.",
				"      com.example.app-main-5:/layout/d.xml:4: error: resource string/d not found.",
				"      com.example.app-main-5:/layout/e.xml:5: error: resource string/e not found.",
				"      com.example.app-main-5:/layout/f.xml:6: error: resource string/f not found.",
				"      com.example.app-main-5:/layout/g.xml:7: error: resource string/g not found.",
				"      com.example.app-main-5:/layout/h.xml:8: error: resource string/h not found.",
				"      com.example.app-main-5:/layout/i.xml:9: error: resource string/i not found.",
				"      com.example.app-main-5:/layout/j.xml:10: error: resource string/j not found.",
				"      ... and 2 more lines",
			},
		},
		{
			name:   "more compile errors than the message holds",
			stderr: manyKotlinErrorsStderr,
			want: []string{
				"/bitrise/src/app/A.kt:1:1: Unresolved reference 'a'.",
				"/bitrise/src/app/B.kt:2:2: Unresolved reference 'b'.",
				"/bitrise/src/app/C.kt:3:3: Unresolved reference 'c'.",
				"... and 1 more compile errors",
			},
		},
		{
			name:   "the same compile error reported twice",
			stderr: duplicateKotlinErrorStderr,
			want:   []string{"/bitrise/src/app/A.kt:1:1: Unresolved reference 'a'."},
		},
		{
			name:   "unknown task",
			stderr: unknownTaskStderr,
			want:   []string{"Task 'assembleFoo' not found in root project 'MyApp'. Some candidates are: 'assembleDebug'."},
		},
		{
			name:   "unknown task, with no cause under the problem",
			stderr: selectionFailedStderr,
			want: []string{
				"Selection failed",
				"  Task 'assembleFoo' not found in root project 'MyApp' and its subprojects.",
			},
		},
		{
			name:   "unresolved dependency",
			stdout: "> Task :app:compileDebugJavaWithJavac FAILED",
			stderr: unresolvedDependencyStderr,
			want: []string{
				"Execution failed for task ':app:compileDebugJavaWithJavac'.",
				"  Could not resolve all files for configuration ':app:debugCompileClasspath'.",
				"    Could not find com.example:missing:1.0.0.",
			},
		},
		{
			name:   "a task line of the other stream inside the block",
			stderr: interleavedTaskLineStderr,
			want: []string{
				"Execution failed for task ':app:processDebugNavigationResources'.",
				"  Could not resolve all files for configuration ':app:debugRuntimeClasspath'.",
				"    Could not find com.example:missing:1.0.0.",
			},
		},
		{
			name:   "failing unit tests",
			stdout: failingTestsStdout,
			stderr: failingTestsStderr,
			want: []string{
				"Execution failed for task ':app:testDebugUnitTest'.",
				"  com.example.ExampleUnitTest > addition_isCorrect FAILED",
				"  org.junit.ComparisonFailure at ExampleUnitTest.java:16",
				"  com.example.ExampleUnitTest > subtraction_isCorrect FAILED",
				"  java.lang.AssertionError at ExampleUnitTest.java:24",
			},
		},
		{
			name:   "failed task only",
			stdout: "> Task :app:externalNativeBuildDebug FAILED",
			want:   []string{"Gradle task :app:externalNativeBuildDebug failed"},
		},
		{
			name:   "nothing to report",
			stdout: "> Task :app:assembleDebug\nBUILD SUCCESSFUL in 12s",
			stderr: "Warning: SDK processing. This version only understands SDK XML versions up to 3",
			want:   nil,
		},
		{
			name:   "ANSI colored output",
			stderr: "\x1b[31;1me: file:///bitrise/src/app/A.kt:1:1 Unresolved reference 'a'.\x1b[0m",
			want:   []string{"/bitrise/src/app/A.kt:1:1: Unresolved reference 'a'."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewMatcher()

			consume(matcher, tt.stdout)
			consume(matcher, tt.stderr)

			require.Equal(t, tt.want, matcher.Reasons())
		})
	}
}

func consume(matcher *Matcher, output string) {
	if output == "" {
		return
	}

	for _, line := range strings.Split(strings.TrimSuffix(output, "\n"), "\n") {
		matcher.ConsumeLine(line)
	}
}
