package gomodulepath

import (
	"fmt"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

const (
	cosmosModulePath     = "github.com/cosmos/cosmos-sdk"
	tendermintModulePath = "github.com/tendermint/tendermint"
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
	packageName := stripNonAlphaNumeric(rootName)
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

// Find search the Go module in the current and parent paths until finding it.
func Find(path string) (parsed Path, appPath string, err error) {
	for len(path) != 0 && path != "." && path != "/" {
		parsed, err = ParseAt(path)
		if errors.Is(err, gomodule.ErrGoModNotFound) {
			path = filepath.Dir(path)
			continue
		}
		return parsed, path, err
	}
	return Path{}, "", errors.Wrap(gomodule.ErrGoModNotFound, "could not locate your app's root dir")
}

// ParseAt parses Go module path of an app resides at path.
func ParseAt(path string) (Path, error) {
	parsed, err := gomodule.ParseAt(path)
	if err != nil {
		return Path{}, err
	}
	if err := validateGoModule(parsed); err != nil {
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

// validateGoModule check if the cosmos-sdk and the tendermint packages are imported.
func validateGoModule(module *modfile.File) error {
	moduleCheck := map[string]bool{
		cosmosModulePath:     true,
		tendermintModulePath: true,
	}
	for _, r := range module.Require {
		delete(moduleCheck, r.Mod.Path)
	}
	for m := range moduleCheck {
		return fmt.Errorf("invalid go module, missing %s package dependency", m)
	}
	return nil
}

func root(path string) string {
	sp := strings.Split(path, "/")
	name := sp[len(sp)-1]
	if semver.IsValid(name) { // omit versions.
		name = sp[len(sp)-2]
	}
	return name
}

func stripNonAlphaNumeric(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return strings.ToLower(reg.ReplaceAllString(name, ""))
}
