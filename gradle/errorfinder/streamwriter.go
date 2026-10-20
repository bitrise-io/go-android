package errorfinder

import "bytes"

// maxLineLength bounds how much of a line is held back while its newline is missing. Failure lines are
// short, so a longer line is not one: the rest of it is dropped instead of buffered.
const maxLineLength = 64 * 1024

// streamWriter feeds one output stream to the collector, line by line.
type streamWriter struct {
	collector *Collector
	pending   []byte
	skipping  bool
}

// Write never fails and always reports the whole chunk as written, so it can not interrupt the
// io.MultiWriter it is part of and cost the build its log output.
func (w *streamWriter) Write(p []byte) (int, error) {
	w.collector.mu.Lock()
	defer w.collector.mu.Unlock()

	n := len(p)
	for len(p) > 0 {
		line, rest, found := bytes.Cut(p, []byte{'\n'})
		w.appendPending(line)
		if !found {
			break
		}

		w.finishLine()
		p = rest
	}

	return n, nil
}

func (w *streamWriter) appendPending(b []byte) {
	if w.skipping {
		return
	}
	if len(w.pending)+len(b) > maxLineLength {
		w.pending = w.pending[:0]
		w.skipping = true

		return
	}

	w.pending = append(w.pending, b...)
}

func (w *streamWriter) flush() {
	if len(w.pending) == 0 && !w.skipping {
		return
	}

	w.finishLine()
}

// finishLine hands the completed line to the parser. A skipped (overlong) line is not a failure line,
// but it still ends an open block, so it is handed over as an empty line.
func (w *streamWriter) finishLine() {
	line := ""
	if !w.skipping {
		line = string(w.pending)
	}
	w.pending = w.pending[:0]
	w.skipping = false

	w.collector.parser.consumeLine(line)
}
