package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

const flagForceReset = "force-reset"
const flagResetOnce = "reset-once"
const flagConfig = "config"

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
	c.Flags().BoolP(flagForceReset, "f", false, "Force reset of the app state on start and every source change")
	c.Flags().BoolP(flagResetOnce, "r", false, "Reset of -the app state on first start")
	c.Flags().StringP(flagConfig, "c", "", "Starport config file (default: ./config.yml)")

	return c
}

func serveHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	// check if custom config is defined
	config, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	if config != "" {
		chainOption = append(chainOption, chain.ConfigPath(config))
	}

	// create the chain
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
	resetOnce, err := cmd.Flags().GetBool(flagResetOnce)
	if err != nil {
		return err
	}
	if resetOnce {
		serveOptions = append(serveOptions, chain.ServeResetOnce())
	}

	return c.Serve(cmd.Context(), serveOptions...)
}
