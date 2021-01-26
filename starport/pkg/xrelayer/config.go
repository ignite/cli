package xrelayer

import (
	"context"
	"os"
	"path/filepath"
	"time"

	relayercmd "github.com/cosmos/relayer/cmd"
	"github.com/cosmos/relayer/relayer"
	"github.com/tendermint/starport/starport/pkg/confile"
)

var (
	// confHome is the home path of relayer.
	confHome = os.ExpandEnv("$HOME/.relayer")

	// confYamlPath is the path of relayer's config.yaml.
	confYamlPath = filepath.Join(confHome, "config/config.yaml")

	// cfile is used to load relayer's config yaml and overwrite any changes.
	cfile = confile.New(confile.DefaultYAMLEncodingCreator, confYamlPath)

	// defaultConf is a default configuration for relayer's config.yml.
	defaultConf = relayercmd.Config{
		Global: relayercmd.GlobalConfig{
			Timeout:        "10s",
			LightCacheSize: 20,
		},
		Chains: relayer.Chains{},
		Paths:  relayer.Paths{},
	}
)

// config returns the representation of config.yml.
// it deals with creating and adding default configs if there wasn't a config.yml before.
func config(ctx context.Context) (relayercmd.Config, error) {
	// ensure that config.yaml exists.
	if _, err := os.Stat(confYamlPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(confYamlPath), os.ModePerm); err != nil {
			return relayercmd.Config{}, err
		}

		if err := cfile.Save(defaultConf); err != nil {
			return relayercmd.Config{}, err
		}
	}

	// load config.yaml
	rconf := relayercmd.Config{}
	if err := cfile.Load(&rconf); err != nil {
		return relayercmd.Config{}, err
	}

	// init loaded configs.
	globalTimeout, err := time.ParseDuration(rconf.Global.Timeout)
	if err != nil {
		return relayercmd.Config{}, newRelayerError("global.timeout is invalid")
	}

	for _, i := range rconf.Chains {
		if err := i.Init(confHome, globalTimeout, false); err != nil {
			return relayercmd.Config{}, newRelayerError("cannot init")
		}
	}

	return rconf, nil
}
