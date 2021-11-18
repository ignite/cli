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

func Parse(ctx context.Context, obj interface{}) (string, error) {
	yml, err := yaml.MarshalContext(ctx, obj)
	if err != nil {
		return "", err
	}
	file, err := parser.ParseBytes(yml, 0)
	if err != nil {
		return "", err
	}

	err = normalizeByteSlice(file, yml, mutableGentxPath)
	if err != nil {
		return "", err
	}

	err = normalizeByteSlice(file, yml, mutableConsPubKeyPath)
	if err != nil {
		return "", err
	}

	return file.String(), nil
}

func normalizeByteSlice(file *ast.File, yml []byte, yamlPath string) error {
	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return err
	}
	var obj []byte
	err = path.Read(strings.NewReader(string(yml)), &obj)
	if !errors.Is(err, yaml.ErrNotFoundNode) {
		return err
	}
	return path.ReplaceWithReader(file, strings.NewReader(string(obj)))
}

func normalizeByteSlice2(file *ast.File, yml []byte, yamlPath string) (*ast.File, []byte, error) {
	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return file, yml, err
	}
	var obj []byte
	err = path.Read(strings.NewReader(string(yml)), &obj)
	if !errors.Is(err, yaml.ErrNotFoundNode) {
		return file, yml, err
	}
	return file, yml, path.ReplaceWithReader(file, strings.NewReader(string(obj)))
}
