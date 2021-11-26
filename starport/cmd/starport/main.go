package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
	starportcmd "github.com/tendermint/starport/starport/cmd"
	"github.com/tendermint/starport/starport/pkg/clictx"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/validation"
)

func main() {
	ctx := clictx.From(context.Background())

	err := starportcmd.New(ctx).ExecuteContext(ctx)
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

func getDefaultConfig(cmd *cobra.Command) (chaincfg.Config, string, error) {
	// need new way of getting path of chain
	appPath := flagGetPath(cmd)
	absPath, err := filepath.Abs(appPath)
	if err != nil {
		return chaincfg.Config{}, "", err
	}

	p, gappPath, err := gomodulepath.Find(absPath)
	if err != nil {
		return chaincfg.Config{}, "", err
	}

	chainId := p.Root

	configPath, err := chaincfg.LocateDefault(gappPath)
	if err != nil {
		return chaincfg.Config{}, "", err
	}

	res, err := chaincfg.ParseFile(configPath)
	return res, chainId, err
}

func flagGetPath(cmd *cobra.Command) (path string) {
	path, _ = cmd.Flags().GetString("path")
	return
}
