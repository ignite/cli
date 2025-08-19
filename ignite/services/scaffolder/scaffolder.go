// Package scaffolder initializes Ignite CLI apps and modifies existing ones
// to add more features in a later time.
package scaffolder

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/version"
)

// Scaffolder is Ignite CLI app scaffolder.
type Scaffolder struct {
	// Version of the chain
	Version cosmosver.Version

	// appPath path of the app.
	appPath string

	// protoDir path of the proto folder.
	protoDir string

	// modpath represents the go module Path of the app.
	modpath gomodulepath.Path

	// runner represents the scaffold xgenny runner.
	runner *xgenny.Runner
}

// New creates a new scaffold app.
func New(context context.Context, appPath, protoDir string) (Scaffolder, error) {
	path, err := filepath.Abs(appPath)
	if err != nil {
		return Scaffolder{}, err
	}

	modpath, path, err := gomodulepath.Find(path)
	if err != nil {
		return Scaffolder{}, err
	}

	ver, err := cosmosver.Detect(path)
	if err != nil {
		return Scaffolder{}, err
	}

	// Make sure that the app was scaffolded with a supported Cosmos SDK version
	if err := version.AssertSupportedCosmosSDKVersion(ver); err != nil {
		return Scaffolder{}, err
	}

	if err := cosmosanalysis.IsChainPath(path); err != nil {
		return Scaffolder{}, err
	}

	s := Scaffolder{
		Version:  ver,
		appPath:  path,
		protoDir: protoDir,
		modpath:  modpath,
		runner:   xgenny.NewRunner(context, path),
	}

	return s, nil
}

func (s Scaffolder) ApplyModifications(options ...xgenny.ApplyOption) (xgenny.SourceModification, error) {
	return s.runner.ApplyModifications(options...)
}

func (s Scaffolder) Tracer() *placeholder.Tracer {
	return s.runner.Tracer()
}

func (s Scaffolder) Run(gens ...*genny.Generator) error {
	return s.runner.Run(gens...)
}

func (s Scaffolder) PostScaffold(ctx context.Context, cacheStorage cache.Storage, skipProto bool) error {
	return PostScaffold(ctx, cacheStorage, s.appPath, s.protoDir, s.modpath.RawPath, skipProto)
}

func PostScaffold(ctx context.Context, cacheStorage cache.Storage, path, protoDir, goModPath string, skipProto bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Errorf("failed to get current working directory: %w", err)
	}

	// go to the app path, in other to use the go tools
	if err := os.Chdir(path); err != nil {
		return errors.Errorf("failed to change directory to %s: %w", path, err)
	}

	if !skipProto {
		// go mod tidy prior and after the proto generation is required.
		if err := gocmd.ModTidy(ctx, path); err != nil {
			return err
		}

		if err := protoc(ctx, cacheStorage, path, protoDir, goModPath); err != nil {
			return err
		}
	}

	if err := gocmd.ModTidy(ctx, path); err != nil {
		return err
	}

	if err := gocmd.Fmt(ctx, path); err != nil {
		return err
	}

	if err := gocmd.GoImports(ctx, path); err != nil {
		return err
	}

	// return to the original working directory
	if err := os.Chdir(wd); err != nil {
		return errors.Errorf("failed to change directory to %s: %w", wd, err)
	}

	return nil
}

func protoc(ctx context.Context, cacheStorage cache.Storage, projectPath, protoDir, goModPath string) error {
	confpath, err := chainconfig.LocateDefault(projectPath)
	if err != nil {
		return err
	}
	conf, err := chainconfig.ParseFile(confpath)
	if err != nil {
		return err
	}

	options := []cosmosgen.Option{
		cosmosgen.UpdateBufModule(),
		cosmosgen.WithGoGeneration(),
	}

	// Generate Typescript client code if it's enabled
	if conf.Client.Typescript.Path != "" { //nolint:staticcheck,nolintlint
		tsClientPath := chainconfig.TSClientPath(*conf)
		if !filepath.IsAbs(tsClientPath) {
			tsClientPath = filepath.Join(projectPath, tsClientPath)
		}

		options = append(options,
			cosmosgen.WithTSClientGeneration(
				cosmosgen.TypescriptModulePath(tsClientPath),
				tsClientPath,
				true,
			),
		)
	}

	if conf.Client.OpenAPI.Path != "" {
		openAPIPath := conf.Client.OpenAPI.Path
		if !filepath.IsAbs(openAPIPath) {
			openAPIPath = filepath.Join(projectPath, openAPIPath)
		}

		options = append(options, cosmosgen.WithOpenAPIGeneration(openAPIPath))
	}

	if err := cosmosgen.Generate(
		ctx,
		cacheStorage,
		projectPath,
		protoDir,
		goModPath,
		chainconfig.DefaultVuePath,
		options...,
	); err != nil {
		return err
	}

	buf, err := cosmosbuf.New(cacheStorage, goModPath)
	if err != nil {
		return err
	}

	return buf.Format(ctx, projectPath)
}
