package cosmosgen

import (
	"path/filepath"
	"strings"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/protopath"
)

var sdkImport = "github.com/cosmos/cosmos-sdk"

func (g *generator) setup() (err error) {
	// Cosmos SDK hosts proto files of own x/ modules and some third party ones needed by itself and
	// blockchain apps. Generate should be aware of these and make them available to the blockchain
	// app that wants to generate code for its own proto.
	//
	// blockchain apps may use different versions of the SDK. following code first makes sure that
	// app's dependencies are download by 'go mod' and cached under the local filesystem.
	// and then, it determines which version of the SDK is used by the app and what is the absolute path
	// of its source code.
	if err := cmdrunner.
		New(cmdrunner.DefaultWorkdir(g.appPath)).
		Run(g.ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
		return err
	}

	// parse the go.mod of the app and extract dependencies.
	modfile, err := gomodule.ParseAt(g.appPath)
	if err != nil {
		return err
	}

	g.deps, err = gomodule.ResolveDependencies(modfile)

	return
}

func (g *generator) resolveInclude(path string) (paths []string, err error) {
	paths = append(paths, filepath.Join(path, g.protoDir))
	for _, p := range g.o.includeDirs {
		paths = append(paths, filepath.Join(path, p))
	}

	includePaths, err := protopath.ResolveDependencyPaths(g.deps,
		protopath.NewModule(sdkImport, append([]string{g.protoDir}, g.o.includeDirs...)...))
	if err != nil {
		return nil, err
	}

	paths = append(paths, includePaths...)
	return paths, nil
}

func (g *generator) discoverModules(path string) ([]module.Module, error) {
	var filteredModules []module.Module

	modules, err := module.Discover(g.ctx, path)
	if err != nil {
		return nil, err
	}

	for _, m := range modules {
		pp := filepath.Join(path, g.protoDir)
		if !strings.HasPrefix(m.Pkg.Path, pp) {
			continue
		}
		filteredModules = append(filteredModules, m)
	}

	return filteredModules, nil
}
