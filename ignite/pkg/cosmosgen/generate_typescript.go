package cosmosgen

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/dirchange"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/nodetime/programs/sta"
	tsproto "github.com/ignite/cli/ignite/pkg/nodetime/programs/ts-proto"
	"github.com/ignite/cli/ignite/pkg/protoc"
)

var (
	dirchangeCacheNamespace = "generate.typescript.dirchange"
	jsOpenAPIOut            = []string{"--openapiv2_out=logtostderr=true,allow_merge=true,json_names_for_fields=false,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:."}
	tsOut                   = []string{"--ts_proto_out=."}
)

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

func (g *generator) generateTS() error {
	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	appModulePath := gomodulepath.ExtractAppPath(chainPath.RawPath)
	data := generatePayload{
		Modules:   g.appModules,
		PackageNS: strings.ReplaceAll(appModulePath, "/", "-"),
	}

	// Third party modules are always required to generate the root
	// template because otherwise it would be generated only with
	// custom modules losing the registration of the third party
	// modules when the root templates are re-generated.
	for _, modules := range g.thirdModules {
		data.Modules = append(data.Modules, modules...)
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

func (g *tsGenerator) generateModuleTemplates() error {
	protocCmd, cleanupProtoc, err := protoc.Command()
	if err != nil {
		return err
	}

	defer cleanupProtoc()

	tsprotoPluginPath, cleanupPlugin, err := tsproto.BinaryPath()
	if err != nil {
		return err
	}

	defer cleanupPlugin()

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

				err = g.generateModuleTemplate(g.g.ctx, protocCmd, staCmd, tsprotoPluginPath, sourcePath, m)
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
	protocCmd protoc.Cmd,
	staCmd sta.Cmd,
	tsprotoPluginPath, appPath string,
	m module.Module,
) error {
	var (
		out      = g.g.o.jsOut(m)
		typesOut = filepath.Join(out, "types")
	)

	includePaths, err := g.g.resolveInclude(appPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(typesOut, 0o766); err != nil {
		return err
	}

	// generate ts-proto types
	err = protoc.Generate(
		ctx,
		typesOut,
		m.Pkg.Path,
		includePaths,
		tsOut,
		protoc.Plugin(tsprotoPluginPath, "--ts_proto_opt=snakeToCamel=true", "--ts_proto_opt=esModuleInterop=true"),
		protoc.Env("NODE_OPTIONS="), // unset nodejs options to avoid unexpected issues with vercel "pkg"
		protoc.WithCommand(protocCmd),
	)
	if err != nil {
		return err
	}

	// generate OpenAPI spec
	tmp, err := os.MkdirTemp("", "gen-js-openapi-module-spec")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	err = protoc.Generate(
		ctx,
		tmp,
		m.Pkg.Path,
		includePaths,
		jsOpenAPIOut,
		protoc.WithCommand(protocCmd),
	)
	if err != nil {
		return err
	}

	// generate the REST client from the OpenAPI spec
	var (
		srcSpec = filepath.Join(tmp, "apidocs.swagger.json")
		outREST = filepath.Join(out, "rest.ts")
	)

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
