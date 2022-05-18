package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite-hq/cli/ignite/pkg/relayer"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

const (
	flagChainFaucet        = "chain-faucet"
	flagChainAddressPrefix = "chain-prefix"
	flagChainAccount       = "chain-account"
	flagChainGasPrice      = "chain-gasprice"
	flagChainGasLimit      = "chain-gaslimit"
	flagSPNGasPrice        = "spn-gasprice"
	flagSPNGasLimit        = "spn-gaslimit"

	defaultGasPrice = "0.0000025"
	defaultGasLimit = 400000
)

// NewNetworkClientConnect connects the monitoring modules of launched chains with SPN
func NewNetworkClientConnect() *cobra.Command {
	c := &cobra.Command{
		Use:   "connect [launch-id] [chain-rpc]",
		Short: "Connect the monitoring modules of launched chains with SPN",
		Args:  cobra.ExactArgs(2),
		RunE:  networkConnectHandler,
	}

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().String(flagSPNGasPrice, defaultGasPrice, "Gas price used for transactions on SPN")
	c.Flags().String(flagChainGasPrice, defaultGasPrice, "Gas price used for transactions on target chain")
	c.Flags().Int64(flagSPNGasLimit, defaultGasLimit, "Gas limit used for transactions on SPN")
	c.Flags().Int64(flagChainGasLimit, defaultGasLimit, "Gas limit used for transactions on target chain")
	c.Flags().String(flagChainAddressPrefix, "", "Address prefix of the target chain")
	c.Flags().String(flagChainAccount, cosmosaccount.DefaultAccount, "Target chain Account")
	c.Flags().String(flagChainFaucet, "", "Faucet address of the target chain")

	return c
}

// ignite network --local connect [launch-id] [target-rpc]
// Flag
// --target-faucet
// Flag values with defaults
// --source-gaslimit
// --target-gaslimit
// --source-gasprice
// --target-gasprice
// --source-account (from)
// --target-account (from)
// Hardcoded flags
// --ordered
// --source-rpc "http://0.0.0.0:26657"
// --source-faucet "http://0.0.0.0:4500"
// --source-port "monitoringc"
// --target-port "monitoringp"
// --source-version "monitoring-1"
// --target-version "monitoring-1"
// --source-prefix "spn"
// --target-prefix "cosmos" (fetch from genesis)

func networkConnectHandler(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		err = handleRelayerAccountErr(err)
	}()

	session := cliui.New()
	defer session.Cleanup()

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
	)
	if err != nil {
		return err
	}

	if err := ca.EnsureDefaultAccount(); err != nil {
		return err
	}

	if err := printSection(session, "Setting up chains"); err != nil {
		return err
	}

	var (
		spnGasPrice, _        = cmd.Flags().GetString(flagSPNGasPrice)
		chainGasPrice, _      = cmd.Flags().GetString(flagChainGasPrice)
		spnGasLimit, _        = cmd.Flags().GetInt64(flagSPNGasLimit)
		chainGasLimit, _      = cmd.Flags().GetInt64(flagChainGasLimit)
		chainAddressPrefix, _ = cmd.Flags().GetString(flagChainAddressPrefix)
		chainAccount, _       = cmd.Flags().GetString(flagChainAccount)
		chainFaucet, _        = cmd.Flags().GetString(flagChainFaucet)
	)

	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}
	chainRPC := args[1]

	session.StartSpinner("Creating network relayer client ID...")
	spnClientID, chainClientID, err := clientCreate(cmd, launchID, chainRPC)
	if err != nil {
		return err
	}

	session.StartSpinner("Fetching chain info...")
	session.Println()

	r := relayer.New(ca)
	// initialize the chains
	spnChain, err := initChain(
		cmd,
		r,
		session,
		relayerSource,
		getFrom(cmd),
		spnNodeAddress,
		spnFaucetAddress,
		spnGasPrice,
		spnGasLimit,
		networktypes.SPN,
		spnClientID,
	)
	if err != nil {
		return err
	}

	targetChain, err := initChain(
		cmd,
		r,
		session,
		relayerTarget,
		chainAccount,
		chainRPC,
		chainFaucet,
		chainGasPrice,
		chainGasLimit,
		chainAddressPrefix,
		chainClientID,
	)
	if err != nil {
		return err
	}

	// TODO remove me
	session.Println(spnChain, targetChain)

	session.StartSpinner("Configuring...")

	session.StopSpinner()
	return nil
}
