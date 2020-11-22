package clispinner

import (
	"time"

	"github.com/briandowns/spinner"
)

var (
	refreshRate = time.Second
	charset     = spinner.CharSets[4]
	color       = "blue"
)

type Spinner struct {
	sp *spinner.Spinner
}

// New creates a new spinner.
func New() *Spinner {
	sp := spinner.New(charset, refreshRate)
	sp.Color(color)
	s := &Spinner{
		sp: sp,
	}
	s.SetText("Initializing...")
	return s
}

// SetText sets the text for spinner.
func (s *Spinner) SetText(text string) *Spinner {
	s.sp.Suffix = " " + text
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
	return s
}
