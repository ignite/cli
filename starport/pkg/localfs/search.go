package localfs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Search searches for files in the fs with given glob pattern by ensuring that
// returned file paths are sorted.
func Search(path, pattern string) ([]string, error) {
	files := make([]string, 0)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}

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
