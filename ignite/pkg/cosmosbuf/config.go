package cosmosbuf

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// BufWork represents the buf.work.yaml file.
type BufWork struct {
	Version     string   `yaml:"version"`
	Directories []string `yaml:"directories"`
}

// ParseBufWork parse the buf.work.yaml file.
func ParseBufWork(path string) (BufWork, error) {
	f, err := os.Open(path)
	if err != nil {
		return BufWork{}, err
	}
	defer f.Close()

	var w BufWork
	return w, yaml.NewDecoder(f).Decode(&w)
}

// ChangeProtoPath change the name of a proto directory path into the buf work file.
func (w *BufWork) ChangeProtoPath(oldPath, newPath string) error {
	for i, path := range w.Directories {
		if path == oldPath {
			w.Directories[i] = newPath
			return nil
		}
	}
	return errors.Errorf("proto path %s not found", oldPath)
}

// HasProtoPath returns true if the proto path exist into the directories slice.
func (w *BufWork) HasProtoPath(path string) bool {
	for _, dirPath := range w.Directories {
		if path == dirPath {
			return true
		}
	}
	return false
}

// Save saves the buf work file.
func (w *BufWork) Save(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(w)
}
