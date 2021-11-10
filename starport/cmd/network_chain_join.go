package starportcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
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

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())

	return c
}

func networkChainJoinHandler(cmd *cobra.Command, args []string) error {
	var (
		gentxPath, _ = cmd.Flags().GetString(flagGentx)
		amountArg, _ = cmd.Flags().GetString(flagAmount)
	)

	amount, err := sdk.ParseCoinsNormalized(amountArg)
	if err != nil {
		return fmt.Errorf("error parsing amount: %s", err.Error())
	}

	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing launchID: %s", err.Error())
	}

	nb, _, endRoutine, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer endRoutine()

	home := getHome(cmd)
	if home == "" {
		var err error
		home, err = network.ChainHome(launchID)
		if err != nil {
			return err
		}
	}

	homeGentxPath, err := checkChainHomeInitialized(home)
	if err != nil {
		return err
	}
	if gentxPath == "" {
		gentxPath = homeGentxPath
	}

	// initialize the blockchain from the launch ID
	sourceOption := network.SourceLaunchID(launchID)
	blockchain, err := nb.Blockchain(cmd.Context(), sourceOption, network.InitializationHomePath(home))
	if err != nil {
		return err
	}

	addr, err := blockchain.GetAccountAddress(cmd.Context(), getFrom(cmd))
	if err != nil {
		return err
	}

	genesis, exist, err := getChainGenesis(gentxPath)
	if err != nil {
		return err
	}
	if exist {
		hasAcc := genesis.HasAccount(addr)
		if !hasAcc {
			exist, err := blockchain.CheckRequestAccount(cmd.Context(), launchID, addr)
			if err != nil {
				return err
			}
			if !exist {
				// TODO handle if the account not exist
			}
		}
	}

	// TODO: If [--amount], then send "create account" message (address has a default)
	_ = amount

	result, err := blockchain.Join(cmd.Context(), launchID, amount)
	if err != nil {
		return err
	}

	// delegatorAddress, validatorAddress, selfDelegation, err := network.ParseGentx(gentxPath)

	fmt.Printf("%s Network joined\n%s", clispinner.OK, result)
	return nil
}

// checkChainHomeInitialized checks if a home with the provided launchID already initialized
func checkChainHomeInitialized(home string) (string, error) {
	gentxPath := filepath.Join(home, "config/gentx/gentx.json")
	_, err := os.Stat(gentxPath)
	if err != nil {
		return home, err
	}
	return gentxPath, err
}

// getChainGenesis return the chain genesis path
func getChainGenesis(home string) (network.ChainGenesis, bool, error) {
	genesisPath := filepath.Join(home, "config/genesis.json")
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return network.ChainGenesis{}, false, nil
	} else if err != nil {
		return network.ChainGenesis{}, false, err
	}
	net, err := network.ParseGenesis(genesisPath)
	if err != nil {
		return network.ChainGenesis{}, false, err
	}
	return net, true, nil
}
