package starportcmd

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagGentx  = "gentx"
	flagAmount = "amount"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [launch-id]",
		Short: "Join to a network as a validator by launch id",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainJoinHandler,
	}

	c.Flags().String(flagGentx, "", "Path to a gentx json file")
	c.Flags().String(flagAmount, "", "If is provided sends the \"create account\" message")

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	// initialize network common methods
	nb, s, shutdown, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer shutdown()

	// parse the amount
	amountArg, _ := cmd.Flags().GetString(flagAmount)
	amount, err := sdk.ParseCoinNormalized(amountArg)
	if err != nil {
		return errors.Wrap(err, "error parsing amount")
	}

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// parse the home path
	home, err := getLaunchIDHome(cmd, launchID)
	if err != nil {
		return err
	}

	// parse the gentx and check if it exist
	gentxPath := getGentxPath(cmd, home)
	if err != nil {
		return errors.Wrap(err, "error parsing gentx path")
	}
	info, gentx, err := network.ParseGentx(gentxPath)
	if err != nil {
		return err
	}

	// get the peer public address for the validator
	peer, err := askPublicAddress(s)
	if err != nil {
		return err
	}

	// create the message to add the validator into the SPN
	valMsg, err := nb.CreateValidatorRequestMsg(
		cmd.Context(),
		launchID,
		peer,
		info.DelegatorAddress,
		gentx,
		info.PubKey,
		info.SelfDelegation,
	)
	if err != nil {
		return err
	}

	// create the message to add the account into the SPN
	accMsg, err := nb.CreateAccountRequestMsg(
		cmd.Context(),
		home,
		launchID,
		amount,
	)
	if err != nil {
		return err
	}

	result, err := nb.Join(valMsg, accMsg)
	if err != nil {
		return err
	}

	s.Stop()
	fmt.Printf("%s Network joined\n%s", clispinner.OK, result)

	return nil
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
