package model

import (
	"fmt"

	"github.com/muesli/reflow/indent"
)

const (
	// ColorSpinner defines the foreground color for the spinner.
	ColorSpinner = "#8878FF"

	// EOL defines the rune for the end of line.
	EOL = '\n'
)

const (
	defaultIndent = 2
)

type (
	// ErrorMsg defines a message for error.
	ErrorMsg struct{ error }
)

// FormatView formats a model view padding and indentation.
func FormatView(view string) string {
	// Indent the view lines
	view = indent.String(view, defaultIndent)
	// Add top and bottom paddings
	return fmt.Sprintf("\n%s\n", view)
}
