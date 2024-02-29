package scaffold

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

var v027 = semver.MustParse("v0.27.0")

type (
	// Scaffold represents a set of commands and prerequisites scaffold commands that are required to run before them.
	Scaffold struct {
		// Name is the unique identifier of the command
		Name string
		// Prerequisites is the names of commands that need to be run before this command set
		Prerequisites []string
		// Commands is the list of scaffold commands that are going to be run
		// The commands will be prefixed with "ignite scaffold" and executed in order
		Commands []string
	}

	Commands map[string]Scaffold
)

type (
	// options represents configuration for the generator.
	options struct {
		cachePath string
		output    string
		commands  Commands
	}
	// Options configures the generator.
	Options func(*options)
)

// newOptions returns a options with default options.
func newOptions() options {
	tmpDir := os.TempDir()
	return options{
		cachePath: filepath.Join(tmpDir, "migration-cache"),
		output:    filepath.Join(tmpDir, "migration"),
		commands:  defaultCommands,
	}
}

// WithOutput set the ignite scaffold output.
func WithOutput(output string) Options {
	return func(o *options) {
		o.output = output
	}
}

// WithCachePath set the ignite scaffold cache path.
func WithCachePath(cachePath string) Options {
	return func(o *options) {
		o.cachePath = cachePath
	}
}

// WithCommandList set the migration docs output.
func WithCommandList(commands Commands) Options {
	return func(o *options) {
		o.commands = commands
	}
}

// Run execute the scaffold commands based in the binary semantic version.
func Run(binary string, ver *semver.Version, options ...Options) (string, error) {
	opts := newOptions()
	for _, apply := range options {
		apply(&opts)
	}

	output, err := filepath.Abs(opts.output)
	if err != nil {
		return "", err
	}
	output = filepath.Join(output, ver.Original())

	for _, c := range opts.commands {
		if err := runCommand(binary, output, c.Name, c.Prerequisites, c.Commands, ver, opts.commands); err != nil {
			return "", err
		}
		if err := applyPostScaffoldExceptions(ver, c.Name, output); err != nil {
			return "", err
		}
	}
	return output, nil
}

func runCommand(
	binary, output, name string,
	prerequisites, scaffoldCommands []string,
	ver *semver.Version,
	commandList Commands,
) error {
	// TODO add cache for duplicated scaffoldCommands.
	for _, p := range prerequisites {
		c, ok := commandList[p]
		if !ok {
			return errors.Errorf("command %s not found", name)
		}
		if err := runCommand(binary, output, name, c.Prerequisites, c.Commands, ver, commandList); err != nil {
			return err
		}
	}

	for _, cmd := range scaffoldCommands {
		if err := executeScaffold(binary, name, cmd, output, ver); err != nil {
			return err
		}
	}
	return nil
}

func executeScaffold(binary, name, cmd, output string, ver *semver.Version) error {
	args := append([]string{binary, "scaffold"}, strings.Fields(cmd)...)
	args = append(args, "--path", filepath.Join(output, name))
	args = applyPreExecuteExceptions(ver, args)

	if err := exec.Exec(context.Background(), args); err != nil {
		return errors.Wrapf(err, "failed to execute ignite scaffold command: %s", cmd)
	}
	return nil
}

// applyPreExecuteExceptions this function we can manipulate command arguments before executing it in
// order to compensate for differences in versions.
func applyPreExecuteExceptions(ver *semver.Version, args []string) []string {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the name of chain at the given --path
	// so we need to append "example" to the path if the command is not "chain"
	if ver.LessThan(v027) && args[2] != "chain" {
		args[len(args)-1] = filepath.Join(args[len(args)-1], "example")
	}
	return args
}

// applyPostScaffoldExceptions this function we can manipulate the output of scaffold commands after
// they have been executed in order to compensate for differences in versions.
func applyPostScaffoldExceptions(ver *semver.Version, name string, output string) error {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the name of chain at the given --path
	// so we need to move the directory to the parent directory.
	if ver.LessThan(v027) {
		if err := os.Rename(filepath.Join(output, name, "example"), filepath.Join(output, "example_tmp")); err != nil {
			return errors.Wrapf(err, "failed to move %s directory to tmp directory", name)
		}

		if err := os.RemoveAll(filepath.Join(output, name)); err != nil {
			return errors.Wrapf(err, "failed to remove %s directory", name)
		}

		if err := os.Rename(filepath.Join(output, "example_tmp"), filepath.Join(output, name)); err != nil {
			return errors.Wrapf(err, "failed to move tmp directory to %s directory", name)
		}
	}

	return nil
}
