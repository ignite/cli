package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	relayerconfig "github.com/ignite-hq/cli/ignite/pkg/relayer/config"
)

// NewRelayerReset returns a new relayer reset command to
// reset all relayer config file
func NewRelayerReset() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset the relayer config",
		RunE:  relayerResetHandler,
	}
}

func relayerResetHandler(cmd *cobra.Command, args []string) (err error) {
	session := cliui.New()
	session.StartSpinner("Resetting relayer...")
	defer session.Cleanup()
	return relayerconfig.Delete()
}
