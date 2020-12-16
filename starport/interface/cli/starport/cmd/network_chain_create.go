package starportcmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xurl"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

const (
	flagRev = "rev"
	spnRev  = "spn"
)

// NewNetworkChainCreate creates a new chain create command to create
// a new network.
func NewNetworkChainCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [chain] [source]",
		Short: "Create a new network",
		RunE:  networkChainCreateHandler,
	}
	c.Flags().String(flagRev, "", "Git revision of the source code")
	return c
}

func networkChainCreateHandler(cmd *cobra.Command, args []string) error {
	// collect required values.
	var (
		chainID string
		source  string
		rev, _  = cmd.Flags().GetString(flagRev)
	)

	if len(args) >= 1 {
		chainID = args[0]
	}

	if len(args) >= 2 {
		source = args[1]
	}

	var questions []cliquiz.Question
	if chainID == "" {
		questions = append(questions, cliquiz.NewQuestion("Chain ID", &chainID, cliquiz.Required()))
	}
	if source == "" {
		questions = append(questions, cliquiz.NewQuestion("Git repository of the chain's source code (local or remote)", &source, cliquiz.Required()))
	}
	if len(questions) > 0 {
		if err := cliquiz.Ask(questions...); err != nil {
			return err
		}
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
	if _, err := nb.ShowChain(cmd.Context(), chainID); err == nil {
		s.Stop()

		return fmt.Errorf("chain with id %q already exists", chainID)
	}

	initChain := func() (*networkbuilder.Blockchain, error) {
		if xurl.IsLocalPath(source) {
			return nb.InitBlockchainFromPath(cmd.Context(), chainID, source, true)
		}
		return nb.InitBlockchainFromURL(cmd.Context(), chainID, source, rev, true)
	}

	initChainWithSPNFallback := func() (*networkbuilder.Blockchain, error) {
		nb, err := initChain()
		if err == chain.ErrCouldntLocateConfig && rev != spnRev {
			fmt.Printf("%s Default branch does not have a config.yml, trying with 'spn' branch...\n", color.New(color.FgYellow).SprintFunc()("ℹ"))
			rev = "spn"
			return initChain()
		}
		return nb, err
	}

	// init the chain.
	blockchain, err := initChainWithSPNFallback()

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

		blockchain, err = initChainWithSPNFallback()
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

	s.SetText("Submiting...")
	s.Start()

	// create blockchain.
	if err := blockchain.Create(cmd.Context()); err != nil {
		return err
	}
	s.Stop()

	fmt.Println("\n🌐  Network submited")
	return nil
}
