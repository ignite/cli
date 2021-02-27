// Package sta provides access to swagger-typescript-api CLI.
package sta

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/nodetime"
)

var placeOnce sync.Once

// Generate generates client code and TS types to outPath from an OpenAPI spec that resides at specPath.
func Generate(ctx context.Context, outPath, specPath, moduleNameIndex string) error {
	var err error

	// places the protobufjs-cli into BinaryPath.
	placeOnce.Do(func() { err = nodetime.PlaceBinary() })

	if err != nil {
		return err
	}

	dir := filepath.Dir(outPath)
	file := filepath.Base(outPath)

	// command constructs the sta command.
	command := []string{
		nodetime.BinaryPath,
		nodetime.CommandSTA,
		"--module-name-index",
		moduleNameIndex,
		"-p",
		specPath,
		"-o",
		dir,
		"-n",
		file,
	}

	// execute the command.
	return cmdrunner.Exec(ctx, command[0], command[1:]...)
}
