package cosmosgen

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/giturl"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/localfs"
	"github.com/tendermint/starport/starport/pkg/nodetime/programs/sta"
	tsproto "github.com/tendermint/starport/starport/pkg/nodetime/programs/ts-proto"
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
)

type tsGenerator struct {
	g *generator
}

func newTSGenerator(g *generator) *tsGenerator {
	return &tsGenerator{
		g: g,
	}
}

func (g *generator) generateTS() error {
	tsg := newTSGenerator(g)

	if err := tsg.generateModules(); err != nil {
		return err
	}

	if err := tsg.generatePiniaStores(); err != nil {
		return err
	}

	if err := tsg.generateRootClasses(); err != nil {
		return err
	}

	return nil
}

func (g *tsGenerator) generatePiniaStores() error {
	gg := &errgroup.Group{}

	add := func(modules []module.Module) {
		for _, m := range modules {
			m := m
			gg.Go(func() error {
				path := filepath.Join(g.g.o.tsClientRootPath, "pinia", m.Pkg.Name)

				if err := os.MkdirAll(path, 0766); err != nil {
					return err
				}
				if err := templateTSClientPinia.Write(path, "", struct{ Module module.Module }{m}); err != nil {
					return err
				}

				return nil
			})
		}
	}

	add(g.g.appModules)

	if g.g.o.jsIncludeThirdParty {
		for _, modules := range g.g.thirdModules {
			add(modules)
		}
	}

	return gg.Wait()
}

func (g *tsGenerator) generateModules() error {
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
func (g *tsGenerator) generateModule(ctx context.Context, tsprotoPluginPath, appPath string, m module.Module) error {
	var (
		out      = filepath.Join(g.g.o.tsClientRootPath, "client", m.Pkg.Name)
		typesOut = filepath.Join(out, "types")
	)

	includePaths, err := g.g.resolveInclude(appPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(typesOut, 0766); err != nil {
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

	pp := filepath.Join(appPath, g.g.protoDir)
	if err := templateTSClientModule.Write(out, pp, struct{ Module module.Module }{m}); err != nil {
		return err
	}

	return nil
}

func (g *tsGenerator) generateRootClasses() error {
	modulePaths, err := localfs.Search(g.g.o.tsClientRootPath, "module.ts")
	if err != nil {
		return err
	}

	chainPath, _, err := gomodulepath.Find(g.g.appPath)
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
		pathrel, err := filepath.Rel(g.g.o.tsClientRootPath, path)
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

	tsClientOut := filepath.Join(g.g.o.tsClientRootPath, "client")
	if err := os.MkdirAll(tsClientOut, 0766); err != nil {
		return err
	}
	if err := templateTSClientRoot.Write(tsClientOut, "", data); err != nil {
		return err
	}

	piniaOut := filepath.Join(g.g.o.tsClientRootPath, "pinia")
	if err := os.MkdirAll(piniaOut, 0766); err != nil {
		return err
	}
	if err := templateTSClientPiniaRoot.Write(piniaOut, "", data); err != nil {
		return err
	}

	return nil
}
