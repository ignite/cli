package xerrors

import (
	"github.com/cockroachdb/errors"
)

func New(msg string) error { return errors.New(msg) }

func Errorf(format string, args ...interface{}) error { return errors.Errorf(format, args...) }

func Wrap(err error, msg string) error { return errors.Wrap(err, msg) }

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func Unwrap(err error) error { return errors.Unwrap(err) }

func Join(errs ...error) error { return errors.Join(errs...) }

func Is(err, reference error) bool { return errors.Is(err, reference) }

func As(err error, target interface{}) bool { return errors.As(err, target) }

func WithStack(err error) error { return errors.WithStack(err) }
