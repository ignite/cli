package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

const (
	flagMetadata = "metadata"
)

// NewNetworkProjectPublish returns a new command to publish a new projects on Ignite.
func NewNetworkProjectPublish() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name] [total-supply]",
		Short: "Create a project",
		Args:  cobra.ExactArgs(2),
		RunE:  networkProjectPublishHandler,
	}
	c.Flags().String(flagMetadata, "", "Add a metadata to the chain")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkProjectPublishHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	totalSupply, err := sdk.ParseCoinsNormalized(args[1])
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	metadata, _ := cmd.Flags().GetString(flagMetadata)
	projectID, err := n.CreateProject(cmd.Context(), args[0], metadata, totalSupply)
	if err != nil {
		return err
	}

	return session.Printf("%s Project ID: %d \n", icons.Bullet, projectID)
}
