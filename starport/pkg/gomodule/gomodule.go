package gomodule

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// ErrGoModNotFound returned when go.mod file cannot be found for an app.
var ErrGoModNotFound = errors.New("go.mod not found")

// ParseAt finds and parses go.mod at app's path.
func ParseAt(path string) (*modfile.File, error) {
	gomod, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrGoModNotFound
		}
		return nil, err
	}
	return modfile.Parse("", gomod, nil)
}

// FilterVersions filters dependencies under require section by their paths.
func FilterVersions(dependencies []module.Version, paths ...string) []module.Version {
	var filtered []module.Version

	for _, dep := range dependencies {
		for _, path := range paths {
			if dep.Path == path {
				filtered = append(filtered, dep)
				break
			}
		}
	}

	return filtered
}

func ResolveDependencies(f *modfile.File) ([]module.Version, error) {
	var versions []module.Version

	isReplacementAdded := func(rv module.Version) bool {
		for _, rep := range f.Replace {
			if rv.Path == rep.Old.Path {
				versions = append(versions, rep.New)

				return true
			}
		}

		return false
	}

	for _, req := range f.Require {
		if !isReplacementAdded(req.Mod) {
			versions = append(versions, req.Mod)
		}
	}

	return versions, nil
}

// LocatePath locates pkg's absolute path managed by 'go mod' on the local filesystem.
func LocatePath(pkg module.Version) (path string, err error) {
	path = filepath.Join(build.Default.GOPATH, "pkg/mod", fmt.Sprintf("%s@%s", pkg.Path, pkg.Version))
	_, err = os.Stat(path)
	return
}
