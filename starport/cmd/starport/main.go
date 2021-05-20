package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	starportcmd "github.com/tendermint/starport/starport/cmd"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/gacli"
	"github.com/tendermint/starport/starport/pkg/validation"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			addMetric(Metric{
				Err: fmt.Errorf("%v", r),
			})
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	gaclient = gacli.New(gaid)
	name, hadLogin := prepLoginName()
	if !hadLogin {
		addMetric(Metric{
			Login:          name,
			IsInstallation: true,
		})
	}
	// if running serve command, don't wait sending metric until the end of
	// execution because it takes a long time.
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		addMetric(Metric{})
	}

	ctx := clictx.From(context.Background())
	err := starportcmd.New().ExecuteContext(ctx)

	if err == context.Canceled {
		addMetric(Metric{
			Err: err,
		})
		fmt.Println("aborted")
		return
	}
	if err != nil {
		var validationErr validation.Error
		if errors.As(err, &validationErr) {
			panic(validationErr.ValidationInfo())
		} else {
			panic(err)
		}
	}
}
