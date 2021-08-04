// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/localfs"
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

// SearchRecursive recursively finds all proto files under path.
func SearchRecursive(path string) ([]string, error) {
	return localfs.Search(path, protoFilePattern)
}
