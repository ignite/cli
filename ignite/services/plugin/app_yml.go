package plugin

import (
	semver "github.com/blang/semver/v4"
)

// AppYML is the structure of app.ignite.yml file.
type AppYML struct {
	Version semver.Version     `yaml:"version"`
	Apps    map[string]AppInfo `yaml:"apps"`
}

// AppInfo is the structure of app info in app.ignite.yml file which only holds
// the description and the relative path of the app.
type AppInfo struct {
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}
