package plugin

// ExtensionsConfig is the structure of ext.ignite.yml file.
type ExtensionsConfig struct {
	Version    uint               `yaml:"version"`
	Extensions map[string]ExtInfo `yaml:"extensions"`
}

// ExtInfo is the structure of app info in ext.ignite.yml file which only holds
// the description and the relative path of the app.
type ExtInfo struct {
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}
