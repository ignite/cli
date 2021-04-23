// Package sperrors holds starport spesific errors.
package sperrors

import "errors"

var (
	// ErrOnlyStargateSupported is returned when underlying chain is not a stargate chain.
	ErrOnlyStargateSupported = errors.New("this version of Cosmos SDK is no longer supported")
)
