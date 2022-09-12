package chain

import (
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/xstrings"
)

// App keeps info about chain.
type App struct {
	Name       string
	Path       string
	ImportPath string
}

// NewAppAt creates an App from the blockchain source code located at path.
func NewAppAt(path string) (App, error) {
	p, appPath, err := gomodulepath.Find(path)
	if err != nil {
		return App{}, err
	}
	return App{
		Path:       appPath,
		Name:       p.Root,
		ImportPath: p.RawPath,
	}, nil
}

// N returns app name without dashes.
func (a App) N() string {
	return xstrings.NoDash(a.Name)
}

// D returns appd name.
func (a App) D() string {
	return a.Name + "d"
}

// ND returns no-dash appd name.
func (a App) ND() string {
	return a.N() + "d"
}

// Root returns the root path of app.
func (a App) Root() string {
	path, _ := filepath.Abs(a.Path)
	return path
}
