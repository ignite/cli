// Package errors provides helpers for error creation, avoiding
// using different packages for errors.
//
// e.g.:
//
//	import "github.com/ignite/cli/v28/ignite/pkg/errors"
//
//	func main() {
//	 err1 := errors.New("error new")
//	 err2 := errors.Errorf("%s: error", foo)
//	 err3 := errors.Wrap(errFoo, errBar)
//	}
package errors

import (
	"github.com/cockroachdb/errors"
)

// New creates an error with a simple error message.
// A stack trace is retained.
func New(msg string) error { return errors.New(msg) }

// Errorf aliases Newf().
func Errorf(format string, args ...interface{}) error { return errors.Errorf(format, args...) }

// Wrap wraps an error with a message prefix. A stack trace is retained.
func Wrap(err error, msg string) error { return errors.Wrap(err, msg) }

// Wrapf wraps an error with a formatted message prefix. A stack
// trace is also retained. If the format is empty, no prefix is added,
// but the extra arguments are still processed for reportable strings.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Unwrap accesses the direct cause of the error if any, otherwise
// returns nil.
func Unwrap(err error) error { return errors.Unwrap(err) }

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if errs contains no non-nil values.
func Join(errs ...error) error { return errors.Join(errs...) }

// Is determines whether one of the causes of the given error or any
// of its causes is equivalent to some reference error.
func Is(err, reference error) bool { return errors.Is(err, reference) }

// As finds the first error in err's chain that matches the type to which target
// points, and if so, sets the target to its value and returns true. An error
// matches a type if it is assignable to the target type, or if it has a method
// As(interface{}) bool such that As(target) returns true. As will panic if target
// is not a non-nil pointer to a type which implements error or is of interface type.
func As(err error, target interface{}) bool { return errors.As(err, target) }

// WithStack annotates err with a stack trace at the point WithStack was called.
func WithStack(err error) error { return errors.WithStack(err) }
