package plugins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ParseDir expects to find a plugin config file in dir. If dir is not a folder,
// an error is returned.
// The plugin config file format can be `plugins.yml` or `plugins.yaml`. If
// found, the file is parsed into a Config and returned. If no file from the
// given names above are found, then an empty config is returned, w/o errors.
func ParseDir(dir string) (*Config, error) {
	// handy function that wraps and prefix err with a common label
	errf := func(err error) (*Config, error) {
		return nil, fmt.Errorf("plugin config parse: %w", err)
	}
	fi, err := os.Stat(dir)
	if err != nil {
		return errf(err)
	}
	if !fi.IsDir() {
		return errf(fmt.Errorf("path %s is not a dir", dir))
	}

	filename, err := locateFile(dir)
	if err != nil {
		return errf(err)
	}
	c := Config{
		path: filename,
	}

	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &c, nil
		}
		return errf(err)
	}
	defer f.Close()

	// if the error is end of file meaning an empty file on read return nil
	if err := yaml.NewDecoder(f).Decode(&c); err != nil && !errors.Is(err, io.EOF) {
		return errf(err)
	}
	return &c, nil
}

var (
	// filenames is a list of recognized names as Ignite's plugins config file.
	filenames       = []string{"plugins.yml", "plugins.yaml"}
	defaultFilename = filenames[0]
)

func locateFile(root string) (string, error) {
	for _, name := range filenames {
		path := filepath.Join(root, name)
		_, err := os.Stat(path)
		if err == nil {
			// file found
			return path, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
	}
	// no file found, return the default config name
	return filepath.Join(root, defaultFilename), nil
}
