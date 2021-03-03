package cosmosgen

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mattn/go-zglob"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/nodetime/sta"
	tsproto "github.com/tendermint/starport/starport/pkg/nodetime/ts-proto"
	"github.com/tendermint/starport/starport/pkg/nodetime/tsc"
	"github.com/tendermint/starport/starport/pkg/protoc"
	"golang.org/x/sync/errgroup"
)

var (
	tsOut = []string{
		"--ts_proto_out=.",
	}

	openAPIOut = []string{
		"--openapiv2_out=logtostderr=true,allow_merge=true:.",
	}

	vuexRootMarker = "vuex-root"
)

type jsGenerator struct {
	g                 *generator
	tsprotoPluginPath string
}

func newJSGenerator(g *generator) (jsGenerator, error) {
	tsprotoPluginPath, err := tsproto.BinaryPath()
	if err != nil {
		return jsGenerator{}, err
	}

	return jsGenerator{
		g:                 g,
		tsprotoPluginPath: tsprotoPluginPath,
	}, nil
}

func (g *generator) generateJS() error {
	jsg, err := newJSGenerator(g)
	if err != nil {
		return err
	}

	if err := jsg.generateModules(); err != nil {
		return err
	}

	if err := jsg.generateVuexModuleLoader(); err != nil {
		return err
	}

	return nil
}

func (g *jsGenerator) generateModules() error {
	// sourcePaths keeps a list of root paths of Go projects (source codes) that might contain
	// Cosmos SDK modules inside.
	sourcePaths := []string{
		g.g.appPath, // user's blockchain. may contain internal modules. it is the first place to look for.
	}

	if g.g.o.jsIncludeThirdParty {
		// go through the Go dependencies (inside go.mod) of each source path, some of them might be hosting
		// Cosmos SDK modules that could be in use by user's blockchain.
		//
		// Cosmos SDK is a dependency of all blockchains, so it's absolute that we'll be discovering all modules of the
		// SDK as well during this process.
		//
		// even if a dependency contains some SDK modules, not all of these modules could be used by user's blockchain.
		// this is fine, we can still generate JS clients for those non modules, it is up to user to use (import in JS)
		// not use generated modules.
		// not used ones will never get resolved inside JS environment and will not ship to production, JS bundlers will avoid.
		//
		// TODO(ilgooz): we can still implement some sort of smart filtering to detect non used modules by the user's blockchain
		// at some point, it is a nice to have.
		for _, dep := range g.g.deps {
			deppath, err := gomodule.LocatePath(dep)
			if err != nil {
				return err
			}
			sourcePaths = append(sourcePaths, deppath)
		}
	}

	gs := &errgroup.Group{}

	// try to discover SDK modules in all source paths.
	for _, sourcePath := range sourcePaths {
		sourcePath := sourcePath

		gs.Go(func() error {
			modules, err := g.g.discoverModules(sourcePath)
			if err != nil {
				return err
			}

			gg, ctx := errgroup.WithContext(g.g.ctx)

			// do code generation for each found module.
			for _, m := range modules {
				m := m

				gg.Go(func() error { return g.generateModule(ctx, g.tsprotoPluginPath, sourcePath, m) })
			}

			return gg.Wait()
		})
	}

	return gs.Wait()
}

// generateModule generates generates JS code for a module.
func (g *jsGenerator) generateModule(ctx context.Context, tsprotoPluginPath, appPath string, m module.Module) error {
	var (
		out          = g.g.o.jsOut(m)
		storeDirPath = filepath.Dir(out)
		typesOut     = filepath.Join(out, "types")
	)

	includePaths, err := g.g.resolveInclude(appPath)
	if err != nil {
		return err
	}

	// reset destination dir.
	if err := os.RemoveAll(out); err != nil {
		return err
	}
	if err := os.MkdirAll(typesOut, 0755); err != nil {
		return err
	}

	// generate ts-proto types.
	err = protoc.Generate(
		g.g.ctx,
		typesOut,
		m.Pkg.Path,
		includePaths,
		tsOut,
		protoc.Plugin(tsprotoPluginPath),
	)
	if err != nil {
		return err
	}

	// generate OpenAPI spec.
	oaitemp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(oaitemp)

	err = protoc.Generate(
		ctx,
		oaitemp,
		m.Pkg.Path,
		includePaths,
		openAPIOut,
	)
	if err != nil {
		return err
	}

	// generate the REST client from the OpenAPI spec.
	var (
		srcspec = filepath.Join(oaitemp, "apidocs.swagger.json")
		outREST = filepath.Join(out, "rest.ts")
	)

	if err := sta.Generate(g.g.ctx, outREST, srcspec, "-1"); err != nil { // -1 removes the route namespace.
		return err
	}

	// generate the js client wrapper.
	outclient := filepath.Join(out, "index.ts")
	f, err := os.OpenFile(outclient, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	pp := filepath.Join(appPath, g.g.protoDir)
	err = templateJSClient(pp).Execute(f, struct{ Module module.Module }{m})
	if err != nil {
		return err
	}

	// generate Vuex if enabled.
	if g.g.o.vuexStoreRootPath != "" {
		storePath := filepath.Join(storeDirPath, "index.ts")
		f, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		err = templateVuexStore(pp).Execute(f, struct{ Module module.Module }{m})
		if err != nil {
			return err
		}

		// mark vuex root dir.
		f, err = os.Create(filepath.Join(storeDirPath, vuexRootMarker))
		if err != nil {
			return err
		}
		f.Close()
	}

	// generate .js and .d.ts files for all ts files.
	return tsc.Generate(g.g.ctx, tscConfig(storeDirPath+"/**/*.ts"))
}

func (g *jsGenerator) generateVuexModuleLoader() error {
	modulePaths, err := zglob.Glob(filepath.Join(g.g.o.vuexStoreRootPath, "/**/"+vuexRootMarker))
	if err != nil {
		return err
	}

	type module struct {
		Name string
		Path string
	}

	var modules []module

	for _, path := range modulePaths {
		pathrel, err := filepath.Rel(g.g.o.vuexStoreRootPath, path)
		if err != nil {
			return err
		}
		pathrel = filepath.Dir(pathrel)
		name := strcase.ToCamel(strings.ReplaceAll(pathrel, "/", "_"))
		modules = append(modules, module{name, pathrel})
	}

	loaderPath := filepath.Join(g.g.o.vuexStoreRootPath, "index.ts")

	f, err := os.OpenFile(loaderPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = templateVuexLoader(g.g.o.vuexStoreRootPath).Execute(f, modules)
	if err != nil {
		return err
	}

	return tsc.Generate(g.g.ctx, tscConfig(loaderPath))
}

func tscConfig(include ...string) tsc.Config {
	return tsc.Config{
		Include: include,
		CompilerOptions: tsc.CompilerOptions{
			Declaration: true,
		},
	}
}
