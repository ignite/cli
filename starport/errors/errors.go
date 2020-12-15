// Package errors keeps Starport errors.
package errors

import "errors"

var (
	// ErrStarportRequiresProtoc returned when protoc is not installed.
	ErrStarportRequiresProtoc = errors.New("starport requires protoc installed.\nPlease, follow instructions on https://grpc.io/docs/protoc-installation")
)
