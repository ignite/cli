package chainconfig

import (
	"github.com/ignite/cli/ignite/chainconfig/common"
)

// ConvertLatest converts a Config to the latest version of Config.
func ConvertLatest(config common.Config) (common.Config, error) {
	var err error
	version := config.Version()

	for version < common.LatestVersion {
		config, err = config.ConvertNext()
		if err != nil {
			return config, err
		}
		version = config.Version()
	}
	return config, err
}
