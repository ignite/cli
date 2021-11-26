package pluginsrpc

import "errors"

var (
	ErrCommandNotFound            = errors.New("Command not found")
	ErrCommandPluginNotRecognized = errors.New("Command plugin not recognized")
	ErrHookPluginNotRecognized    = errors.New("Hook plugin not recognized")
)
