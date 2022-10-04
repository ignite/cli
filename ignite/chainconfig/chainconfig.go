package chainconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// DefaultTSClientPath defines the default relative path to use when generating the TS client.
	// The path is relative to the app's directory.
	DefaultTSClientPath = "ts-client"

	// LatestVersion defines the latest version of the config.
	LatestVersion config.Version = 1

	// Versions holds config types for the supported versions.
	Versions = map[config.Version]config.Converter{
		0: &v0.Config{},
		1: &v1.Config{},
	}
)

// Config defines the latest config.
type Config = v1.Config

// DefaultConfig returns a config for the latest version initialized with default values.
func DefaultConfig() *Config {
	return v1.DefaultConfig()
}

// FaucetHost returns the faucet host to use.
func FaucetHost(cfg *Config) string {
	// We keep supporting Port option for backward compatibility
	// TODO: drop this option in the future
	host := cfg.Faucet.Host
	if cfg.Faucet.Port != 0 {
		host = fmt.Sprintf(":%d", cfg.Faucet.Port)
	}

	return host
}

// TSClientPath returns the relative path to the Typescript client directory.
// Path is relative to the app's directory.
func TSClientPath(conf *Config) string {
	if path := strings.TrimSpace(conf.Client.Typescript.Path); path != "" {
		return filepath.Clean(path)
	}

	return DefaultTSClientPath
}

// CreateConfigDir creates config directory if it is not created yet.
func CreateConfigDir() error {
	path, err := ConfigDirPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(path, 0o755)
}

// LocateDefault locates the default path for the config file.
// Returns ErrConfigNotFound when no config file found.
func LocateDefault(root string) (path string, err error) {
	for _, name := range ConfigFileNames {
		path = filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", ErrConfigNotFound
}
