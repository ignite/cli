package chainconfig

import (
	"errors"
	"fmt"

	"github.com/ignite/cli/ignite/chainconfig/config"
)

// ErrConfigNotFound indicates that the config.yml can't be found.
var ErrConfigNotFound = errors.New("could not locate a config.yml in your chain")

// ValidationError is returned when a configuration is invalid.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config is not valid: %s", e.Message)
}

// UnsupportedVersionError is returned when the version of the config is not supported.
type UnsupportedVersionError struct {
	Version config.Version
}

func (e *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("config version %d is not supported", e.Version)
}
