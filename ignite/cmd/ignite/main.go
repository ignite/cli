package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"google.golang.org/grpc/status"

	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	chainconfig "github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/internal/analytics"
	"github.com/ignite/cli/v28/ignite/pkg/clictx"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/validation"
	"github.com/ignite/cli/v28/ignite/pkg/xstrings"
)

const exitCodeOK, exitCodeError = 0, 1

func main() {
	os.Exit(run())
}

func run() int {
	ctx := clictx.From(context.Background())
	cmd, cleanUp, err := ignitecmd.New(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		return exitCodeError
	}
	defer cleanUp()

	// find command and send to analytics
	subCmd, _, err := cmd.Find(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		return exitCodeError
	}
	var wg sync.WaitGroup
	analytics.SendMetric(&wg, subCmd)
	analytics.EnableSentry(ctx, &wg)

	err = cmd.ExecuteContext(ctx)
	if err != nil {
		err = ensureError(err)
	}

	if errors.Is(ctx.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
		fmt.Println("aborted")
		return exitCodeOK
	}

	if err != nil {
		var (
			validationErr validation.Error
			versionErr    chainconfig.VersionError
			msg           string
		)

		if errors.As(err, &validationErr) {
			msg = validationErr.ValidationInfo()
		} else {
			msg = err.Error()
		}

		// Make sure the error message starts with an upper case character
		msg = xstrings.ToUpperFirst(msg)

		fmt.Printf("%s %s\n", icons.NotOK, colors.Error(msg))

		if errors.As(err, &versionErr) {
			fmt.Println("Use a more recent CLI version or upgrade blockchain app's config")
		}

		return exitCodeError
	}

	// waits for analytics to finish
	wg.Wait()

	return exitCodeOK
}

func ensureError(err error) error {
	// Extract gRPC error status.
	// These errors are returned by the plugins.
	s, ok := status.FromError(err)
	if !ok {
		// The error is not a gRPC error
		return err
	}

	// Get the error message
	cause := s.Proto().GetMessage()
	if cause == "" {
		return err
	}

	// Restore context canceled errors
	if cause == context.Canceled.Error() {
		return context.Canceled
	}

	// Use the gRPC description as error to avoid printing
	// extra gRPC error information like code or prefix.
	return errors.New(cause)
}
