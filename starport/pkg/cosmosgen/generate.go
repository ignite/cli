package cosmosgen

import (
	"context"
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	"github.com/tendermint/starport/starport/pkg/nodetime/sta"
	tsproto "github.com/tendermint/starport/starport/pkg/nodetime/ts-proto"
	"github.com/tendermint/starport/starport/pkg/nodetime/tsc"
	"github.com/tendermint/starport/starport/pkg/protoanalysis"
	"github.com/tendermint/starport/starport/pkg/protoc"
	"github.com/tendermint/starport/starport/pkg/protopath"
	gomodmodule "golang.org/x/mod/module"
	"golang.org/x/sync/errgroup"
)

var (
	goOuts = []string{
		"--gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
		"--grpc-gateway_out=logtostderr=true:.",
	}

	tsOut = []string{
		"--ts_proto_out=.",
	}

	openAPIOut = []string{
		"--openapiv2_out=logtostderr=true,allow_merge=true:.",
	}

	sdkImport = "github.com/cosmos/cosmos-sdk"
)

//go:embed templates/*
var templates embed.FS

// tpl holds the js client template which is for wrapping the generated protobufjs types and rest client,
// utilizing cosmjs' type registry, tx signing & broadcasting through exported, high level txClient() and queryClient() funcs.
func tpl(protoPath string) *template.Template {
	return template.Must(
		template.New("client.ts.tpl").
			Funcs(template.FuncMap{
				"camelCase": strcase.ToLowerCamel,
				"resolveFile": func(fullPath string) string {
					rel, _ := filepath.Rel(protoPath, fullPath)
					rel = strings.TrimSuffix(rel, ".proto")
					return rel
				},
			}).
			ParseFS(templates, "templates/client.ts.tpl"),
	)
}

type generateOptions struct {
	includeDirs         []string
	gomodPath           string
	jsOut               func(module.Module) string
	jsIncludeThirdParty bool
}

// TODO add WithInstall.

// Option configures code generation.
type Option func(*generateOptions)

// WithJSGeneration adds JS code generation. out hook is called for each module to
// retrieve the path that should be used to place generated js code inside for a given module.
// if includeThirdPartyModules set to true, code generation will be made for the 3rd party modules
// used by the app -including the SDK- as well.
func WithJSGeneration(includeThirdPartyModules bool, out func(module.Module) (path string)) Option {
	return func(o *generateOptions) {
		o.jsOut = out
		o.jsIncludeThirdParty = includeThirdPartyModules
	}
}

// WithGoGeneration adds Go code generation.
func WithGoGeneration(gomodPath string) Option {
	return func(o *generateOptions) {
		o.gomodPath = gomodPath
	}
}

// IncludeDirs configures the third party proto dirs that used by app's proto.
// relative to the projectPath.
func IncludeDirs(dirs []string) Option {
	return func(o *generateOptions) {
		o.includeDirs = dirs
	}
}

// generator generates code for sdk and sdk apps.
type generator struct {
	ctx      context.Context
	appPath  string
	protoDir string
	o        *generateOptions
	deps     []gomodmodule.Version
}

// Generate generates code from protoDir of an SDK app residing at appPath with given options.
// protoDir must be relative to the projectPath.
func Generate(ctx context.Context, appPath, protoDir string, options ...Option) error {
	g := &generator{
		ctx:      ctx,
		appPath:  appPath,
		protoDir: protoDir,
		o:        &generateOptions{},
	}

	for _, apply := range options {
		apply(g.o)
	}

	if err := g.setup(); err != nil {
		return err
	}

	if g.o.gomodPath != "" {
		if err := g.generateGo(); err != nil {
			return err
		}
	}

	// js generation requires Go types to be existent in the source code.
	// so it needs to run after Go code gen.
	if g.o.jsOut != nil {
		if err := g.generateJS(); err != nil {
			return err
		}
	}

	return nil
}

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

func (g *generator) generateGo() error {
	includePaths, err := g.resolveInclude(g.appPath)
	if err != nil {
		return err
	}

	// created a temporary dir to locate generated code under which later only some of them will be moved to the
	// app's source code. this also prevents having leftover files in the app's source code or its parent dir -when
	// command executed directly there- in case of an interrupt.
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	// discover proto packages in the app.
	pp := filepath.Join(g.appPath, g.protoDir)
	pkgs, err := protoanalysis.DiscoverPackages(pp)
	if err != nil {
		return err
	}

	// code generate for each module.
	for _, pkg := range pkgs {
		if err := protoc.Generate(g.ctx, tmp, pkg.Path, includePaths, goOuts); err != nil {
			return err
		}
	}

	// move generated code for the app under the relative locations in its source code.
	generatedPath := filepath.Join(tmp, g.o.gomodPath)

	_, err = os.Stat(generatedPath)
	if err == nil {
		err = copy.Copy(generatedPath, g.appPath)
		return errors.Wrap(err, "cannot copy path")
	}
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (g *generator) generateJS() error {
	tsprotoPluginPath, err := tsproto.BinaryPath()
	if err != nil {
		return err
	}

	// generate generates JS code for a module.
	generate := func(ctx context.Context, appPath string, m module.Module) error {
		var (
			out      = g.o.jsOut(m)
			typesOut = filepath.Join(out, "types")
		)

		includePaths, err := g.resolveInclude(appPath)
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
			g.ctx,
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

		if err := sta.Generate(g.ctx, outREST, srcspec, "-1"); err != nil { // -1 removes the route namespace.
			return err
		}

		// generate the js client wrapper.
		outclient := filepath.Join(out, "index.ts")
		f, err := os.OpenFile(outclient, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		pp := filepath.Join(appPath, g.protoDir)
		err = tpl(pp).Execute(f, struct{ Module module.Module }{m})
		if err != nil {
			return err
		}

		// generate .js and .d.ts files for all ts files.
		err = tsc.Generate(g.ctx, tsc.Config{
			Include: []string{out + "/**/*.ts"},
			CompilerOptions: tsc.CompilerOptions{
				Declaration: true,
			},
		})

		return err
	}

	// sourcePaths keeps a list of root paths of Go projects (source codes) that might contain
	// Cosmos SDK modules inside.
	sourcePaths := []string{
		g.appPath, // user's blockchain. may contain internal modules. it is the first place to look for.
	}

	if g.o.jsIncludeThirdParty {
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
		for _, dep := range g.deps {
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
			modules, err := g.discoverModules(sourcePath)
			if err != nil {
				return err
			}

			gg, ctx := errgroup.WithContext(g.ctx)

			// do code generation for each found module.
			for _, m := range modules {
				m := m

				gg.Go(func() error { return generate(ctx, sourcePath, m) })
			}

			return gg.Wait()
		})
	}

	return gs.Wait()
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

	modules, err := module.Discover(path)
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
