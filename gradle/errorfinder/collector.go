// Package errorfinder collects the lines of a Gradle run's output that say why the build failed, so a
// step can report them as its failure reason instead of a bare exit status.
package errorfinder

import (
	"io"
	"strings"
	"sync"
)

const (
	// maxErrorBytes leaves room for the step's own wrapping layers inside the 1 KB cap of
	// BITRISE_FAILED_STEP_ERROR_MESSAGE.
	maxErrorBytes = 800

	truncatedMarker = "... (truncated)"
)

// Collector keeps the lines of a Gradle run's output that say why it failed.
//
// Gradle splits those lines over both streams: the FAILURE block and the compiler errors go to stderr,
// the failed task and failing test lines to stdout. os/exec copies the two streams on separate
// goroutines, so each stream gets its own writer, holding its own partial line, over shared state.
type Collector struct {
	mu     sync.Mutex
	parser parser

	stdout *streamWriter
	stderr *streamWriter
}

// NewCollector ...
func NewCollector() *Collector {
	collector := &Collector{}
	collector.stdout = &streamWriter{collector: collector}
	collector.stderr = &streamWriter{collector: collector}

	return collector
}

// StdoutWriter returns the writer the command's stdout has to be teed into.
func (c *Collector) StdoutWriter() io.Writer {
	return c.stdout
}

// StderrWriter returns the writer the command's stderr has to be teed into.
func (c *Collector) StderrWriter() io.Writer {
	return c.stderr
}

// Errors returns the failure reason found so far, or nothing when no line matched, in which case the
// caller keeps its own generic message rather than showing a guess.
//
// The end of the output is treated as the end of a line, so it is safe to call more than once.
func (c *Collector) Errors() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stdout.flush()
	c.stderr.flush()

	return capLines(c.parser.errors())
}

func capLines(lines []string) []string {
	var capped []string
	remaining := maxErrorBytes

	for _, line := range lines {
		if len(line) > remaining {
			if len(capped) == 0 {
				capped = append(capped, strings.ToValidUTF8(line[:remaining], ""))
			}
			capped = append(capped, truncatedMarker)

			break
		}

		capped = append(capped, line)
		remaining -= len(line) + 1
	}

	return capped
}
