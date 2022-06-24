package ignitecmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/gitpod"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networkchain"
	"github.com/ignite/cli/ignite/services/network/networktypes"
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

	spnNodeAddressNightly   = "https://rpc.nightly.starport.network:443"
	spnFaucetAddressNightly = "https://faucet.nightly.starport.network"

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
		NewNetworkClient(),
	)

	return c
}

var cosmos *cosmosclient.Client

type (
	NetworkBuilderOption func(builder *NetworkBuilder)

	NetworkBuilder struct {
		AccountRegistry cosmosaccount.Registry

		ev  events.Bus
		cmd *cobra.Command
		cc  cosmosclient.Client
	}
)

func CollectEvents(ev events.Bus) NetworkBuilderOption {
	return func(builder *NetworkBuilder) {
		builder.ev = ev
	}
}

func flagSetSPNAccountPrefixes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagAddressPrefix, networktypes.SPN, "Account address prefix")
	return fs
}

func newNetworkBuilder(cmd *cobra.Command, options ...NetworkBuilderOption) (NetworkBuilder, error) {
	var (
		err error
		n   = NetworkBuilder{cmd: cmd}
	)

	if n.cc, err = getNetworkCosmosClient(cmd); err != nil {
		return NetworkBuilder{}, err
	}

	n.AccountRegistry = n.cc.AccountRegistry

	for _, apply := range options {
		apply(&n)
	}
	return n, nil
}

func (n NetworkBuilder) Chain(source networkchain.SourceOption, options ...networkchain.Option) (*networkchain.Chain, error) {
	if home := getHome(n.cmd); home != "" {
		options = append(options, networkchain.WithHome(home))
	}

	options = append(options, networkchain.CollectEvents(n.ev))

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

	options = append(options, network.CollectEvents(n.ev))

	return network.New(*cosmos, account, options...), nil
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
