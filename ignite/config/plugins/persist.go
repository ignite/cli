package plugins

import (
	"os"

	"gopkg.in/yaml.v2"
)

// persist writes a plugin configuration file to a specified file.
// the configuration state that is passed in will be the new state of the file
// before writing the new definition to disk, a truncatte and seek operation
// are performed to assure the file contnets will be overriden.
func persist(config *Config, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}

	// zero out the file for a new config write
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	// seek to the beginning to assure there are no trailing properties
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	// check that plugins are in fact defined within the config
	// if there is an empty array of plugins the encoding will be `{}`
	if len(config.Plugins) > 0 {
		return yaml.NewEncoder(f).Encode(&config)
	}

	return nil
}
