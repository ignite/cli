package ignitecmd

import (
	"sync"

	"github.com/ignite-hq/cli/ignite/pkg/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/pkg/events"
	"github.com/ignite-hq/cli/ignite/pkg/gitpod"
	"github.com/ignite-hq/cli/ignite/services/network"
	"github.com/ignite-hq/cli/ignite/services/network/networkchain"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
)

var (
	nightly bool
	local   bool

	spnNodeAddress   string
	spnFaucetAddress string
)

const (
	flagNightly = "nightly"
	flagLocal   = "local"

	flagSPNNodeAddress   = "spn-node-address"
	flagSPNFaucetAddress = "spn-faucet-address"

	spnNodeAddressNightly   = "https://rpc.nightly.ignite.com:443"
	spnFaucetAddressNightly = "https://faucet.nightly.ignite.com"

	spnNodeAddressLocal   = "http://0.0.0.0:26657"
	spnFaucetAddressLocal = "http://0.0.0.0:4500"
)

// NewNetwork creates a new network command that holds some other sub commands
// related to creating a new network collaboratively.
func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:     "network [command]",
		Aliases: []string{"n"},
		Short:   "Launch a blockchain network in production",
		Args:    cobra.ExactArgs(1),
		Hidden:  true,
	}

	// configure flags.
	c.PersistentFlags().BoolVar(&local, flagLocal, false, "Use local SPN network")
	c.PersistentFlags().BoolVar(&nightly, flagNightly, false, "Use nightly SPN network")
	c.PersistentFlags().StringVar(&spnNodeAddress, flagSPNNodeAddress, spnNodeAddressNightly, "SPN node address")
	c.PersistentFlags().StringVar(&spnFaucetAddress, flagSPNFaucetAddress, spnFaucetAddressNightly, "SPN faucet address")

	// add sub commands.
	c.AddCommand(
		NewNetworkChain(),
		NewNetworkCampaign(),
		NewNetworkRequest(),
		NewNetworkReward(),
	)

	return c
}

var cosmos *cosmosclient.Client

type NetworkBuilder struct {
	AccountRegistry cosmosaccount.Registry
	Spinner         *clispinner.Spinner

	ev  events.Bus
	wg  *sync.WaitGroup
	cmd *cobra.Command
	cc  cosmosclient.Client
}

func newNetworkBuilder(cmd *cobra.Command) (NetworkBuilder, error) {
	var err error

	n := NetworkBuilder{
		Spinner: clispinner.New(),
		ev:      events.NewBus(),
		wg:      &sync.WaitGroup{},
		cmd:     cmd,
	}

	n.wg.Add(1)
	go printEvents(n.wg, n.ev, n.Spinner)

	if n.cc, err = getNetworkCosmosClient(cmd); err != nil {
		n.Cleanup()
		return NetworkBuilder{}, err
	}

	n.AccountRegistry = n.cc.AccountRegistry

	return n, nil
}

func (n NetworkBuilder) Chain(source networkchain.SourceOption, options ...networkchain.Option) (*networkchain.Chain, error) {
	options = append(options, networkchain.CollectEvents(n.ev))

	if home := getHome(n.cmd); home != "" {
		options = append(options, networkchain.WithHome(home))
	}

	return networkchain.New(n.cmd.Context(), n.AccountRegistry, source, options...)
}

func (n NetworkBuilder) Network(options ...network.Option) (network.Network, error) {
	var (
		err     error
		from    = getFrom(n.cmd)
		account = cosmosaccount.Account{}
	)
	if from != "" {
		account, err = cosmos.AccountRegistry.GetByName(getFrom(n.cmd))
		if err != nil {
			return network.Network{}, errors.Wrap(err, "make sure that this account exists, use 'ignite account -h' to manage accounts")
		}
	}

	options = append(options,
		network.CollectEvents(n.ev),
		network.WithCampaignQueryClient(campaigntypes.NewQueryClient(cosmos.Context())),
		network.WithLaunchQueryClient(launchtypes.NewQueryClient(cosmos.Context())),
		network.WithProfileQueryClient(profiletypes.NewQueryClient(cosmos.Context())),
		network.WithRewardQueryClient(rewardtypes.NewQueryClient(cosmos.Context())),
	)

	return network.New(*cosmos, account, options...), nil
}

func (n NetworkBuilder) Cleanup() {
	n.Spinner.Stop()
	n.ev.Shutdown()
	n.wg.Wait()
}

func getNetworkCosmosClient(cmd *cobra.Command) (cosmosclient.Client, error) {
	// check preconfigured networks
	if nightly && local {
		return cosmosclient.Client{}, errors.New("local and nightly networks can't both be specified in the same command, specify local or nightly")
	}
	if local {
		spnNodeAddress = spnNodeAddressLocal
		spnFaucetAddress = spnFaucetAddressLocal
	} else if nightly {
		spnNodeAddress = spnNodeAddressNightly
		spnFaucetAddress = spnFaucetAddressNightly
	}

	cosmosOptions := []cosmosclient.Option{
		cosmosclient.WithHome(cosmosaccount.KeyringHome),
		cosmosclient.WithNodeAddress(spnNodeAddress),
		cosmosclient.WithAddressPrefix(networktypes.SPN),
		cosmosclient.WithUseFaucet(spnFaucetAddress, networktypes.SPNDenom, 5),
		cosmosclient.WithKeyringServiceName(cosmosaccount.KeyringServiceName),
	}

	keyringBackend := getKeyringBackend(cmd)
	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	//
	// when not on Gitpod, OS keyring backend is used which only asks password once.
	if gitpod.IsOnGitpod() {
		keyringBackend = cosmosaccount.KeyringTest
	}
	if keyringBackend != "" {
		cosmosOptions = append(cosmosOptions, cosmosclient.WithKeyringBackend(keyringBackend))
	}

	// init cosmos client only once on start in order to spnclient to
	// reuse unlocked keyring in the following steps.
	if cosmos == nil {
		client, err := cosmosclient.New(cmd.Context(), cosmosOptions...)
		if err != nil {
			return cosmosclient.Client{}, err
		}
		cosmos = &client
	}

	if err := cosmos.AccountRegistry.EnsureDefaultAccount(); err != nil {
		return cosmosclient.Client{}, err
	}

	return *cosmos, nil
}
