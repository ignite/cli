// Package protoc provides high level access to protoc command.
package protoc

import (
	"context"
	"os"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
)

// Option configures Generate configs.
type Option func(*configs)

// configs holds Generate configs.
type configs struct {
	pluginPath string
}

// Plugin configures a plugin for code generation.
func Plugin(path string) Option {
	return func(c *configs) {
		c.pluginPath = path
	}
}

// Generate generates code into outDir from protoPath and its includePaths by using plugins provided with protocOuts.
func Generate(ctx context.Context, outDir, protoPath string, includePaths, protocOuts []string, options ...Option) error {
	c := &configs{}
	for _, o := range options {
		o(c)
	}

	// start preparing the protoc command for execution.
	command := []string{
		"protoc",
	}

	// add plugin if set.
	if c.pluginPath != "" {
		command = append(command, "--plugin", c.pluginPath)
	}

	// append third party proto locations to the command.
	for _, importPath := range includePaths {
		// skip if a third party proto source actually doesn't exist on the filesystem.
		if _, err := os.Stat(importPath); os.IsNotExist(err) {
			continue
		}
		command = append(command, "-I", importPath)
	}

	// find out the list of proto files under the app and generate code for them.
	files, err := protoanalysis.SearchProto(protoPath)
	if err != nil {
		return err
	}

	// run command for each protocOuts.
	for _, out := range protocOuts {
		command := append(command, out)
		command = append(command, files...)

		if err := exec.Exec(ctx, command,
			exec.StepOption(step.Workdir(outDir)),
			exec.IncludeStdLogsToError(),
		); err != nil {
			return err
		}
	}

	return nil
}
