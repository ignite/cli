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
	"github.com/tendermint/starport/starport/pkg/gentx"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagGentx = "gentx"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [launch-id] [amount]",
		Short: "Join to a network as a validator by launch id",
		Args:  cobra.ExactArgs(2),
		RunE:  networkChainJoinHandler,
	}
	c.Flags().String(flagGentx, "", "Path to a gentx json file")
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

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing launchID")
	}
	if launchID == 0 {
		return errors.New("launch ID must be greater than 0")
	}

	// parse the amount
	amount, err := sdk.ParseCoinNormalized(args[1])
	if err != nil {
		return errors.Wrap(err, "error parsing amount")
	}

	// parse the home path
	home := getHome(cmd)
	if home == "" {
		home, err = network.ChainHome(launchID)
		if err != nil {
			return err
		}
	}

	// parse the gentx and check if it exists
	gentxPath, _ := cmd.Flags().GetString(flagGentx)
	customGentx := true
	if gentxPath == "" {
		customGentx = false
		gentxPath = network.Gentx(home)
	}

	if err != nil {
		return errors.Wrap(err, "error parsing gentx path")
	}
	info, gentx, err := gentx.ParseGentx(gentxPath)
	if err != nil {
		return err
	}

	// get the gentx account for the validator to sign the tx
	valAcc, err := getAccountByAddress(cmd, info.DelegatorAddress)
	if err != nil {
		return err
	}

	// get the peer public address for the validator
	peer, err := askPublicAddress(s)
	if err != nil {
		return err
	}

	// create the message to add the account if needed
	accReqID, accAutoApproved, created, err := nb.CreateAccountRequestMsg(cmd.Context(),
		home,
		customGentx,
		launchID,
		amount)
	if err != nil {
		return err
	}
	s.Stop()
	if accAutoApproved {
		fmt.Printf("%s Account %s added to the network by the coordinator!\n",
			clispinner.OK, getFrom(cmd))
	} else if created {
		fmt.Printf("%s Request %d to add account to the network has been submitted!\n",
			clispinner.OK, accReqID)
	}
	s.Start()

	// create the message to add the validator
	valReqID, valAutoApproved, err := nb.Join(cmd.Context(),
		launchID,
		peer,
		valAcc.Name,
		info.DelegatorAddress,
		gentx,
		info.PubKey,
		info.SelfDelegation,
	)
	if err != nil {
		return err
	}

	s.Stop()
	if valAutoApproved {
		fmt.Printf("%s Validator %s added to the network by the coordinator!\n",
			clispinner.OK, info.DelegatorAddress)
	} else {
		fmt.Printf("%s Request %d to join the network as a validator has been submitted!\n",
			clispinner.OK, valReqID)
	}

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
