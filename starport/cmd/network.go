package starportcmd

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/services/network"
	"github.com/tendermint/starport/starport/services/network/networkchain"
	"github.com/tendermint/starport/starport/services/network/networktypes"
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

	spnNodeAddressAlpha   = "https://rpc.alpha.starport.network:443"
	spnFaucetAddressAlpha = "https://faucet.alpha.starport.network"

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
	c.PersistentFlags().StringVar(&spnNodeAddress, flagSPNNodeAddress, spnNodeAddressAlpha, "SPN node address")
	c.PersistentFlags().StringVar(&spnFaucetAddress, flagSPNFaucetAddress, spnFaucetAddressAlpha, "SPN faucet address")

	// add sub commands.
	c.AddCommand(
		NewNetworkChain(),
		NewNetworkCampaign(),
		NewNetworkRequest(),
		NewNetworkProfile(),
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
	options = append(options, network.CollectEvents(n.ev))

	account, err := cosmos.AccountRegistry.GetByName(getFrom(n.cmd))
	if err != nil {
		return network.Network{}, errors.Wrap(err, "make sure that this account exists, use 'starport account -h' to manage accounts")
	}

	return network.New(*cosmos, account, options...)
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
	cosmosOptions = append(cosmosOptions, cosmosclient.WithKeyringBackend(keyringBackend))

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
