package cosmosgen

import (
	"context"
	"os"
	"path/filepath"
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

	if g.o.jsIncludeThirdParty {
		for _, modules := range g.thirdModules {
			data.Modules = append(data.Modules, modules...)
		}
	}

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
				changed, err := dirchange.HasDirChecksumChanged(dirCache, cacheKey, sourcePath, paths...)
				if err != nil {
					return err
				}

				if !changed {
					return nil
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

	if g.g.o.jsIncludeThirdParty {
		for sourcePath, modules := range g.g.thirdModules {
			add(sourcePath, modules)
		}
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
		protoc.Plugin(tsprotoPluginPath, "--ts_proto_opt=snakeToCamel=false"),
		protoc.Env("NODE_OPTIONS="), // unset nodejs options to avoid unexpected issues with vercel "pkg"
		protoc.WithCommand(protocCmd),
	)
	if err != nil {
		return err
	}

	// generate OpenAPI spec
	oaitemp, err := os.MkdirTemp("", "gen-js-openapi-module-spec")
	if err != nil {
		return err
	}

	defer os.RemoveAll(oaitemp)

	err = protoc.Generate(
		ctx,
		oaitemp,
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
		srcspec = filepath.Join(oaitemp, "apidocs.swagger.json")
		outREST = filepath.Join(out, "rest.ts")
	)

	if err := sta.Generate(ctx, outREST, srcspec, sta.WithCommand(staCmd)); err != nil {
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
	outDir := filepath.Join(g.g.o.tsClientRootPath)
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	return templateTSClientRoot.Write(outDir, "", p)
}
