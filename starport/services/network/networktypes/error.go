package networktypes

import "fmt"

// ErrInvalidRequest is an error returned in methods manipulating requests when they are invalid
type ErrInvalidRequest struct {
	requestID uint64
}

// Error implements error
func (err ErrInvalidRequest) Error() string {
	return fmt.Sprintf("request %d is invalid", err.requestID)
}

// NewErrInvalidRequest returns a new ErrInvalidRequest
func NewErrInvalidRequest(requestID uint64) ErrInvalidRequest {
	return ErrInvalidRequest{requestID: requestID}
}
