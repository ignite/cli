package model

const (
	// ColorSpinner defines the foreground color for the spinner.
	ColorSpinner = "#8878FF"

	// EOL defines the rune for the end of line.
	EOL = '\n'
)

type (
	// ErrorMsg defines a message for error.
	ErrorMsg struct{ error }
)
