package cosmoserror

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal       = errors.New("some invariants expected by the underlying system has been broken")
	ErrInvalidRequest = errors.New("invalid GRPC request argument or object not found")
)

func Unwrap(err error) error {
	s, ok := status.FromError(err)
	if ok {
		switch s.Code() {
		case codes.InvalidArgument:
			return ErrInvalidRequest
		case codes.Internal:
			return ErrInternal
		}
	}
	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		return unwrapped
	}
	return err
}
