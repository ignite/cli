package docs

import "embed"

// Docs are Starport docs.
//go:embed *.md */*.md
var Docs embed.FS

// MigrationDocs are Starport version migration docs.
//go:embed */migration/*.md
var MigrationDocs embed.FS
