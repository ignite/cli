package starportcmd

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network"
)

const (
	flagRecover  = "recover"
	flagMnemonic = "mnemomic"
	flagKeyName  = "key-name"
	flagOut      = "out"
)

// NewNetworkChainInit returns a new command to initialize a chain from a published chain ID
func NewNetworkChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init [launch-id]",
		Short: "Initialize a chain from a published chain ID",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainInitHandler,
	}

	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().Bool(flagRecover, false, "Recover chain account from a mnemonic")
	c.Flags().String(flagMnemonic, "", "Mnemonic for recovered account")
	c.Flags().String(flagKeyName, "", "key name for the chain account")

	return c
}

func networkChainInitHandler(cmd *cobra.Command, args []string) error {
	//var (
	//	recover, _     = cmd.Flags().GetBool(flagRecover)
	//	mnemonic, _       = cmd.Flags().GetString(flagMnemonic)
	//	keyName, _ = cmd.Flags().GetString(flagKeyName)
	//)

	s := clispinner.New()
	defer s.Stop()

	var (
		wg sync.WaitGroup
		ev = events.NewBus()
	)
	wg.Add(1)

	defer wg.Wait()
	defer ev.Shutdown()

	go printEvents(&wg, ev, s)

	// parse launch ID
	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing launchID: %s", err.Error())
	}

	nb, err := newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		return err
	}

	// initialize the blockchain from the launch ID
	initOptions := initOptionWithHomeFlag(cmd, []network.InitOption{network.MustNotInitializedBefore()})
	sourceOption := network.SourceLaunchID(launchID)
	_, err = nb.Blockchain(cmd.Context(), sourceOption, initOptions...)
	if err != nil {
		return err
	}

	return nil
}
