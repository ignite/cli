package uilog

import (
	"io"
	"os"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/lineprefixer"
	"github.com/ignite/cli/ignite/pkg/cliui/prefixgen"
	"github.com/ignite/cli/ignite/pkg/xio"
)

const (
	defaultVerboseLabel      = "ignite"
	defaultVerboseLabelColor = colors.Red
)

// Verbosity enumerates possible verbosity levels for CLI output.
type Verbosity uint8

const (
	VerbositySilent = iota
	VerbosityDefault
	VerbosityVerbose
)

// Outputer defines an interface for logging output creation.
type Outputer interface {
	// NewOutput returns a new logging output.
	NewOutput(label, color string) Output

	// Verbosity returns the current verbosity level for the logging output.
	Verbosity() Verbosity
}

// Output stores writers for standard output and error.
type Output struct {
	verbosity Verbosity
	stdout    io.WriteCloser
	stderr    io.WriteCloser
}

// Stdout returns the standard output writer.
func (o Output) Stdout() io.WriteCloser {
	return o.stdout
}

// Stderr returns the standard error writer.
func (o Output) Stderr() io.WriteCloser {
	return o.stderr
}

// Verbosity returns the log output verbosity.
func (o Output) Verbosity() Verbosity {
	return o.verbosity
}

type option struct {
	stdout            io.WriteCloser
	stderr            io.WriteCloser
	verbosity         Verbosity
	verboseLabel      string
	verboseLabelColor string
}

// Option configures log output options.
type Option func(*option)

// Verbose changes the log output to be prefixed with "ignite".
func Verbose() Option {
	return func(o *option) {
		o.verbosity = VerbosityVerbose
		o.verboseLabel = defaultVerboseLabel
		o.verboseLabelColor = defaultVerboseLabelColor
	}
}

// CustomVerbose changes the log output to be prefixed with a custom label.
func CustomVerbose(label, color string) Option {
	return func(o *option) {
		o.verbosity = VerbosityVerbose
		o.verboseLabel = label
		o.verboseLabelColor = color
	}
}

// Silent creates a log output that doesn't print any of the written lines.
func Silent() Option {
	return func(o *option) {
		o.verbosity = VerbositySilent
	}
}

// WithStdout sets a custom writer to use instead of the default `os.Stdout`.
func WithStdout(r io.WriteCloser) Option {
	return func(o *option) {
		o.stdout = r
	}
}

// WithStderr sets a custom writer to use instead of the default `os.Stderr`.
func WithStderr(r io.WriteCloser) Option {
	return func(o *option) {
		o.stderr = r
	}
}

// NewOutput creates a new log output.
// By default, the new output uses the default OS stdout and stderr to
// initialize the outputs with a default verbosity that doesn't change
// the output.
func NewOutput(options ...Option) (out Output) {
	o := option{
		verbosity: VerbosityDefault,
		stdout:    os.Stdout,
		stderr:    os.Stderr,
	}

	for _, apply := range options {
		apply(&o)
	}

	out.verbosity = o.verbosity

	switch o.verbosity {
	case VerbositySilent:
		out.stdout = xio.NopWriteCloser(io.Discard)
		out.stderr = xio.NopWriteCloser(io.Discard)
	case VerbosityVerbose:
		// Function to add a custom prefix to each log output
		prefixer := func(w io.Writer) *lineprefixer.Writer {
			options := prefixgen.Common(prefixgen.Color(o.verboseLabelColor))
			prefix := prefixgen.New(o.verboseLabel, options...).Gen()

			return lineprefixer.NewWriter(w, func() string { return prefix })
		}

		out.stdout = xio.NopWriteCloser(prefixer(o.stdout))
		out.stderr = xio.NopWriteCloser(prefixer(o.stderr))
	default:
		out.stdout = o.stdout
		out.stderr = o.stderr
	}

	return out
}
