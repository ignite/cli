package app

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/tendermint/starport/starport/pkg/cosmosanalysis"
)

var appImplementation = []string{
	"RegisterAPIRoutes",
	"RegisterTxService",
	"RegisterTendermintService",
}

// CheckKeeper checks for the existence of the keeper with the provided name in the app structure
func CheckKeeper(path, keeperName string) error {
	// find app type
	appImpl, err := cosmosanalysis.FindImplementation(path, appImplementation)
	if err != nil {
		return err
	}
	if len(appImpl) != 1 {
		return errors.New("app.go should contain a single app")
	}
	appTypeName := appImpl[0]

	// Inspect the module for app struct
	var found bool
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, path, nil, 0)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				// look for struct methods.
				appType, ok := n.(*ast.TypeSpec)
				if !ok || appType.Name.Name != appTypeName {
					return true
				}

				appStruct, ok := appType.Type.(*ast.StructType)
				if !ok {
					return true
				}

				// Search for the keeper specific field
				for _, field := range appStruct.Fields.List {
					for _, fieldName := range field.Names {
						if fieldName.Name == keeperName {
							found = true
						}
					}
				}

				return false
			})
		}
	}

	if !found {
		return fmt.Errorf("app doesn't contain %s", keeperName)
	}
	return nil
}
