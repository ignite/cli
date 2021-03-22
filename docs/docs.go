package docs

import "embed"

// Docs are Starport docs.
//go:embed */*.md
var Docs embed.FS
