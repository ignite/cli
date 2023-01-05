// Package protoc provides high level access to protoc command.
package protoc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/localfs"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/ignite/pkg/protoc/data"
)

// Option configures Generate configs.
type Option func(*configs)

// configs holds Generate configs.
type configs struct {
	pluginPath             string
	isGeneratedDepsEnabled bool
	pluginOptions          []string
	env                    []string
	command                Cmd
}

// Plugin configures a plugin for code generation.
func Plugin(path string, options ...string) Option {
	return func(c *configs) {
		c.pluginPath = path
		c.pluginOptions = options
	}
}

// GenerateDependencies enables code generation for the proto files that your protofile depends on.
// use this if your protoc plugin does not give you an option to enable the same feature.
func GenerateDependencies() Option {
	return func(c *configs) {
		c.isGeneratedDepsEnabled = true
	}
}

// Env assigns environment values during the code generation.
func Env(v ...string) Option {
	return func(c *configs) {
		c.env = v
	}
}

// WithCommand assigns a protoc command to use for code generation.
// This allows to use a single protoc binary in multiple code generation calls.
// Otherwise, `Generate` creates a new protoc binary each time it is called.
func WithCommand(command Cmd) Option {
	return func(c *configs) {
		c.command = command
	}
}

// Cmd contains the information necessary to execute the protoc command.
type Cmd struct {
	command  []string
	includes []string
}

// Command returns the strings to execute the `protoc` command.
func (c Cmd) Command() []string {
	return c.command
}

// Includes returns the proto files import paths.
func (c Cmd) Includes() []string {
	return c.includes
}

// Command sets the protoc binary up and returns the command needed to execute c.
func Command() (command Cmd, cleanup func(), err error) {
	path, cleanupProto, err := localfs.SaveBytesTemp(data.Binary(), "protoc", 0o755)
	if err != nil {
		return Cmd{}, nil, err
	}

	include, cleanupInclude, err := localfs.SaveTemp(data.Include())
	if err != nil {
		cleanupProto()
		return Cmd{}, nil, err
	}

	cleanup = func() {
		cleanupProto()
		cleanupInclude()
	}

	command = Cmd{
		command:  []string{path, "-I", include},
		includes: []string{include},
	}

	return command, cleanup, nil
}

// Generate generates code into outDir from protoPath and its includePaths by using plugins provided with protocOuts.
func Generate(ctx context.Context, outDir, protoPath string, includePaths, protocOuts []string, options ...Option) error {
	c := configs{}

	for _, o := range options {
		o(&c)
	}

	// init the string to run the protoc command and the proto files import path
	command := c.command.Command()
	includes := c.command.Includes()

	if command == nil {
		cmd, cleanup, err := Command()
		if err != nil {
			return err
		}

		defer cleanup()

		command = cmd.Command()
		includes = cmd.Includes()
	}

	// add plugin if set.
	if c.pluginPath != "" {
		command = append(command, "--plugin", c.pluginPath)
	}
	var existentIncludePaths []string

	// skip if a third party proto source actually doesn't exist on the filesystem.
	for _, path := range includePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		existentIncludePaths = append(existentIncludePaths, path)
	}

	// append third party proto locations to the command.
	for _, importPath := range existentIncludePaths {
		command = append(command, "-I", importPath)
	}

	// find out the list of proto files to generate code for and perform code generation.
	files, err := discoverFiles(ctx, c, protoPath, append(includes, existentIncludePaths...), protoanalysis.NewCache())
	if err != nil {
		return err
	}

	// run command for each protocOuts.
	for _, out := range protocOuts {
		command := append(command, out)
		command = append(command, files...)
		command = append(command, c.pluginOptions...)

		execOpts := []exec.Option{
			exec.StepOption(step.Workdir(outDir)),
			exec.IncludeStdLogsToError(),
		}
		if c.env != nil {
			execOpts = append(execOpts, exec.StepOption(step.Env(c.env...)))
		}

		if err := exec.Exec(ctx, command, execOpts...); err != nil {
			return err
		}
	}

	return nil
}

// discoverFiles discovers .proto files to do code generation for. .proto files of the app
// (everything under protoPath) will always be a part of the discovered files.
//
// when .proto files of the app depends on another proto package under includePaths (dependencies), those
// may need to be discovered as well. some protoc plugins already do this discovery internally but
// for the ones that don't, it needs to be handled here if GenerateDependencies() is enabled.
func discoverFiles(ctx context.Context, c configs, protoPath string, includePaths []string, cache protoanalysis.Cache) (
	discovered []string, err error,
) {
	packages, err := protoanalysis.Parse(ctx, cache, protoPath)
	if err != nil {
		return nil, err
	}

	discovered = packages.Files().Paths()

	if !c.isGeneratedDepsEnabled {
		return discovered, nil
	}

	for _, file := range packages.Files() {
		d, err := searchFile(file, protoPath, includePaths)
		if err != nil {
			return nil, err
		}
		discovered = append(discovered, d...)
	}

	return discovered, nil
}

func searchFile(file protoanalysis.File, protoPath string, includePaths []string) (discovered []string, err error) {
	dir := filepath.Dir(file.Path)

	for _, dep := range file.Dependencies {
		// try to locate imported .proto file relative to the this .proto file.
		guessedPath := filepath.Join(dir, dep)
		_, err := os.Stat(guessedPath)
		if err == nil {
			discovered = append(discovered, guessedPath)
			continue
		}
		if !os.IsNotExist(err) {
			return nil, err
		}

		// otherwise, search by absolute path in includePaths.
		var found bool
		for _, included := range includePaths {
			guessedPath := filepath.Join(included, dep)
			_, err := os.Stat(guessedPath)
			if err == nil {
				// found the dependency.
				// if it's under protoPath, it is already discovered so, skip it.
				if !strings.HasPrefix(guessedPath, protoPath) {
					discovered = append(discovered, guessedPath)

					// perform a complete search on this one to discover its dependencies as well.
					depFile, err := protoanalysis.ParseFile(guessedPath)
					if err != nil {
						return nil, err
					}
					d, err := searchFile(depFile, protoPath, includePaths)
					if err != nil {
						return nil, err
					}
					discovered = append(discovered, d...)
				}

				found = true
				break
			}
			if !os.IsNotExist(err) {
				return nil, err
			}
		}

		if !found {
			return nil, fmt.Errorf("cannot locate dependency %q for %q", dep, file.Path)
		}
	}

	return discovered, nil
}
