package plugins

import "errors"

// ErrConfigNotFound indicates that the plugins.yml can't be found.
var ErrConfigNotFound = errors.New("could not locate a plugins.yml")
