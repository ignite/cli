package app

import (
	"errors"
	"fmt"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis"
	"go/ast"
	"go/parser"
	"go/token"
)

var appImplementation = []string{
	"RegisterAPIRoutes",
	"RegisterGRPCServer",
	"RegisterTxService",
	"RegisterTendermintService",
}


// CheckKeeper checks for the existence of the keeper with the provided name in the app structure
func CheckKeeper(appPath, keeperName string) error {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, appPath, nil, 0)
	if err != nil {
		return err
	}

	// find app type
	appImpl, err := cosmosanalysis.FindImplementation(appPath, appImplementation)
	if err != nil {
		return err
	}
	if len(appImpl) != 1 {
		return errors.New("app.go should contain a single app")
	}
	appTypeName := appImpl[0]

	// Inspect the file for app struct
	var found bool
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

	if !found {
		return fmt.Errorf("app doesn't contain %s", keeperName)
	}
	return nil
}
