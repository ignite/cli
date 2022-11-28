package plugins

import (
	"io"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// DefaultConfig returns a config with default values.
func DefaultConfig() *Config {
	c := Config{}
	return &c
}

type Config struct {
	Plugins []Plugin `yaml:"plugins,omitempty"`
}

// Plugin keeps plugin name and location.
type Plugin struct {
	// Path holds the location of the plugin.
	// A path can be local, in that case it must start with a `/`.
	// A remote path on the other hand, is an URL to a public remote git
	// repository. For example:
	//
	// path: github.com/foo/bar
	//
	// It can contain a path inside that repository, if for instance the repo
	// contains multiple plugins, For example:
	//
	// path: github.com/foo/bar/plugin1
	//
	// It can also specify a tag or a branch, by adding a `@` and the branch/tag
	// name at the end of the path. For example:
	//
	// path: github.com/foo/bar/plugin1@v42
	Path string `yaml:"path"`
	// With holds arguments passed to the plugin interface
	With map[string]string `yaml:"with"`
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() (*Config, error) {
	copy := Config{}
	if err := mergo.Merge(&copy, c, mergo.WithAppendSlice); err != nil {
		return nil, err
	}

	return &copy, nil
}

// Decode decodes the config file values from YAML.
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}
