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
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/add"
)

const (
	wasmImport = "github.com/CosmWasm/wasmd"
	apppkg     = "app"
)


func NewAdd() *cobra.Command {

validArgs := []string{ "wasm" }

	c := &cobra.Command{
		Use:   "add [feature]",
		Short: "Adds features like wasm sto your blockchain",
		Long: "Use starport add wasm to add support for webassembly smart contracts to your blockchain.",
		Args:  cobra.ExactValidArgs(1),
		ValidArgs: validArgs,
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
		return errors.New("CosmWasm is already added.")
	}
	path, err := gomodulepath.Parse(getModule(appPath))
	if err != nil {
		return err
	}
	g, _ := add.New(&add.Options{
		Feature: args[0],
		AppName: path.Package,
	})
	run := genny.WetRunner(context.Background())
	run.With(g)
	run.Run()
	return nil
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
