// Package gomodulepath implements functions for the manipulation of Go module paths.
// Paths are typically defined as a domain name and a path containing the user and
// repository names, e.g. "github.com/username/reponame", but Go also allows other module
// names like "domain.com/name", "name", "namespace/name", or similar variants.
package gomodulepath

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/ignite/cli/ignite/pkg/gomodule"
)

// Path represents a Go module's path.
type Path struct {
	// Path is Go module's full path.
	// e.g.: github.com/ignite/cli.
	RawPath string

	// Root is the root directory name of Go module.
	// e.g.: cli for github.com/ignite/cli.
	Root string

	// Package is the default package name for the Go module that can be used
	// to host main functionality of the module.
	// e.g.: cli for github.com/ignite/cli.
	Package string
}

// Parse parses rawpath into a module Path.
func Parse(rawpath string) (Path, error) {
	if err := validateRawPath(rawpath); err != nil {
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

// ParseAt parses Go module path of an app resides at path.
func ParseAt(path string) (Path, error) {
	parsed, err := gomodule.ParseAt(path)
	if err != nil {
		return Path{}, err
	}
	return Parse(parsed.Module.Mod.Path)
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

// ExtractAppPath extracts the app module path from a Go module path.
func ExtractAppPath(path string) string {
	if path == "" {
		return ""
	}

	items := strings.Split(path, "/")

	// Remove the first path item if it is assumed to be a domain name
	if len(items) > 1 && strings.Contains(items[0], ".") {
		items = items[1:]
	}

	count := len(items)
	if count == 1 {
		// The Go module path is a single name
		return items[0]
	}

	// The last two items in the path define the namespace and app name
	return strings.Join(items[count-2:], "/")
}

func hasDomainNamePrefix(path string) bool {
	if path == "" {
		return false
	}

	// TODO: should we use a regexp instead of the simplistic check ?
	name, _, _ := strings.Cut(path, "/")
	return strings.Contains(name, ".")
}

func validateRawPath(path string) error {
	// A raw path should be either a URI, a single name or a path
	if hasDomainNamePrefix(path) {
		return validateURIPath(path)
	}
	return validateNamePath(path)
}

func validateURIPath(path string) error {
	if err := module.CheckPath(path); err != nil {
		return fmt.Errorf("app name is an invalid go module name: %w", err)
	}
	return nil
}

func validateNamePath(path string) error {
	if err := module.CheckImportPath(path); err != nil {
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
