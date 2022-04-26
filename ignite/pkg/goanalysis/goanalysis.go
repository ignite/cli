// Package goanalysis provides a toolset for statically analysing Go applications
package goanalysis

import (
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const (
	mainPackage     = "main"
	goFileExtension = ".go"
)

var (
	// ErrMultipleMainPackagesFound is returned when multiple main packages found while expecting only one.
	ErrMultipleMainPackagesFound = errors.New("multiple main packages found")
)

// DiscoverMain finds main Go packages under path.
func DiscoverMain(path string) (pkgPaths []string, err error) {
	uniquePaths := make(map[string]struct{})

	err = filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() || !strings.HasSuffix(filePath, goFileExtension) {
			return err
		}

		parsed, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.PackageClauseOnly)
		if err != nil {
			return err
		}

		if mainPackage == parsed.Name.Name {
			dir := filepath.Dir(filePath)
			uniquePaths[dir] = struct{}{}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for path := range uniquePaths {
		pkgPaths = append(pkgPaths, path)
	}

	return pkgPaths, nil
}

// DiscoverOneMain tries to find only one main Go package under path.
func DiscoverOneMain(path string) (pkgPath string, err error) {
	pkgPaths, err := DiscoverMain(path)
	if err != nil {
		return "", err
	}

	count := len(pkgPaths)
	if count == 0 {
		return "", errors.New("main package cannot be found")
	}
	if count > 1 {
		return "", ErrMultipleMainPackagesFound
	}

	return pkgPaths[0], nil
}

// FindImportedPackages finds the imported packages in a Go file and returns a map
// with package name, import path pair.
func FindImportedPackages(name string) (map[string]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}

	packages := make(map[string]string) // name -> import
	for _, imp := range f.Imports {
		var importName string
		if imp.Name != nil {
			importName = imp.Name.Name
		} else {
			importParts := strings.Split(imp.Path.Value, "/")
			importName = importParts[len(importParts)-1]
		}

		name := strings.Trim(importName, "\"")
		packages[name] = strings.Trim(imp.Path.Value, "\"")
	}

	return packages, nil
}
