package protobufjs

import (
	"bytes"
	"context"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/nodetime"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

var placeOnce sync.Once

// Generate generates static protobuf.js types for given proto where includePaths holds dependency protos.
// TODO add ts generation. protobufjs supports this but by executing jsdoc command with node dynamically,
// it doesn't work with bundled node pkg. things needs to be reconstructed.
func Generate(ctx context.Context, outDir, outName, protoPath string, includePaths []string) error {
	var err error

	// places the protobufjs-cli into BinaryPath.
	placeOnce.Do(func() { err = nodetime.PlaceBinary() })

	if err != nil {
		return err
	}

	runcmd := func(command []string) error {
		errb := &bytes.Buffer{}

		err = cmdrunner.
			New(
				cmdrunner.DefaultStderr(errb)).
			Run(ctx,
				step.New(step.Exec(command[0], command[1:]...)))

		return errors.Wrap(err, errb.String())
	}

	jsOutPath := filepath.Join(outDir, outName+".js")

	// construct js gen command for the actual code generation.
	command := []string{
		nodetime.BinaryPath,
		nodetime.CommandPBJS,
		"-t",
		"static-module",
		"-w",
		"es6",
		"-o",
		jsOutPath,
	}

	// add proto dependency paths to that.
	for _, includePath := range includePaths {
		command = append(
			command,
			"-p",
			includePath,
		)
	}

	// add target proto path to that.
	command = append(command, protoanalysis.GlobPattern(protoPath))

	// run the js command.
	return runcmd(command)
}
