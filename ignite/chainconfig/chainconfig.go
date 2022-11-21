package chainconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

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

	// DefaultVuePath defines the default relative path to use when scaffolding a Vue app.
	// The path is relative to the app's directory.
	DefaultVuePath = "vue"

	// DefaultReactPath defines the default relative path to use when scaffolding a React app.
	// The path is relative to the app's directory.
	DefaultReactPath = "react"

	// DefaultVuexPath defines the default relative path to use when generating Vuex stores for a Vue app.
	// The path is relative to the app's directory.
	DefaultVuexPath = "vue/src/store"

	// DefaultComposablesPath defines the default relative path to use when generating useQuery composables for a Vue app.
	// The path is relative to the app's directory.
	DefaultComposablesPath = "vue/src/composables"

	// DefaultHooksPath defines the default relative path to use when generating useQuery hooks for a React app.
	// The path is relative to the app's directory.
	DefaultHooksPath = "react/src/hooks"

	// DefaultOpenAPIPath defines the default relative path to use when when generating an OpenAPI schema.
	// The path is relative to the app's directory.
	DefaultOpenAPIPath = "docs/static/openapi.yml"

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

// Plugin defines the latest plugin config
type Plugin = v1.Plugin

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
func TSClientPath(conf Config) string {
	if path := strings.TrimSpace(conf.Client.Typescript.Path); path != "" {
		return filepath.Clean(path)
	}

	return DefaultTSClientPath
}

// VuexPath returns the relative path to the Vuex stores directory.
// Path is relative to the app's directory.
func VuexPath(conf *Config) string {
	//nolint:staticcheck //ignore SA1019 until vuex config option is removed
	if path := strings.TrimSpace(conf.Client.Vuex.Path); path != "" {
		return filepath.Clean(path)
	}

	return DefaultVuexPath
}

// ComposablesPath returns the relative path to the Vue useQuery composables directory.
// Path is relative to the app's directory.
func ComposablesPath(conf *Config) string {
	if path := strings.TrimSpace(conf.Client.Composables.Path); path != "" {
		return filepath.Clean(path)
	}

	return DefaultComposablesPath
}

// HooksPath returns the relative path to the React useQuery hooks directory.
// Path is relative to the app's directory.
func HooksPath(conf *Config) string {
	if path := strings.TrimSpace(conf.Client.Hooks.Path); path != "" {
		return filepath.Clean(path)
	}

	return DefaultHooksPath
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

// CheckVersion checks that the config version is the latest
// and if not a VersionError is returned.
func CheckVersion(configFile io.Reader) error {
	version, err := ReadConfigVersion(configFile)
	if err != nil {
		return err
	}

	if version != LatestVersion {
		return VersionError{version}
	}

	return nil
}

// Save saves a config to a YAML file.
func Save(c Config, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}

	defer file.Close()

	return yaml.NewEncoder(file).Encode(c)
}
