package starportcmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

const (
	flagChainID = "chain"
	flagSource  = "source"
)

func NewNetworkChainCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "Create a new network",
		RunE:  networkChainCreateHandler,
	}
	c.Flags().String(flagChainID, "", "Chain ID")
	c.Flags().String(flagSource, "", "Git repository of the chain's source code (local or remote)")
	return c
}

func networkChainCreateHandler(cmd *cobra.Command, args []string) error {
	values, err := cliquiz.ValuesFromFlagsOrAsk(cmd.Flags(), "Please, provide:\n", []string{
		flagChainID,
		flagSource,
	})
	if err != nil {
		return err
	}

	s := clispinner.New()
	defer s.Stop()

	ev := events.NewBus()
	go printEvents(ev, s)

	nb, err := newNetworkBuilder(networkbuilder.CollectEvents(ev))
	if err != nil {
		return err
	}

	// check if chain already exists on SPN.
	if _, err := nb.ShowChain(cmd.Context(), values[flagChainID]); err == nil {
		s.Stop()

		return fmt.Errorf("chain with id %q already exists", values[flagChainID])
	}

	initChain := func() (*networkbuilder.Blockchain, error) {
		if xurl.IsLocalPath(values[flagSource]) {
			return nb.InitBlockchainFromPath(cmd.Context(), values[flagChainID], values[flagSource], true)
		}
		return nb.InitBlockchainFromURL(cmd.Context(), values[flagChainID], values[flagSource], "", true)
	}

	// init the chain.
	blockchain, err := initChain()

	// ask to delete data dir for the chain if already exists on the fs.
	var e *networkbuilder.DataDirExistsError
	if errors.As(err, &e) {
		s.Stop()

		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Data directory for %q blockchain already exists: %s. Would you like to overwrite it",
				e.ID,
				e.Home,
			),
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err != nil {
			fmt.Println("said no")
			return nil
		}

		if err := os.RemoveAll(e.Home); err != nil {
			return err
		}

		s.Start()

		blockchain, err = initChain()
	}

	s.Stop()

	if err == context.Canceled {
		fmt.Println("aborted")
		return nil
	}
	if err != nil {
		return err
	}
	defer blockchain.Cleanup()

	info, err := blockchain.Info()
	if err != nil {
		return err
	}

	// ask to confirm Genesis.
	prettyGenesis, err := info.Genesis.Pretty()
	if err != nil {
		return err
	}

	fmt.Printf("\nGenesis: \n\n%s\n\n", prettyGenesis)

	prompt := promptui.Prompt{
		Label:     "Proceed with the Genesis configuration above",
		IsConfirm: true,
	}
	if _, err := prompt.Run(); err != nil {
		fmt.Println("said no")
		return nil
	}

	// create blockchain.
	if err := blockchain.Create(cmd.Context()); err != nil {
		return err
	}

	fmt.Println("\n🌐 Network submited")
	return nil
}
