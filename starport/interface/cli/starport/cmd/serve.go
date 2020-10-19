package starportcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/services/chain"
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
	path, err := gomodulepath.Parse(getModule(appPath))
	if err != nil {
		return err
	}
	app := chain.App{
		Name: path.Root,
		Path: appPath,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		cancel()
	}()

	s, err := chain.New(app, verbose)
	if err != nil {
		return err
	}
	err = s.Serve(ctx)
	if err == context.Canceled {
		fmt.Println("aborted")
		return nil
	}
	return err
}
