// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// ErrImportNotFound returned when proto file import cannot be found.
var ErrImportNotFound = errors.New("proto import not found")

const protoFilePattern = "*.proto"

// Parse parses proto packages by finding them with given glob pattern.
func Parse(ctx context.Context, cache *Cache, path string) (Packages, error) {
	if cache != nil {
		if packages, ok := cache.Get(path); ok {
			return packages, nil
		}
	}

	parsed, err := parse(ctx, path, protoFilePattern)
	if err != nil {
		return nil, err
	}

	var packages Packages

	for _, pp := range parsed {
		packages = append(packages, build(*pp))
	}

	if cache != nil {
		cache.Add(path, packages)
	}

	return packages, nil
}

// ParseFile parses a proto file at path.
func ParseFile(path string) (File, error) {
	packages, err := Parse(context.Background(), nil, path)
	if err != nil {
		return File{}, err
	}
	files := packages.Files()
	if len(files) != 1 {
		return File{}, errors.New("path does not point to single file or it cannot be found")
	}
	return files[0], nil
}

// HasMessages checks if the proto package under path contains messages with given names.
func HasMessages(ctx context.Context, path string, names ...string) error {
	pkgs, err := Parse(ctx, NewCache(), path)
	if err != nil {
		return err
	}

	hasName := func(name string) error {
		for _, pkg := range pkgs {
			for _, msg := range pkg.Messages {
				if msg.Name == name {
					return nil
				}
			}
		}
		return fmt.Errorf("invalid proto message name %s", name)
	}

	for _, name := range names {
		if err := hasName(name); err != nil {
			return err
		}
	}
	return nil
}

// IsImported checks if the proto package under path imports list of dependencies.
func IsImported(path string, dependencies ...string) error {
	f, err := ParseFile(path)
	if err != nil {
		return err
	}

	for _, wantDep := range dependencies {
		found := false
		for _, fileDep := range f.Dependencies {
			if wantDep == fileDep {
				found = true
				break
			}
		}
		if !found {
			return errors.Wrap(ErrImportNotFound, fmt.Sprintf(
				"invalid proto dependency %s for file %s", wantDep, path),
			)
		}
	}
	return nil
}
