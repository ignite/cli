// Package goanalysis provides a toolset for statically analysing Go applications
package goanalysis

import (
	"errors"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	mainPackage     = "main"
	goFileExtension = ".go"
)

// DiscoverMain finds the main application package path
func DiscoverMain(appPath string) (string, error) {
	var (
		fset     = token.NewFileSet()
		mainPath = ""
	)
	err := filepath.Walk(appPath, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() || !strings.HasSuffix(filePath, goFileExtension) {
			return err
		}
		parsed, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
		if err != nil {
			return err
		}
		name := parsed.Name.Name
		if name == mainPackage {
			mainPath = filepath.Dir(filePath)
			return io.EOF
		}
		return nil
	})
	if err == io.EOF {
		return mainPath, nil
	} else if err != nil {
		return "", err
	}
	return "", errors.New("main package not found")
}
