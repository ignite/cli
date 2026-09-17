package os

import "time"

type (
	stdoutWriter struct{}
	stderrWriter struct{}
)

var (
	Stdout = &stdoutWriter{}
	Stderr = &stderrWriter{}
)

func (w *stdoutWriter) Write(p []byte) (n int, err error) {
	return write(p, false), nil
}

func (w *stderrWriter) Write(p []byte) (n int, err error) {
	return write(p, true), nil
}

func Sleep(duration time.Duration) {
	sleep(int64(duration))
}

// native bindings
func write(p []byte, isStderr bool) int
func sleep(duration int64)
