package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
)

// NewNetworkValidatorShow creates a command to show validator information
func NewNetworkValidatorShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [address]",
		Short: "Show a validator profile",
		RunE:  networkValidatorShowHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkValidatorShowHandler(cmd *cobra.Command, args []string) error {
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

	validator, err := n.Validator(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	// convert the request object to YAML to be more readable
	// and convert the byte array fields to string.
	validatorYaml, err := yaml.Marshal(cmd.Context(), struct {
		Identity string
		Details  string
		Website  string
		Security string
	}{
		validator.Identity,
		validator.Details,
		validator.Website,
		validator.SecurityContact,
	})
	if err != nil {
		return err
	}

	return session.Println(validatorYaml)
}
