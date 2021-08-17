package localfs

import (
	"io/fs"
	"os"
)

// MkdirAllReset is same as os.MkdirAll except it deletes path before creating it.
func MkdirAllReset(path string, perm fs.FileMode) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.MkdirAll(path, perm)
}
