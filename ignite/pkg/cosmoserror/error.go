package cosmoserror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
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
	if err := errors.Unwrap(err); err != nil {
		return err
	}
	return err
}

// IsNotFound returns if the given error is "not found".
func IsNotFound(err error) bool {
	return errors.Is(Unwrap(err), ErrNotFound)
}
