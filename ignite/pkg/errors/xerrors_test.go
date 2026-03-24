package errors

import (
	"context"
	stdErrors "errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type customErr struct {
	msg string
}

func (e customErr) Error() string { return e.msg }

func TestBasicHelpers(t *testing.T) {
	err := New("boom")
	require.EqualError(t, err, "boom")

	err = Errorf("value: %d", 10)
	require.EqualError(t, err, "value: 10")
}

func TestWrapHelpers(t *testing.T) {
	base := stdErrors.New("base")

	require.Nil(t, Wrap(nil, "prefix"))
	require.Nil(t, Wrapf(nil, "prefix %s", "x"))

	wrapped := Wrap(base, "prefix")
	require.Error(t, wrapped)
	require.True(t, Is(wrapped, base))

	wrapped = Wrapf(base, "prefix %s", "x")
	require.Error(t, wrapped)
	require.True(t, Is(wrapped, base))
}

func TestJoinUnwrapAs(t *testing.T) {
	e1 := customErr{msg: "one"}
	e2 := stdErrors.New("two")
	j := Join(e1, e2)
	require.Error(t, j)
	require.True(t, Is(j, e2))

	var target customErr
	require.True(t, As(j, &target))
	require.Equal(t, "one", target.msg)

	wrapped := Wrap(e2, "prefix")
	require.True(t, Is(Unwrap(wrapped), e2))
}

func TestWithStack(t *testing.T) {
	base := stdErrors.New("base")
	require.True(t, Is(WithStack(base), base))
}

func TestShouldCaptureException(t *testing.T) {
	t.Run("canceled errors are ignored", func(t *testing.T) {
		require.False(t, shouldCaptureException(context.Canceled))
		require.False(t, shouldCaptureException(Wrap(context.Canceled, "prefix")))
		require.False(t, shouldCaptureException(Wrapf(context.Canceled, "prefix %s", "x")))
		require.False(t, shouldCaptureException(WithStack(context.Canceled)))
	})

	t.Run("other errors are reported", func(t *testing.T) {
		require.True(t, shouldCaptureException(stdErrors.New("boom")))
	})
}
