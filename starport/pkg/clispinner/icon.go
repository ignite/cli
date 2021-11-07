package clispinner

import "github.com/fatih/color"

var (
	// OK is an OK mark.
	OK = color.New(color.FgGreen).SprintFunc()("✔")
)
