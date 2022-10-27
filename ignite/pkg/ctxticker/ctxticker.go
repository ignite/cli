package ctxticker

import (
	"context"
	"time"
)

// Do calls fn every d until ctx canceled or fn returns with a non-nil error.
func Do(ctx context.Context, d time.Duration, fn func() error) error {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := fn(); err != nil {
				return err
			}
		}
	}
}

// DoNow is same as Do except it makes +1 call to fn on start.
func DoNow(ctx context.Context, d time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		return err
	}
	return Do(ctx, d, fn)
}
