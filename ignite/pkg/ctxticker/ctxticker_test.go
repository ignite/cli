package ctxticker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoNow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var callCount int

	require.Error(t, context.Canceled, DoNow(ctx, time.Millisecond, func() error {
		if callCount == 3 {
			cancel()
			return nil
		}
		callCount++
		return nil
	}))

	require.True(t, callCount >= 3)
}

func TestDoNowError(t *testing.T) {
	errDone := errors.New("done")
	var callCount int

	require.Error(t, errDone, DoNow(context.Background(), time.Millisecond, func() error {
		if callCount == 3 {
			return errDone
		}
		callCount++
		return nil
	}))

	require.True(t, callCount >= 3)
}
