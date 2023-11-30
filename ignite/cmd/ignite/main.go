package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	ignitecmd "github.com/ignite/cli/ignite/cmd"
	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/gacli"
	"github.com/ignite/cli/ignite/pkg/validation"
	"github.com/ignite/cli/ignite/pkg/xstrings"
)

func main() {
	gaclient = gacli.New()
	os.Exit(run())
}

func run() int {
	const exitCodeOK, exitCodeError = 0, 1
	var wg sync.WaitGroup

	osArgs := strings.Join(os.Args, " ")

	defer func() {
		if r := recover(); r != nil {
			sendMetric(&wg, metric{err: fmt.Errorf("%v", r), command: osArgs})
			fmt.Println(r)
			os.Exit(exitCodeError)
		}
	}()

	if len(os.Args) > 1 {
		sendMetric(&wg, metric{command: osArgs})
	}

	ctx := clictx.From(context.Background())
	cmd, cleanUp, err := ignitecmd.New(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		return exitCodeError
	}
	defer cleanUp()

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
