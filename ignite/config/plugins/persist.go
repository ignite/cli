package plugins

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Persist(config Config, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(&config)
}
