package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	ignitecmd "github.com/ignite/cli/ignite/cmd"
	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/validation"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := clictx.From(context.Background())

	cmd := ignitecmd.New()

	// Load plugins if any
	if err := ignitecmd.LoadPlugins(ctx, cmd); err != nil {
		fmt.Printf("Error while loading chain's plugins: %v\n", err)
		return 1
	}
	defer ignitecmd.UnloadPlugins()

	err := cmd.ExecuteContext(ctx)

	if ctx.Err() == context.Canceled || err == context.Canceled {
		fmt.Println("aborted")
		return 0
	}

	if err != nil {
		var validationErr validation.Error

		if errors.As(err, &validationErr) {
			fmt.Println(validationErr.ValidationInfo())
		} else {
			fmt.Println(err)
		}
		return 1
	}
	return 0
}
