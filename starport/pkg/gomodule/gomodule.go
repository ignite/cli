package gomodule

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// ParseAt finds and parses go.mod at app's path.
func ParseAt(path string) (*modfile.File, error) {
	gomod, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, err
	}
	return modfile.Parse("", gomod, nil)
}

// FilterRequire filters dependencies under require section by their paths.
func FilterRequire(dependencies []*modfile.Require, paths ...string) []*modfile.Require {
	var filtered []*modfile.Require

	for _, dep := range dependencies {
		for _, path := range paths {
			if dep.Mod.Path == path {
				filtered = append(filtered, dep)
				break
			}
		}
	}

	return filtered
}

// LocatePath locates pkg's absolute path managed by 'go mod' on the local filesystem.
func LocatePath(pkg module.Version) (path string, err error) {
	path = filepath.Join(build.Default.GOPATH, "pkg/mod", fmt.Sprintf("%s@%s", pkg.Path, pkg.Version))
	_, err = os.Stat(path)
	return
}
