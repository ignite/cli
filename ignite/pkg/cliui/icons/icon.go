package icons

import "github.com/fatih/color"

var (
	Earth = "🌍"

	// OK is an OK mark.
	OK = color.New(color.FgGreen).SprintFunc()("✔")
	// NotOK is a red cross mark
	NotOK = color.New(color.FgRed).SprintFunc()("✘")
	// Bullet is a bullet mark
	Bullet = color.New(color.FgYellow).SprintFunc()("⋆")
	// Info is an info mark
	Info = color.New(color.FgYellow).SprintFunc()("𝓲")
)
