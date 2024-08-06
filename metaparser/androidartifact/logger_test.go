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

func (l *testLogger) Infof(format string, args ...interface{})   { fmt.Printf(format, args...) }
func (l *testLogger) Printf(format string, args ...interface{})  { fmt.Printf(format, args...) }
func (l *testLogger) Donef(format string, args ...interface{})   { fmt.Printf(format, args...) }
func (l *testLogger) Debugf(format string, args ...interface{})  { fmt.Printf(format, args...) }
func (l *testLogger) Errorf(format string, args ...interface{})  { fmt.Printf(format, args...) }
func (l *testLogger) TInfof(format string, args ...interface{})  { fmt.Printf(format, args...) }
func (l *testLogger) TWarnf(format string, args ...interface{})  { fmt.Printf(format, args...) }
func (l *testLogger) TPrintf(format string, args ...interface{}) { fmt.Printf(format, args...) }
func (l *testLogger) TDonef(format string, args ...interface{})  { fmt.Printf(format, args...) }
func (l *testLogger) TDebugf(format string, args ...interface{}) { fmt.Printf(format, args...) }
func (l *testLogger) TErrorf(format string, args ...interface{}) { fmt.Printf(format, args...) }
func (l *testLogger) Println()                                   { fmt.Println() }
func (l *testLogger) EnableDebugLog(_ bool)                      {}
