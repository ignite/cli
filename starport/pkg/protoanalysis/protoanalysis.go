// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	protoFolder      = "proto"
	protoFilePattern = "*.proto"
)

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

// FindMessage finds a message from a
// specific module declared into the app
func FindMessage(
	module,
	structName string,
	staticTypes map[string]string,
) ([]string, error) {
	path := filepath.Join(protoFolder, module)
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	pkgs, err := Parse(context.Background(), NewCache(), path)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		for _, msg := range pkg.Messages {
			if msg.Name != structName {
				continue
			}
			fields := make([]string, 0)
			for fieldName, fieldType := range msg.Fields {
				if staticType, ok := staticTypes[fieldType]; ok {
					fieldType = staticType
				}
				fieldString := fmt.Sprintf("%s%s:%s", structName, fieldName, fieldType)
				fields = append(fields, fieldString)
			}
			return fields, nil
		}
	}
	return nil, fmt.Errorf("struct '%s' not found", structName)
}
