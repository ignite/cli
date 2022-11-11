package clispinner

import (
	"io"
	"time"

	"github.com/briandowns/spinner"
)

// DefaultText defines the default spinner text.
const DefaultText = "Initializing..."

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
func New(options ...Option) *Spinner {
	o := Options{}
	for _, apply := range options {
		apply(&o)
	}

	text := o.text
	if text == "" {
		text = DefaultText
	}

	spOptions := []spinner.Option{
		spinner.WithColor(spinnerColor),
		spinner.WithSuffix(" " + text),
	}

	if o.writer != nil {
		spOptions = append(spOptions, spinner.WithWriter(o.writer))
	}

	return &Spinner{
		sp: spinner.New(charset, refreshRate, spOptions...),
	}
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
