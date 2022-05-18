package cliui

import (
	"io"

	"github.com/docker/docker/pkg/ioutils"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/lineprefixer"
	"github.com/ignite-hq/cli/ignite/pkg/cliui/prefixgen"
)

const (
	defaultLogStreamLabel = "ignite"
	defaultLogStreamColor = 91
)

// LogStreamer specifies that object could create new LogStream objects
type LogStreamer interface {
	NewLogStream(label string, color uint8) (logStream LogStream)
}

// LogStream API of Session which provides ability to write logs to io.WriteCloser type object
type LogStream struct {
	stdout io.WriteCloser
	stderr io.WriteCloser
}

// Stdout returns LogStream stdout writer
func (ls LogStream) Stdout() io.WriteCloser {
	return ls.stdout
}

// Stderr returns LogStream stderr writer
func (ls LogStream) Stderr() io.WriteCloser {
	return ls.stderr
}

// NewLogStream creates new LogStream object bound to the Session instance
func (s Session) NewLogStream(label string, color uint8) (logStream LogStream) {
	prefixed := func(w io.Writer) *lineprefixer.Writer {
		options := prefixgen.Common(prefixgen.Color(color))
		prefixStr := prefixgen.New(label, options...).Gen()
		return lineprefixer.NewWriter(w, func() string { return prefixStr })
	}

	verbosity := s.verbosity
	if s.isDefaultLogStreamInitialised && verbosity != VerbosityVerbose {
		verbosity = VerbositySilent
	}
	s.isDefaultLogStreamInitialised = true

	switch verbosity {
	case VerbositySilent:
		logStream.stdout = ioutils.NopWriteCloser(io.Discard)
		logStream.stderr = ioutils.NopWriteCloser(io.Discard)
	case VerbosityVerbose:
		logStream.stdout = prefixed(s.stdout)
		logStream.stderr = prefixed(s.stderr)
	default:
		logStream.stdout = s.stdout
		logStream.stderr = s.stderr
	}

	return
}
