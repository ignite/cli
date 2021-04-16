// Package sperrors holds starport spesific errors.
package sperrors

import "errors"

var (
	// ErrOnlyStargateSupported is returned when underlying chain is not a stargate chain.
	ErrOnlyStargateSupported = errors.New("only stargate chains are supported")
)
