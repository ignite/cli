package chain

import (
	"io"

	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/config/chain/version"
)

// Build time check for the latest config version type.
// This is required to be sure that conversion to latest
// doesn't break when a new config version is added without
// updating the references to the previous version.
var _ = Versions[LatestVersion].(*Config)

// ConvertLatest converts a config to the latest version.
func ConvertLatest(c version.Converter) (_ *Config, err error) {
	for c.GetVersion() < LatestVersion {
		c, err = c.ConvertNext()
		if err != nil {
			return nil, err
		}
	}

	// Cast to the latest version type.
	// This is safe because there is a build time check that makes sure
	// the type for the latest config version is the right one here.
	return c.(*Config), nil
}

// MigrateLatest migrates a config file to the latest version.
func MigrateLatest(current io.Reader, latest io.Writer) error {
	cfg, err := Parse(current)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(latest).Encode(cfg)
}
