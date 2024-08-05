package androidartifact

import "fmt"

type testLogger struct{}

func (l *testLogger) Warnf(format string, args ...interface{}) { fmt.Printf(format, args...) }
func (l *testLogger) AABParseWarnf(_, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (l *testLogger) APKParseWarnf(_, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
