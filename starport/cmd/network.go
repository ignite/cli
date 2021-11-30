package starportcmd

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/services/network"
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
	c.PersistentFlags().StringVar(&spnFaucetAddress, flagSPNFaucetAddress, spnFaucetAddressAlpha, "SPN Faucet address")

	// add sub commands.
	c.AddCommand(NewNetworkChain())

	return c
}

var cosmos *cosmosclient.Client

// initializeNetwork initializes event bus, CLIn components such as spinner and returns a new network builder
func initializeNetwork(cmd *cobra.Command) (
	nb *network.Builder,
	spinner *clispinner.Spinner,
	cleanup func(),
	err error,
) {
	var (
		wg sync.WaitGroup
		s  = clispinner.New()
		ev = events.NewBus()
	)
	wg.Add(1)
	go printEvents(&wg, ev, s)

	shutdown := func() {
		s.Stop()
		ev.Shutdown()
		wg.Wait()
	}

	nb, err = newNetwork(cmd, network.CollectEvents(ev))
	if err != nil {
		shutdown()
	}
	return nb, s, shutdown, err
}

// newNetwork returns a new network builder initialized with command flag
func newNetwork(cmd *cobra.Command, options ...network.Option) (*network.Builder, error) {
	// check preconfigured networks
	if nightly && local {
		return nil, errors.New("local and nightly networks can't be specified in the same command")
	}
	if local {
		spnNodeAddress = spnNodeAddressLocal
		spnFaucetAddress = spnFaucetAddressLocal
	} else if nightly {
		spnNodeAddress = spnNodeAddressNightly
		spnFaucetAddress = spnFaucetAddressNightly
	}

	cosmosOptions := []cosmosclient.Option{
		cosmosclient.WithHome(getHome(cmd)),
		cosmosclient.WithNodeAddress(spnNodeAddress),
		cosmosclient.WithAddressPrefix(network.SPNAddressPrefix),
		cosmosclient.WithUseFaucet(spnFaucetAddress, "", 0),
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
			return nil, err
		}
		cosmos = &client
	}

	if err := cosmos.AccountRegistry.EnsureDefaultAccount(); err != nil {
		return nil, err
	}

	account, err := cosmos.AccountRegistry.GetByName(getFrom(cmd))
	if err != nil {
		return nil, errors.Wrap(err, "make sure that this account exists, use 'starport account -h' to manage accounts")
	}

	return network.New(*cosmos, account, options...)
}

func printSection(title string) {
	fmt.Printf("---------------------------------------------\n%s\n---------------------------------------------\n\n", title)
}
