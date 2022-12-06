package version

import (
	"fmt"
	"io"
)

// Version defines the type for the config version number.
type Version uint

func (v Version) String() string {
	return fmt.Sprintf("v%d", v)
}

// Converter defines the interface required to migrate configurations to newer versions.
type Converter interface {
	// Clone clones the config by returning a new copy of the current one.
	Clone() (Converter, error)

	// SetDefaults assigns default values to empty config fields.
	SetDefaults() error

	// GetVersion returns the config version.
	GetVersion() Version

	// ConvertNext converts the config to the next version.
	ConvertNext() (Converter, error)

	// Decode decodes the config file from YAML and updates its values.
	Decode(io.Reader) error
}
