package diff

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// computeFS computes the unified diffs between the origin and modified filesystems.
// but ignores files that match the given globs.
func computeFS(origin, modified fs.FS, ignoreGlobs ...string) ([]gotextdiff.Unified, error) {
	compiledGlobs, err := compileGlobs(ignoreGlobs)
	if err != nil {
		return nil, err
	}

	marked := make(map[string]struct{})
	unified := make([]gotextdiff.Unified, 0)
	err = fs.WalkDir(origin, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Errorf("failed to walk origin: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		if matchGlobs(compiledGlobs, path) {
			return nil
		}

		marked[path] = struct{}{}
		data, err := fs.ReadFile(origin, path)
		if err != nil {
			return errors.Errorf("failed to read file %q from origin: %w", path, err)
		}
		originFile := string(data)

		data, err = fs.ReadFile(modified, path)
		if !os.IsNotExist(err) && err != nil {
			return errors.Errorf("failed to read file %q from modified: %w", path, err)
		}
		modifiedFile := string(data)

		edits := myers.ComputeEdits(span.URIFromURI(fmt.Sprintf("file://%s", path)), originFile, modifiedFile)
		if len(edits) > 0 {
			unified = append(unified, gotextdiff.ToUnified(path, path, originFile, edits))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(modified, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Errorf("failed to walk modified: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		if _, ok := marked[path]; ok {
			return nil
		}

		if matchGlobs(compiledGlobs, path) {
			return nil
		}

		originFile := ""
		data, err := fs.ReadFile(modified, path)
		if err != nil {
			return errors.Errorf("failed to read file %q from modified: %w", path, err)
		}
		modifiedFile := string(data)

		edits := myers.ComputeEdits(span.URIFromURI(fmt.Sprintf("file://%s", path)), originFile, modifiedFile)
		if len(edits) > 0 {
			unified = append(unified, gotextdiff.ToUnified(path, path, originFile, edits))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return unified, nil
}

func compileGlobs(globs []string) ([]glob.Glob, error) {
	var compiledGlobs []glob.Glob
	for _, g := range globs {
		compiledGlob, err := glob.Compile(g, filepath.Separator)
		if err != nil {
			return nil, errors.Errorf("failed to compile glob %q: %w", g, err)
		}
		compiledGlobs = append(compiledGlobs, compiledGlob)
	}
	return compiledGlobs, nil
}

func matchGlobs(globs []glob.Glob, path string) bool {
	for _, g := range globs {
		if g.Match(path) {
			return true
		}
	}
	return false
}
