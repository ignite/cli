package common

// Config is the interface defining all the common methods for the ConfigYaml struct across all supported versions
type Config interface {
	Clone() Config

	// Version returns the version of the Config
	Version() Version

	// ConvertNext converts the instance of Config from the current version to the next version.
	ConvertNext() (Config, error)
}
