package xembed

import (
	"embed"
	"io/fs"
	"path/filepath"
)

// FileList list all files into an embed.FS in a provider path.
func FileList(efs embed.FS, path string) ([]string, error) {
	return fileList(efs, path, path)
}

func fileList(efs embed.FS, path, currentDir string) ([]string, error) {
	dir, err := fs.ReadDir(efs, currentDir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, f := range dir {
		if !f.IsDir() {
			relPath, err := filepath.Rel(path, filepath.Join(currentDir, f.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, relPath)
			continue
		}

		newDir := filepath.Join(currentDir, f.Name())
		dirFiles, err := fileList(efs, path, newDir)
		if err != nil {
			return nil, err
		}
		files = append(files, dirFiles...)
	}
	return files, nil
}
