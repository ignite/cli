package starportcmd

import (
	"context"
	"fmt"
	"io"
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

	chainAccSummaryHeader = []string{"Genesis Account", "Coins"}
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

	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())

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
	case chainShowGenesis:
		return printChainGenesis(c)
	case chainShowInfo:
		return printChainInfo(c, launchID)
	case chainShowAccounts:
		return printChainAccounts(cmd.Context(), n, launchID, os.Stdout)
	case chainShowPeers:
		return printChainPeers(cmd.Context(), n, launchID)
	}
	return nil
}

func printChainGenesis(c network.Chain) error {
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}
	if _, err = os.Stat(genesisPath); os.IsNotExist(err) {
		return fmt.Errorf("chain genesis not initialized: %s", genesisPath)
	}
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return err
	}
	fmt.Println(string(genesisFile))
	return nil
}

func printChainInfo(c *networkchain.Chain, launchID uint64) error {
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
	return nil
}

func printChainAccounts(ctx context.Context, n network.Network, launchID uint64, out io.Writer) error {
	genesisInformation, err := n.GenesisInformation(ctx, launchID)
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
	return entrywriter.MustWrite(out, chainAccSummaryHeader, genesisAccEntries...)
}

func printChainPeers(ctx context.Context, n network.Network, launchID uint64) error {
	genesisInformation, err := n.GenesisInformation(ctx, launchID)
	if err != nil {
		return err
	}

	peers := make([]string, 0)
	for _, acc := range genesisInformation.GenesisValidators {
		peers = append(peers, acc.Peer)
	}
	fmt.Printf("Persistent Peers: %s", strings.Join(peers, ","))
	return nil
}
