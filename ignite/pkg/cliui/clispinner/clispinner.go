package clispinner

import (
	"io"
	"os"

	"golang.org/x/term"
)

// DefaultText defines the default spinner text.
const DefaultText = "Initializing..."

type (
	Spinner interface {
		SetText(text string) Spinner
		SetPrefix(text string) Spinner
		SetCharset(charset []string) Spinner
		SetColor(color string) Spinner
		Start() Spinner
		Stop() Spinner
		IsActive() bool
		Writer() io.Writer
	}

	Option func(*Options)

	Options struct {
		writer  io.Writer
		text    string
		charset []string
	}
)

// WithWriter configures an output for a spinner.
func WithWriter(w io.Writer) Option {
	return func(options *Options) {
		options.writer = w
	}
}

// WithText configures the spinner text.
func WithText(text string) Option {
	return func(options *Options) {
		options.text = text
	}
}

// WithCharset configures the spinner charset.
func WithCharset(charset []string) Option {
	return func(options *Options) {
		options.charset = charset
	}
}

// New creates a new spinner.
func New(options ...Option) Spinner {
	o := Options{}
	for _, apply := range options {
		apply(&o)
	}

	if isRunningInTerminal(o.writer) {
		return newTermSpinner(o)
	}
	return newSimpleSpinner(o)
}

// isRunningInTerminal check if the writer file descriptor is a terminal.
func isRunningInTerminal(w io.Writer) bool {
	if w == nil {
		return term.IsTerminal(int(os.Stdout.Fd()))
	}
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}
