package validation

// Error must be implemented by errors that provide validation info.
type Error interface {
	error
	ValidationInfo() string
}
