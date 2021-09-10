// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

const protoFilePattern = "*.proto"

type Cache map[string]Packages // proto dir path-proto packages pair.

func NewCache() Cache {
	return make(Cache)
}

// Parse parses proto packages by finding them with given glob pattern.
func Parse(ctx context.Context, cache Cache, path string) (Packages, error) {
	if cache != nil {
		if packages, ok := cache[path]; ok {
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
		cache[path] = packages
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

// CheckTypes check if the proto files into a folder contains a list of message names
func CheckTypes(path string, names []string) error {
	pkgs, err := Parse(context.Background(), NewCache(), path)
	if err != nil {
		return err
	}

	for _, name := range names {
		if err := checkMsgName(pkgs, name); err != nil {
			return err
		}
	}
	return nil
}

// checkMsgName check if a message name exist into the package list
func checkMsgName(pkgs Packages, name string) error {
	for _, pkg := range pkgs {
		for _, msg := range pkg.Messages {
			if msg.Name == name {
				return nil
			}
		}
	}
	return fmt.Errorf("invalid proto message name %s", name)
}
