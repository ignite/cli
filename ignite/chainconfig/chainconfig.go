package chainconfig

import (
	"errors"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0 "github.com/ignite-hq/cli/ignite/chainconfig/v0"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
	"github.com/ignite-hq/cli/ignite/pkg/xfilepath"
)

var (
	// LatestVersion defines the latest version of the config.
	LatestVersion config.Version = 1

	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New(
		"could not locate a config.yml in your chain. please follow the link for" +
			"how-to: https://github.com/ignite-hq/cli/blob/develop/docs/configure/index.md")

	// ConfigDirPath returns the path of configuration directory of Ignite.
	ConfigDirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

	// ConfigFileNames is a list of recognized names as for Ignite's config file.
	ConfigFileNames = []string{"config.yml", "config.yaml"}

	// DefaultConfig defines the default config without the validators.
	DefaultConfig = &v1.Config{
		BaseConfig: config.BaseConfig{
			Build: config.Build{
				Proto: config.Proto{
					Path: "proto",
					ThirdPartyPaths: []string{
						"third_party/proto",
						"proto_vendor",
					},
				},
			},
			Faucet: config.Faucet{
				Host: "0.0.0.0:4500",
			},
		},
	}

	// Migration defines the version as the key and the config instance as the value
	Migration = map[config.Version]config.Config{
		0: &v0.Config{},
		1: &v1.Config{},
	}
)
