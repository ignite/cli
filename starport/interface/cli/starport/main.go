package main

import (
	"fmt"
	"os"

	starportcmd "github.com/tendermint/starport/starport/interface/cli/starport/cmd"
	"github.com/tendermint/starport/starport/pkg/gacli"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			addMetric(Metric{
				Err: fmt.Errorf("%s", r),
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
	err := starportcmd.New().Execute()
	addMetric(Metric{
		Err: err,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
