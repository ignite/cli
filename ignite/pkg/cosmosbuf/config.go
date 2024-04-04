package cosmosbuf

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const workFilename = "buf.work.yaml"

// BufWork represents the buf.work.yaml file.
type BufWork struct {
	path        string   `yaml:"-"`
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

	var w BufWork
	return w, yaml.NewDecoder(f).Decode(&w)
}

// AddProtoPath change the name of a proto directory path into the buf work file.
func (w BufWork) AddProtoPath(newPath string) error {
	w.Directories = append(w.Directories, newPath)
	return w.save()
}

// HasProtoPath returns true if the proto path exist into the directories slice.
func (w BufWork) HasProtoPath(path string) bool {
	for _, dirPath := range w.Directories {
		if path == dirPath {
			return true
		}
	}
	return false
}

// save saves the buf work file.
func (w BufWork) save() error {
	file, err := os.OpenFile(w.path, os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(w)
}
