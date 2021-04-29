package xrelayer

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	tsrelayer "github.com/tendermint/starport/starport/pkg/nodetime/ts-relayer"
	"github.com/tendermint/starport/starport/pkg/tendermintrpc"
)

// faucetTimeout used to set a timeout while transferring coins from a faucet.
const faucetTimeout = time.Second * 20

const (
	TransferPort      = "transfer"
	TransferVersion   = "ics20-1"
	OrderingUnordered = "unordered"
	OrderingOrdered   = "ordered"
)

// Chain represents a chain in relayer.
type Chain struct {
	// ID is id of the chain.
	ID string

	// rpcAddress is the node address of tm.
	rpcAddress string

	// faucetAddress is the faucet address to get tokens for relayer accounts.
	faucetAddress string

	// tmclient used to interact with tm apis.
	tmclient tendermintrpc.Client

	// gasPrice is the gas price used when sending transactions to the chain
	gasPrice string
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

// WithGasPrice gives the gas price to use to send transactions to the chain
func WithGasPrice(gasPrice string) Option {
	return func(c *Chain) {
		c.gasPrice = gasPrice
	}
}

// NewChain creates a new chain on relayer or uses the existing matching chain.
func NewChain(ctx context.Context, rpcAddress string, options ...Option) (*Chain, error) {
	c := &Chain{
		rpcAddress: rpcAddress,
		tmclient:   tendermintrpc.New(rpcAddress),
	}

	// apply user options.
	for _, o := range options {
		o(c)
	}

	return c, c.ensureChainSetup(ctx)
}

// Account retrieves the default account on chain.
func (c *Chain) Account(ctx context.Context) (Account, error) {
	var account Account
	err := tsrelayer.Call(ctx, "getDefaultAccount", c.ID, &account)
	return account, err
}

// TryFaucet tries to receive tokens from the faucet. user given faucet address is
// used when it's available. otherwise, TryFaucet tries to guess the address for
// faucet to retrieve coins. a non-nil error is returned when coin retrieval is unsuccessful.
func (c *Chain) TryFaucet(ctx context.Context) error {
	// find faucet url. can be the user given, otherwise it is the guessed one.
	u, err := c.findFaucetURL(ctx)
	if err != nil {
		return err
	}

	// retrieve the default relayer account to ask tokens for.
	account, err := c.Account(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, faucetTimeout)
	defer cancel()

	fc := cosmosfaucet.NewClient(u.String())

	resp, err := fc.Transfer(ctx, cosmosfaucet.TransferRequest{
		AccountAddress: account.Address,
	})
	if err != nil {
		return errors.Wrap(err, "faucet is not operational")
	}
	if resp.Error != "" {
		return fmt.Errorf("faucet is not operational: %s", resp.Error)
	}
	for _, transfer := range resp.Transfers {
		if transfer.Error != "" {
			return fmt.Errorf("faucet is not operational: %s", transfer.Error)
		}
	}

	return nil
}

// Balance returns the balance for default account in the chain.
func (c *Chain) Balance(ctx context.Context) (sdk.Coins, error) {
	var coins sdk.Coins
	err := tsrelayer.Call(ctx, "getDefaultAccountBalance", c.ID, &coins)
	return coins, err
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
func (c *Chain) Connect(ctx context.Context, dst *Chain, options ...ChannelOption) (Path, error) {
	channelOptions := newChannelOptions()

	for _, apply := range options {
		apply(&channelOptions)
	}

	var path Path
	err := tsrelayer.Call(ctx, "createPath", []interface{}{c.ID, dst.ID, channelOptions}, &path)
	return path, err
}

// findFaucetURL finds faucet address by returning the address if given by the user
// otherwise, tries to guess it.
func (c *Chain) findFaucetURL(ctx context.Context) (*url.URL, error) {
	// use if there is a user given faucet address.
	if c.faucetAddress != "" {
		return url.Parse(c.faucetAddress)
	}

	// guess faucet address otherwise.
	guessedURLs, err := c.guessFaucetURLs()
	if err != nil {
		return nil, err
	}

	for _, u := range guessedURLs {
		// check if the potential faucet server accepts connections.
		address := u.Host
		if u.Scheme == "https" {
			address += ":443"
		}
		if _, err := net.DialTimeout("tcp", address, time.Second); err != nil {
			continue
		}

		// ensure that this is a real faucet server.
		info, err := cosmosfaucet.NewClient(u.String()).FaucetInfo(ctx)
		if err != nil || info.ChainID != c.ID || !info.IsAFaucet {
			continue
		}

		return u, nil
	}

	return nil, errors.New("no faucet available, please send coins to the address")
}

// guess tries to guess all possible faucet addresses.
func (c *Chain) guessFaucetURLs() ([]*url.URL, error) {
	u, err := url.Parse(c.rpcAddress)
	if err != nil {
		return nil, err
	}

	var guessedURLs []*url.URL

	possibilities := []struct {
		port         string
		subname      string
		nameSperator string
	}{
		{"4500", "", "."},
		{"", "faucet", "."},
		{"", "4500", "-"}, // Gitpod uses port number as sub domain name.
	}

	// creating guesses addresses by basing RPC address.
	for _, poss := range possibilities {
		guess, _ := url.Parse(u.String())                  // copy the original url.
		for _, scheme := range []string{"http", "https"} { // do for both schemes.
			guess, _ := url.Parse(guess.String()) // copy guess.
			guess.Scheme = scheme

			// try with port numbers.
			if poss.port != "" {
				guess.Host = fmt.Sprintf("%s:%s", u.Hostname(), "4500")
				guessedURLs = append(guessedURLs, guess)
				continue
			}

			// try with subnames.
			if poss.subname != "" {
				bases := []string{
					// try with appending subname to the default name.
					// e.g.: faucet.my.domain.
					u.Hostname(),
				}

				// try with replacing the subname for 1 level.
				// e.g.: faucet.domain.
				sp := strings.SplitN(u.Hostname(), poss.nameSperator, 2)
				if len(sp) == 2 {
					bases = append(bases, sp[1])
				}
				for _, basename := range bases {
					guess, _ := url.Parse(guess.String()) // copy guess.
					guess.Host = fmt.Sprintf("%s%s%s", poss.subname, poss.nameSperator, basename)
					guessedURLs = append(guessedURLs, guess)
				}
			}
		}
	}

	return guessedURLs, nil
}

// ensureChainSetup sets up the new or existing chain.
func (c *Chain) ensureChainSetup(ctx context.Context) error {
	var reply struct {
		ID string `json:"id"`
	}
	err := tsrelayer.Call(ctx, "ensureSetupChain", c.rpcAddress, &reply)
	c.ID = reply.ID
	return err
}
