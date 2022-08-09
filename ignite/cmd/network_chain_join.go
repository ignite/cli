package ignitecmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networkchain"
)

const (
	flagGentx           = "gentx"
	flagAmount          = "amount"
	flagPeerAddress     = "peer-address"
	flagHidePeerAddress = "hide-peer-address"
	flagNodeID          = "node-id"
)

// NewNetworkChainJoin creates a new chain join command to join
// to a network as a validator.
func NewNetworkChainJoin() *cobra.Command {
	c := &cobra.Command{
		Use:   "join [launch-id]",
		Short: "Request to join a network as a validator",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainJoinHandler,
	}

	c.Flags().String(flagGentx, "", "Path to a gentx json file")
	c.Flags().String(flagAmount, "", "Amount of coins for account request")
	c.Flags().String(flagPeerAddress, "", "Specify validator node's peer address ip and port")
	c.Flags().String(flagNodeID, "", "Specify validator node's id")
	c.Flags().Bool(flagHidePeerAddress, false, "Join without peer address")
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetCheckDependencies())

	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	var (
		amount, _            = cmd.Flags().GetString(flagAmount)
		customPeerAddress, _ = cmd.Flags().GetString(flagPeerAddress)
		nodeID, _            = cmd.Flags().GetString(flagNodeID)
		hidePeerAddress, _   = cmd.Flags().GetBool(flagHidePeerAddress)
		gentxPath, _         = cmd.Flags().GetString(flagGentx)
	)

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID.
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
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

	var networkOptions []networkchain.Option

	if flagGetCheckDependencies(cmd) {
		networkOptions = append(networkOptions, networkchain.CheckDependencies())
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch), networkOptions...)
	if err != nil {
		return err
	}

	var joinOptions []network.JoinOption

	// parse the peer
	if !hidePeerAddress {
		var publicAddr string
		if customPeerAddress != "" {
			publicAddr = customPeerAddress
		} else {
			publicAddr, err = askPublicAddress(cmd, session)
			if err != nil {
				return err
			}
		}
		joinOptions = append(joinOptions, network.WithPublicAddress(publicAddr))
	}

	if nodeID != "" {
		joinOptions = append(joinOptions, network.WithNodeID(nodeID))
	}

	// use the default gentx path from chain home if not provided
	if gentxPath == "" {
		gentxPath, err = c.DefaultGentxPath()
		if err != nil {
			return err
		}
	} else {
		// if a custom gentx is provided, we initialize the chain home in order to check accounts
		if err := c.Init(cmd.Context(), cacheStorage); err != nil {
			return err
		}
	}

	if amount != "" {
		// parse the amount.
		amountCoins, err := sdk.ParseCoinsNormalized(amount)
		if err != nil {
			return errors.Wrap(err, "error parsing amount")
		}
		joinOptions = append(joinOptions, network.WithAccountRequest(amountCoins))
	} else {
		if !getYes(cmd) {
			question := fmt.Sprintf(
				"You haven't set the --%s flag and therefore an account request won't be submitted. Do you confirm",
				flagAmount,
			)
			if err := session.AskConfirm(question); err != nil {
				if errors.Is(err, cliui.ErrRejectConfirmation) {
					return session.PrintSaidNo()
				}
				return err
			}
		}

		_ = session.Printf("%s %s\n", icons.Info, "Account request won't be submitted")
	}

	// create the message to add the validator.
	return n.Join(cmd.Context(), c, launchID, gentxPath, joinOptions...)
}

// askPublicAddress prepare questions to interactively ask for a publicAddress
// when peer isn't provided and not running through chisel proxy.
func askPublicAddress(cmd *cobra.Command, session cliui.Session) (publicAddress string, err error) {
	var options []cliquiz.Option
	// even if GetIp fails we won't handle the error because we don't want to interrupt a join process.
	// just in case if GetIp fails user should enter his address manually
	ip, err := ipify.GetIp()
	if err == nil {
		options = append(options, cliquiz.Suggestion(fmt.Sprintf("%s:26656", ip)))
	}

	publicAddress, err = session.AskString("Peer's address", options...)
	if err != nil {
		return "", err
	}
	if publicAddress == "" && !getYes(cmd) {
		question := "You didn't specify your peer's address. Do you want to proceed without it"
		if err := session.AskConfirm(question); err != nil {
			return "", err
		}
	}

	return publicAddress, nil
}
