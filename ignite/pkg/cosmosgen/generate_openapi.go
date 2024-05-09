package cosmosgen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	swaggercombine "github.com/ignite/cli/v29/ignite/pkg/swagger-combine"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
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
		conf     = swaggercombine.New("HTTP API Console", g.gomodPath)
	)
	defer func() {
		for _, dir := range specDirs {
			os.RemoveAll(dir)
		}
	}()

	specCache := cache.New[[]byte](g.cacheStorage, specCacheNamespace)

	var hasAnySpecChanged bool

	// gen generates a spec for a module where it's source code resides at src.
	// and adds needed swaggercombine configure for it.
	gen := func(appPath, protoPath string) (err error) {
		name := extractName(appPath)
		dir, err := os.MkdirTemp("", "gen-openapi-module-spec")
		if err != nil {
			return err
		}

		checksum, err := dirchange.ChecksumFromPaths(appPath, protoPath)
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
		if err = g.buf.Generate(ctx, filepath.Join(appPath, g.protoDir), dir, g.openAPITemplate(), "module.proto"); err != nil {
			return err
		}

		specs, err := xos.FindFiles(dir, xos.JSONFile)
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
	if err := gen(g.appPath, g.protoDir); err != nil {
		return err
	}

	for src := range g.thirdModules {
		if err := gen(src, ""); err != nil {
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

	// combine specs into one and save to out.
	if err := conf.Combine(out); err != nil {
		return err
	}

	return dirchange.SaveDirChecksum(specCache, out, g.appPath, out)
}

// generateModuleOpenAPISpec generates a spec for a module where it's source code resides at src.
// and adds needed swaggercombine configure for it.
func (g *generator) generateModuleOpenAPISpec(ctx context.Context, m module.Module, out string) error {
	var (
		specDirs []string
		title    = "HTTP API Console " + m.Pkg.Name
		conf     = swaggercombine.New(title, g.gomodPath)
	)
	defer func() {
		for _, dir := range specDirs {
			os.RemoveAll(dir)
		}
	}()

	// generate specs for each module and persist them in the file system
	// after add their path and config to swaggercombine.Config so we can combine them
	// into a single spec.
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
		if err := conf.AddSpec(strcase.ToCamel(m.Pkg.Name), spec, false); err != nil {
			return err
		}
	}
	specDirs = append(specDirs, dir)

	// combine specs into one and save to out.
	return conf.Combine(out)
}

// extractName takes a full path and returns the name.
func extractName(path string) string {
	// Extract the last part of the path
	lastPart := filepath.Base(path)
	// If there is a version suffix (e.g., @v0.50), remove it
	name := strings.Split(lastPart, "@")[0]
	return strcase.ToCamel(name)
}
