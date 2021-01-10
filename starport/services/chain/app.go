package chain

import (
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/gomodulepath"
)

// App keeps info about chain.
type App struct {
	ChainID    string
	Name       string
	Path       string
	ImportPath string
}

// NewAppAt creates an App from the blockchain source code located at path.
func NewAppAt(path string) (App, error) {
	p, err := gomodulepath.ParseAt(path)
	if err != nil {
		return App{}, err
	}
	return App{
		Path:       path,
		Name:       p.Root,
		ImportPath: p.RawPath,
	}, nil
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
