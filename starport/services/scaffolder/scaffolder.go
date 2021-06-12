// Package scaffolder initializes Starport apps and modifies existing ones
// to add more features in a later time.
package scaffolder

import (
	"context"
	"errors"
	"os"
	"strings"

	conf "github.com/tendermint/starport/starport/chainconf"
	sperrors "github.com/tendermint/starport/starport/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gocmd"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/services"
)

// Scaffolder is Starport app scaffolder.
type Scaffolder struct {
	// path is app's path on the filesystem.
	path string

	// options to configure scaffolding.
	options *scaffoldingOptions

	// version of the chain
	version cosmosver.Version
}

// Option configures scaffolding.
type Option func(*scaffoldingOptions)

// scaffoldingOptions keeps set of options to apply scaffolding.
type scaffoldingOptions struct {
	addressPrefix string
}

func newOptions(options ...Option) *scaffoldingOptions {
	opts := &scaffoldingOptions{}
	opts.apply(options...)
	return opts
}

func (s *scaffoldingOptions) apply(options ...Option) {
	for _, o := range options {
		o(s)
	}
}

// AddressPrefix configures address prefix for the app.
func AddressPrefix(prefix string) Option {
	return func(o *scaffoldingOptions) {
		o.addressPrefix = strings.ToLower(prefix)
	}
}

// New initializes a new Scaffolder for app at path.
func New(path string, options ...Option) (*Scaffolder, error) {
	s := &Scaffolder{
		path:    path,
		options: newOptions(options...),
	}

	// determine the chain version.
	var err error
	s.version, err = cosmosver.Detect(path)
	if err != nil && !errors.Is(err, gomodule.ErrGoModNotFound) {
		return nil, err
	}
	if err == nil && !s.version.Major().Is(cosmosver.Stargate) {
		return nil, sperrors.ErrOnlyStargateSupported
	}

	return s, nil
}

func owner(modulePath string) string {
	return strings.Split(modulePath, "/")[1]
}

func (s *Scaffolder) finish(path, gomodPath string) error {
	if err := protoc(path, gomodPath); err != nil {
		return err
	}
	if err := tidy(path); err != nil {
		return err
	}
	return fmtProject(path)
}

func protoc(projectPath, gomodPath string) error {
	if err := cosmosgen.InstallDependencies(context.Background(), projectPath); err != nil {
		return err
	}

	confpath, err := conf.LocateDefault(projectPath)
	if err != nil {
		return err
	}
	conf, err := conf.ParseFile(confpath)
	if err != nil {
		return err
	}

	options := services.CodegenOptions(projectPath, gomodPath, false, conf)

	return cosmosgen.Generate(context.Background(), projectPath, conf.Build.Proto.Path, options...)
}

func tidy(path string) error {
	return cmdrunner.
		New(
			cmdrunner.DefaultStderr(os.Stderr),
			cmdrunner.DefaultWorkdir(path),
		).
		Run(context.Background(),
			step.New(
				step.Exec(gocmd.Name(), "mod", "tidy"),
			),
		)
}

func fmtProject(path string) error {
	return cmdrunner.
		New(
			cmdrunner.DefaultStderr(os.Stderr),
			cmdrunner.DefaultWorkdir(path),
		).
		Run(context.Background(),
			step.New(
				step.Exec(
					gocmd.Name(),
					"fmt",
					"./...",
				),
			),
		)
}
