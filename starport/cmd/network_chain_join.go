package starportcmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/ignite-hq/cli/starport/pkg/cliquiz"
	"github.com/ignite-hq/cli/starport/pkg/clispinner"
	"github.com/ignite-hq/cli/starport/pkg/xchisel"
	"github.com/ignite-hq/cli/starport/services/network"
	"github.com/ignite-hq/cli/starport/services/network/networkchain"
)

const (
	flagGentx = "gentx"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [launch-id] [amount]",
		Short: "Request to join a network as a validator",
		Args:  cobra.ExactArgs(2),
		RunE:  networkChainJoinHandler,
	}
	c.Flags().String(flagGentx, "", "Path to a gentx json file")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID.
	launchID, err := network.ParseLaunchID(args[0])
	if err != nil {
		return err
	}

	// parse the amount.
	amount, err := sdk.ParseCoinNormalized(args[1])
	if err != nil {
		return errors.Wrap(err, "error parsing amount")
	}

	gentxPath, _ := cmd.Flags().GetString(flagGentx)

	// get the peer public address for the validator.
	publicAddr, err := askPublicAddress(nb.Spinner)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	// create the message to add the validator.
	return n.Join(cmd.Context(), c, launchID, amount, publicAddr, gentxPath)
}

// askPublicAddress prepare questions to interactively ask for a publicAddress
// when peer isn't provided and not running through chisel proxy.
func askPublicAddress(s *clispinner.Spinner) (publicAddress string, err error) {
	s.Stop()
	defer s.Start()

	options := []cliquiz.Option{
		cliquiz.Required(),
	}
	if !xchisel.IsEnabled() {
		ip, _ := ipify.GetIp()
		if err == nil {
			options = append(options, cliquiz.DefaultAnswer(fmt.Sprintf("%s:26656", ip)))
		}
	}
	questions := []cliquiz.Question{cliquiz.NewQuestion(
		"Peer's address",
		&publicAddress,
		options...,
	)}
	return publicAddress, cliquiz.Ask(questions...)
}
