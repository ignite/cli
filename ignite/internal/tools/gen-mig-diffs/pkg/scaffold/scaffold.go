package scaffold

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/ignite/cli/v29/ignite/internal/tools/gen-mig-diffs/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/randstr"
)

var v027 = semver.MustParse("v0.27.0")

type (
	// Scaffold holder the Scaffold logic.
	Scaffold struct {
		Output      string
		binary      string
		version     *semver.Version
		cache       *cache.Cache
		cachePath   string
		commandList Commands
	}

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
	tmpDir := filepath.Join(os.TempDir(), randstr.Runes(4))
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

	c, err := cache.New(opts.cachePath)
	if err != nil {
		return nil, err
	}

	if err := opts.commands.Validate(); err != nil {
		return nil, err
	}

	return &Scaffold{
		binary:      binary,
		version:     ver,
		cache:       c,
		cachePath:   opts.cachePath,
		Output:      filepath.Join(output, ver.Original()),
		commandList: opts.commands,
	}, nil
}

// Run execute the scaffold command based in the binary semantic version.
func (s *Scaffold) Run(ctx context.Context) error {
	if err := os.RemoveAll(s.Output); err != nil {
		return errors.Wrapf(err, "failed to remove the scaffold output directory: %s", s.Output)
	}

	for _, command := range s.commandList {
		if err := s.runCommand(ctx, command.Name, command); err != nil {
			return err
		}
		if err := applyPostScaffoldExceptions(s.version, command.Name, s.Output); err != nil {
			return err
		}
	}
	return nil
}

// Cleanup cleanup all temporary directories.
func (s *Scaffold) Cleanup() error {
	if err := os.RemoveAll(s.cachePath); err != nil {
		return err
	}
	return os.RemoveAll(s.Output)
}

func (s *Scaffold) runCommand(ctx context.Context, name string, command Command) error {
	path := filepath.Join(s.Output, name)
	if command.Prerequisite != "" {
		reqCmd, err := s.commandList.Get(command.Prerequisite)
		if err != nil {
			return errors.Wrapf(err, "pre-requisite command %s from %s not found", command.Prerequisite, name)
		}

		if s.cache.Has(command.Prerequisite) {
			if err := s.cache.Get(command.Prerequisite, path); err != nil {
				return errors.Wrapf(err, "failed to get cache key %s", command.Prerequisite)
			}
		} else {
			if err := s.runCommand(ctx, name, reqCmd); err != nil {
				return err
			}
		}
	}

	for _, cmd := range command.Commands {
		if err := s.executeScaffold(ctx, cmd, path); err != nil {
			return err
		}
	}
	return s.cache.Save(command.Name, path)
}

func (s *Scaffold) executeScaffold(ctx context.Context, cmd, path string) error {
	args := append([]string{s.binary, "scaffold"}, strings.Fields(cmd)...)
	args = append(args, "--path", path)
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
	// name of chain at the given '--path', so we need to append "example" to the path if the
	// command is not "chain".
	if ver.LessThan(v027) && args[2] != "chain" {
		args[len(args)-1] = filepath.Join(args[len(args)-1], "example")
	}
	return args
}

// applyPostScaffoldExceptions this function we can manipulate the Output of scaffold command after
// they have been executed in order to compensate for differences in versions.
func applyPostScaffoldExceptions(ver *semver.Version, name string, output string) error {
	// In versions <0.27.0, "scaffold chain" command always creates a new directory with the name of
	// chain at the given '--path', so we need to move the directory to the parent directory.
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
