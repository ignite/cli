package localfs

import (
	"os"
	"path/filepath"

	"io/fs"
)

// SaveTemp saves file system f to a temporary path in the local file system
// and returns that path.
func SaveTemp(f fs.FS) (path string, err error) {
	path, err = os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	return path, Save(f, path)
}

// Save saves file system f to path in the local file system.
func Save(f fs.FS, path string) error {
	return fs.WalkDir(f, ".", func(wpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		out := filepath.Join(path, wpath)

		if d.IsDir() {
			return os.MkdirAll(out, 0744)
		}

		content, err := fs.ReadFile(f, wpath)
		if err != nil {
			return err
		}

		return os.WriteFile(out, content, 0644)
	})
}
