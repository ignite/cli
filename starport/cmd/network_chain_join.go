package starportcmd

import (
	"errors"
	"fmt"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
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
	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	var (
		gentxPath, _ = cmd.Flags().GetString(flagGentx)
		amount, _    = cmd.Flags().GetString(flagAmount)
	)

	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing launchID: %s", err.Error())
	}

	var (
		wg sync.WaitGroup
		ev = events.NewBus()
		s  = clispinner.New()
	)
	defer s.Stop()
	wg.Add(1)
	defer wg.Wait()
	defer ev.Shutdown()

	go printEvents(&wg, ev, s)

	if gentxPath == "" {
		chainHome, exist, err := checkChainHomeExist(launchID)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("the chain home not exist")
		}
		gentxPath = chainHome
	}

	nb, err := newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		return err
	}

	// initialize the blockchain from the launch ID
	initOptions := initOptionWithHomeFlag(cmd, []network.InitOption{})
	sourceOption := network.SourceLaunchID(launchID)
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	if err := blockchain.Init(cmd.Context()); err != nil {
		return err
	}

	gentx, err := network.ParseGentx(gentxPath)

	fmt.Printf("%s Network joined\n", clispinner.OK)
	return nil
}
