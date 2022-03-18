package starportcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/yaml"
)

// NewNetworkProfile returns a new command to show the address profile info on Starport Network.
func NewNetworkProfile() *cobra.Command {
	c := &cobra.Command{
		Use:   "profile",
		Short: "Show the address profile info",
		Args:  cobra.NoArgs,
		RunE:  networkProfileHandler,
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkProfileHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	profile, err := n.Profile(cmd.Context())
	if err != nil {
		return err
	}

	profileInfo, err := yaml.Marshal(cmd.Context(), profile)
	if err != nil {
		return err
	}

	nb.Cleanup()
	fmt.Print(profileInfo)
	return nil
}
