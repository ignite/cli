package xos

import (
	"os"
	"path/filepath"
)

func RemoveAllUnderHome(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(home, path))
}
