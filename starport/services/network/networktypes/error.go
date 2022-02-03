package networktypes

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrInvalidRequest is an error returned in methods manipulating requests when they are invalid
type ErrInvalidRequest struct {
	requestID uint64
}

// Error implements error
func (err ErrInvalidRequest) Error() string {
	return fmt.Sprintf("request %d is invalid", err.requestID)
}

// NewWrappedErrInvalidRequest returns a wrapped ErrInvalidRequest
func NewWrappedErrInvalidRequest(requestID uint64, message string) error {
	return errors.Wrap(ErrInvalidRequest{requestID: requestID}, message)
}
