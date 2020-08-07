package starportcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/xos"
	starportserve "github.com/tendermint/starport/starport/services/serve"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
)

var appPath string

func NewServe() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Launches a reloading server",
		Args:  cobra.ExactArgs(0),
		RunE:  serveHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func serveHandler(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	appName, _ := getAppAndModule(appPath)
	app := starportserve.App{
		Name: appName,
		Path: appPath,
	}

	confFile, err := xos.OpenFirst(starportconf.FileNames...)
	if err != nil {
		return errors.Wrap(err, "config file cannot be found")
	}
	defer confFile.Close()
	conf, err := starportconf.Parse(confFile)
	if err != nil {
		return errors.Wrap(err, "config file is not valid")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		cancel()
	}()

	err = starportserve.Serve(ctx, app, conf, verbose)
	if err == context.Canceled {
		fmt.Println("aborted")
		return nil
	}
	return err
}
