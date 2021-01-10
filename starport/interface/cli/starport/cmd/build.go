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
	// Check if custom home is provided
	_, _, err := getHomeFlags(cmd)
	if err != nil {
		return err
	}
	// TODO: fill the command

	c, err := chain.New(appPath, chain.LogLevel(logLevel(cmd)), chain.KeyringBackend(chaincmd.KeyringBackendTest))
	if err != nil {
		return err
	}
	return c.Build(cmd.Context())
}
