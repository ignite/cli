package protopath

import (
	"context"
	"errors"
	"path/filepath"

	"golang.org/x/mod/module"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/gomodule"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

var (
	globalInclude = xfilepath.List(
		// this one should be already known by naked protoc execution, but adding it anyway to making sure.
		xfilepath.JoinFromHome(xfilepath.Path("local/include")),
		// this one is the suggested installation path for placing default proto by
		// https://grpc.io/docs/protoc-installation/.
		xfilepath.JoinFromHome(xfilepath.Path(".local/include")),
	)
)

// Module represents a go module that hosts dependency proto paths.
type Module struct {
	importPath string
	include    []string
}

// NewModule cretes a new go module representation to look for protoPaths.
func NewModule(importPath string, protoPaths ...string) Module {
	return Module{
		importPath: importPath,
		include:    protoPaths,
	}
}

// ResolveDependencyPaths resolves dependency proto paths (include/-I) for modules over given r inside go modules.
// r should be the list of required packages of the target go app. it is used to resolve exact versions
// of the go modules that used by the target app.
// global dependencies are also included to paths.
func ResolveDependencyPaths(ctx context.Context, cacheStorage cache.Storage, src string, versions []module.Version, modules ...Module) (paths []string, err error) {
	globalInclude, err := globalInclude()
	if err != nil {
		return nil, err
	}

	paths = append(paths, globalInclude...)

	var importPaths []string

	for _, module := range modules {
		importPaths = append(importPaths, module.importPath)
	}

	vs := gomodule.FilterVersions(versions, importPaths...)

	if len(vs) != len(modules) {
		return nil, errors.New("go.mod has missing proto modules")
	}

	for i, v := range vs {
		path, err := gomodule.LocatePath(ctx, cacheStorage, src, v)
		if err != nil {
			return nil, err
		}

		module := modules[i]
		for _, relpath := range module.include {
			paths = append(paths, filepath.Join(path, relpath))
		}
	}

	return paths, nil
}
