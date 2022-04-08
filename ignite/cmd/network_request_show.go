package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/ignite-hq/cli/ignite/pkg/yaml"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/yaml"
	"github.com/tendermint/starport/starport/services/network"
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
	// initialize network common methods
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

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
		"$.content.content.genesisValidator.genTx",
		"$.content.content.genesisValidator.consPubKey",
	)
	if err != nil {
		return err
	}

	nb.Spinner.Stop()
	fmt.Println(requestYaml)
	return nil
}
