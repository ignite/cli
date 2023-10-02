package apps

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/pkg/gomodule"
)

type Config struct {
	path string

	// Apps holds the list of installed apps.
	Apps []App `yaml:"apps"`
}

// App keeps app name and location.
type App struct {
	// Path holds the location of the app.
	// A path can be local, in that case it must start with a `/`.
	// A remote path on the other hand, is an URL to a public remote git
	// repository. For example:
	//
	// path: github.com/foo/bar
	//
	// It can contain a path inside that repository, if for instance the repo
	// contains multiple apps, for example:
	//
	// path: github.com/foo/bar/app1
	//
	// It can also specify a tag or a branch, by adding a `@` and the branch/tag
	// name at the end of the path. For example:
	//
	// path: github.com/foo/bar/app1@v42
	Path string `yaml:"path"`

	// With holds arguments passed to the app interface
	With map[string]string `yaml:"with,omitempty"`

	// Global holds whether the app is installed globally
	// (default: $HOME/.ignite/apps/igniteapps.yml) or locally for a chain.
	Global bool `yaml:"-"`
}

// RemoveDuplicates takes a list of Apps and returns a new list with only unique values.
// Local apps take precedence over global apps if duplicate paths exist.
// Duplicates are compared regardless of version.
func RemoveDuplicates(apps []App) (unique []App) {
	// struct to track app configs
	type check struct {
		hasPath   bool
		global    bool
		prevIndex int
	}

	keys := make(map[string]check)
	for i, app := range apps {
		c := keys[app.CanonicalPath()]
		if !c.hasPath {
			keys[app.CanonicalPath()] = check{
				hasPath:   true,
				global:    app.Global,
				prevIndex: i,
			}
			unique = append(unique, app)
		} else if c.hasPath && !app.Global && c.global { // overwrite global app if local duplicate exists
			unique[c.prevIndex] = app
		}
	}

	return unique
}

// IsGlobal returns whether the app is installed globally or locally for a chain.
func (a App) IsGlobal() bool {
	return a.Global
}

// IsLocalPath returns true if the app path is a local directory.
func (a App) IsLocalPath() bool {
	return strings.HasPrefix(a.Path, "/")
}

// HasPath verifies if an app has the given path regardless of version.
// For example github.com/foo/bar@v1 and github.com/foo/bar@v2 have the
// same path so "true" will be returned.
func (a App) HasPath(path string) bool {
	if path == "" {
		return false
	}
	if a.Path == path {
		return true
	}
	appPath := a.CanonicalPath()
	path = strings.Split(path, "@")[0]
	return appPath == path
}

// CanonicalPath returns the canonical path of an app (excludes version ref).
func (a App) CanonicalPath() string {
	return strings.Split(a.Path, "@")[0]
}

// Path return the path of the config file.
func (c Config) Path() string {
	return c.path
}

// Save persists a config yaml to a specified path on disk.
// Must be writable.
func (c *Config) Save() error {
	errf := func(err error) error {
		return fmt.Errorf("app config save: %w", err)
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

// HasApp returns true if c contains an app with given path.
// Returns also true if there's a local app with the module name equal to
// that path.
func (c Config) HasApp(path string) bool {
	return slices.ContainsFunc(c.Apps, func(app App) bool {
		if app.HasPath(path) {
			return true
		}
		if app.IsLocalPath() {
			// check local app go.mod to see if module name match app path
			gm, err := gomodule.ParseAt(app.Path)
			if err != nil {
				// Skip if we can't parse gomod
				return false
			}
			return App{Path: gm.Module.Mod.Path}.HasPath(path)
		}
		return false
	})
}
