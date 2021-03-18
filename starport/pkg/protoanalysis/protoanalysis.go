// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"

	"github.com/mattn/go-zglob"
)

// Parse parses proto packages by finding them with given glob pattern.
func Parse(ctx context.Context, pattern string) ([]Package, error) {
	p := newParser()

	if err := p.parse(ctx, pattern); err != nil {
		return nil, err
	}

	var packages []Package

	for _, pp := range p.packages {
		packages = append(packages, newBuilder(*pp).build())
	}

	return packages, nil
}

// SearchRecursive recursively finds all proto files under path.
func SearchRecursive(dir string) ([]string, error) {
	return zglob.Glob(PatternRecursive(dir))
}

// PatternRecursive returns a recursive glob search pattern to find all proto files under path.
func PatternRecursive(dir string) string { return dir + "/**/*.proto" }
