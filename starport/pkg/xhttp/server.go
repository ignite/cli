package xhttp

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ShutdownTimeout is the timeout for waiting all requests to complete.
const ShutdownTimeout = time.Minute

// Serve starts s server and shutdowns it once the ctx is cancelled.
func Serve(ctx context.Context, s *http.Server) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		s.Shutdown(shutdownCtx)
	}()

	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
