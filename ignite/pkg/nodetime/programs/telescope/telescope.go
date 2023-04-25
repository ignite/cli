// Package telescope provides access to @osmosis-labs/telescope CLI.
package telescope

import (
	"context"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/nodetime"
)

// Option configures Generate configs.
type Option func(*configs)

// Configs holds Generate configs.
type configs struct {
	command Cmd
}

// WithCommand assigns a telescope command to use for code generation.
// This allows to use a single nodetime telescope binary in multiple code generation
// calls. Otherwise, `Generate` creates a new generator binary each time it is called.
func WithCommand(command Cmd) Option {
	return func(c *configs) {
		c.command = command
	}
}

// Cmd contains the information necessary to execute the telescope command.
type Cmd struct {
	command []string
}

// Command returns the strings to execute the telescope command.
func (c Cmd) Command() []string {
	return c.command
}

// Command sets the telescope binary up and returns the command needed to execute it.
func Command() (command Cmd, cleanup func(), err error) {
	c, cleanup, err := nodetime.Command(nodetime.CommandTelescope)
	command = Cmd{c}
	return
}

// Generate generates client code and TS types to outPath from an array of proto paths.
func Generate(ctx context.Context, outPath string, includePaths []string, options ...Option) error {
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

	// command constructs the telescope command.
	command = append(command, []string{
		"transpile",
		"--includeRPCClients",
		"--includeLCDClients",
		"--includeAminos",
		"--outPath",
		dir,
	}...)

	for i := 0; i < len(includePaths); i++ {
		command = append(command, []string{
			"--protoDirs",
			includePaths[i],
		}...)
	}
	// execute the command.
	return exec.Exec(ctx, command, exec.IncludeStdLogsToError())
}
