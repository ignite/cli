package conversion

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/common"
)

var (
	// LatestVersion defines the latest version of the config.
	LatestVersion = 1
)

// ConvertNext converts a Config to the next version of Config.
func ConvertNext(config common.Config) (common.Config, error) {
	return config.ConvertNext()
}

// ConvertLatest converts a Config to the latest version of Config.
func ConvertLatest(config common.Config) (common.Config, error) {
	var err error
	version := config.GetVersion()

	for version < LatestVersion {
		config, err = ConvertNext(config)
		if err != nil {
			return config, err
		}
		version = config.GetVersion()
	}
	return config, err
}
