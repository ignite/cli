package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

// NewBuild returns a new build command to build a blockchain app.
func NewBuild() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Builds an app and installs binaries",
		Args:  cobra.ExactArgs(0),
		RunE:  buildHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func buildHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	// Check if custom home is provided
	home, cliHome, err := getHomeFlags(cmd)
	if err != nil {
		return err
	}
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	if cliHome != "" {
		chainOption = append(chainOption, chain.CLIHomePath(cliHome))
	}

	c, err := chain.New(appPath, chainOption...)
	if err != nil {
		return err
	}
	return c.Build(cmd.Context())
}
