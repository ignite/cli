// Package sperrors holds starport spesific errors.
package sperrors

import "errors"

// ErrOnlyStargateSupported is returned when underlying chain is not a stargate chain.
var ErrOnlyStargateSupported = errors.New("this version of Cosmos SDK is no longer supported")
