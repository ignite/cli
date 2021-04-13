// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/localfs"
)

// Parse parses proto packages by finding them with given glob pattern.
func Parse(ctx context.Context, pattern string) ([]Package, error) {
	parsed, err := parse(ctx, pattern)
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
func SearchRecursive(dir string) ([]string, error) {
	return localfs.Search(PatternRecursive(dir))
}

// PatternRecursive returns a recursive glob search pattern to find all proto files under path.
func PatternRecursive(dir string) string { return dir + "/**/*.proto" }
