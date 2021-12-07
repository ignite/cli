package starportcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/entrywriter"
	"github.com/tendermint/starport/starport/services/network"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

type ShowType string

const (
	chainShowInfo     ShowType = "info"
	chainShowGenesis  ShowType = "genesis"
	chainShowAccounts ShowType = "accounts"
	chainShowPeers    ShowType = "peers"
)

var (
	showTypes = map[ShowType]struct{}{
		chainShowInfo:     {},
		chainShowGenesis:  {},
		chainShowAccounts: {},
		chainShowPeers:    {},
	}

	chainAccSummaryHeader        = []string{"Address", "Coin"}
	chainValidatorSummaryHeader  = []string{"Peer", "Gentx"}
	chainVestingAccSummaryHeader = []string{"Address", "Vesting", "End Time", "Starting Balance"}
)

// NewNetworkChainShow creates a new chain show command to show
// a chain on SPN.
func NewNetworkChainShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show [info|genesis|accounts|peers] [launch-id]",
		Short: "Show details of a chain",
		RunE:  networkChainShowHandler,
		Args:  cobra.ExactArgs(2),
	}
	return c
}

func networkChainShowHandler(cmd *cobra.Command, args []string) error {
	showType := ShowType(args[0])
	if _, ok := showTypes[showType]; !ok {
		cmd.Usage()
		return fmt.Errorf("invalid arg %s", showType)
	}

	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return err
	}
	defer nb.Cleanup()

	// parse launch ID.
	launchID, err := network.ParseLaunchID(args[1])
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

	switch showType {
	case chainShowAccounts:
		genesisInformation, err := n.GenesisInformation(cmd.Context(), launchID)
		if err != nil {
			return err
		}

		genesisAccEntries := make([][]string, 0)
		for _, acc := range genesisInformation.GenesisAccounts {
			genesisAccEntries = append(genesisAccEntries, []string{
				acc.Address,
				acc.Coins,
			})
		}
		err = entrywriter.MustWrite(os.Stdout, chainAccSummaryHeader, genesisAccEntries...)
		if err != nil {
			return err
		}

		genesisVestingAccEntries := make([][]string, 0)
		for _, acc := range genesisInformation.VestingAccounts {
			genesisVestingAccEntries = append(genesisVestingAccEntries, []string{
				acc.Address,
				acc.Vesting,
				fmt.Sprintf("%d", acc.EndTime),
				acc.StartingBalance,
			})
		}
		err = entrywriter.MustWrite(os.Stdout, chainVestingAccSummaryHeader, genesisVestingAccEntries...)
		if err != nil {
			return err
		}

		genesisValidatorEntries := make([][]string, 0)
		for _, acc := range genesisInformation.GenesisValidators {
			genesisValidatorEntries = append(genesisValidatorEntries, []string{
				acc.Peer,
				string(acc.Gentx),
			})
		}
		err = entrywriter.MustWrite(os.Stdout, chainValidatorSummaryHeader, genesisValidatorEntries...)
		if err != nil {
			return err
		}
	case chainShowGenesis:
		genesisPath, err := c.GenesisPath()
		if err != nil {
			return err
		}
		genesisFile, err := os.ReadFile(genesisPath)
		if err != nil {
			return err
		}
		fmt.Println(string(genesisFile))
	case chainShowInfo:
		home, err := c.Home()
		if err != nil {
			return err
		}
		id, err := c.ID()
		if err != nil {
			return err
		}

		fmt.Printf(`Chain Info:
	-Launch ID: %d
	-Chain ID: %s
	-Name: %s
	-Source URL: %s
	-Hash: %s
	-Home Path: %s`,
			launchID,
			id,
			c.Name(),
			c.SourceURL(),
			c.SourceHash(),
			home,
		)
	case chainShowPeers:
		peers, err := c.Peers()
		if err != nil {
			return err
		}

		var b strings.Builder
		fmt.Fprintf(&b, "Persistent Peers:\n")
		for i, peer := range peers {
			fmt.Fprintf(&b, "%d - %s\n", i, peer)
		}
		fmt.Println(b.String())
	default:
	}

	return nil
}
