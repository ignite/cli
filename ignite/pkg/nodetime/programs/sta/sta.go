// Package sta provides access to swagger-typescript-api CLI.
package sta

import (
	"context"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/nodetime"
)

// Option configures Generate configs.
type Option func(*configs)

// configs holds Generate configs.
type configs struct {
	command Cmd
}

// WithCommand assigns a typescript API generator command to use for code generation.
// This allows to use a single nodetime STA generator binary in multiple code generation
// calls. Otherwise `Generate` creates a new generator binary each time it is called.
func WithCommand(command Cmd) Option {
	return func(c *configs) {
		c.command = command
	}
}

// Cmd contains the information necessary to execute the typescript API generator command.
type Cmd struct {
	command []string
}

// Command returns the strings to execute the typescript API generator command.
func (c Cmd) Command() []string {
	return c.command
}

// Command sets the typescript API generator binary up and returns the command needed to execute it.
func Command() (command Cmd, cleanup func(), err error) {
	c, cleanup, err := nodetime.Command(nodetime.CommandSTA)
	command = Cmd{c}
	return
}

// Generate generates client code and TS types to outPath from an OpenAPI spec that resides at specPath.
func Generate(ctx context.Context, outPath, specPath, moduleNameIndex string, options ...Option) error {
	c := configs{}

	for _, o := range options {
		o(&c)
	}

	command := c.command.Command()
	if command == nil {
		cmd, cleanup, err := Command()
		if err != nil {
			return err
		}

		defer cleanup()

		command = cmd.Command()
	}

	dir := filepath.Dir(outPath)
	file := filepath.Base(outPath)

	// command constructs the sta command.
	command = append(command, []string{
		"--module-name-index",
		moduleNameIndex,
		"-p",
		specPath,
		"-o",
		dir,
		"-n",
		file,
	}...)

	// execute the command.
	return exec.Exec(ctx, command, exec.IncludeStdLogsToError())
}
