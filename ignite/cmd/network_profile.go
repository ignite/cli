package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
)

// NewNetworkProfile returns a new command to show the address profile info on Starport Network.
func NewNetworkProfile() *cobra.Command {
	c := &cobra.Command{
		Use:   "profile [campaign-id]",
		Short: "Show the address profile info",
		Args:  cobra.RangeArgs(0, 1),
		RunE:  networkProfileHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkProfileHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	var (
		campaignID uint64
	)
	if len(args) > 0 {
		campaignID, err = network.ParseID(args[0])
		if err != nil {
			return err
		}
	}

	profile, err := n.Profile(cmd.Context(), campaignID)
	if err != nil {
		return err
	}

	profileInfo, err := yaml.Marshal(cmd.Context(), profile)
	if err != nil {
		return err
	}
	return session.Println(profileInfo)
}
