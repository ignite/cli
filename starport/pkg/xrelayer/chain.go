package xrelayer

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/tendermintrpc"
)

const (
	//faucetTimeout used to set a timeout while transferring coins from a faucet.
	faucetTimeout = time.Second * 20
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
}

// Account represents an account in relayer.
type Account struct {
	// Address of the account.
	Address string
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

	return c, c.setupChain(ctx)
}

// Account retrieves the default account on chain.
func (c *Chain) Account(ctx context.Context) (Account, error) {
	conf, err := config(ctx, false)
	if err != nil {
		return Account{}, err
	}

	rchain, err := conf.Chains.Get(c.ID)
	if err != nil {
		return Account{}, err
	}

	key, err := rchain.Keybase.Key(defaultAccountKey)
	if err != nil {
		return Account{}, err
	}

	return Account{
		Address: key.GetAddress().String(),
	}, nil
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

// Balance returns the balance for default key in chain
func (c *Chain) Balance(ctx context.Context) (sdk.Coins, error) {
	conf, err := config(ctx, false)
	if err != nil {
		return nil, err
	}

	rchain, err := conf.Chains.Get(c.ID)
	if err != nil {
		return nil, err
	}

	return rchain.QueryBalance(rchain.Key)
}

// Connect connects dst chain to c chain. it returns the path id on success otherwise,
// returns with a non-nil error.
func (c *Chain) Connect(ctx context.Context, dst *Chain) (id string, err error) {
	id = fmt.Sprintf("%s-%s", c.ID, dst.ID)

	conf, err := config(ctx, false)
	if err != nil {
		return "", err
	}

	// construct and add path to paths array.
	conf.Paths[id] = &relayer.Path{
		Strategy: relayer.NewNaiveStrategy(),
		Src: &relayer.PathEnd{
			ChainID: c.ID,
			PortID:  "transfer",
			Order:   "unordered",
			Version: "ics20-1",
		},
		Dst: &relayer.PathEnd{
			ChainID: dst.ID,
			PortID:  "transfer",
			Order:   "unordered",
			Version: "ics20-1",
		},
	}

	// save the config.
	if err := cfile.Save(conf); err != nil {
		return "", err
	}

	// init light clients.
	for _, c := range []*Chain{c, dst} {
		rchain, err := conf.Chains.Get(c.ID)
		if err != nil {
			return "", err
		}

		db, _, err := rchain.NewLightDB()
		if err != nil {
			return "", err
		}

		if _, err := rchain.LightClientWithoutTrust(db); err != nil {
			db.Close()
			return "", err
		}
		db.Close()
	}

	return id, nil
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

// setupChain sets up the new or existing chain.
func (c *Chain) setupChain(ctx context.Context) error {
	if err := c.determineAndSetID(ctx); err != nil {
		return err
	}
	if err := c.ensureAddedToRelayer(ctx); err != nil {
		return err
	}
	if err := c.determineAndSetAccount(ctx); err != nil {
		return err
	}
	return nil
}

// determineAndSetID determines chain's id and uses it.
func (c *Chain) determineAndSetID(ctx context.Context) error {
	genesis, err := c.tmclient.GetGenesis(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot fetch chain info")
	}
	c.ID = genesis.ChainID
	return nil
}

// ensureAddedToRelayer ensures that chain added to relayer's database.
func (c *Chain) ensureAddedToRelayer(ctx context.Context) error {
	conf, err := config(ctx, false)
	if err != nil {
		return err
	}

	if _, err := conf.Chains.Get(c.ID); err != nil { // not configured err
		rchain := &relayer.Chain{
			Key:            "testkey",
			ChainID:        c.ID,
			RPCAddr:        c.rpcAddress,
			AccountPrefix:  "cosmos",
			GasAdjustment:  1.5,
			TrustingPeriod: "336h",
		}

		if err := conf.AddChain(rchain); err != nil {
			return err
		}

		if err := cfile.Save(conf); err != nil {
			return err
		}
	}

	return nil
}

const (
	defaultAccountKey        = "testkey"
	defaultCoinType   uint32 = 118
)

// determineAndSetAccount determines and sets the default account for relayer if
// it wasn't exists.
func (c *Chain) determineAndSetAccount(ctx context.Context) error {
	conf, err := config(ctx, false)
	if err != nil {
		return err
	}

	rchain, err := conf.Chains.Get(c.ID)
	if err != nil {
		return err
	}

	keys, err := rchain.Keybase.List()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if key.GetName() == defaultAccountKey {
			return nil
		}
	}

	mnemonic, err := relayer.CreateMnemonic()
	if err != nil {
		return err
	}

	_, err = rchain.Keybase.NewAccount(defaultAccountKey, mnemonic, "", hd.CreateHDPath(defaultCoinType, 0, 0).String(), hd.Secp256k1)
	return err
}
