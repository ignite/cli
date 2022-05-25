package ignitecmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite-hq/cli/ignite/pkg/relayer"
	relayerconf "github.com/ignite-hq/cli/ignite/pkg/relayer/config"
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

	defaultGasPrice = "0.00025"
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
	c.Flags().String(flagSPNGasPrice, defaultGasPrice+networktypes.SPNDenom, "Gas price used for transactions on SPN")
	// TODO fetch the stake coin from chain genesis/config
	c.Flags().String(flagChainGasPrice, defaultGasPrice+"stake", "Gas price used for transactions on target chain")
	c.Flags().Int64(flagSPNGasLimit, defaultGasLimit, "Gas limit used for transactions on SPN")
	c.Flags().Int64(flagChainGasLimit, defaultGasLimit, "Gas limit used for transactions on target chain")
	c.Flags().String(flagChainAddressPrefix, "cosmos", "Address prefix of the target chain")
	c.Flags().String(flagChainAccount, cosmosaccount.DefaultAccount, "Target chain Account")
	c.Flags().String(flagChainFaucet, "", "Faucet address of the target chain")
	c.Flags().String(flagSPNChainID, networktypes.SPNChainID, "Chain ID of SPN")

	return c
}

func networkConnectHandler(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		err = handleRelayerAccountErr(err)
	}()

	session := cliui.New()
	defer session.Cleanup()

	session.StartSpinner("Setting up chains...")

	ca, err := cosmosaccount.New(
		cosmosaccount.WithKeyringBackend(getKeyringBackend(cmd)),
	)
	if err != nil {
		return err
	}

	if err := ca.EnsureDefaultAccount(); err != nil {
		return err
	}

	var (
		spnGasPrice, _   = cmd.Flags().GetString(flagSPNGasPrice)
		chainGasPrice, _ = cmd.Flags().GetString(flagChainGasPrice)
		spnGasLimit, _   = cmd.Flags().GetInt64(flagSPNGasLimit)
		chainGasLimit, _ = cmd.Flags().GetInt64(flagChainGasLimit)
		// TODO fetch from genesis
		chainAddressPrefix, _ = cmd.Flags().GetString(flagChainAddressPrefix)
		chainAccount, _       = cmd.Flags().GetString(flagChainAccount)
		chainFaucet, _        = cmd.Flags().GetString(flagChainFaucet)
		// TODO fetch from node state
		spnChainID, _ = cmd.Flags().GetString(flagSPNChainID)
	)

	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}
	// TODO fetch from chain peer
	chainRPC := args[1]

	session.StartSpinner("Creating network relayer client ID...")
	chain, spn, err := clientCreate(cmd, launchID, chainRPC, spnChainID)
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
		spn.ClientID,
	)
	if err != nil {
		return err
	}
	spnChain.ID = spn.ChainID

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
		chain.ClientID,
	)
	if err != nil {
		return err
	}
	targetChain.ID = chain.ChainID

	session.StartSpinner("Creating links between chains...")

	needsLink, pathID, cfg := spnRelayerConfig(*spnChain, *targetChain, spn, chain)
	if needsLink {
		cfg, err = r.Link(cmd.Context(), cfg, pathID)
		if err != nil {
			return err
		}
	}

	session.StopSpinner()
	if err := printSection(session, "Paths"); err != nil {
		return err
	}

	session.StartSpinner("Loading...")

	path, err := cfg.PathByID(pathID)
	if err != nil {
		return err
	}

	session.StopSpinner()

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "%s:\n", path.ID)
	fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", path.Src.ChainID, path.Src.PortID, path.Src.ChannelID)
	fmt.Fprintf(w, "   \t%s\t>\t(port: %s)\t(channel: %s)\n", path.Dst.ChainID, path.Dst.PortID, path.Dst.ChannelID)
	fmt.Fprintln(w)
	w.Flush()
	session.Print(buf.String())

	if err := printSection(session, "Listening and relaying packets between chains..."); err != nil {
		return err
	}

	_, err = r.Start(cmd.Context(), cfg, pathID)
	return err
}

func spnRelayerConfig(
	srcChain,
	dstChain relayer.Chain,
	srcChannel,
	dstChannel networktypes.Relayer,
) (bool, string, relayerconf.Config) {
	needsLink := !(srcChannel.ConnectionID != "" &&
		srcChannel.Channel.ChannelID != "" &&
		dstChannel.ConnectionID != "" &&
		dstChannel.Channel.ChannelID != "")
	pathID := relayer.PathID(srcChain.ID, dstChain.ID)
	return needsLink, pathID, relayerconf.Config{
		Version: relayerconf.SupportVersion,
		Chains:  []relayerconf.Chain{srcChain.Config(), dstChain.Config()},
		Paths: []relayerconf.Path{
			{
				ID:       pathID,
				Ordering: relayer.OrderingOrdered,
				Src: relayerconf.PathEnd{
					ChainID:      srcChain.ID,
					PortID:       networktypes.SPNPortID,
					Version:      networktypes.SPNVersion,
					ConnectionID: srcChannel.ConnectionID,
					ChannelID:    srcChannel.Channel.ChannelID,
				},
				Dst: relayerconf.PathEnd{
					ChainID:      dstChain.ID,
					PortID:       networktypes.ChainPortID,
					Version:      networktypes.SPNVersion,
					ConnectionID: dstChannel.ConnectionID,
					ChannelID:    dstChannel.Channel.ChannelID,
				},
			},
		},
	}
}
