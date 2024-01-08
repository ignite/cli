package xos

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	JSONFile  = "json"
	ProtoFile = "proto"
)

func FindFiles(directory, extension string) ([]string, error) {
	files := make([]string, 0)
	return files, filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if filepath.Ext(path) == fmt.Sprintf(".%s", extension) {
				files = append(files, path)
			}
		}
		return nil
	})
}

// FileExists check if a file from a given path exists.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
