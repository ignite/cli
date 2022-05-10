package clispinner

import (
	"io"
	"time"

	"github.com/briandowns/spinner"
)

var (
	refreshRate  = time.Millisecond * 200
	charset      = spinner.CharSets[4]
	spinnerColor = "blue"
)

type Spinner struct {
	sp *spinner.Spinner
}

type (
	Option func(*Options)

	Options struct {
		writer io.Writer
	}
)

// WithWriter configures an output for a spinner
func WithWriter(w io.Writer) Option {
	return func(options *Options) {
		options.writer = w
	}
}

// New creates a new spinner.
func New(options ...Option) *Spinner {
	o := Options{}
	for _, apply := range options {
		apply(&o)
	}

	underlyingSpinnerOptions := []spinner.Option{}
	if o.writer != nil {
		underlyingSpinnerOptions = append(underlyingSpinnerOptions, spinner.WithWriter(o.writer))
	}

	sp := spinner.New(charset, refreshRate, underlyingSpinnerOptions...)

	sp.Color(spinnerColor)
	s := &Spinner{
		sp: sp,
	}
	return s.SetText("Initializing...")
}

// SetText sets the text for spinner.
func (s *Spinner) SetText(text string) *Spinner {
	s.sp.Lock()
	s.sp.Suffix = " " + text
	s.sp.Unlock()
	return s
}

// SetPrefix sets the prefix for spinner.
func (s *Spinner) SetPrefix(text string) *Spinner {
	s.sp.Lock()
	s.sp.Prefix = text + " "
	s.sp.Unlock()
	return s
}

// SetCharset sets the prefix for spinner.
func (s *Spinner) SetCharset(charset []string) *Spinner {
	s.sp.UpdateCharSet(charset)
	return s
}

// SetColor sets the prefix for spinner.
func (s *Spinner) SetColor(color string) *Spinner {
	s.sp.Color(color)
	return s
}

// Start starts spinning.
func (s *Spinner) Start() *Spinner {
	s.sp.Start()
	return s
}

// Stop stops spinning.
func (s *Spinner) Stop() *Spinner {
	s.sp.Stop()
	s.sp.Prefix = ""
	s.sp.Color(spinnerColor)
	s.sp.UpdateCharSet(charset)
	s.sp.Stop()
	return s
}

func (s *Spinner) IsActive() bool {
	return s.sp.Active()
}
