package docs

import "embed"

// Docs are Starport docs.
//go:embed *.md */*.md **/*/*.md
var Docs embed.FS
