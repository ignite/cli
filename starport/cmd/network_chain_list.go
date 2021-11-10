package starportcmd

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network"
)

var LaunchSummaryHeader = []string{"Launch ID", "Chain ID", "Source"}

// LaunchSummary holds summarized information about a chain launch
type LaunchSummary struct {
	LaunchID string
	ChainID  string
	Source   string
}

// NewNetworkChainList returns a new command to list all published chains on Starport Network
func NewNetworkChainList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "List published chains",
		Args:  cobra.NoArgs,
		RunE:  networkChainListHandler,
	}
	c.Flags().String(flagFrom, cosmosaccount.DefaultAccount, "Account name to use for sending transactions to SPN")
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func networkChainListHandler(cmd *cobra.Command, args []string) error {
	// TODO: use routine from network init PR
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

	nb, err := newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		return err
	}

	chains, err := nb.ChainLaunches(cmd.Context())
	if err != nil {
		return err
	}
	sums := launchSummaries(chains)

	s.Stop()
	renderLaunchSummaries(sums, os.Stdout)

	return nil
}

// launchSummaries returns the list of launch summaries from the list of chain launches
func launchSummaries(chains []launchtypes.Chain) (sums []LaunchSummary) {
	for _, chain := range chains {
		sums = append(sums, LaunchSummary{
			LaunchID: fmt.Sprintf("%d", chain.LaunchID),
			ChainID:  chain.GenesisChainID,
			Source:   chain.SourceURL,
		})
	}
	return sums
}

// renderLaunchSummaries writes into the provided out, the list of summarized launches
func renderLaunchSummaries(launchSummaries []LaunchSummary, out io.Writer) {
	launchTable := tablewriter.NewWriter(out)
	launchTable.SetHeader(LaunchSummaryHeader)

	for _, sum := range launchSummaries {
		launchTable.Append([]string{sum.LaunchID, sum.ChainID, sum.Source})
	}

	launchTable.Render()
}
