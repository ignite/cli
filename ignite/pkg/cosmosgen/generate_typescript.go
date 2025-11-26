package cosmosgen

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v29/ignite/internal/buf"
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
)

var (
	bufTokenEnvName = "BUF_TOKEN"

	dirchangeCacheNamespace = "generate.typescript.dirchange"

	protocGenTSProtoBin = "protoc-gen-ts_proto"

	msgBufAuth = "Note: Buf is limits remote plugin requests from unauthenticated users on 'buf.build'. Intensively using this function will get you rate limited. Authenticate with 'buf registry login' to avoid this (https://buf.build/docs/generate/auth-required)."
)

const localTSProtoTmpl = `version: v1
plugins:
  - plugin: ts_proto
    out: .
    opt:
      - logtostderr=true
      - allow_merge=true
      - json_names_for_fields=false
      - ts_proto_opt=snakeToCamel=true
      - ts_proto_opt=esModuleInterop=true
      - ts_proto_out=.
`

type tsGenerator struct {
	g              *generator
	tsTemplateFile string
	isLocalProto   bool

	// hasLocalBufToken indicates whether the user had already a local Buf token.
	hasLocalBufToken bool
}

type generatePayload struct {
	Modules   []module.Module
	PackageNS string
}

func newTSGenerator(g *generator) *tsGenerator {
	tsg := &tsGenerator{g: g}
	if _, err := exec.LookPath(protocGenTSProtoBin); err == nil {
		tsg.isLocalProto = true
	}

	if !tsg.isLocalProto {
		if os.Getenv(bufTokenEnvName) == "" {
			token, err := buf.FetchToken()
			if err != nil {
				log.Printf("No '%s' binary found in PATH, using remote buf plugin for Typescript generation. %s\n", protocGenTSProtoBin, msgBufAuth)
			} else {
				os.Setenv(bufTokenEnvName, token)
			}
		} else {
			tsg.hasLocalBufToken = true
		}
	}

	return tsg
}

func (g *tsGenerator) tsTemplate() (string, error) {
	if !g.isLocalProto {
		return g.g.tsTemplate(), nil
	}
	if g.tsTemplateFile != "" {
		return g.tsTemplateFile, nil
	}
	f, err := os.CreateTemp("", "buf-gen-ts-*.yaml")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(localTSProtoTmpl); err != nil {
		return "", err
	}
	g.tsTemplateFile = f.Name()
	return g.tsTemplateFile, nil
}

func (g *tsGenerator) cleanup() {
	if g.tsTemplateFile != "" {
		os.Remove(g.tsTemplateFile)
	}

	// unset ignite buf token from env
	if !g.hasLocalBufToken {
		os.Unsetenv(bufTokenEnvName)
	}
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
	defer tsg.cleanup()
	if err := tsg.generateModuleTemplates(ctx); err != nil {
		return err
	}

	// add third party modules to for the root template.
	for _, modules := range g.thirdModules {
		data.Modules = append(data.Modules, modules...)
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

	gg := &errgroup.Group{}
	for _, m := range g.g.appModules {
		gg.Go(func() error {
			return add(g.g.appPath, m)
		})
	}

	// Always generate third party modules; This is required because not generating them might
	// lead to issues with the module registration in the root template. The root template must
	// always be generated with 3rd party modules which means that if a new 3rd party module
	// is available and not generated it would lead to the registration of a new not generated
	// 3rd party module.
	for sourcePath, modules := range g.g.thirdModules {
		for _, m := range modules {
			gg.Go(func() error {
				return add(sourcePath, m)
			})
		}
	}

	return gg.Wait()
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

	if module.IsCosmosSDKPackage(appPath) {
		protoPath = filepath.Join(g.g.sdkDir, "proto")
	}

	// check if directory exists
	if _, err := os.Stat(protoPath); os.IsNotExist(err) {
		var err error
		protoPath, err = findInnerProtoFolder(appPath)
		if err != nil {
			// if proto directory does not exist, we just skip it
			log.Print(err.Error())
			return nil
		}
	}

	tsTemplate, err := g.tsTemplate()
	if err != nil {
		return err
	}

	// code generate for each module.
	if err := g.g.buf.Generate(
		ctx,
		protoPath,
		typesOut,
		tsTemplate,
		cosmosbuf.IncludeWKT(),
		cosmosbuf.WithModuleName(m.Pkg.Name),
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
		module.Module
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
