package scaffolder

import (
	"context"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/add"
)

const (
	wasmImport = "github.com/CosmWasm/wasmd"
	apppkg     = "app"
)

// AddModule adds sepecified module with name to the scaffolded app.
func (s *Scaffolder) AddModule(name string) error {
	version, err := s.version()
	if err != nil {
		return err
	}
	if version == cosmosver.Stargate {
		return errors.New("adding modules currently is not supported on Stargate")
	}
	ok, err := isWasmAdded(s.path)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("CosmWasm is already added.")
	}
	path, err := gomodulepath.ParseFile(s.path)
	if err != nil {
		return err
	}
	g, err := add.New(&add.Options{
		Feature: name,
		AppName: path.Package,
	})
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	return run.Run()
}

func isWasmAdded(appPath string) (bool, error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, apppkg))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ImportsOnly)
	if err != nil {
		return false, err
	}
	for _, pkg := range all {
		for _, f := range pkg.Files {
			for _, imp := range f.Imports {
				if strings.Contains(imp.Path.Value, wasmImport) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
