package chain

import (
	"os"
	"path/filepath"
	"strings"
)

// App keeps info about scaffold.
type App struct {
	Name       string
	Path       string
	ImportPath string
}

// n returns app name without dashes.
func (a App) n() string {
	return strings.ReplaceAll(a.Name, "-", "")
}

// d returns appd name.
func (a App) d() string {
	return a.Name + "d"
}

// cli return appcli name.
func (a App) cli() string {
	return a.Name + "cli"
}

// nd returns no-dash appd name.
func (a App) nd() string {
	return a.n() + "d"
}

// ncli returns no-dash appcli name.
func (a App) ncli() string {
	return a.n() + "cli"
}

// home returns the home config path of app.
func (a App) home() string {
	h, _ := os.UserHomeDir()
	return filepath.Join(h, "."+a.Name)
}

// root returns the root path of app.
func (a App) root() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, a.Path)
}
