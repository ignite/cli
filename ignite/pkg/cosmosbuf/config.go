package cosmosbuf

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const bufYamlFilename = "buf.yaml"

type (
	// BufWork represents the buf.yaml file.
	BufWork struct {
		appPath  string   `yaml:"-"`
		filePath string   `yaml:"-"`
		Version  string   `yaml:"version"`
		Modules  []Module `yaml:"modules"`
	}
	// Module represents the buf.yaml module.
	Module struct {
		Path string `yaml:"path"`
		Name string `yaml:"name"`
	}
)

// ParseBufConfig parse the buf.yaml file at app path.
func ParseBufConfig(appPath string) (BufWork, error) {
	path := filepath.Join(appPath, bufYamlFilename)

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
	for _, module := range w.Modules {
		protoDir := filepath.Join(w.appPath, module.Path)
		if _, err := os.Stat(protoDir); os.IsNotExist(err) {
			missingPaths = append(missingPaths, module.Path)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return missingPaths, nil
}

// AddProtoDir add a proto directory path from the buf work file.
func (w BufWork) AddProtoDir(newPath string) error {
	w.Modules = append(w.Modules, Module{Path: newPath})
	return w.save()
}

// RemoveProtoDirs remove a list a proto directory paths from the buf work file.
func (w BufWork) RemoveProtoDirs(paths ...string) error {
	for _, path := range paths {
		for i, module := range w.Modules {
			if module.Path == path {
				w.Modules = append(w.Modules[:i], w.Modules[i+1:]...)
				break
			}
		}
	}
	return w.save()
}

// HasProtoDir returns true if the proto path exist into the directories slice.
func (w BufWork) HasProtoDir(path string) bool {
	for _, module := range w.Modules {
		if path == module.Path {
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
