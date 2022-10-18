package ignitecmd

import (
	"errors"

	"github.com/spf13/cobra"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

// NewNetworkCoordinatorSet creates a command to set an information in a coordinator profile
func NewNetworkCoordinatorSet() *cobra.Command {
	c := &cobra.Command{
		Use:   "set details|identity|website [value]",
		Short: "Set an information in a coordinator profile",
		Long: `Coordinators on Ignite can set a profile containing a description for the coordinator.
The coordinator set command allows to set information for the coordinator.
The following information can be set:
- details: general information about the coordinator.
- identity: a piece of information to verify the identity of the coordinator with a system like Keybase or Veramo.
- website: website of the coordinator.
`,
		RunE: networkCoordinatorSetHandler,
		Args: cobra.ExactArgs(2),
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkCoordinatorSetHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	var description profiletypes.CoordinatorDescription
	switch args[0] {
	case "details":
		description.Details = args[1]
	case "identity":
		description.Identity = args[1]
	case "website":
		description.Website = args[1]
	default:
		return errors.New("invalid attribute, must provide details, identity, website or security")
	}

	if err := n.SetCoordinatorDescription(cmd.Context(), description); err != nil {
		return err
	}

	return session.Printf("%s Coordinator updated \n", icons.OK)
}
