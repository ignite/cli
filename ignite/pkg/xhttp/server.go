package xhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// ShutdownTimeout is the timeout for waiting all requests to complete.
const ShutdownTimeout = time.Minute

// Serve starts s server and shutdowns it once the ctx is cancelled.
func Serve(ctx context.Context, s *http.Server) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		_ = s.Shutdown(shutdownCtx)
	}()

	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
