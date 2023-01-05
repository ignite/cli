package looseerrgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Wait waits until g.Wait() returns or ctx canceled, whichever occurs first.
// returned error is context.Canceled if ctx canceled otherwise the error returned by g.Wait().
//
// this is useful when errgroup cannot be used with errgroup.WithContext which happens if executed
// func does not support cancellation.
func Wait(ctx context.Context, g *errgroup.Group) error {
	doneC := make(chan struct{})

	go func() { g.Wait(); close(doneC) }()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-doneC:
		return g.Wait()
	}
}
