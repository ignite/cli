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
