package starportcmd

import (
	"context"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/templates/add"
)

const (
	wasmImport = "github.com/CosmWasm/wasmd"
	apppkg     = "app"
)

func NewAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [feature]",
		Short: "Adds a feature to a project.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  addHandler,
	}
	return c
}

func addHandler(cmd *cobra.Command, args []string) error {
	ok, err := isWasmAdded(appPath)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("wasm alread addded")
	}
	appName, _ := getAppAndModule(appPath)
	g, _ := add.New(&add.Options{
		Feature: args[0],
		AppName: appName,
	})
	run := genny.WetRunner(context.Background())
	run.With(g)
	run.Run()
	return nil
	// fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
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
