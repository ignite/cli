package cosmosgen

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/dirchange"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/nodetime/programs/sta"
	swaggercombine "github.com/ignite/cli/ignite/pkg/nodetime/programs/swagger-combine"
	"github.com/ignite/cli/ignite/pkg/xos"
)

var dirchangeCacheNamespace = "generate.typescript.dirchange"

type tsGenerator struct {
	g *generator
}

type generatePayload struct {
	Modules         []module.Module
	PackageNS       string
	IsConsumerChain bool
}

func newTSGenerator(g *generator) *tsGenerator {
	return &tsGenerator{g}
}

func (g *generator) generateTS() error {
	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	data := generatePayload{
		Modules:         g.appModules,
		PackageNS:       strings.ReplaceAll(appModulePath, "/", "-"),
		IsConsumerChain: false,
	}

	// Third party modules are always required to generate the root
	// template because otherwise it would be generated only with
	// custom modules losing the registration of the third party
	// modules when the root templates are re-generated.
	for _, modules := range g.thirdModules {
		data.Modules = append(data.Modules, modules...)
		for _, m := range modules {
			if strings.HasPrefix(m.Pkg.Name, "interchain_security.ccv.consumer") {
				data.IsConsumerChain = true
			}
		}
	}
	// Make sure the modules are always sorted to keep the import
	// and module registration order consistent so the generated
	// files are not changed.
	sort.SliceStable(data.Modules, func(i, j int) bool {
		return data.Modules[i].Pkg.Name < data.Modules[j].Pkg.Name
	})

	tsg := newTSGenerator(g)
	if err := tsg.generateModuleTemplates(); err != nil {
		return err
	}

	return tsg.generateRootTemplates(data)
}

func (g *tsGenerator) tsTemplate() string {
	return filepath.Join(g.g.appPath, g.g.protoDir, "buf.gen.ts.yaml")
}

func (g *tsGenerator) generateModuleTemplates() error {
	staCmd, cleanupSTA, err := sta.Command()
	if err != nil {
		return err
	}

	defer cleanupSTA()

	gg := &errgroup.Group{}
	dirCache := cache.New[[]byte](g.g.cacheStorage, dirchangeCacheNamespace)
	add := func(sourcePath string, modules []module.Module) {
		for _, m := range modules {
			m := m

			gg.Go(func() error {
				cacheKey := m.Pkg.Path
				paths := append([]string{m.Pkg.Path, g.g.o.jsOut(m)}, g.g.o.includeDirs...)

				// Always generate module templates by default unless cache is enabled, in which
				// case the module template is generated when one or more files were changed in
				// the module since the last generation.
				if g.g.o.useCache {
					changed, err := dirchange.HasDirChecksumChanged(dirCache, cacheKey, sourcePath, paths...)
					if err != nil {
						return err
					}

					if !changed {
						return nil
					}
				}

				err = g.generateModuleTemplate(g.g.ctx, staCmd, sourcePath, m)
				if err != nil {
					return err
				}

				return dirchange.SaveDirChecksum(dirCache, cacheKey, sourcePath, paths...)
			})
		}
	}

	add(g.g.appPath, g.g.appModules)

	// Always generate third party modules; This is required because not generating them might
	// lead to issues with the module registration in the root template. The root template must
	// always be generated with 3rd party modules which means that if a new 3rd party module
	// is available and not generated it would lead to the registration of a new not generated
	// 3rd party module.
	for sourcePath, modules := range g.g.thirdModules {
		add(sourcePath, modules)
	}

	return gg.Wait()
}

func (g *tsGenerator) generateModuleTemplate(
	ctx context.Context,
	staCmd sta.Cmd,
	appPath string,
	m module.Module,
) error {
	var (
		out      = g.g.o.jsOut(m)
		typesOut = filepath.Join(out, "types")
		conf     = swaggercombine.Config{
			Swagger: "2.0",
			Info: swaggercombine.Info{
				Title: "HTTP API Console",
			},
		}
	)

	if err := os.MkdirAll(typesOut, 0o766); err != nil {
		return err
	}

	// generate ts-proto types
	if err := g.g.buf.Generate(
		ctx,
		m.Pkg.Path,
		typesOut,
		g.tsTemplate(),
		"module.proto",
	); err != nil {
		return err
	}

	// generate OpenAPI spec
	tmp, err := os.MkdirTemp("", "gen-js-openapi-module-spec")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	if err := g.g.buf.Generate(
		ctx,
		m.Pkg.Path,
		tmp,
		g.g.openAPITemplate(),
		"module.proto",
	); err != nil {
		return err
	}

	// combine all swagger files
	specs, err := xos.FindFiles(tmp, xos.JSONFile)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if err := conf.AddSpec(strcase.ToCamel(m.Pkg.Name), spec); err != nil {
			return err
		}
	}

	// combine specs into one and save to out.
	srcSpec := filepath.Join(tmp, "apidocs.swagger.json")
	if err := swaggercombine.Combine(ctx, conf, srcSpec); err != nil {
		return err
	}

	// generate the REST client from the OpenAPI spec
	outREST := filepath.Join(out, "rest.ts")
	if err := sta.Generate(ctx, outREST, srcSpec, sta.WithCommand(staCmd)); err != nil {
		return err
	}

	pp := filepath.Join(appPath, g.g.protoDir)

	return templateTSClientModule.Write(out, pp, struct {
		Module module.Module
	}{
		Module: m,
	})
}

func (g *tsGenerator) generateRootTemplates(p generatePayload) error {
	outDir := g.g.o.tsClientRootPath
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientRoot.Write(outDir, "", p)
}
