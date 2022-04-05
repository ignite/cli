package starportcmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ignite-hq/cli/starport/pkg/cosmosutil"
	"github.com/ignite-hq/cli/starport/pkg/entrywriter"
	"github.com/ignite-hq/cli/starport/pkg/yaml"
	"github.com/ignite-hq/cli/starport/services/network"
	"github.com/ignite-hq/cli/starport/services/network/networkchain"
	"github.com/ignite-hq/cli/starport/services/network/networktypes"
	"github.com/spf13/cobra"
)

var (
	chainGenesisValSummaryHeader = []string{"Genesis Validator", "Self Delegation", "Peer"}
	chainGenesisAccSummaryHeader = []string{"Genesis Account", "Coins"}
	chainVestingAccSummaryHeader = []string{"Vesting Account", "Total Balance", "Vesting", "EndTime"}
)

// NewNetworkChainShow creates a new chain show
// command to show a chain details on SPN.
func NewNetworkChainShow() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "Show details of a chain",
	}
	c.AddCommand(
		newNetworkChainShowInfo(),
		newNetworkChainShowGenesis(),
		newNetworkChainShowAccounts(),
		newNetworkChainShowValidators(),
		newNetworkChainShowPeers(),
	)
	c.PersistentFlags().AddFlagSet(flagNetworkFrom())
	c.PersistentFlags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkChainLaunch(cmd *cobra.Command, args []string) (NetworkBuilder, uint64, error) {
	nb, err := newNetworkBuilder(cmd)
	if err != nil {
		return nb, 0, err
	}
	// parse launch ID.
	launchID, err := network.ParseLaunchID(args[0])
	if err != nil {
		return nb, launchID, err
	}
	return nb, launchID, err
}

func newNetworkChainShowInfo() *cobra.Command {
	c := &cobra.Command{
		Use:   "info [launch-id]",
		Short: "Show info details of the chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nb, launchID, err := networkChainLaunch(cmd, args)
			if err != nil {
				return err
			}
			defer nb.Cleanup()
			n, err := nb.Network()
			if err != nil {
				return err
			}

			chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
			if err != nil {
				return err
			}

			var genesis []byte
			if chainLaunch.GenesisURL != "" {
				genesis, _, err = cosmosutil.GenesisAndHashFromURL(cmd.Context(), chainLaunch.GenesisURL)
				if err != nil {
					return err
				}
			}
			chainInfo := struct {
				Chain   networktypes.ChainLaunch `json:"Chain"`
				Genesis []byte                   `json:"Genesis"`
			}{
				Chain:   chainLaunch,
				Genesis: genesis,
			}
			info, err := yaml.Marshal(cmd.Context(), chainInfo, "$.Genesis")
			if err != nil {
				return err
			}
			nb.Spinner.Stop()
			fmt.Print(info)
			return nil
		},
	}
	return c
}

func newNetworkChainShowGenesis() *cobra.Command {
	c := &cobra.Command{
		Use:   "genesis [launch-id]",
		Short: "Show the chain genesis file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nb, launchID, err := networkChainLaunch(cmd, args)
			if err != nil {
				return err
			}
			defer nb.Cleanup()
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
			genesisPath, err := c.GenesisPath()
			if err != nil {
				return err
			}

			// check if the genesis already exist
			if _, err = os.Stat(genesisPath); os.IsNotExist(err) {
				// fetch the information to construct genesis
				genesisInformation, err := n.GenesisInformation(cmd.Context(), launchID)
				if err != nil {
					return err
				}

				// create the chain into a temp dir
				home := filepath.Join(os.TempDir(), "spn/temp", chainLaunch.ChainID)
				c.SetHome(home)
				defer os.RemoveAll(home)

				err = c.Prepare(cmd.Context(), genesisInformation)
				if err != nil {
					return err
				}

				// get the new genesis path
				genesisPath, err = c.GenesisPath()
				if err != nil {
					return err
				}
			}
			genesisFile, err := os.ReadFile(genesisPath)
			if err != nil {
				return err
			}

			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, genesisFile, "", "    "); err != nil {
				return err
			}
			nb.Spinner.Stop()
			fmt.Printf("Genesis: \n%s", prettyJSON.String())
			return nil
		},
	}
	return c
}

func newNetworkChainShowAccounts() *cobra.Command {
	c := &cobra.Command{
		Use:   "accounts [launch-id]",
		Short: "Show all vesting and genesis accounts of the chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nb, launchID, err := networkChainLaunch(cmd, args)
			if err != nil {
				return err
			}
			defer nb.Cleanup()
			n, err := nb.Network()
			if err != nil {
				return err
			}

			accountSummary := bytes.NewBufferString("")

			// get all chain genesis accounts
			genesisAccs, err := n.GenesisAccounts(cmd.Context(), launchID)
			if err != nil {
				return err
			}
			genesisAccEntries := make([][]string, 0)
			for _, acc := range genesisAccs {
				genesisAccEntries = append(genesisAccEntries, []string{
					acc.Address,
					acc.Coins,
				})
			}
			if len(genesisAccEntries) > 0 {
				if err = entrywriter.MustWrite(
					accountSummary,
					chainGenesisAccSummaryHeader,
					genesisAccEntries...,
				); err != nil {
					return err
				}
			}

			// get all chain vesting accounts
			vestingAccs, err := n.VestingAccounts(cmd.Context(), launchID)
			if err != nil {
				return err
			}
			genesisVestingAccEntries := make([][]string, 0)
			for _, acc := range vestingAccs {
				genesisVestingAccEntries = append(genesisVestingAccEntries, []string{
					acc.Address,
					acc.TotalBalance,
					acc.Vesting,
					strconv.FormatInt(acc.EndTime, 10),
				})
			}
			if len(genesisVestingAccEntries) > 0 {
				if err = entrywriter.MustWrite(
					accountSummary,
					chainVestingAccSummaryHeader,
					genesisVestingAccEntries...,
				); err != nil {
					return err
				}
			}
			nb.Spinner.Stop()
			fmt.Print(accountSummary.String())
			return nil
		},
	}
	return c
}

func newNetworkChainShowValidators() *cobra.Command {
	c := &cobra.Command{
		Use:   "validators [launch-id]",
		Short: "Show all validators of the chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nb, launchID, err := networkChainLaunch(cmd, args)
			if err != nil {
				return err
			}
			defer nb.Cleanup()
			n, err := nb.Network()
			if err != nil {
				return err
			}

			validatorSummary := bytes.NewBufferString("")
			validators, err := n.GenesisValidators(cmd.Context(), launchID)
			if err != nil {
				return err
			}
			validatorEntries := make([][]string, 0)
			for _, acc := range validators {
				validatorEntries = append(validatorEntries, []string{
					acc.Address,
					acc.SelfDelegation.String(),
					acc.Peer,
				})
			}
			if len(validatorEntries) > 0 {
				if err = entrywriter.MustWrite(
					validatorSummary,
					chainGenesisValSummaryHeader,
					validatorEntries...,
				); err != nil {
					return err
				}
			}
			nb.Spinner.Stop()
			fmt.Print(validatorSummary.String())
			return nil
		},
	}
	return c
}

func newNetworkChainShowPeers() *cobra.Command {
	c := &cobra.Command{
		Use:   "peers [launch-id]",
		Short: "Show peers list of the chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nb, launchID, err := networkChainLaunch(cmd, args)
			if err != nil {
				return err
			}
			defer nb.Cleanup()
			n, err := nb.Network()
			if err != nil {
				return err
			}
			genVals, err := n.GenesisValidators(cmd.Context(), launchID)
			if err != nil {
				return err
			}

			peers := make([]string, 0)
			for _, acc := range genVals {
				peers = append(peers, acc.Peer)
			}
			nb.Spinner.Stop()
			if len(peers) > 0 {
				fmt.Printf("Peers: %s\n", strings.Join(peers, ","))
			} else {
				fmt.Print("empty peer list")
			}
			return nil
		},
	}
	return c
}
