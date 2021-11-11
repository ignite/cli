package starportcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rdegges/go-ipify"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/cliquiz"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/xchisel"
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

	amount, err := sdk.ParseCoinNormalized(amountArg)
	if err != nil {
		return fmt.Errorf("error parsing amount: %s", err.Error())
	}

	launchID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing launchID: %s", err.Error())
	}

	// initialize network common methods
	nb, _, endRoutine, err := initializeNetwork(cmd)
	if err != nil {
		return err
	}
	defer endRoutine()

	// get the pear public address for the validator
	publicAddress, err := askPublicAddress()
	if err != nil {
		return err
	}

	home := getHome(cmd)
	if home == "" {
		var err error
		home, err = network.ChainHome(launchID)
		if err != nil {
			return err
		}
	}

	// add the custom gentx if provided
	if gentxPath != "" {
		gentxPath = filepath.Join(home, "config/gentx/gentx.json")
	}

	info, gentx, err := network.ParseGentx(gentxPath)
	if err != nil {
		return err
	}

	// send message to add the validator into the SPN
	result, err := nb.Join(
		cmd.Context(),
		launchID,
		info.ValidatorAddress,
		publicAddress,
		gentx,
		info.PubKey,
		amount,
	)
	if err != nil {
		return err
	}

	addr := ""

	genesis, exist, err := getChainGenesis(gentxPath)
	if err != nil {
		return err
	}
	if exist {
		hasAcc := genesis.HasAccount(addr)
		if !hasAcc {
			exist, err := nb.CheckRequestAccount(cmd.Context(), launchID, addr)
			if err != nil {
				return err
			}
			if !exist {
				return fmt.Errorf("account already exist %s", addr)
			}
		}
	}

	fmt.Printf("%s Network joined\n%s", clispinner.OK, result)
	return nil
}

// askPublicAddress prepare questions to interactively ask for a publicAddress when peer isn't provided
// and not running through chisel proxy.
func askPublicAddress() (publicAddress string, err error) {
	options := []cliquiz.Option{
		cliquiz.Required(),
	}
	if !xchisel.IsEnabled() {
		ip, _ := ipify.GetIp()
		if err == nil {
			options = append(options, cliquiz.DefaultAnswer(fmt.Sprintf("%s:26656", ip)))
		}
	}
	questions := []cliquiz.Question{cliquiz.NewQuestion(
		"Peer's address",
		&publicAddress,
		options...,
	)}
	return publicAddress, cliquiz.Ask(questions...)
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
