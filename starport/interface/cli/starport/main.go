package main

import (
	"fmt"
	"os"

	starportcmd "github.com/tendermint/starport/starport/interface/cli/starport/cmd"
	"github.com/tendermint/starport/starport/internal/version"
	"github.com/tendermint/starport/starport/pkg/analyticsutil"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			addMetric(Metric{
				Err: fmt.Errorf("%s", r),
			})
			analyticsc.Close()
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	analyticsc = analyticsutil.New(analyticsEndpoint, analyticsKey)
	// TODO add version of new installation.
	name, hadLogin := prepLoginName()
	analyticsc.Login(name, version.Version)
	if !hadLogin {
		addMetric(Metric{
			Login:          name,
			IsInstallation: true,
		})
	}
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		addMetric(Metric{})
	}
	err := starportcmd.New().Execute()
	addMetric(Metric{
		Err: err,
	})
	analyticsc.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
