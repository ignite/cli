package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	ignitecmd "github.com/ignite/cli/ignite/cmd"
	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/validation"
	"github.com/ignite/cli/ignite/pkg/xstrings"
)

func main() {
	os.Exit(run())
}

func run() int {
	const (
		exitCodeOK    = 0
		exitCodeError = 1
	)
	ctx := clictx.From(context.Background())

	cmd := ignitecmd.New()

	// Load plugins if any
	if err := ignitecmd.LoadPlugins(ctx, cmd); err != nil {
		fmt.Printf("Error while loading plugins: %v\n", err)
		return exitCodeError
	}
	defer ignitecmd.UnloadPlugins()

	err := cmd.ExecuteContext(ctx)

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
	return exitCodeOK
}
