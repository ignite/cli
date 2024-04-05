package cosmosbuf

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const workFilename = "buf.work.yaml"

// BufWork represents the buf.work.yaml file.
type BufWork struct {
	appPath     string   `yaml:"-"`
	filePath    string   `yaml:"-"`
	Version     string   `yaml:"version"`
	Directories []string `yaml:"directories"`
}

// ParseBufWork parse the buf.work.yaml file at app path.
func ParseBufWork(appPath string) (BufWork, error) {
	path := filepath.Join(appPath, workFilename)

	f, err := os.Open(path)
	if err != nil {
		return BufWork{}, err
	}
	defer f.Close()

	w := BufWork{appPath: appPath, filePath: path}
	return w, yaml.NewDecoder(f).Decode(&w)
}

// MissingDirectories check if the directories inside the buf work exist.
func (w BufWork) MissingDirectories() ([]string, error) {
	missingPaths := make([]string, 0)
	for _, dir := range w.Directories {
		protoDir := filepath.Join(w.appPath, dir)
		if _, err := os.Stat(protoDir); os.IsNotExist(err) {
			missingPaths = append(missingPaths, dir)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return missingPaths, nil
}

// AddProtoDir add a proto directory path from the buf work file.
func (w BufWork) AddProtoDir(newPath string) error {
	w.Directories = append(w.Directories, newPath)
	return w.save()
}

// RemoveProtoDirs remove a list a proto directory paths from the buf work file.
func (w BufWork) RemoveProtoDirs(paths ...string) error {
	for _, path := range paths {
		for i, dir := range w.Directories {
			if dir == path {
				w.Directories = append(w.Directories[:i], w.Directories[i+1:]...)
				break
			}
		}
	}
	return w.save()
}

// HasProtoDir returns true if the proto path exist into the directories slice.
func (w BufWork) HasProtoDir(path string) bool {
	for _, dirPath := range w.Directories {
		if path == dirPath {
			return true
		}
	}
	return false
}

// save saves the buf work file.
func (w BufWork) save() error {
	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(&w)
}
