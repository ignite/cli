package cosmosgen

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
)

var dirchangeCacheNamespace = "generate.typescript.dirchange"

type tsGenerator struct {
	g *generator
}

type generatePayload struct {
	Modules   []module.Module
	PackageNS string
}

func newTSGenerator(g *generator) *tsGenerator {
	return &tsGenerator{g}
}

func (g *generator) tsTemplate() string {
	return filepath.Join(g.appPath, g.protoDir, "buf.gen.ts.yaml")
}

func (g *generator) generateTS(ctx context.Context) error {
	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	data := generatePayload{
		Modules:   g.appModules,
		PackageNS: strings.ReplaceAll(appModulePath, "/", "-"),
	}

	// Make sure the modules are always sorted to keep the import
	// and module registration order consistent so the generated
	// files are not changed.
	sort.SliceStable(data.Modules, func(i, j int) bool {
		return data.Modules[i].Pkg.Name < data.Modules[j].Pkg.Name
	})

	tsg := newTSGenerator(g)
	if err := tsg.generateModuleTemplates(ctx); err != nil {
		return err
	}

	return tsg.generateRootTemplates(data)
}

func (g *tsGenerator) generateModuleTemplates(ctx context.Context) error {
	dirCache := cache.New[[]byte](g.g.cacheStorage, dirchangeCacheNamespace)
	add := func(sourcePath string, m module.Module) error {
		cacheKey := m.Pkg.Path
		paths := []string{m.Pkg.Path, g.g.opts.jsOut(m)}

		// Always generate module templates by default unless cache is enabled, in which
		// case the module template is generated when one or more files were changed in
		// the module since the last generation.
		if g.g.opts.useCache {
			changed, err := dirchange.HasDirChecksumChanged(dirCache, cacheKey, sourcePath, paths...)
			if err != nil {
				return err
			}

			if !changed {
				return nil
			}
		}

		if err := g.generateModuleTemplate(ctx, sourcePath, m); err != nil {
			return err
		}

		return dirchange.SaveDirChecksum(dirCache, cacheKey, sourcePath, paths...)
	}

	for _, m := range g.g.appModules {
		if err := add(g.g.appPath, m); err != nil {
			return err
		}
	}

	// Always generate third party modules; This is required because not generating them might
	// lead to issues with the module registration in the root template. The root template must
	// always be generated with 3rd party modules which means that if a new 3rd party module
	// is available and not generated it would lead to the registration of a new not generated
	// 3rd party module.
	for sourcePath, modules := range g.g.thirdModules {
		for _, m := range modules {
			if err := add(sourcePath, m); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *tsGenerator) generateModuleTemplate(
	ctx context.Context,
	appPath string,
	m module.Module,
) error {
	var (
		out      = g.g.opts.jsOut(m)
		typesOut = filepath.Join(out, "types")
	)

	if err := os.MkdirAll(typesOut, 0o766); err != nil {
		return err
	}
	if err := generateRouteNameFile(typesOut); err != nil {
		return err
	}

	// All "cosmossdk.io" module packages must use SDK's
	// proto path which is where the proto files are stored.
	protoPath := filepath.Join(appPath, g.g.protoDir) // use module app path
	if module.IsCosmosSDKModulePkg(appPath) {
		protoPath = filepath.Join(g.g.sdkDir, "proto")
	}

	// code generate for each module.
	if err := g.g.buf.Generate(
		ctx,
		protoPath,
		typesOut,
		g.g.tsTemplate(),
		cosmosbuf.IncludeWKT(),
	); err != nil {
		return err
	}

	// Generate the module template
	if err := templateTSClientModule.Write(out, protoPath, struct {
		Module module.Module
	}{
		Module: m,
	}); err != nil {
		return err
	}

	// Generate the rest API template (using axios)
	return templateTSClientRest.Write(out, protoPath, struct {
		Module module.Module
	}{
		Module: m,
	})
}

func (g *tsGenerator) generateRootTemplates(p generatePayload) error {
	outDir := g.g.opts.tsClientRootPath
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientRoot.Write(outDir, "", p)
}
