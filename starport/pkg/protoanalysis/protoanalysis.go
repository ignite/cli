// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/localfs"
)

const protoFilePattern = "*.proto"

// Parse parses proto packages by finding them with given glob pattern.
func Parse(ctx context.Context, path string) ([]Package, error) {
	parsed, err := parse(ctx, path, protoFilePattern)
	if err != nil {
		return nil, err
	}

	var packages []Package

	for _, pp := range parsed {
		packages = append(packages, build(*pp))
	}

	return packages, nil
}

// SearchRecursive recursively finds all proto files under path.
func SearchRecursive(path string) ([]string, error) {
	return localfs.Search(path, protoFilePattern)
}
