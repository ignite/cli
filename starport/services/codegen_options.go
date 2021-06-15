package services

import (
	"path/filepath"

	starportconf "github.com/tendermint/starport/starport/chainconf"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/cosmosgen"
	"github.com/tendermint/starport/starport/pkg/giturl"
)

// CodegenOptions gets client code generation options.
func CodegenOptions(
	projecPath string,
	goModPath string,
	enableThirdPartyModuleCodegen bool,
	conf starportconf.Config,
) []cosmosgen.Option {
	options := []cosmosgen.Option{
		cosmosgen.WithGoGeneration(goModPath),
		cosmosgen.IncludeDirs(conf.Build.Proto.ThirdPartyPaths),
	}

	// generate Vuex code as well if it is enabled.
	if conf.Client.Vuex.Path != "" {
		storeRootPath := filepath.Join(projecPath, conf.Client.Vuex.Path, "generated")
		options = append(options,
			cosmosgen.WithVuexGeneration(
				enableThirdPartyModuleCodegen,
				getModulePathFn(storeRootPath),
				storeRootPath,
			),
		)
	}

	// generate ts client code as well if it is enabled.
	// NB: Vuex generates JS/TS client code as well but in addition to a vuex store.
	// Path options will conflict with each other.
	if conf.Client.Typescript.Path != "" {
		tsRootPath := filepath.Join(projecPath, conf.Client.Typescript.Path)
		options = append(options,
			cosmosgen.WithTSGeneration(
				enableThirdPartyModuleCodegen,
				getModulePathFn(tsRootPath),
			),
		)

	}

	if conf.Client.OpenAPI.Path != "" {
		options = append(options, cosmosgen.WithOpenAPIGeneration(conf.Client.OpenAPI.Path))
	}

	return options
}

func getModulePathFn(path string) func(m module.Module) string {
	return func(m module.Module) string {
		parsedGitURL, _ := giturl.Parse(m.Pkg.GoImportName)
		return filepath.Join(path, parsedGitURL.UserAndRepo(), m.Pkg.Name, "module")
	}
}
