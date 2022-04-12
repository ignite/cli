package docs

import "embed"

// Docs are Ignite CLI docs.
//go:embed *.md */*.md
var Docs embed.FS
