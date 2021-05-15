package goanalysis

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
)

type Type struct {
	Name string
}

func FindTypes(srcDirPath string, names []string) ([]Type, error) {
	namesSet := make(map[string]struct{})
	for _, name := range names {
		namesSet[name] = struct{}{}
	}
	var typesFound []Type

	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, srcDirPath, func(fs.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {

		// Iterate through the files in the package
		for _, file := range pkg.Files {

			// Walk through the AST
			ast.Inspect(file, func(n ast.Node) bool {

				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					return true
				}

				if _, ok := namesSet[typeSpec.Name.Name]; ok {
					typesFound = append(typesFound, Type{Name: typeSpec.Name.Name})
				}
				return true
			})
		}
	}

	return typesFound, nil
}
