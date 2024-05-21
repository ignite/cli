package cosmosgen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
<<<<<<< HEAD
	"sort"
=======
	"strings"
>>>>>>> 0b412628 (feat: improve buf rate limit (#4133))

	"github.com/blang/semver/v4"
	"github.com/iancoleman/strcase"

<<<<<<< HEAD
	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v28/ignite/pkg/dirchange"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/nodetime"
	swaggercombine "github.com/ignite/cli/v28/ignite/pkg/nodetime/programs/swagger-combine"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
=======
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	swaggercombine "github.com/ignite/cli/v29/ignite/pkg/swagger-combine"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
>>>>>>> 0b412628 (feat: improve buf rate limit (#4133))
)

const (
	specCacheNamespace = "generate.openapi.spec"
	specFilename       = "swagger.config.json"
)

func (g *generator) openAPITemplate() string {
	return filepath.Join(g.appPath, g.protoDir, "buf.gen.swagger.yaml")
}

func (g *generator) openAPITemplateForSTA() string {
	return filepath.Join(g.appPath, g.protoDir, "buf.gen.sta.yaml")
}

func (g *generator) generateOpenAPISpec(ctx context.Context) error {
	var (
		specDirs []string
<<<<<<< HEAD
		conf     = swaggercombine.Config{
			Swagger: "2.0",
			Info: swaggercombine.Info{
				Title: "HTTP API Console",
			},
		}
=======
		conf     = swaggercombine.New("HTTP API Console", g.goModPath)
>>>>>>> 0b412628 (feat: improve buf rate limit (#4133))
	)
	command, cleanup, err := nodetime.Command(nodetime.CommandSwaggerCombine)
	if err != nil {
		return err
	}
	defer cleanup()
	defer func() {
		for _, dir := range specDirs {
			os.RemoveAll(dir)
		}
	}()

	specCache := cache.New[[]byte](g.cacheStorage, specCacheNamespace)

	var hasAnySpecChanged bool

	// gen generates a spec for a module where it's source code resides at src.
	// and adds needed swaggercombine configure for it.
	gen := func(appPath, protoDir, name string) error {
		name = strcase.ToCamel(name)
		protoPath := filepath.Join(appPath, protoDir)

		dir, err := os.MkdirTemp("", "gen-openapi-module-spec")
		if err != nil {
			return err
		}

		checksum, err := dirchange.ChecksumFromPaths(appPath, protoDir)
		if err != nil {
			return err
		}
		cacheKey := fmt.Sprintf("%x", checksum)
		existingSpec, err := specCache.Get(cacheKey)
		if err != nil && !errors.Is(err, cache.ErrorNotFound) {
			return err
		}

		if !errors.Is(err, cache.ErrorNotFound) {
			specPath := filepath.Join(dir, specFilename)
			if err := os.WriteFile(specPath, existingSpec, 0o644); err != nil {
				return err
			}
			return conf.AddSpec(name, specPath, true)
		}

		hasAnySpecChanged = true
		if err = g.buf.Generate(
			ctx,
			protoPath,
			dir,
			g.openAPITemplate(),
			cosmosbuf.ExcludeFiles(
				"*/module.proto",
				"*/testutil/*",
				"*/testdata/*",
				"*/cosmos/orm/*",
				"*/cosmos/reflection/*",
				"*/cosmos/app/v1alpha1/*",
				"*/cosmos/tx/config/v1/config.proto",
				"*/cosmos/msg/textual/v1/textual.proto",
			),
			cosmosbuf.FileByFile(),
		); err != nil {
			return errors.Wrapf(err, "failed to generate openapi spec %s, probally you need to exclude some proto files", protoPath)
		}

		specs, err := xos.FindFilesExtension(dir, xos.JSONFile)
		if err != nil {
			return err
		}

		for _, spec := range specs {
			f, err := os.ReadFile(spec)
			if err != nil {
				return err
			}
			if err := specCache.Put(cacheKey, f); err != nil {
				return err
			}
			if err := conf.AddSpec(name, spec, true); err != nil {
				return err
			}
		}
		specDirs = append(specDirs, dir)

		return nil
	}

	// generate specs for each module and persist them in the file system
	// after add their path and config to swaggercombine.Config so we can combine them
	// into a single spec.

	// protoc openapi generator acts weird on concurrent run, so do not use goroutines here.
	if err := gen(g.appPath, g.protoDir, g.goModPath); err != nil {
		return err
	}

	doneMods := make(map[string]struct{})
	for _, modules := range g.thirdModules {
		if len(modules) == 0 {
			continue
		}
		var (
			m    = modules[0]
			path = extractRootModulePath(m.Pkg.Path)
		)

		if _, ok := doneMods[path]; ok {
			continue
		}
		doneMods[path] = struct{}{}

		if err := gen(path, "", m.Name); err != nil {
			return err
		}
	}

	out := g.opts.specOut

	if !hasAnySpecChanged {
		// In case the generated output has been changed
		changed, err := dirchange.HasDirChecksumChanged(specCache, out, g.appPath, out)
		if err != nil {
			return err
		}

		if !changed {
			return nil
		}
	}

	sort.Slice(conf.APIs, func(a, b int) bool { return conf.APIs[a].ID < conf.APIs[b].ID })

	// ensure out dir exists.
	outDir := filepath.Dir(out)
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}

	// combine specs into one and save to out.
	if err := swaggercombine.Combine(ctx, conf, command, out); err != nil {
		return err
	}

	return dirchange.SaveDirChecksum(specCache, out, g.appPath, out)
}

func (g *generator) generateModuleOpenAPISpec(ctx context.Context, m module.Module, out string) error {
	var (
		specDirs []string
<<<<<<< HEAD
		conf     = swaggercombine.Config{
			Swagger: "2.0",
			Info: swaggercombine.Info{
				Title: "HTTP API Console " + m.Pkg.Name,
			},
		}
=======
		title    = "HTTP API Console " + m.Pkg.Name
		conf     = swaggercombine.New(title, g.goModPath)
>>>>>>> 0b412628 (feat: improve buf rate limit (#4133))
	)
	command, cleanup, err := nodetime.Command(nodetime.CommandSwaggerCombine)
	if err != nil {
		return err
	}
	defer cleanup()

	defer func() {
		for _, dir := range specDirs {
			os.RemoveAll(dir)
		}
	}()

	// gen generates a spec for a module where it's source code resides at src.
	// and adds needed swaggercombine configure for it.
	gen := func(m module.Module) (err error) {
		dir, err := os.MkdirTemp("", "gen-openapi-module-spec")
		if err != nil {
			return err
		}

		err = g.buf.Generate(ctx, m.Pkg.Path, dir, g.openAPITemplateForSTA(), "module.proto")
		if err != nil {
			return err
		}

		specs, err := xos.FindFiles(dir, xos.JSONFile)
		if err != nil {
			return err
		}

		for _, spec := range specs {
			if err != nil {
				return err
			}
			if err := conf.AddSpec(strcase.ToCamel(m.Pkg.Name), spec, false); err != nil {
				return err
			}
		}
		specDirs = append(specDirs, dir)

		return nil
	}

	// generate specs for each module and persist them in the file system
	// after add their path and config to swaggercombine.Config so we can combine them
	// into a single spec.

<<<<<<< HEAD
	add := func(modules []module.Module) error {
		for _, m := range modules {
			if err := gen(m); err != nil {
				return err
			}
=======
	err = g.buf.Generate(ctx, m.Pkg.Path, dir, g.openAPITemplateForSTA(), cosmosbuf.ExcludeFiles("*/module.proto"))
	if err != nil {
		return err
	}

	specs, err := xos.FindFilesExtension(dir, xos.JSONFile)
	if err != nil {
		return err
	}

	for _, spec := range specs {
		if err := conf.AddSpec(strcase.ToCamel(m.Pkg.Name), spec, false); err != nil {
			return err
>>>>>>> 0b412628 (feat: improve buf rate limit (#4133))
		}
		return nil
	}

	// protoc openapi generator acts weird on concurrent run, so do not use goroutines here.
	if err := add([]module.Module{m}); err != nil {
		return err
	}

	sort.Slice(conf.APIs, func(a, b int) bool { return conf.APIs[a].ID < conf.APIs[b].ID })

	// ensure out dir exists.
	outDir := filepath.Dir(out)
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}
	// combine specs into one and save to out.
	return swaggercombine.Combine(ctx, conf, command, out)
}

func extractRootModulePath(fullPath string) string {
	var (
		segments   = strings.Split(fullPath, "/")
		modulePath = "/"
	)

	for _, segment := range segments {
		modulePath = filepath.Join(modulePath, segment)
		segmentName := strings.Split(segment, "@")
		if len(segmentName) > 1 {
			if _, err := semver.ParseTolerant(segmentName[1]); err == nil {
				return modulePath
			}
		}
	}
	return fullPath
}
