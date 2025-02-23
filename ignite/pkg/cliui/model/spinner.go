package cliuimodel

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// ColorSpinner defines the foreground color for the spinner.
const ColorSpinner = "#3465A4"

// Spinner defines the spinner model animation.
var Spinner = spinner.Spinner{
	Frames: []string{"◢ ", "◣ ", "◤ ", "◥ "},
	FPS:    time.Second / 5,
}

// NewSpinner returns a new spinner model.
func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = Spinner
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSpinner))
	return s
}
