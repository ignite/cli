package ignitecmd

import (
	"errors"

	"github.com/spf13/cobra"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

// NewNetworkValidatorSet creates a command to set an information in a validator profile
func NewNetworkValidatorSet() *cobra.Command {
	c := &cobra.Command{
		Use:   "set details|identity|website|security [value]",
		Short: "Set an information in a validator profile",
		Long: `Validators on Ignite can set a profile containing a description for the validator.
The validator set command allows to set information for the validator.
The following information can be set:
- details: general information about the validator.
- identity: piece of information to verify identity of the validator with a system like Keybase of Veramo.
- website: website of the validator.
- security: security contact for the validator.
`,
		RunE: networkValidatorSetHandler,
		Args: cobra.ExactArgs(2),
	}
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkValidatorSetHandler(cmd *cobra.Command, args []string) error {
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

	var validator profiletypes.Validator
	switch args[0] {
	case "details":
		validator.Description.Details = args[1]
	case "identity":
		validator.Description.Identity = args[1]
	case "website":
		validator.Description.Website = args[1]
	case "security":
		validator.Description.SecurityContact = args[1]
	default:
		return errors.New("invalid attribute, must provide details, identity, website or security")
	}

	if err := n.SetValidatorDescription(cmd.Context(), validator); err != nil {
		return err
	}

	return session.Printf("%s Validator updated \n", icons.OK)
}
