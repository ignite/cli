// Package gomodulepath implements functions for the manipulation of Go module paths.
// A module’s path is the prefix for package paths within a module.
// Paths are typically defined as a domain name and a path containing the user and
// repository names, e.g. "github.com/username/reponame", but Go also allows other module
// names like "domain.com/reponame", "reponame", "reponame/pkgname", or similar variants.
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

	"github.com/tendermint/starport/starport/pkg/giturl"
	"github.com/tendermint/starport/starport/pkg/gomodule"
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

// HasDomainNamePrefix checks is a Go module path has a domain name as prefix.
// The first path element is assumed to be a domain when it contains a ".".
func HasDomainNamePrefix(path string) bool {
	if path == "" {
		return false
	}

	// TODO: should we use a regexp instead of the simplistic check ?
	name, _, _ := strings.Cut(path, "/")
	return strings.Contains(name, ".")
}

// ExtractUserAndRepoNames extracts the user and repository names from a go module path.
// Path must be a URL path, e.g. github.com/tendermint/starport, in which case the user
// value is "tendermint" and the repo name is "starport", or it can be a local path,
// e.g. "starport" or "starport/x" which results in the same "starport" value for the
// user and repo names.
func ExtractUserAndRepoNames(path string) (string, string, error) {
	if HasDomainNamePrefix(path) {
		// Use the path element as user and repo names when the path
		// contains a short URL as value, e.g. domain.com/starport
		if v := strings.SplitN(path, "/", 3); len(v) == 2 {
			return v[1], v[1], nil
		}

		g, err := giturl.Parse(path)
		if err != nil {
			return "", "", err
		}

		return g.User, g.Repo, nil
	}

	// The first path element defines the user and repository names
	name, _, _ := strings.Cut(path, "/")

	return name, name, nil
}

func validateRawPath(path string) error {
	// A raw path should be either a URL, a single name or a path
	if HasDomainNamePrefix(path) {
		return validateURLPath(path)
	}
	return validateNamePath(path)
}

func validateURLPath(path string) error {
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
