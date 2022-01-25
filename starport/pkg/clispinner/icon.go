package clispinner

import "github.com/fatih/color"

var (
	// OK is an OK mark.
	OK = color.New(color.FgGreen).SprintFunc()("✔")
	// NotOK is a red cross mark
	NotOK  = color.New(color.FgRed).SprintFunc()("✘")
	Bullet = color.New(color.FgYellow).SprintFunc()("⋆")
)
