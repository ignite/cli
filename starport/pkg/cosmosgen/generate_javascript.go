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
	"github.com/tendermint/starport/starport/pkg/giturl"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/nodetime/sta"
	tsproto "github.com/tendermint/starport/starport/pkg/nodetime/ts-proto"
	"github.com/tendermint/starport/starport/pkg/nodetime/tsc"
	"github.com/tendermint/starport/starport/pkg/protoc"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"golang.org/x/sync/errgroup"
)

var (
	tsOut = []string{
		"--ts_proto_out=.",
	}

	jsOpenAPIOut = []string{
		"--openapiv2_out=logtostderr=true,allow_merge=true,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
	}

	vuexRootMarker = "vuex-root"
)

type jsGenerator struct {
	g *generator
}

func newJSGenerator(g *generator) *jsGenerator {
	return &jsGenerator{
		g: g,
	}
}

func (g *generator) generateJS() error {
	jsg := newJSGenerator(g)

	if err := jsg.generateModules(); err != nil {
		return err
	}

	if g.o.vuexStoreRootPath != "" {
		if err := jsg.generateVuexModuleLoader(); err != nil {
			return err
		}
	}

	return nil
}

func (g *jsGenerator) generateModules() error {
	tsprotoPluginPath, cleanup, err := tsproto.BinaryPath()
	if err != nil {
		return err
	}
	defer cleanup()

	gg := &errgroup.Group{}

	add := func(sourcePath string, modules []module.Module) {
		for _, m := range modules {
			m := m
			gg.Go(func() error { return g.generateModule(g.g.ctx, tsprotoPluginPath, sourcePath, m) })
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
	oaitemp, err := ioutil.TempDir("", "gen-js-openapi-module-spec")
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
	pp := filepath.Join(appPath, g.g.protoDir)
	if err := templateJSClient.Write(out, pp, struct{ Module module.Module }{m}); err != nil {
		return err
	}

	// generate Vuex if enabled.
	if g.g.o.vuexStoreRootPath != "" {
		err = templateVuexStore.Write(storeDirPath, pp, struct{ Module module.Module }{m})
		if err != nil {
			return err
		}
	}
	// generate .js and .d.ts files for all ts files.
	return tsc.Generate(g.g.ctx, tscConfig(g.g.o.typescript, storeDirPath+"/**/*.ts"))
}

func (g *jsGenerator) generateVuexModuleLoader() error {
	modulePaths, err := zglob.Glob(filepath.Join(g.g.o.vuexStoreRootPath, "/**/"+vuexRootMarker))
	if err != nil {
		return err
	}

	chainPath, err := gomodulepath.ParseAt(g.g.appPath)
	if err != nil {
		return err
	}

	chainURL, err := giturl.Parse(chainPath.RawPath)
	if err != nil {
		return err
	}

	type module struct {
		Name     string
		Path     string
		FullName string
		FullPath string
	}

	data := struct {
		Modules []module
		User    string
		Repo    string
	}{
		User: chainURL.User,
		Repo: chainURL.Repo,
	}

	for _, path := range modulePaths {
		pathrel, err := filepath.Rel(g.g.o.vuexStoreRootPath, path)
		if err != nil {
			return err
		}

		var (
			fullPath = filepath.Dir(pathrel)
			fullName = xstrings.FormatUsername(strcase.ToCamel(strings.ReplaceAll(fullPath, "/", "_")))
			path     = filepath.Base(fullPath)
			name     = strcase.ToCamel(path)
		)
		data.Modules = append(data.Modules, module{
			Name:     name,
			Path:     path,
			FullName: fullName,
			FullPath: fullPath,
		})
	}

	loaderPath := filepath.Join(g.g.o.vuexStoreRootPath, "index.ts")

	if err := templateVuexRoot.Write(g.g.o.vuexStoreRootPath, "", data); err != nil {
		return err
	}

	return tsc.Generate(g.g.ctx, tscConfig(g.g.o.typescript, loaderPath))
}

func tscConfig(noEmit bool, include ...string) tsc.Config {
	return tsc.Config{
		Include: include,
		CompilerOptions: tsc.CompilerOptions{
			Declaration: true,
			NoEmit:      noEmit,
		},
	}
}
