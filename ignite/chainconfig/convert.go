package chainconfig

import (
	"os"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
	"gopkg.in/yaml.v2"
)

// Build time check for the latest config version type.
// This is required to be sure that conversion to latest
// doesn't break when a new config version is added without
// updating the references to the previous version.
var _ = Versions[LatestVersion].(*v1.Config)

// ConvertLatest converts a config to the latest version.
func ConvertLatest(c config.Converter) (*v1.Config, error) {
	var err error

	for c.GetVersion() < LatestVersion {
		c, err = c.ConvertNext()
		if err != nil {
			return nil, err
		}
	}

	// Cast to the latest version type.
	// This is safe because there is a build time check that makes sure
	// the type for the latest config version is the right one here.
	return c.(*v1.Config), nil
}

// MigrateLatest migrates a config file to the latest version.
// TODO: Change to receive an io.Reader and an io.Writer instead of saving to file
func MigrateLatest(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	cfg, err := Parse(file)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(file).Encode(cfg)
}
