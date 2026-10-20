package errorfinder

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func write(t *testing.T, writer io.Writer, s string) {
	t.Helper()

	if s == "" {
		return
	}

	n, err := writer.Write([]byte(s))
	require.NoError(t, err)
	require.Equal(t, len(s), n)
}

const kotlinCompileStdout = `> Task :app:preBuild UP-TO-DATE
> Task :app:compileDebugKotlin FAILED
`

const kotlinCompileStderr = `e: file:///bitrise/src/app/src/main/java/com/example/Foo.kt:4:27 Unresolved reference 'bar'.

FAILURE: Build failed with an exception.

* What went wrong:
Execution failed for task ':app:compileDebugKotlin'.
> A failure occurred while executing org.jetbrains.kotlin.compilerRunner.GradleCompilerRunnerWithWorkers$GradleKotlinCompilerWorkAction
   > Compilation error. See log for more details

* Try:
> Run with --stacktrace option to get the stack trace.

* Get more help at https://help.gradle.org

BUILD FAILED in 5s
`

const javacCompileStderr = `/bitrise/src/app/src/main/java/com/example/Bar.java:4: error: cannot find symbol
        int x = y;
                ^
  symbol:   variable y
  location: class Bar
1 error

FAILURE: Build failed with an exception.

* What went wrong:
Execution failed for task ':app:compileDebugJavaWithJavac'.
> Compilation failed; see the compiler error output for details.

* Try:
> Run with --stacktrace option to get the stack trace.
`

const manyKotlinErrorsStderr = `e: file:///bitrise/src/app/A.kt:1:1 Unresolved reference 'a'.
e: file:///bitrise/src/app/B.kt:2:2 Unresolved reference 'b'.
e: file:///bitrise/src/app/C.kt:3:3 Unresolved reference 'c'.
e: file:///bitrise/src/app/D.kt:4:4 Unresolved reference 'd'.
`

const duplicateKotlinErrorStderr = `e: file:///bitrise/src/app/A.kt:1:1 Unresolved reference 'a'.
e: file:///bitrise/src/app/A.kt:1:1 Unresolved reference 'a'.
`

const unknownTaskStderr = `
FAILURE: Build failed with an exception.

* What went wrong:
Task 'assembleFoo' not found in root project 'MyApp'. Some candidates are: 'assembleDebug'.

* Try:
> Run gradlew tasks to get a list of available tasks.
`

const unresolvedDependencyStderr = `
FAILURE: Build failed with an exception.

* What went wrong:
Execution failed for task ':app:compileDebugJavaWithJavac'.
> Could not resolve all files for configuration ':app:debugCompileClasspath'.
   > Could not find com.example:missing:1.0.0.
     Searched in the following locations:
       - https://repo.maven.apache.org/maven2/com/example/missing/1.0.0/missing-1.0.0.pom
     Required by:
         project :app

* Try:
> Run with --stacktrace option to get the stack trace.
`

const failingTestsStdout = `> Task :app:testDebugUnitTest

com.example.ExampleUnitTest > addition_isCorrect FAILED
    org.junit.ComparisonFailure at ExampleUnitTest.java:16

com.example.ExampleUnitTest > subtraction_isCorrect FAILED
    java.lang.AssertionError at ExampleUnitTest.java:24

2 tests completed, 2 failed

> Task :app:testDebugUnitTest FAILED
`

const failingTestsStderr = `
FAILURE: Build failed with an exception.

* What went wrong:
Execution failed for task ':app:testDebugUnitTest'.
> There were failing tests. See the report at: file:///bitrise/src/app/build/reports/tests/testDebugUnitTest/index.html

* Try:
> Run with --scan to get full insights.
`
