package relayer

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/imdario/mergo"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
	relayerconfig "github.com/ignite/cli/ignite/pkg/relayer/config"
)

const (
	TransferPort      = "transfer"
	TransferVersion   = "ics20-1"
	OrderingUnordered = "ORDER_UNORDERED"
	OrderingOrdered   = "ORDER_ORDERED"
)

var errEndpointExistsWithDifferentChainID = errors.New("rpc endpoint already exists with a different chain id")

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

	// gasPrice is the gas price used when sending transactions to the chain
	gasPrice string

	// gasLimit is the gas limit used when sending transactions to the chain
	gasLimit int64

	// addressPrefix is the address prefix of the chain.
	addressPrefix string

	// clientID is the client id of the chain for relayer connection.
	clientID string

	r Relayer
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
		c.gasPrice = gasPrice
	}
}

// WithGasLimit gives the gas limit to use to send ibc transactions to the chain.
func WithGasLimit(limit int64) Option {
	return func(c *Chain) {
		c.gasLimit = limit
	}
}

// WithAddressPrefix configures the account key prefix used on the chain.
func WithAddressPrefix(addressPrefix string) Option {
	return func(c *Chain) {
		c.addressPrefix = addressPrefix
	}
}

// WithClientID configures the chain client id.
func WithClientID(clientID string) Option {
	return func(c *Chain) {
		c.clientID = clientID
	}
}

// NewChain creates a new chain on relayer or uses the existing matching chain.
func (r Relayer) NewChain(accountName, rpcAddress string, options ...Option) (
	*Chain, cosmosaccount.Account, error,
) {
	c := &Chain{
		accountName: accountName,
		rpcAddress:  fixRPCAddress(rpcAddress),
		r:           r,
	}

	// apply user options.
	for _, o := range options {
		o(c)
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

	addr, err := acc.Address(c.addressPrefix)
	if err != nil {
		return nil, err
	}

	if err = cosmosfaucet.TryRetrieve(ctx, c.ID, c.rpcAddress, c.faucetAddress, addr); err != nil {
		return nil, err
	}
	return c.r.balance(ctx, c.rpcAddress, c.accountName, c.addressPrefix)
}

func (c *Chain) Config() relayerconfig.Chain {
	return relayerconfig.Chain{
		ID:            c.ID,
		Account:       c.accountName,
		AddressPrefix: c.addressPrefix,
		RPCAddress:    c.rpcAddress,
		GasPrice:      c.gasPrice,
		GasLimit:      c.gasLimit,
		ClientID:      c.clientID,
	}
}

// channelOptions represents options for configuring the IBC channel between two chains.
type channelOptions struct {
	sourcePort    string
	sourceVersion string
	targetPort    string
	targetVersion string
	ordering      string
}

// newChannelOptions returns default channel options.
func newChannelOptions() channelOptions {
	return channelOptions{
		sourcePort:    TransferPort,
		sourceVersion: TransferVersion,
		targetPort:    TransferPort,
		targetVersion: TransferVersion,
		ordering:      OrderingUnordered,
	}
}

// ChannelOption is used to configure relayer IBC connection.
type ChannelOption func(*channelOptions)

// SourcePort configures the source port of the new channel.
func SourcePort(port string) ChannelOption {
	return func(c *channelOptions) {
		c.sourcePort = port
	}
}

// TargetPort configures the target port of the new channel.
func TargetPort(port string) ChannelOption {
	return func(c *channelOptions) {
		c.targetPort = port
	}
}

// SourceVersion configures the source version of the new channel.
func SourceVersion(version string) ChannelOption {
	return func(c *channelOptions) {
		c.sourceVersion = version
	}
}

// TargetVersion configures the target version of the new channel.
func TargetVersion(version string) ChannelOption {
	return func(c *channelOptions) {
		c.targetVersion = version
	}
}

// Ordered sets the new channel as ordered.
func Ordered() ChannelOption {
	return func(c *channelOptions) {
		c.ordering = OrderingOrdered
	}
}

// Connect connects dst chain to c chain and creates a path in between in offline mode.
// it returns the path id on success otherwise, returns with a non-nil error.
func (c *Chain) Connect(dst *Chain, options ...ChannelOption) (id string, err error) {
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
	pathID := PathID(c.ID, dst.ID)
	var suffix string
	i := 2
	for {
		guess := pathID + suffix
		if _, err := conf.PathByID(guess); err != nil { // guess is unique.
			pathID = guess
			break
		}
		suffix = fmt.Sprintf("-%d", i)
		i++
	}

	confPath := relayerconfig.Path{
		ID:       pathID,
		Ordering: channelOptions.ordering,
		Src: relayerconfig.PathEnd{
			ChainID: c.ID,
			PortID:  channelOptions.sourcePort,
			Version: channelOptions.sourceVersion,
		},
		Dst: relayerconfig.PathEnd{
			ChainID: dst.ID,
			PortID:  channelOptions.targetPort,
			Version: channelOptions.targetVersion,
		},
	}

	conf.Paths = append(conf.Paths, confPath)

	if err := relayerconfig.Save(conf); err != nil {
		return "", err
	}

	return pathID, nil
}

// EnsureChainSetup sets up the new or existing chain.
func (c *Chain) EnsureChainSetup(ctx context.Context) error {
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(c.rpcAddress))
	if err != nil {
		return err
	}
	status, err := client.RPC.Status(ctx)
	if err != nil {
		return err
	}
	c.ID = status.NodeInfo.Network

	confChain := c.Config()
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

// PathID creates path name from chain ids.
func PathID(srcChainID, dstChainID string) string {
	return fmt.Sprintf("%s-%s", srcChainID, dstChainID)
}
