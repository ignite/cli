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
	ctx := clictx.From(context.Background())

	err := ignitecmd.New().ExecuteContext(ctx)

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
