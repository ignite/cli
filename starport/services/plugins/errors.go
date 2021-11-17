package plugins

import "errors"

var (
	ErrCommandNotFound            = errors.New("Command not found")
	ErrCommandPluginNotRecognized = errors.New("Command plugin not recognized")
)
