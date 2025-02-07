package xos

import (
	"io"
	"os"

	"github.com/creack/pty"
)

// StdOutRedirect manages the redirection of stdout/stderr with a PTY.
// The returned writer can be used to send data, and the cleanup function handles resource cleanup.
func StdOutRedirect(writer io.Writer) (*os.File, func() error, error) {
	// Create a new pseudo-terminal.
	aux := writer
	pty, tty, err := pty.Open()
	if err != nil {
		return nil, nil, err
	}

	// Redirect os.Stdout and os.Stderr to the TTY.
	writer = tty

	// Start a goroutine to forward the output from the PTY to the real terminal
	go func() {
		_, _ = io.Copy(aux, pty) // Copy all output from the PTY to the real terminal
	}()

	cleanup := func() error {
		if err := pty.Close(); err != nil {
			return err
		}
		return tty.Close()
	}

	return tty, cleanup, nil
}
