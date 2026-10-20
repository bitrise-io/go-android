package errorfinder

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollector_Errors(t *testing.T) {
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
				"/bitrise/src/app/src/main/java/com/example/Foo.kt:4:27: Unresolved reference 'bar'.",
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
				"/bitrise/src/app/src/main/java/com/example/Bar.java:4: cannot find symbol",
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
			name:   "unresolved dependency",
			stdout: "> Task :app:compileDebugJavaWithJavac FAILED",
			stderr: unresolvedDependencyStderr,
			want: []string{
				"Execution failed for task ':app:compileDebugJavaWithJavac'.",
				"> Could not resolve all files for configuration ':app:debugCompileClasspath'.",
				"> Could not find com.example:missing:1.0.0.",
			},
		},
		{
			name:   "failing unit tests",
			stdout: failingTestsStdout,
			stderr: failingTestsStderr,
			want: []string{
				"Execution failed for task ':app:testDebugUnitTest'.",
				"com.example.ExampleUnitTest > addition_isCorrect FAILED",
				"  org.junit.ComparisonFailure at ExampleUnitTest.java:16",
				"com.example.ExampleUnitTest > subtraction_isCorrect FAILED",
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
		{
			name:   "Windows style line endings",
			stderr: "e: file:///bitrise/src/app/A.kt:1:1 Unresolved reference 'a'.\r\n",
			want:   []string{"/bitrise/src/app/A.kt:1:1: Unresolved reference 'a'."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCollector()

			write(t, collector.StdoutWriter(), tt.stdout)
			write(t, collector.StderrWriter(), tt.stderr)

			require.Equal(t, tt.want, collector.Errors())
		})
	}
}

func TestCollector_Errors_cappedToStayUsefulInTheStepError(t *testing.T) {
	collector := NewCollector()

	write(t, collector.StderrWriter(), "* What went wrong:\n")
	for i := range 10 {
		write(t, collector.StderrWriter(), fmt.Sprintf("> cause %d: %s\n", i, strings.Repeat("detail ", 30)))
	}

	errors := collector.Errors()

	require.Equal(t, truncatedMarker, errors[len(errors)-1])
	require.LessOrEqual(t, len(strings.Join(errors, "\n")), maxErrorBytes+len(truncatedMarker)+1)
}

func TestCollector_Errors_singleLineLongerThanTheCap(t *testing.T) {
	collector := NewCollector()

	write(t, collector.StderrWriter(), "* What went wrong:\n"+strings.Repeat("é", maxErrorBytes)+"\n")

	errors := collector.Errors()

	require.Len(t, errors, 2)
	require.Equal(t, strings.Repeat("é", maxErrorBytes/2), errors[0])
	require.Equal(t, truncatedMarker, errors[1])
}
