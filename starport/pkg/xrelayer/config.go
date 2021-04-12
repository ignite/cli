package xrelayer

import (
	"context"
	"os"
	"path/filepath"
	"time"

	relayercmd "github.com/cosmos/relayer/cmd"
	"github.com/cosmos/relayer/relayer"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/tendermintlogger"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

var (
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

// confHome returns the home path of relayer
func confHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".relayer"), nil
}

// confYamlPath returns the path of relayer's config.yaml
func confYamlPath() (string, error) {
	confHome, err := confHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(confHome, "config", "config.yaml"), nil
}

// confFile returns the file used to load relayer's config yaml and overwrite any changes
func confFile() (*confile.ConfigFile, error) {
	confYamlPath, err := confYamlPath()
	if err != nil {
		return nil, err
	}

	return confile.New(confile.DefaultYAMLEncodingCreator, confYamlPath), nil
}

// config returns the representation of config.yml.
// it deals with creating and adding default configs if there wasn't a config.yml before.
func config(_ context.Context, enableLogs bool) (relayercmd.Config, error) {
	confHome, err := confHome()
	if err != nil {
		return relayercmd.Config{}, err
	}
	confYamlPath, err := confYamlPath()
	if err != nil {
		return relayercmd.Config{}, err
	}
	confFile, err := confFile()
	if err != nil {
		return relayercmd.Config{}, err
	}

	// ensure that config.yaml exists.
	if _, err := os.Stat(confYamlPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(confYamlPath), os.ModePerm); err != nil {
			return relayercmd.Config{}, err
		}

		if err := confFile.Save(defaultConf); err != nil {
			return relayercmd.Config{}, err
		}
	}

	// load config.yaml
	rconf := relayercmd.Config{}
	if err := confFile.Load(&rconf); err != nil {
		return relayercmd.Config{}, err
	}

	// init loaded configs.
	globalTimeout, err := time.ParseDuration(rconf.Global.Timeout)
	if err != nil {
		return relayercmd.Config{}, errors.New("relayer's global.timeout is invalid")
	}

	var logger tmlog.Logger
	if !enableLogs {
		logger = tendermintlogger.DiscardLogger{}
	}

	for _, i := range rconf.Chains {
		if err := i.Init(confHome, globalTimeout, logger, false); err != nil {
			return relayercmd.Config{}, errors.New("cannot init relayer")
		}
	}

	return rconf, nil
}
