package cosmosgen

import (
	"github.com/ignite-hq/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite-hq/cli/ignite/pkg/giturl"
	"github.com/ignite-hq/cli/ignite/pkg/gomodulepath"
	"github.com/ignite-hq/cli/ignite/pkg/nodetime/programs/sta"
	tsproto "github.com/ignite-hq/cli/ignite/pkg/nodetime/programs/ts-proto"
	"github.com/ignite-hq/cli/ignite/pkg/protoc"

	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

var (
	tsOut = []string{
		"--ts_proto_out=.",
	}

	jsOpenAPIOut = []string{
		"--openapiv2_out=logtostderr=true,allow_merge=true,json_names_for_fields=false,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
	}
)

type tsGenerator struct {
	g *generator
}

type generatePayload struct {
	Modules []module.Module
	User    string ``
	Repo    string ``
}

func newTSGenerator(g *generator) *tsGenerator {
	return &tsGenerator{
		g: g,
	}
}

func (g *generator) generateTS() error {
	tsg := newTSGenerator(g)

	chainPath, _, err := gomodulepath.Find(g.appPath)
	if err != nil {
		return err
	}

	chainInfo, err := giturl.Parse(chainPath.RawPath)
	if err != nil {
		return err
	}

	data := generatePayload{
		Modules: g.appModules,
		User:    chainInfo.User,
		Repo:    chainInfo.Repo,
	}

	if g.o.jsIncludeThirdParty {
		for _, modules := range g.thirdModules {
			data.Modules = append(data.Modules, modules...)
		}
	}

	if err := tsg.generateModuleTemplates(); err != nil {
		return err
	}

	if err := tsg.generateVueTemplates(data); err != nil {
		return err
	}

	if err := tsg.generateRootTemplates(data); err != nil {
		return err
	}

	return nil
}

func (g *tsGenerator) generateModuleTemplates() error {
	tsprotoPluginPath, cleanup, err := tsproto.BinaryPath()
	if err != nil {
		return err
	}
	defer cleanup()

	gg := &errgroup.Group{}

	generate := func(sourcePath string, modules []module.Module) {
		for _, m := range modules {
			m := m
			gg.Go(func() error {
				var (
					out      = filepath.Join(g.g.o.tsClientRootPath, "client", m.Pkg.Name)
					typesOut = filepath.Join(out, "types")
					ctx      = g.g.ctx
					appPath  = sourcePath
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
					protoc.Plugin(tsprotoPluginPath, "--ts_proto_opt=snakeToCamel=false"),
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
			})
		}
	}

	generate(g.g.appPath, g.g.appModules)

	if g.g.o.jsIncludeThirdParty {
		for sourcePath, modules := range g.g.thirdModules {
			generate(sourcePath, modules)
		}
	}

	return gg.Wait()
}

func (g *tsGenerator) generateVueTemplates(payload generatePayload) error {
	gg := &errgroup.Group{}

	generate := func() {
		for _, m := range payload.Modules {
			m := m

			gg.Go(func() error {
				vueAPIOut := filepath.Join(g.g.o.tsClientRootPath, "vue", m.Pkg.Name)

				if err := os.MkdirAll(vueAPIOut, 0766); err != nil {
					return err
				}

				if err := templateTSClientVue.Write(vueAPIOut, "", struct {
					Module module.Module
					User   string
					Repo   string
				}{
					Module: m,
					User:   payload.User,
					Repo:   payload.Repo,
				}); err != nil {
					return err
				}

				return nil
			})
		}
	}

	generate()

	return gg.Wait()
}

func (g *tsGenerator) generateRootTemplates(payload generatePayload) error {
	tsClientOut := filepath.Join(g.g.o.tsClientRootPath, "client")
	if err := os.MkdirAll(tsClientOut, 0766); err != nil {
		return err
	}
	if err := templateTSClientRoot.Write(tsClientOut, "", payload); err != nil {
		return err
	}

	vueOut := filepath.Join(g.g.o.tsClientRootPath, "vue")
	if err := os.MkdirAll(vueOut, 0766); err != nil {
		return err
	}
	if err := templateTSClientVueRoot.Write(vueOut, "", payload); err != nil {
		return err
	}

	return nil
}
