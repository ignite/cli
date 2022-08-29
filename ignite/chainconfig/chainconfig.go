package chainconfig

import (
	"github.com/ignite/cli/ignite/chainconfig/config"
	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

var (
	// ConfigDirPath returns the path of configuration directory of Ignite.
	ConfigDirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

	// ConfigFileNames is a list of recognized names as for Ignite's config file.
	ConfigFileNames = []string{"config.yml", "config.yaml"}

	// LatestVersion defines the latest version of the config.
	LatestVersion config.Version = 1

	// Versions holds config types for the supported versions.
	Versions = map[config.Version]config.Converter{
		0: &v0.Config{},
		1: &v1.Config{},
	}
)

// DefaultConfig returns a config for the latest version initialized with default values.
func DefaultConfig() *v1.Config {
	return v1.DefaultConfig()
}
