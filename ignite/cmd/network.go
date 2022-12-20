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
)

const (
	flagNightly = "nightly"
	flagLocal   = "local"

	flagSPNNodeAddress   = "spn-node-address"
	flagSPNFaucetAddress = "spn-faucet-address"

	spnNodeAddressNightly   = "https://rpc.devnet.ignite.com:443"
	spnFaucetAddressNightly = "https://faucet.devnet.ignite.com:443"

	spnNodeAddressLocal   = "http://0.0.0.0:26661"
	spnFaucetAddressLocal = "http://0.0.0.0:4502"
)

// NewNetwork creates a new network command that holds some other sub commands
// related to creating a new network collaboratively.
func NewNetwork() *cobra.Command {
	c := &cobra.Command{
		Use:     "network [command]",
		Aliases: []string{"n"},
		Short:   "Launch a blockchain in production",
		Long: `
Ignite Network commands allow to coordinate the launch of sovereign Cosmos blockchains.

To launch a Cosmos blockchain you need someone to be a coordinator and others to
be validators. These are just roles, anyone can be a coordinator or a validator.
A coordinator publishes information about a chain to be launched on the Ignite
blockchain, approves validator requests and coordinates the launch. Validators
send requests to join a chain and start their nodes when a blockchain is ready
for launch.

To publish the information about your chain as a coordinator run the following
command (the URL should point to a repository with a Cosmos SDK chain):

	ignite network chain publish github.com/ignite/example

This command will return a launch identifier you will be using in the following
commands. Let's say this identifier is 42.

Next, ask validators to initialize their nodes and request to join the network
as validators. For a testnet you can use the default values suggested by the
CLI.

	ignite network chain init 42

	ignite network chain join 42 --amount 95000000stake

As a coordinator list all validator requests:

	ignite network request list 42

Approve validator requests:

	ignite network request approve 42 1,2

Once you've approved all validators you need in the validator set, announce that
the chain is ready for launch:

	ignite network chain launch 42

Validators can now prepare their nodes for launch:

	ignite network chain prepare 42

The output of this command will show a command that a validator would use to
launch their node, for example “exampled --home ~/.example”. After enough
validators launch their nodes, a blockchain will be live.
`,
		Args: cobra.ExactArgs(1),
	}

	// configure flags.
	c.PersistentFlags().BoolVar(&local, flagLocal, false, "Use local SPN network")
	c.PersistentFlags().BoolVar(&nightly, flagNightly, false, "Use nightly SPN network")
	// Includes Flags for Node and Faucet Address
	c.PersistentFlags().AddFlagSet(flagSetSpnAddresses())

	// add sub commands.
	c.AddCommand(
		NewNetworkChain(),
		NewNetworkProject(),
		NewNetworkRequest(),
		NewNetworkReward(),
		NewNetworkValidator(),
		NewNetworkProfile(),
		NewNetworkCoordinator(),
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

	NetworkAddresses struct {
		NodeAddress   string
		FaucetAddress string
	}
)

func CollectEvents(ev events.Bus) NetworkBuilderOption {
	return func(builder *NetworkBuilder) {
		builder.ev = ev
	}
}

func flagSetSPNAccountPrefixes() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagAddressPrefix, networktypes.SPN, "account address prefix")
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
	spn, err := getSpnAddresses(cmd)
	if err != nil {
		return cosmosclient.Client{}, err
	}

	cosmosOptions := []cosmosclient.Option{
		cosmosclient.WithHome(cosmosaccount.KeyringHome),
		cosmosclient.WithNodeAddress(spn.NodeAddress),
		cosmosclient.WithAddressPrefix(networktypes.SPN),
		cosmosclient.WithUseFaucet(spn.FaucetAddress, networktypes.SPNDenom, 5),
		cosmosclient.WithKeyringServiceName(cosmosaccount.KeyringServiceName),
		cosmosclient.WithKeyringDir(getKeyringDir(cmd)),
		cosmosclient.WithGas(cosmosclient.GasAuto),
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

func flagSetSpnAddresses() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagSPNNodeAddress, spnNodeAddressNightly, "SPN node address")
	fs.String(flagSPNFaucetAddress, spnFaucetAddressNightly, "SPN faucet address")
	return fs
}

func getSpnAddresses(cmd *cobra.Command) (NetworkAddresses, error) {
	// check preconfigured networks
	if nightly && local {
		return NetworkAddresses{}, errors.New("local and nightly networks can't both be specified in the same command, specify local or nightly")
	}
	if nightly {
		return NetworkAddresses{spnNodeAddressNightly, spnFaucetAddressNightly}, nil
	}
	if local {
		return NetworkAddresses{spnNodeAddressLocal, spnFaucetAddressLocal}, nil
	}

	spnNodeAddress, err := cmd.Flags().GetString(flagSPNNodeAddress)
	if err != nil {
		return NetworkAddresses{}, err
	}

	spnFaucetAddress, err := cmd.Flags().GetString(flagSPNFaucetAddress)
	if err != nil {
		return NetworkAddresses{}, err
	}
	return NetworkAddresses{spnNodeAddress, spnFaucetAddress}, nil
}
