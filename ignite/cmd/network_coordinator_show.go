package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
)

// NewNetworkCoordinatorShow creates a command to show coordinator information
func NewNetworkCoordinatorShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [address]",
		Short: "Show a coordinator profile",
		RunE:  networkCoordinatorShowHandler,
		Args:  cobra.ExactArgs(1),
	}
	return c
}

func networkCoordinatorShowHandler(cmd *cobra.Command, args []string) error {
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

	coordinator, err := n.Coordinator(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	// convert the request object to YAML to be more readable
	// and convert the byte array fields to string.
	coordinatorYaml, err := yaml.Marshal(cmd.Context(), struct {
		Identity string
		Details  string
		Website  string
	}{
		coordinator.Identity,
		coordinator.Details,
		coordinator.Website,
	})
	if err != nil {
		return err
	}

	return session.Println(coordinatorYaml)
}
