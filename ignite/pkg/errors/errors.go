package errors

// ValidationError must be implemented by errors that provide validation info.
type ValidationError interface {
	error
	ValidationInfo() string
}
