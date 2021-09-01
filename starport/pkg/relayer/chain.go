package relayer

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/imdario/mergo"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	relayerconfig "github.com/tendermint/starport/starport/pkg/relayer/config"
)

const (
	TransferPort      = "transfer"
	TransferVersion   = "ics20-1"
	OrderingUnordered = "ORDER_UNORDERED"
	OrderingOrdered   = "ORDER_ORDERED"
)

var (
	errEndpointExistsWithDifferentChainID = errors.New("rpc endpoint already exists with a different chain id")
)

// Chain represents a chain in relayer.
type Chain struct {
	// ID is id of the chain.
	ID string

	// accountName is account used on the chain.
	accountName string

	// rpcAddress is the node address of tm.
	rpcAddress string

	// faucetAddress is the faucet address to get tokens for relayer accounts.
	faucetAddress string

	r Relayer

	// options are used to set up the chain.
	options chainOptions
}

// chainOptions holds options to be used setting up the chain.
type chainOptions struct {
	// GasPrice is the gas price used when sending transactions to the chain
	GasPrice string `json:"gasPrice"`

	// GasLimit is the gas limit used when sending transactions to the chain
	GasLimit int64 `json:"gasLimit"`

	// AddressPrefix is the address prefix of the chain.
	AddressPrefix string `json:"addressPrefix"`
}

// Account represents an account in relayer.
type Account struct {
	// Address of the account.
	Address string `json:"address"`
}

// Option is used to configure Chain.
type Option func(*Chain)

// WithFaucet provides a faucet address for chain to get tokens from.
// when it isn't provided.
func WithFaucet(address string) Option {
	return func(c *Chain) {
		c.faucetAddress = address
	}
}

// WithGasPrice gives the gas price to use to send ibc transactions to the chain.
func WithGasPrice(gasPrice string) Option {
	return func(c *Chain) {
		c.options.GasPrice = gasPrice
	}
}

// WithGasLimit gives the gas limit to use to send ibc transactions to the chain.
func WithGasLimit(limit int64) Option {
	return func(c *Chain) {
		c.options.GasLimit = limit
	}
}

// WithAddressPrefix configures the account key prefix used on the chain.
func WithAddressPrefix(addressPrefix string) Option {
	return func(c *Chain) {
		c.options.AddressPrefix = addressPrefix
	}
}

// NewChain creates a new chain on relayer or uses the existing matching chain.
func (r Relayer) NewChain(ctx context.Context, accountName, rpcAddress string, options ...Option) (
	*Chain, cosmosaccount.Account, error) {
	c := &Chain{
		accountName: accountName,
		rpcAddress:  fixRPCAddress(rpcAddress),
		r:           r,
	}

	// apply user options.
	for _, o := range options {
		o(c)
	}

	if err := c.ensureChainSetup(ctx); err != nil {
		return nil, cosmosaccount.Account{}, err
	}

	account, err := r.ca.GetByName(accountName)
	if err != nil {
		return nil, cosmosaccount.Account{}, err
	}

	return c, account, nil
}

// TryRetrieve tries to receive some coins to the account and returns the total balance.
func (c *Chain) TryRetrieve(ctx context.Context) (sdk.Coins, error) {
	acc, err := c.r.ca.GetByName(c.accountName)
	if err != nil {
		return nil, err
	}

	addr := acc.Address(c.options.AddressPrefix)

	err = cosmosfaucet.TryRetrieve(ctx, c.ID, c.rpcAddress, c.faucetAddress, addr)
	if err != nil {
		return nil, err
	}
	return c.r.balance(ctx, c.rpcAddress, c.accountName, c.options.AddressPrefix)
}

// channelOptions represents options for configuring the IBC channel between two chains
type channelOptions struct {
	SourcePort    string `json:"sourcePort"`
	SourceVersion string `json:"sourceVersion"`
	TargetPort    string `json:"targetPort"`
	TargetVersion string `json:"targetVersion"`
	Ordering      string `json:"ordering"`
}

// newChannelOptions returns default channel options
func newChannelOptions() channelOptions {
	return channelOptions{
		SourcePort:    TransferPort,
		SourceVersion: TransferVersion,
		TargetPort:    TransferPort,
		TargetVersion: TransferVersion,
		Ordering:      OrderingUnordered,
	}
}

// ChannelOption is used to configure relayer IBC connection
type ChannelOption func(*channelOptions)

// SourcePort configures the source port of the new channel
func SourcePort(port string) ChannelOption {
	return func(c *channelOptions) {
		c.SourcePort = port
	}
}

// TargetPort configures the target port of the new channel
func TargetPort(port string) ChannelOption {
	return func(c *channelOptions) {
		c.TargetPort = port
	}
}

// SourceVersion configures the source version of the new channel
func SourceVersion(version string) ChannelOption {
	return func(c *channelOptions) {
		c.SourceVersion = version
	}
}

// TargetVersion configures the target version of the new channel
func TargetVersion(version string) ChannelOption {
	return func(c *channelOptions) {
		c.TargetVersion = version
	}
}

// Ordered sets the new channel as ordered
func Ordered() ChannelOption {
	return func(c *channelOptions) {
		c.Ordering = OrderingOrdered
	}
}

// Connect connects dst chain to c chain and creates a path in between in offline mode.
// it returns the path id on success otherwise, returns with a non-nil error.
func (c *Chain) Connect(ctx context.Context, dst *Chain, options ...ChannelOption) (id string, err error) {
	channelOptions := newChannelOptions()

	for _, apply := range options {
		apply(&channelOptions)
	}

	conf, err := relayerconfig.Get()
	if err != nil {
		return "", err
	}

	// determine a unique path name from chain ids with incremental numbers. e.g.:
	// - src-dst
	// - src-dst-2
	pathID := fmt.Sprintf("%s-%s", c.ID, dst.ID)
	var suffix string
	i := 2
	for {
		guess := pathID + suffix
		if _, err := conf.PathByID(guess); err != nil { // guess is inique.
			pathID = guess
			break
		}
		suffix = fmt.Sprintf("-%d", i)
		i++
	}

	confPath := relayerconfig.Path{
		ID:       pathID,
		Ordering: channelOptions.Ordering,
		Src: relayerconfig.PathEnd{
			ChainID: c.ID,
			PortID:  channelOptions.SourcePort,
			Version: channelOptions.SourceVersion,
		},
		Dst: relayerconfig.PathEnd{
			ChainID: dst.ID,
			PortID:  channelOptions.TargetPort,
			Version: channelOptions.TargetVersion,
		},
	}

	conf.Paths = append(conf.Paths, confPath)

	if err := relayerconfig.Save(conf); err != nil {
		return "", err
	}

	return pathID, nil
}

// ensureChainSetup sets up the new or existing chain.
func (c *Chain) ensureChainSetup(ctx context.Context) error {
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(c.rpcAddress))
	if err != nil {
		return err
	}
	status, err := client.RPC.Status(ctx)
	if err != nil {
		return err
	}
	c.ID = status.NodeInfo.Network

	confChain := relayerconfig.Chain{
		ID:            c.ID,
		Account:       c.accountName,
		AddressPrefix: c.options.AddressPrefix,
		RPCAddress:    c.rpcAddress,
		GasPrice:      c.options.GasPrice,
		GasLimit:      c.options.GasLimit,
	}

	conf, err := relayerconfig.Get()
	if err != nil {
		return err
	}

	var found bool

	for i, chain := range conf.Chains {
		if chain.ID == c.ID {
			if chain.RPCAddress != c.rpcAddress {
				return errEndpointExistsWithDifferentChainID
			}

			if err := mergo.Merge(&conf.Chains[i], confChain, mergo.WithOverride); err != nil {
				return err
			}

			found = true
			break
		}
	}

	if !found {
		conf.Chains = append(conf.Chains, confChain)
	}

	return relayerconfig.Save(conf)
}
