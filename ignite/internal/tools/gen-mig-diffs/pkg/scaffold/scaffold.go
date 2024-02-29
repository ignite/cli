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
	// Scaffold holder the Scaffold logic.
	Scaffold struct {
		Output      string
		binary      string
		version     *semver.Version
		cache       *cache
		cachePath   string
		commandList Commands
	}
	// Command represents a set of commandList and prerequisites scaffold commandList that are required to run before them.
	Command struct {
		// Prerequisites is the names of commandList that need to be run before this command set
		Prerequisites []string
		// Commands is the list of scaffold commandList that are going to be run
		// The commandList will be prefixed with "ignite scaffold" and executed in order
		Commands []string
	}

	Commands map[string]Command

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

// WithOutput set the ignite scaffold Output.
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

// WithCommandList set the migration docs Output.
func WithCommandList(commands Commands) Options {
	return func(o *options) {
		o.commands = commands
	}
}

// New returns a new Scaffold.
func New(binary string, ver *semver.Version, options ...Options) (*Scaffold, error) {
	opts := newOptions()
	for _, apply := range options {
		apply(&opts)
	}

	output, err := filepath.Abs(opts.output)
	if err != nil {
		return nil, err
	}
	return &Scaffold{
		binary:      binary,
		version:     ver,
		cache:       newCache(opts.cachePath),
		cachePath:   opts.cachePath,
		Output:      filepath.Join(output, ver.Original()),
		commandList: opts.commands,
	}, nil
}

// Run execute the scaffold commandList based in the binary semantic version.
func (s *Scaffold) Run(ctx context.Context) error {
	for name, c := range s.commandList {
		if err := s.runCommand(ctx, name, c.Prerequisites, c.Commands); err != nil {
			return err
		}
		if err := applyPostScaffoldExceptions(s.version, name, s.Output); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scaffold) runCommand(
	ctx context.Context,
	name string,
	prerequisites, scaffoldCommands []string,
) error {
	// TODO add cache for duplicated scaffoldCommands.
	for _, name := range prerequisites {
		c, ok := s.commandList[name]
		if !ok {
			return errors.Errorf("command %s not found", name)
		}
		if err := s.runCommand(ctx, name, c.Prerequisites, c.Commands); err != nil {
			return err
		}
	}

	for _, cmd := range scaffoldCommands {
		if err := s.executeScaffold(ctx, name, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scaffold) executeScaffold(ctx context.Context, name, cmd string) error {
	args := append([]string{s.binary, "scaffold"}, strings.Fields(cmd)...)
	args = append(args, "--path", filepath.Join(s.Output, name))
	args = applyPreExecuteExceptions(s.version, args)

	if err := exec.Exec(ctx, args); err != nil {
		return errors.Wrapf(err, "failed to execute ignite scaffold command: %s", cmd)
	}
	return nil
}

// applyPreExecuteExceptions this function we can manipulate command arguments before executing it in
// order to compensate for differences in versions.
func applyPreExecuteExceptions(ver *semver.Version, args []string) []string {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the
	// name of chain at the given '--path' so we need to append "example" to the path if the
	// command is not "chain".
	if ver.LessThan(v027) && args[2] != "chain" {
		args[len(args)-1] = filepath.Join(args[len(args)-1], "example")
	}
	return args
}

// applyPostScaffoldExceptions this function we can manipulate the Output of scaffold commandList after
// they have been executed in order to compensate for differences in versions.
func applyPostScaffoldExceptions(ver *semver.Version, name string, output string) error {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the name of
	// chain at the given --path so we need to move the directory to the parent directory.
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
