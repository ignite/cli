package cli

import _ "embed"

//go:embed ignite_apps.md
var igniteAppsDoc []byte

// IgniteAppsDoc returns the Markdown contents of `ignite_apps.md`.
func IgniteAppsDoc() []byte {
	return igniteAppsDoc
}
