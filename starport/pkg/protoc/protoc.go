// Package protoc provides high level access to protoc command.
package protoc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner/exec"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/localfs"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
	"github.com/tendermint/starport/starport/pkg/protoc/data"
)

// Option configures Generate configs.
type Option func(*configs)

// configs holds Generate configs.
type configs struct {
	pluginPath             string
	isGeneratedDepsEnabled bool
}

// Plugin configures a plugin for code generation.
func Plugin(path string) Option {
	return func(c *configs) {
		c.pluginPath = path
	}
}

// GenerateDependencies enables code generation for the proto files that your protofile depends on.
// use this if your protoc plugin does not give you an option to enable the same feature.
func GenerateDependencies() Option {
	return func(c *configs) {
		c.isGeneratedDepsEnabled = true
	}
}

type Cmd struct {
	Command  []string
	Included []string
}

// Command sets the protoc binary up and returns the command needed to execute c.
func Command() (command Cmd, cleanup func(), err error) {
	path, cleanupProto, err := localfs.SaveBytesTemp(data.Binary(), "protoc", 0755)
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
		Command:  []string{path, "-I", include},
		Included: []string{include},
	}

	return command, cleanup, nil
}

// Generate generates code into outDir from protoPath and its includePaths by using plugins provided with protocOuts.
func Generate(ctx context.Context, outDir, protoPath string, includePaths, protocOuts []string, options ...Option) error {
	c := configs{}
	for _, o := range options {
		o(&c)
	}

	cmd, cleanup, err := Command()
	if err != nil {
		return err
	}
	defer cleanup()

	var command []string

	// add plugin if set.
	if c.pluginPath != "" {
		command = append(cmd.Command, "--plugin", c.pluginPath)
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
	files, err := discoverFiles(ctx, c, protoPath, append(cmd.Included, existentIncludePaths...), protoanalysis.NewCache())
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

// discoverFiles discovers .proto files to do code generation for. .proto files of the app
// (everything under protoPath) will always be a part of the disovered files.
//
// when .proto files of the app depends on another proto package under includePaths (dependencies), those
// ones may need to be discovered as well. some protoc plugins already do this discovery internally but
// for the ones that don't, it needs to be handled here if GenerateDependencies() is enabled.
func discoverFiles(ctx context.Context, c configs, protoPath string, includePaths []string, cache protoanalysis.Cache) (
	discovered []string, err error) {
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
		if _, err := os.Stat(guessedPath); err == nil {
			discovered = append(discovered, guessedPath)
			continue
		}

		// otherwise, search by absolute path in includePaths.
		var found bool
		for _, included := range includePaths {

			guessedPath := filepath.Join(included, dep)
			if _, err := os.Stat(guessedPath); err == nil {
				// found the dependency.
				// if it's under protoPath, it is already discovered so, skip it.
				if !strings.HasPrefix(guessedPath, protoPath) {
					discovered = append(discovered, guessedPath)

					// perform a complete search on this one as well to discover its dependencies as well.
					d, err := searchFile(protoanalysis.File{Path: dep}, protoPath, includePaths)
					if err != nil {
						return nil, err
					}
					discovered = append(discovered, d...)
				}

				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("cannot locate dependency %q for %q", dep, file.Path)
		}
	}

	return
}
