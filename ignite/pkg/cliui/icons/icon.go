package icons

import (
	"github.com/ignite-hq/cli/ignite/pkg/cliui/colors"
)

var (
	// OK is an OK mark.
	OK = colors.Success("✔")
	// NotOK is a red cross mark
	NotOK = colors.Error("✘")
	// Bullet is a bullet mark
	Bullet = colors.Info("⋆")
	// Info is an info mark
	Info = colors.Info("𝓲")
)
