package conversion

import (
	"github.com/ignite/cli/ignite/chainconfig/common"
	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
)

// ConvertLatest converts a Config to the latest version of Config.
func ConvertLatest(config common.Config) (*v1.Config, error) {
	var err error
	version := config.Version()

	for version < common.LatestVersion {
		config, err = config.ConvertNext()
		if err != nil {
			return config.(*v1.Config), err
		}
		version = config.Version()
	}
	return config.(*v1.Config), err
}
