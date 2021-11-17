package yaml

import (
	"context"
	"errors"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

var (
	mutablePaths = []string{
		"$.content.content.genesisValidator.genTx",
		"$.content.content.genesisValidator.consPubKey",
	}
)

const (
	mutableGentxPath      = "$.content.content.genesisValidator.genTx"
	mutableConsPubKeyPath = "$.content.content.genesisValidator.consPubKey"
)

func Parse(ctx context.Context, obj interface{}) (string, error) {
	yml, err := yaml.MarshalContext(ctx, obj)
	if err != nil {
		return "", err
	}
	file, err := parser.ParseBytes(yml, 0)
	if err != nil {
		return "", err
	}

	file, err = normalizeByteSlice(file, string(yml), mutableGentxPath)
	if err != nil {
		return "", err
	}

	file, err = normalizeByteSlice(file, string(yml), mutableConsPubKeyPath)
	if err != nil {
		return "", err
	}

	return file.String(), nil
	//for _, path := range mutablePaths {
	//	yml, err = normalizeByteSlice(file, string(yml), path)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
}

func normalizeByteSlice(file *ast.File, yml, yamlPath string) (*ast.File, error) {
	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return file, err
	}
	var obj []byte
	err = path.Read(strings.NewReader(yml), &obj)
	if !errors.Is(err, yaml.ErrNotFoundNode) {
		return file, err
	}
	return file, path.ReplaceWithReader(file, strings.NewReader(string(obj)))
}

func ParseNormalized(ctx context.Context, obj interface{}) (string, error) {
	requestYaml, err := yaml.MarshalContext(ctx, obj)
	if err != nil {
		return "", err
	}
	file, err := parser.ParseBytes(requestYaml, 0)
	if err != nil {
		return "", err
	}
	pathGentx, err := yaml.PathString(mutableGentxPath)
	if err != nil {
		return "", err
	}
	var gentx []byte
	err = pathGentx.Read(strings.NewReader(string(requestYaml)), &gentx)
	if !errors.Is(err, yaml.ErrNotFoundNode) {
		return "", err
	}
	if err := pathGentx.ReplaceWithReader(file, strings.NewReader(string(gentx))); err != nil {
		return "", err
	}

	pathConsPubKey, err := yaml.PathString(mutableConsPubKeyPath)
	if err != nil {
		return "", err
	}
	var consPubKey []byte
	err = pathConsPubKey.Read(strings.NewReader(string(requestYaml)), &consPubKey)
	if !errors.Is(err, yaml.ErrNotFoundNode) {
		return "", err
	}

	if err := pathConsPubKey.ReplaceWithReader(file, strings.NewReader(string(consPubKey))); err != nil {
		return "", err
	}
	return file.String(), nil
}
