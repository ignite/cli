package xos

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	JSONFile = "json"
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
