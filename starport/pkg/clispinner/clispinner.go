package clispinner

import (
	"time"

	"github.com/briandowns/spinner"
)

var (
	refreshRate = 100 * time.Millisecond
	charset     = spinner.CharSets[42]
	color       = "blue"
)

type Spinner struct {
	sp *spinner.Spinner
}

// New creates a new spinner.
func New() *Spinner {
	sp := spinner.New(charset, refreshRate)
	sp.Color(color)
	return &Spinner{
		sp: sp,
	}
}

// SetText sets the text for spinner.
func (s *Spinner) SetText(text string) {
	s.sp.Suffix = " " + text
}

// Start starts spinning.
func (s *Spinner) Start() { s.sp.Start() }

// Stop stops spinning.
func (s *Spinner) Stop() { s.sp.Stop() }
