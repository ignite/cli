package chainconfig

import "github.com/ignite-hq/cli/ignite/chainconfig/config"

// ConvertLatest converts a Config to the latest version of Config.
func ConvertLatest(cfg config.Converter) (config.Converter, error) {
	var err error

	version := cfg.GetVersion()

	for version < LatestVersion {
		cfg, err = cfg.ConvertNext()
		if err != nil {
			return nil, err
		}

		version = cfg.GetVersion()
	}

	return cfg, nil
}
