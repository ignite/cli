package starportcmd

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/yaml"
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
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkRequestShowHandler(cmd *cobra.Command, args []string) error {
	// initialize network common methods
	nb, s, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// parse request ID
	requestID, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing requestID")
	}

	request, err := nb.FetchRequest(cmd.Context(), launchID, requestID)
	if err != nil {
		return err
	}

	// convert the request object to yaml

	requestYaml, err := yaml.Parse(cmd.Context(), request)

	s.Stop()
	fmt.Println(requestYaml)

	return nil
}
