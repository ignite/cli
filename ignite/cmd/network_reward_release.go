package ignitecmd

import (
	"bytes"
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/relayer"
	relayerconf "github.com/ignite/cli/ignite/pkg/relayer/config"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

const (
	flagTestnetFaucet        = "testnet-faucet"
	flagTestnetAddressPrefix = "testnet-prefix"
	flagTestnetAccount       = "testnet-account"
	flagTestnetGasPrice      = "testnet-gasprice"
	flagTestnetGasLimit      = "testnet-gaslimit"
	flagSPNGasPrice          = "spn-gasprice"
	flagSPNGasLimit          = "spn-gaslimit"
	flagCreateClientOnly     = "create-client-only"

	defaultTestnetGasPrice = "0.0000025stake"
	defaultSPNGasPrice     = "0.0000025" + networktypes.SPNDenom
	defaultGasLimit        = 400000
)

// NewNetworkRewardRelease connects the monitoring modules of launched
// chains with SPN and distribute rewards with chain Relayer.
func NewNetworkRewardRelease() *cobra.Command {
	c := &cobra.Command{
		Use:   "release [launch-id] [chain-rpc]",
		Short: "Connect the monitoring modules of launched chains with SPN",
		Args:  cobra.ExactArgs(2),
		RunE:  networkRewardRelease,
	}

	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().String(flagSPNGasPrice, defaultSPNGasPrice, "gas price used for transactions on SPN")
	c.Flags().String(flagTestnetGasPrice, defaultTestnetGasPrice, "gas price used for transactions on testnet chain")
	c.Flags().Int64(flagSPNGasLimit, defaultGasLimit, "gas limit used for transactions on SPN")
	c.Flags().Int64(flagTestnetGasLimit, defaultGasLimit, "gas limit used for transactions on testnet chain")
	c.Flags().String(flagTestnetAddressPrefix, cosmosaccount.AccountPrefixCosmos, "address prefix of the testnet chain")
	c.Flags().String(flagTestnetAccount, cosmosaccount.DefaultAccount, "testnet chain account")
	c.Flags().String(flagTestnetFaucet, "", "faucet address of the testnet chain")
	c.Flags().Bool(flagCreateClientOnly, false, "only create the network client id")

	return c
}

func networkRewardRelease(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		err = handleRelayerAccountErr(err)
	}()

	session := cliui.New(cliui.StartSpinnerWithText("Setting up chains..."))
	defer session.End()

	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}
	chainRPC := xurl.HTTPEnsurePort(args[1])

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}
	spnChainID, err := n.ChainID(cmd.Context())
	if err != nil {
		return err
	}

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
		createClientOnly, _ = cmd.Flags().GetBool(flagCreateClientOnly)
		spnGasPrice, _      = cmd.Flags().GetString(flagSPNGasPrice)
		testnetGasPrice, _  = cmd.Flags().GetString(flagTestnetGasPrice)
		spnGasLimit, _      = cmd.Flags().GetInt64(flagSPNGasLimit)
		testnetGasLimit, _  = cmd.Flags().GetInt64(flagTestnetGasLimit)
		// TODO fetch from genesis
		testnetAddressPrefix, _ = cmd.Flags().GetString(flagTestnetAddressPrefix)
		testnetAccount, _       = cmd.Flags().GetString(flagTestnetAccount)
		testnetFaucet, _        = cmd.Flags().GetString(flagTestnetFaucet)
	)

	session.StartSpinner("Creating network relayer client ID...")
	chain, spn, err := createClient(cmd, n, session, launchID, chainRPC, spnChainID)
	if err != nil {
		return err
	}
	if createClientOnly {
		return nil
	}

	session.StartSpinner("Fetching chain info...")
	session.Println()

	spnAddresses, err := getSpnAddresses(cmd)
	if err != nil {
		return err
	}

	r := relayer.New(ca)
	// initialize the chains
	spnChain, err := initChain(
		cmd,
		r,
		session,
		relayerSource,
		getFrom(cmd),
		spnAddresses.NodeAddress,
		spnAddresses.FaucetAddress,
		spnGasPrice,
		spnGasLimit,
		networktypes.SPN,
		spn.ClientID,
	)
	if err != nil {
		return err
	}
	spnChain.ID = spn.ChainID

	testnetChain, err := initChain(
		cmd,
		r,
		session,
		relayerTarget,
		testnetAccount,
		chainRPC,
		testnetFaucet,
		testnetGasPrice,
		testnetGasLimit,
		testnetAddressPrefix,
		chain.ClientID,
	)
	if err != nil {
		return err
	}
	testnetChain.ID = chain.ChainID

	session.StartSpinner("Creating links between chains...")

	pathID, cfg, err := spnRelayerConfig(*spnChain, *testnetChain, spn, chain)
	if err != nil {
		return err
	}
	if spn.ChannelID == "" {
		cfg, err = r.Link(cmd.Context(), cfg, pathID)
		if err != nil {
			return err
		}
	}

	if err := printSection(session, "Paths"); err != nil {
		return err
	}

	session.StartSpinner("Loading...")

	path, err := cfg.PathByID(pathID)
	if err != nil {
		return err
	}

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

	return r.Start(cmd.Context(), cfg, pathID, nil)
}

func createClient(
	cmd *cobra.Command,
	n network.Network,
	session *cliui.Session,
	launchID uint64,
	nodeAPI,
	spnChainID string,
) (networktypes.RewardIBCInfo, networktypes.RewardIBCInfo, error) {
	nodeClient, err := cosmosclient.New(cmd.Context(), cosmosclient.WithNodeAddress(nodeAPI))
	if err != nil {
		return networktypes.RewardIBCInfo{}, networktypes.RewardIBCInfo{}, err
	}
	node := network.NewNode(nodeClient)

	chainRelayer, err := node.RewardIBCInfo(cmd.Context())
	if err != nil {
		return networktypes.RewardIBCInfo{}, networktypes.RewardIBCInfo{}, err
	}

	rewardsInfo, chainID, unboundingTime, err := node.RewardsInfo(cmd.Context())
	if err != nil {
		return networktypes.RewardIBCInfo{}, networktypes.RewardIBCInfo{}, err
	}

	spnRelayer, err := n.RewardIBCInfo(cmd.Context(), launchID)
	if errors.Is(err, network.ErrObjectNotFound) {
		spnRelayer.ClientID, err = n.CreateClient(cmd.Context(), launchID, unboundingTime, rewardsInfo)
	}
	if err != nil {
		return networktypes.RewardIBCInfo{}, networktypes.RewardIBCInfo{}, err
	}

	chainRelayer.ChainID = chainID
	spnRelayer.ChainID = spnChainID

	session.Printf(
		"%s Network client: %s\n",
		icons.Info,
		spnRelayer.ClientID,
	)
	printRelayerOptions(session, spnRelayer.ConnectionID, spnRelayer.ChainID, "connection")
	printRelayerOptions(session, spnRelayer.ChannelID, spnRelayer.ChainID, "channel")

	session.Printf(
		"%s Testnet chain %s client: %s\n",
		icons.Info,
		chainRelayer.ChainID,
		chainRelayer.ClientID,
	)
	printRelayerOptions(session, chainRelayer.ConnectionID, chainRelayer.ChainID, "connection")
	printRelayerOptions(session, chainRelayer.ChannelID, chainRelayer.ChainID, "channel")
	return chainRelayer, spnRelayer, err
}

func printRelayerOptions(session *cliui.Session, obj, chainID, option string) {
	if obj != "" {
		session.Printf("%s The chain %s already have a %s: %s\n",
			icons.Bullet,
			chainID,
			option,
			obj,
		)
	}
}

func spnRelayerConfig(
	srcChain,
	dstChain relayer.Chain,
	srcChannel,
	dstChannel networktypes.RewardIBCInfo,
) (string, relayerconf.Config, error) {
	var (
		pathID = relayer.PathID(srcChain.ID, dstChain.ID)
		conf   = relayerconf.Config{
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
						ChannelID:    srcChannel.ChannelID,
					},
					Dst: relayerconf.PathEnd{
						ChainID:      dstChain.ID,
						PortID:       networktypes.ChainPortID,
						Version:      networktypes.SPNVersion,
						ConnectionID: dstChannel.ConnectionID,
						ChannelID:    dstChannel.ChannelID,
					},
				},
			},
		}
	)
	switch {
	case srcChannel.ConnectionID != "" &&
		srcChannel.ChannelID != "" &&
		dstChannel.ConnectionID != "" &&
		dstChannel.ChannelID != "":
		return pathID, conf, nil
	case srcChannel.ConnectionID == "" &&
		srcChannel.ChannelID == "" &&
		dstChannel.ConnectionID == "" &&
		dstChannel.ChannelID == "":
		return pathID, conf, nil
	}
	return pathID, conf, errors.New("connection was already established and is missing in one of the chains")
}
