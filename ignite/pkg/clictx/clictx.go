package clictx

import (
	"context"
	"os"
	"os/signal"
)

// From creates a new context from ctx that is canceled when an exit signal received.
func From(ctx context.Context) context.Context {
	var (
		ctxend, cancel = context.WithCancel(ctx)
		quit           = make(chan os.Signal, 1)
	)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		cancel()
	}()
	return ctxend
}

// Do runs fn and waits for its result unless ctx is canceled.
// Returns fn result or canceled context error.
func Do(ctx context.Context, fn func() error) error {
	errc := make(chan error)
	go func() { errc <- fn() }()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
