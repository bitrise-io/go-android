package errorfinder

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamWriter_linesSplitAcrossWrites(t *testing.T) {
	collector := NewCollector()

	for _, chunk := range splitEvery(kotlinCompileStderr, 7) {
		write(t, collector.StderrWriter(), chunk)
	}

	require.Equal(t, []string{
		"Execution failed for task ':app:compileDebugKotlin'.",
		"/bitrise/src/app/src/main/java/com/example/Foo.kt:4:27: Unresolved reference 'bar'.",
	}, collector.Errors())
}

func TestStreamWriter_lastLineWithoutNewline(t *testing.T) {
	collector := NewCollector()

	write(t, collector.StderrWriter(), "e: file:///bitrise/src/app/A.kt:1:1 Unresolved reference 'a'.")

	require.Equal(t, []string{"/bitrise/src/app/A.kt:1:1: Unresolved reference 'a'."}, collector.Errors())
	require.Equal(t, []string{"/bitrise/src/app/A.kt:1:1: Unresolved reference 'a'."}, collector.Errors())
}

func TestStreamWriter_overlongLineDropped(t *testing.T) {
	collector := NewCollector()

	write(t, collector.StderrWriter(), "e: file:///bitrise/src/app/A.kt:1:1 "+strings.Repeat("x", maxLineLength)+"\n")
	write(t, collector.StderrWriter(), unknownTaskStderr)

	require.Equal(t, []string{"Task 'assembleFoo' not found in root project 'MyApp'. Some candidates are: 'assembleDebug'."}, collector.Errors())
}

func TestStreamWriter_bothStreamsWrittenAtOnce(t *testing.T) {
	collector := NewCollector()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		write(t, collector.StdoutWriter(), strings.Repeat(kotlinCompileStdout, 20))
	}()
	go func() {
		defer wg.Done()
		write(t, collector.StderrWriter(), strings.Repeat(kotlinCompileStderr, 20))
	}()

	wg.Wait()

	require.Equal(t, []string{
		"Execution failed for task ':app:compileDebugKotlin'.",
		"/bitrise/src/app/src/main/java/com/example/Foo.kt:4:27: Unresolved reference 'bar'.",
	}, collector.Errors())
}

func splitEvery(s string, n int) []string {
	var chunks []string
	for len(s) > n {
		chunks = append(chunks, s[:n])
		s = s[n:]
	}

	return append(chunks, s)
}
