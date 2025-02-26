package clispinner

import (
	"io"
	"time"

	"github.com/briandowns/spinner"
)

var (
	terminalCharset     = spinner.CharSets[4]
	terminalRefreshRate = time.Millisecond * 200
	terminalColor       = "blue"
)

type TermSpinner struct {
	sp      *spinner.Spinner
	charset []string
}

// newTermSpinner creates a new terminal spinner.
func newTermSpinner(o Options) *TermSpinner {
	text := o.text
	if text == "" {
		text = DefaultText
	}

	charset := o.charset
	if len(charset) == 0 {
		charset = terminalCharset
	}

	spOptions := []spinner.Option{
		spinner.WithColor(terminalColor),
		spinner.WithSuffix(" " + text),
	}

	if o.writer != nil {
		spOptions = append(spOptions, spinner.WithWriter(o.writer))
	}

	return &TermSpinner{
		sp:      spinner.New(charset, terminalRefreshRate, spOptions...),
		charset: charset,
	}
}

// SetText sets the text for spinner.
func (s *TermSpinner) SetText(text string) Spinner {
	s.sp.Lock()
	s.sp.Suffix = " " + text
	s.sp.Unlock()
	return s
}

// SetPrefix sets the prefix for spinner.
func (s *TermSpinner) SetPrefix(text string) Spinner {
	s.sp.Lock()
	s.sp.Prefix = text + " "
	s.sp.Unlock()
	return s
}

// SetCharset sets the prefix for spinner.
func (s *TermSpinner) SetCharset(charset []string) Spinner {
	s.sp.UpdateCharSet(charset)
	return s
}

// SetColor sets the prefix for spinner.
func (s *TermSpinner) SetColor(color string) Spinner {
	_ = s.sp.Color(color)
	return s
}

// Start starts spinning.
func (s *TermSpinner) Start() Spinner {
	s.sp.Start()
	return s
}

// Stop stops spinning.
func (s *TermSpinner) Stop() Spinner {
	s.sp.Stop()
	s.sp.Prefix = ""
	_ = s.sp.Color(terminalColor)
	s.sp.UpdateCharSet(s.charset)
	s.sp.Stop()
	return s
}

// IsActive returns whether the spinner is currently active.
func (s *TermSpinner) IsActive() bool {
	return s.sp.Active()
}

// Writer returns the spinner writer.
func (s *TermSpinner) Writer() io.Writer {
	return s.sp.Writer
}
