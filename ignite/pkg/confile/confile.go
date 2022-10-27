// Package confile is helper to load and overwrite configuration files.
package confile

import (
	"os"
	"path/filepath"
)

// ConfigFile represents a configuration file.
type ConfigFile struct {
	creator EncodingCreator
	path    string
}

// New starts a new ConfigFile by using creator as underlying EncodingCreator to encode and
// decode config file that presents or will present on path.
func New(creator EncodingCreator, path string) *ConfigFile {
	return &ConfigFile{
		creator: creator,
		path:    path,
	}
}

// Load loads content of config file into v if file exist on path.
// otherwise nothing loaded into v and no error is returned.
func (c *ConfigFile) Load(v interface{}) error {
	file, err := os.Open(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	return c.creator.Create(file).Decode(v)
}

// Save saves v into config file by overwriting the previous content it also creates the
// config file if it wasn't exist.
func (c *ConfigFile) Save(v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	return c.creator.Create(file).Encode(v)
}
