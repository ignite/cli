package yaml

import (
	"context"
	"errors"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

// Marshal converts an object to a string in a YAML format and transforms
// the byte slice fields from the path to string to be more readable.
func Marshal(ctx context.Context, obj interface{}, paths ...string) (string, error) {
	requestYaml, err := yaml.MarshalContext(ctx, obj)
	if err != nil {
		return "", err
	}
	file, err := parser.ParseBytes(requestYaml, 0)
	if err != nil {
		return "", err
	}

	// normalize the structure converting the byte slice fields to string
	for _, path := range paths {
		pathString, err := yaml.PathString(path)
		if err != nil {
			return "", err
		}
		var byteSlice []byte
		err = pathString.Read(strings.NewReader(string(requestYaml)), &byteSlice)
		if err != nil && !errors.Is(err, yaml.ErrNotFoundNode) {
			return "", err
		}
		if err := pathString.ReplaceWithReader(file,
			strings.NewReader(string(byteSlice)),
		); err != nil {
			return "", err
		}
	}

	return file.String(), nil
}
