package gomodulepath

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"golang.org/x/mod/module"
)

// Path represents a Go module's path.
type Path struct {
	// Path is Go module's full path.
	// e.g.: github.com/tendermint/starport.
	RawPath string

	// Root is the root directory name of Go module.
	// e.g.: starport for github.com/tendermint/starport.
	Root string

	// Package is the default package name for the Go module that can be used
	// to host main functionality of the module.
	// e.g.: starport for github.com/tendermint/starport.
	Package string
}

// Parse parses rawpath into a module Path.
func Parse(rawpath string) (Path, error) {
	if err := validateModulePath(rawpath); err != nil {
		return Path{}, err
	}
	rootName := root(rawpath)
	// package name cannot contain "-" so gracefully remove them
	// if they present.
	packageName := strings.ReplaceAll(rootName, "-", "")
	if err := validatePackageName(packageName); err != nil {
		return Path{}, err
	}
	p := Path{
		RawPath: rawpath,
		Root:    rootName,
		Package: packageName,
	}
	return p, nil
}

// ParseFile parses Go module path of an app resides at path.
func ParseFile(path string) (Path, error) {
	parsed, err := gomodule.ParseAt(path)
	if err != nil {
		return Path{}, err
	}
	return Parse(parsed.Module.Mod.Path)
}

func validateModulePath(path string) error {
	if err := module.CheckPath(path); err != nil {
		return fmt.Errorf("app name is an invalid go module name: %w", err)
	}
	return nil
}

func validatePackageName(name string) error {
	fset := token.NewFileSet()
	src := fmt.Sprintf("package %s", name)
	if _, err := parser.ParseFile(fset, "", src, parser.PackageClauseOnly); err != nil {
		// parser error is very low level here so let's hide it from the user
		// completely.
		return errors.New("app name is an invalid go package name")
	}
	return nil
}

func root(path string) string {
	sp := strings.Split(path, "/")
	return sp[len(sp)-1]
}
