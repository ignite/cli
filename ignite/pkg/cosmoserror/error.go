package cosmoserror

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal       = errors.New("internal error")
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
)

// Unwrap checks if an error contains a given grpc error code and returns the corresponding simple error type.
//
//nolint:exhaustive
func Unwrap(err error) error {
	wrapped := err
	for err != nil {
		s, ok := status.FromError(err)
		if ok {
			switch s.Code() {
			case codes.NotFound:
				return ErrNotFound
			case codes.InvalidArgument:
				return ErrInvalidRequest
			case codes.Internal:
				return ErrInternal
			}
		}
		err = errors.Unwrap(err)
	}
	unwrapped := errors.Unwrap(wrapped)
	if unwrapped != nil {
		return unwrapped
	}
	return wrapped
}

// IsNotFound returns if the given error is "not found".
func IsNotFound(err error) bool {
	return errors.Is(Unwrap(err), ErrNotFound)
}
