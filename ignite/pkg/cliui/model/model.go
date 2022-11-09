package cliuimodel

import (
	"fmt"

	"github.com/muesli/reflow/indent"
)

const (
	// EOL defines the rune for the end of line.
	EOL = '\n'
)

const (
	defaultIndent = 2
)

type (
	// ErrorMsg defines a message for errors.
	ErrorMsg struct {
		Error error
	}

	// QuitMsg defines a message for stopping the command.
	QuitMsg struct{}
)

// FormatView formats a model view padding and indentation.
func FormatView(view string) string {
	// Indent the view lines
	view = indent.String(view, defaultIndent)
	// Add top and bottom paddings
	return fmt.Sprintf("\n%s\n", view)
}
