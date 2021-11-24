package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	starportcmd "github.com/tendermint/starport/starport/cmd"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/validation"
	"github.com/tendermint/starport/starport/services/plugins"
)

func main() {
	ctx := clictx.From(context.Background())

	// Check if this actually preruns, idk if it is right now
	starportCmd := starportcmd.New(ctx)
	// Set PersistentPreRunE here (https://pkg.go.dev/github.com/spf13/cobra#Command)
	// Find a way to have this run on setup, not on every usage of the command
	starportCmd.PersistentPreRunE = plugins.PersistentPreRunE
	err := starportCmd.ExecuteContext(ctx)
	if ctx.Err() == context.Canceled || err == context.Canceled {
		fmt.Println("aborted")
		return
	}

	if err != nil {
		var validationErr validation.Error

		if errors.As(err, &validationErr) {
			fmt.Println(validationErr.ValidationInfo())
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}
