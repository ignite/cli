package cosmosgen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/iancoleman/strcase"
	"github.com/ignite-hq/cli/ignite/pkg/cache"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite-hq/cli/ignite/pkg/dirchange"
	swaggercombine "github.com/ignite-hq/cli/ignite/pkg/nodetime/programs/swagger-combine"
	"github.com/ignite-hq/cli/ignite/pkg/protoc"
)

var openAPIOut = []string{
	"--openapiv2_out=logtostderr=true,allow_merge=true,json_names_for_fields=false,fqn_for_openapi_name=true,simple_operation_ids=true,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:.",
}

func generateOpenAPISpec(g *generator) error {
	out := filepath.Join(g.appPath, g.o.specOut)

	var (
		specDirs []string
		conf     = swaggercombine.Config{
			Swagger: "2.0",
			Info: swaggercombine.Info{
				Title: "HTTP API Console",
			},
		}
	)

	defer func() {
		for _, dir := range specDirs {
			os.RemoveAll(dir)
		}
	}()

	specCache := cache.New[[]byte](g.cacheStorage, "generate.openapi")

	anyChanges := false

	// gen generates a spec for a module where it's source code resides at src.
	// and adds needed swaggercombine configure for it.
	gen := func(src string, m module.Module) (err error) {
		dir, err := os.MkdirTemp("", "gen-openapi-module-spec")
		if err != nil {
			return err
		}
		specPath := filepath.Join(dir, "apidocs.swagger.json")

		checksum, err := dirchange.ChecksumFromPaths(src, []string{m.Pkg.Path})
		if err != nil {
			return err
		}
		cacheKey := fmt.Sprintf("%x", checksum)
		existing, found, err := specCache.Get(cacheKey)
		if err != nil {
			return err
		}

		if found {
			if err := os.WriteFile(specPath, existing, 0644); err != nil {
				return err
			}
		} else {
			anyChanges = true
			include, err := g.resolveInclude(src)
			if err != nil {
				return err
			}

			err = protoc.Generate(
				g.ctx,
				dir,
				m.Pkg.Path,
				include,
				openAPIOut,
			)
			if err != nil {
				return err
			}

			f, err := os.ReadFile(specPath)
			if err != nil {
				return err
			}
			if err := specCache.Put(cacheKey, f); err != nil {
				return err
			}
		}

		specDirs = append(specDirs, dir)

		return conf.AddSpec(strcase.ToCamel(m.Pkg.Name), specPath)
	}

	// generate specs for each module and persist them in the file system
	// after add their path and config to swaggercombine.Config so we can combine them
	// into a single spec.

	add := func(src string, modules []module.Module) error {
		for _, m := range modules {
			m := m
			if err := gen(src, m); err != nil {
				return err
			}
		}
		return nil
	}

	// protoc openapi generator acts weird on conccurrent run, so do not use goroutines here.
	if err := add(g.appPath, g.appModules); err != nil {
		return err
	}

	for src, modules := range g.thirdModules {
		if err := add(src, modules); err != nil {
			return err
		}
	}

	if !anyChanges {
		// In case the generated output has been changed
		changed, err := dirchange.HasDirChecksumChanged(g.appPath, []string{out}, specCache, out)
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
	if err := os.MkdirAll(outDir, 0766); err != nil {
		return err
	}

	// combine specs into one and save to out.
	if err := swaggercombine.Combine(g.ctx, conf, out); err != nil {
		return err
	}

	return dirchange.SaveDirChecksum(g.appPath, []string{out}, specCache, out)
}
