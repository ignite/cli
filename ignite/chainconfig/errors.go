package chainconfig

import "fmt"

// ValidationError is returned when a configuration is invalid.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config is not valid: %s", e.Message)
}

// UnsupportedVersionError is returned when the version of the config is not supported.
type UnsupportedVersionError struct {
	Message string
}

func (e *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("the version of the config is unsupported: %s", e.Message)
}

// UnknownInputError is returned when the input of Parse is unknown.
type UnknownInputError struct {
	Message string
}

func (e *UnknownInputError) Error() string {
	return fmt.Sprintf("the version of the config is unsupported: %s", e.Message)
}
