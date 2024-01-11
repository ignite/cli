package plugin

// AppsConfig is the structure of app.ignite.yml file.
type AppsConfig struct {
	Version uint               `yaml:"version"`
	Apps    map[string]AppInfo `yaml:"apps"`
}

// AppInfo is the structure of app info in app.ignite.yml file which only holds
// the description and the relative path of the app.
type AppInfo struct {
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}
