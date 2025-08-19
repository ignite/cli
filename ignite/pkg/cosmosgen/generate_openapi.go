package cosmosgen

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
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

func (g *generator) generateOpenAPISpec(ctx context.Context) error {
	var (
		specDirs = make([]string, 0)
		conf     = swaggercombine.New("HTTP API Console", g.goModPath)
	)
	defer func() {
		for _, dir := range specDirs {
			_ = os.RemoveAll(dir)
		}
	}()

	specCache := cache.New[[]byte](g.cacheStorage, specCacheNamespace)

	var hasAnySpecChanged bool

	// gen generates a spec for a module where it's source code resides at src.
	// and adds needed swaggercombine configure for it.
	gen := func(appPath, protoDir, name string) error {
		name = strcase.ToCamel(name)
		protoPath := filepath.Join(appPath, protoDir)

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

		dir, err := os.MkdirTemp("", "gen-openapi-module-spec")
		if err != nil {
			return err
		}

		specDirs = append(specDirs, dir)

		var noChecksum bool
		checksum, err := dirchange.ChecksumFromPaths(appPath, protoDir)
		if errors.Is(err, dirchange.ErrNoFile) {
			noChecksum = true
		} else if err != nil {
			return err
		}

		cacheKey := fmt.Sprintf("%x", checksum)
		if !noChecksum {
			existingSpec, err := specCache.Get(cacheKey)
			if err != nil && !errors.Is(err, cache.ErrorNotFound) {
				return err
			}

			if !errors.Is(err, cache.ErrorNotFound) {
				specPath := filepath.Join(dir, specFilename)
				if err := os.WriteFile(specPath, existingSpec, 0o600); err != nil {
					return err
				}
				return conf.AddSpec(name, specPath, true)
			}
		}

		hasAnySpecChanged = true
		if err = g.buf.Generate(
			ctx,
			protoPath,
			dir,
			g.openAPITemplate(),
			cosmosbuf.ExcludeFiles(
				"*/osmosis-labs/fee-abstraction/*",
				"*/module.proto",
				"*/testutil/*",
				"*/testdata/*",
				"*/cosmos/orm/*",
				"*/cosmos/reflection/*",
				"*/cosmos/app/v1alpha1/*",
				"*/cosmos/tx/config/v1/config.proto",
				"*/cosmos/msg/textual/v1/textual.proto",
				"*/cosmos/vesting/v1beta1/vesting.proto",
			),
			cosmosbuf.FileByFile(),
		); err != nil {
			return errors.Wrapf(err, "failed to generate openapi spec %s, probably you need to exclude some proto files", protoPath)
		}

		specs, err := xos.FindFiles(dir, xos.WithExtension(xos.JSONFile))
		if err != nil {
			return err
		}

		for _, spec := range specs {
			f, err := os.ReadFile(spec)
			if err != nil {
				return err
			}

			// if no checksum, the cacheKey is wrong, so we do not save it
			if !noChecksum {
				if err := specCache.Put(cacheKey, f); err != nil {
					return err
				}
			}

			if err := conf.AddSpec(name, spec, true); err != nil {
				return err
			}
		}

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
		for _, m := range modules {
			path := extractRootModulePath(m.Pkg.Path)

			if _, ok := doneMods[path]; ok {
				continue
			}
			doneMods[path] = struct{}{}

			if err := gen(path, "proto", m.Name); err != nil {
				return err
			}
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
