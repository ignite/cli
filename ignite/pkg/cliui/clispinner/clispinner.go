package clispinner

import (
	"io"
	"time"

	"github.com/theckman/yacspin"
)

// DefaultText defines the default spinner text.
const DefaultText = "Initializing..."

var (
	refreshRate  = time.Millisecond * 200
	charset      = yacspin.CharSets[4]
	spinnerColor = "blue"
)

type Spinner struct {
	sp *yacspin.Spinner
}

type (
	Option func(*Options)

	Options struct {
		writer io.Writer
		text   string
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

// New creates a new spinner.
func New(options ...Option) (*Spinner, error) {
	o := Options{}
	for _, apply := range options {
		apply(&o)
	}

	text := o.text
	if text == "" {
		text = DefaultText
	}

	cfg := yacspin.Config{
		Frequency:  refreshRate,
		CharSet:    yacspin.CharSets[59],
		Message:    text,
		Colors:     []string{spinnerColor},
		StopColors: []string{"fgGreen"},
	}

	if o.writer != nil {
		cfg.Writer = o.writer
	}

	sp, err := yacspin.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Spinner{
		sp: sp,
	}, nil
}

// SetText sets the text for spinner.
func (s *Spinner) SetText(text string) *Spinner {
	s.sp.Message(" " + text)
	return s
}

// SetPrefix sets the prefix for spinner.
func (s *Spinner) SetPrefix(text string) *Spinner {
	s.sp.Prefix(text + " ")
	return s
}

// SetCharset sets the prefix for spinner.
func (s *Spinner) SetCharset(charset []string) *Spinner {
	_ = s.sp.CharSet(charset)
	return s
}

// SetColor sets the prefix for spinner.
func (s *Spinner) SetColor(color string) *Spinner {
	_ = s.sp.Colors(color)
	return s
}

// Start starts spinning.
func (s *Spinner) Start() *Spinner {
	_ = s.sp.Start()
	return s
}

// Stop stops spinning.
func (s *Spinner) Stop() *Spinner {
	_ = s.sp.Stop()
	s.sp.Prefix("")
	_ = s.sp.Colors(spinnerColor)
	_ = s.sp.CharSet(charset)
	_ = s.sp.Stop()
	return s
}

func (s *Spinner) IsActive() bool {
	return s.sp.Status() != yacspin.SpinnerRunning && s.sp.Status() != yacspin.SpinnerStarting
}
