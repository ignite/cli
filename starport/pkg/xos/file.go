package xos

import (
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
)

// OpenFirst finds and opens the first found file within names.
func OpenFirst(names ...string) (file *os.File, err error) {
	for _, name := range names {
		file, err = os.Open(name)
		if err == nil {
			break
		}
	}
	return file, err
}

// DirList returns a list of dirs of matching pattern.
func DirList(pattern string) ([]string, error) {
	files, err := zglob.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// get a list of dirs wihout duplicates
	dirsNoDup := make(map[string]bool)

	for _, file := range files {
		dirsNoDup[filepath.Dir(file)] = true
	}

	var dirs []string

	for dir := range dirsNoDup {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

// PrefixPathToList adds prefix to the each path in paths and returns the
// newly created path slice.
func PrefixPathToList(paths []string, prefix string) []string {
	var p []string
	for _, path := range paths {
		p = append(p, filepath.Join(prefix, path))
	}
	return p
}
