package localfs

import (
	"sort"

	"github.com/mattn/go-zglob"
)

// Search searches for files in the fs with given glob pattern by ensuring that
// returned file paths are sorted.
func Search(pattern string) (paths []string, err error) {
	files, err := zglob.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}
