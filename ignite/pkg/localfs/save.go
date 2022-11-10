package localfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// SaveTemp saves file system f to a temporary path in the local file system
// and returns that path.
func SaveTemp(f fs.FS) (path string, cleanup func(), err error) {
	path, err = os.MkdirTemp("", "")
	if err != nil {
		return
	}

	cleanup = func() { os.RemoveAll(path) }

	defer func() {
		if err != nil {
			cleanup()
		}
	}()

	err = Save(f, path)

	return
}

// SaveBytesTemp saves data bytes to a temporary file location at path.
func SaveBytesTemp(data []byte, prefix string, perm os.FileMode) (path string, cleanup func(), err error) {
	f, err := os.CreateTemp("", prefix)
	if err != nil {
		return
	}
	defer f.Close()

	path = f.Name()
	cleanup = func() { os.Remove(path) }

	defer func() {
		if err != nil {
			cleanup()
		}
	}()

	if _, err = f.Write(data); err != nil {
		return
	}

	err = os.Chmod(path, perm)

	return
}

// Save saves file system f to path in the local file system.
func Save(f fs.FS, path string) error {
	return fs.WalkDir(f, ".", func(wpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		out := filepath.Join(path, wpath)

		if d.IsDir() {
			return os.MkdirAll(out, 0o744)
		}

		content, err := fs.ReadFile(f, wpath)
		if err != nil {
			return err
		}

		return os.WriteFile(out, content, 0o644)
	})
}
