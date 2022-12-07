package chain

import (
	"errors"
	"fmt"

	"github.com/ignite/cli/ignite/config/chain/version"
)

// ErrConfigNotFound indicates that the config.yml can't be found.
var ErrConfigNotFound = errors.New("could not locate a config.yml in your chain")

// ValidationError is returned when a configuration is invalid.
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config is not valid: %s", e.Message)
}

// UnsupportedVersionError is returned when the version of the config is not supported.
type UnsupportedVersionError struct {
	Version version.Version
}

func (e UnsupportedVersionError) Error() string {
	return fmt.Sprintf("config version %s is not supported", e.Version)
}

// VersionError is returned when config version doesn't match with the version CLI supports.
type VersionError struct {
	Version version.Version
}

func (e VersionError) Error() string {
	if LatestVersion > e.Version {
		return fmt.Sprintf(
			"blockchain app uses a previous config version %s and CLI expects %s",
			e.Version,
			LatestVersion,
		)
	}

	return fmt.Sprintf(
		"blockchain app uses a newer config version %s and CLI expects %s",
		e.Version,
		LatestVersion,
	)
}
