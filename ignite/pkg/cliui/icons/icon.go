package icons

import (
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
)

var (
	Earth   = "🌍"
	CD      = "💿"
	User    = "👤"
	Command = "❯⎯"
	Hook    = "🪝"

	// OK is an OK mark.
	OK = colors.SprintFunc(colors.Green)("✔")
	// NotOK is a red cross mark.
	NotOK = colors.SprintFunc(colors.Red)("✘")
	// Bullet is a bullet mark.
	Bullet = colors.SprintFunc(colors.Yellow)("⋆")
	// Info is an info mark.
	Info = colors.SprintFunc(colors.Yellow)("𝓲")
)
