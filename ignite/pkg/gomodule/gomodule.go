package gomodule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

const pathCacheNamespace = "gomodule.path"

var (
	// ErrGoModNotFound returned when go.mod file cannot be found for an app.
	ErrGoModNotFound = errors.New("go.mod not found")

	// ErrModuleNotFound is returned when a Go module is not found.
	ErrModuleNotFound = errors.New("module not found")
)

// Version is an alias to the module version type.
type Version = module.Version

// Module contains Go module info.
type Module struct {
	// Path is the Go module path.
	Path string

	// Version is the module version.
	Version string

	// Dir is the absolute path to the Go module.
	Dir string
}

// ParseAt finds and parses go.mod at app's path.
func ParseAt(path string) (*modfile.File, error) {
	gomod, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrGoModNotFound
		}
		return nil, err
	}
	return modfile.Parse("", gomod, nil)
}

// FilterVersions filters dependencies under require section by their paths.
func FilterVersions(dependencies []Version, paths ...string) []Version {
	var filtered []Version

	for _, dep := range dependencies {
		if slices.Contains(paths, dep.Path) {
			filtered = append(filtered, dep)
		}
	}

	return filtered
}

// ResolveDependencies resolves dependencies from go.mod file.
// It replaces direct dependencies with their replacements.
func ResolveDependencies(f *modfile.File, includeIndirect bool) ([]Version, error) {
	var versions []Version

	isReplacementAdded := func(rv Version) bool {
		for _, rep := range f.Replace {
			if rv.Path == rep.Old.Path {
				versions = append(versions, rep.New)

				return true
			}
		}

		return false
	}

	for _, req := range f.Require {
		if req.Indirect && !includeIndirect {
			continue
		}
		if !isReplacementAdded(req.Mod) {
			versions = append(versions, req.Mod)
		}
	}

	return versions, nil
}

// LocatePath locates pkg's absolute path managed by 'go mod' on the local filesystem.
func LocatePath(ctx context.Context, cacheStorage cache.Storage, src string, pkg Version) (path string, err error) {
	// can be a local package.
	if pkg.Version == "" { // indicates that this is a local package.
		if filepath.IsAbs(pkg.Path) {
			return pkg.Path, nil
		}
		return filepath.Join(src, pkg.Path), nil
	}

	pathCache := cache.New[string](cacheStorage, pathCacheNamespace)
	cacheKey := cache.Key(pkg.Path, pkg.Version)
	path, err = pathCache.Get(cacheKey)
	if err != nil && !errors.Is(err, cache.ErrorNotFound) {
		return "", err
	}
	if !errors.Is(err, cache.ErrorNotFound) {
		return path, nil
	}

	// otherwise, it is hosted.
	m, err := FindModule(ctx, src, pkg.String())
	if err != nil {
		return "", err
	}

	if err = pathCache.Put(cacheKey, m.Dir); err != nil {
		return "", err
	}
	return m.Dir, nil
}

// SplitPath splits a Go import path into an URI path and version.
// Version is an empty string when the path doesn't contain a version suffix.
// Versioned paths use the "path@version" format.
func SplitPath(path string) (string, string) {
	if len(path) == 0 || path[0] == '@' {
		return "", ""
	}

	parts := strings.SplitN(path, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

// JoinPath joins a Go import path URI to a version.
// The result path have the "path@version" format.
func JoinPath(path, version string) string {
	if path == "" {
		return ""
	}

	if version == "" {
		return path
	}

	return fmt.Sprintf("%s@%s", path, version)
}

// FindModule returns the Go module info for an import path.
// The module is searched within the dependencies of the module defined in root dir.
// If a local module path is passed, it returns the local module info.
func FindModule(ctx context.Context, rootDir, path string) (Module, error) {
	// can be a local module.
	if filepath.IsAbs(path) || strings.HasPrefix(path, ".") { // indicates that this is a local module.
		return Module{
			Path:    path,
			Version: "",
			Dir:     path,
		}, nil
	}

	var stdout bytes.Buffer
	err := gocmd.ModDownload(ctx, rootDir, true, exec.StepOption(step.Stdout(&stdout)))
	if err != nil {
		return Module{}, err
	}

	dec := json.NewDecoder(&stdout)
	p, version := SplitPath(path)

	for dec.More() {
		var m Module
		if err := dec.Decode(&m); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Module{}, err
		}

		if m.Path == p && (version == "" || version == m.Version) {
			return m, nil
		}
	}

	return Module{}, errors.Errorf("%w: %s", ErrModuleNotFound, path)
}
