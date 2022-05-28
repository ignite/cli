package conversion

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

var (
	// LatestVersion defines the latest version of the config.
	LatestVersion common.Version = 1
)

// ConvertLatest converts a Config to the latest version of Config.
func ConvertLatest(config common.Config) (*v1.Config, error) {
	var err error
	version := config.Version()

	for version < LatestVersion {
		config, err = config.ConvertNext()
		if err != nil {
			return config.(*v1.Config), err
		}
		version = config.Version()
	}
	return config.(*v1.Config), err
}
