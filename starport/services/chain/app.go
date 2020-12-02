package chain

import (
	"path/filepath"
	"strings"
)

// App keeps info about scaffold.
type App struct {
	Name       string
	Path       string
	ImportPath string
}

// N returns app name without dashes.
func (a App) N() string {
	return strings.ReplaceAll(a.Name, "-", "")
}

// D returns appd name.
func (a App) D() string {
	return a.Name + "d"
}

// CLI return appcli name.
func (a App) CLI() string {
	return a.Name + "cli"
}

// ND returns no-dash appd name.
func (a App) ND() string {
	return a.N() + "d"
}

// NCLI returns no-dash appcli name.
func (a App) NCLI() string {
	return a.N() + "cli"
}

// Root returns the root path of app.
func (a App) Root() string {
	path, _ := filepath.Abs(a.Path)
	return path
}
