package ignitecmd

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/yaml"
	"github.com/ignite/cli/ignite/services/network"
)

// NewNetworkRequestShow creates a new request show command to show
// requests details for a chain
func NewNetworkRequestShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [launch-id] [request-id]",
		Short: "Show pending requests details",
		RunE:  networkRequestShowHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkRequestShowHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// parse request ID
	requestID, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing requestID")
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	request, err := n.Request(cmd.Context(), launchID, requestID)
	if err != nil {
		return err
	}

	// convert the request object to YAML to be more readable
	// and convert the byte array fields to string.
	requestYaml, err := yaml.Marshal(cmd.Context(), request,
		"$.Content.content.genesisValidator.genTx",
		"$.Content.content.genesisValidator.consPubKey",
	)
	if err != nil {
		return err
	}

	session.StopSpinner()

	return session.Println(requestYaml)
}
