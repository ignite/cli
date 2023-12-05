package main

import (
	"context"
	"fmt"
	"os"
	"sync"

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

func main() {
	os.Exit(run())
}

func run() int {
	const exitCodeOK, exitCodeError = 0, 1
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

	err = cmd.ExecuteContext(ctx)
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

	wg.Wait() // waits for all metrics to be sent

	return exitCodeOK
}
