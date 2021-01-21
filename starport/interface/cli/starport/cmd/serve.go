package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

const flagForceReset = "force-reset"

var appPath string

// NewServe creates a new serve command to serve a blockchain.
func NewServe() *cobra.Command {
	c := &cobra.Command{
		Use:   "serve",
		Short: "Launches a reloading server",
		Args:  cobra.ExactArgs(0),
		RunE:  serveHandler,
	}
	c.Flags().AddFlagSet(flagSetHomes())
	c.Flags().StringVarP(&appPath, "path", "p", "", "Path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	c.Flags().BoolP(flagForceReset, "r", false, "Force reset of the app state")
	return c
}

func serveHandler(cmd *cobra.Command, args []string) error {
	// create the chain
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}
	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	// serve the chain
	var serveOptions []chain.ServeOption
	forceUpdate, err := cmd.Flags().GetBool(flagForceReset)
	if err != nil {
		return err
	}
	if forceUpdate {
		serveOptions = append(serveOptions, chain.ServeForceReset())
	}

	return c.Serve(cmd.Context(), serveOptions...)
}
