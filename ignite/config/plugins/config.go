package plugins

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// path to the config file
	path string

	Plugins []Plugin `yaml:"plugins"`
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
	// contains multiple plugins, for example:
	//
	// path: github.com/foo/bar/plugin1
	//
	// It can also specify a tag or a branch, by adding a `@` and the branch/tag
	// name at the end of the path. For example:
	//
	// path: github.com/foo/bar/plugin1@v42
	Path string `yaml:"path"`
	// With holds arguments passed to the plugin interface
	With map[string]string `yaml:"with,omitempty"`
	// Global holds whether the plugin is installed globally
	// (default: $HOME/.ignite/plugins/plugins.yml) or locally for a chain.
	Global bool `yaml:"-"`
}

func (p Plugin) HasPath(path string) bool {
	if path == "" {
		return false
	}
	if p.Path == path {
		return true
	}
	pluginPath := strings.Split(p.Path, "@")[0]
	path = strings.Split(path, "@")[0]
	return pluginPath == path
}

// Path return the path of the config file.
func (c Config) Path() string {
	return c.path
}

// Save persists a config yaml to a specified path on disk.
// Must be writable.
func (c *Config) Save() error {
	errf := func(err error) error {
		return fmt.Errorf("plugin config save: %w", err)
	}
	if c.path == "" {
		return errf(errors.New("empty path"))
	}
	file, err := os.Create(c.path)
	if err != nil {
		return errf(err)
	}
	defer file.Close()
	if err := yaml.NewEncoder(file).Encode(c); err != nil {
		return errf(err)
	}
	return nil
}

// HasPlugin returns true if c contains a plugin with given path.
func (c Config) HasPlugin(path string) bool {
	return slices.ContainsFunc(c.Plugins, func(cp Plugin) bool {
		return cp.HasPath(path)
	})
}
