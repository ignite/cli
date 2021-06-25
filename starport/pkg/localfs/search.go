package localfs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Search searches for files in the fs with given glob pattern by ensuring that
// returned file paths are sorted.
func Search(path, pattern string) (paths []string, err error) {
	files, err := glob(path, pattern)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// glob returns the names of all files matching pattern or nil
// if there is no matching file in all sub-paths recursively. The
// syntax of patterns is the same as in Match. The pattern may
// describe hierarchical names such as /usr/*/bin/ed (assuming
// the Separator is '/').
func glob(path, pattern string) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		// skip hidden folders
		if f.IsDir() && strings.HasPrefix(base, ".") {
			return filepath.SkipDir
		}
		// avoid check directories
		if f.IsDir() {
			return nil
		}
		// check if the file name and pattern matches
		matched, err := filepath.Match(pattern, base)
		if err != nil {
			return err
		}
		if matched {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}
