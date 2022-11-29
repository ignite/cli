package chain

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Version defines the type for the config version number.
type Version uint

func (v Version) String() string {
	return fmt.Sprintf("v%d", v)
}

// Converter defines the interface required to migrate configurations to newer versions.
type Converter interface {
	// Clone clones the config by returning a new copy of the current one.
	Clone() (Converter, error)

	// SetDefaults assigns default values to empty config fields.
	SetDefaults() error

	// GetVersion returns the config version.
	GetVersion() Version

	// ConvertNext converts the config to the next version.
	ConvertNext() (Converter, error)

	// Decode decodes the config file from YAML and updates its values.
	Decode(io.Reader) error
}

// Build time check for the latest config version type.
// This is required to be sure that conversion to latest
// doesn't break when a new config version is added without
// updating the references to the previous version.
var _ = Versions[LatestVersion].(*ChainConfig)

// ConvertLatest converts a config to the latest version.
func ConvertLatest(c Converter) (_ *ChainConfig, err error) {
	for c.GetVersion() < LatestVersion {
		c, err = c.ConvertNext()
		if err != nil {
			return nil, err
		}
	}

	// Cast to the latest version type.
	// This is safe because there is a build time check that makes sure
	// the type for the latest config version is the right one here.
	return c.(*ChainConfig), nil
}

// MigrateLatest migrates a config file to the latest version.
func MigrateLatest(current io.Reader, latest io.Writer) error {
	cfg, err := Parse(current)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(latest).Encode(cfg)
}
