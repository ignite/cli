// Package cosmosfaucet is a faucet to request tokens for sdk accounts.
package cosmosfaucet

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chaincmdrunner "github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

const (
	// DefaultAccountName is the default account to transfer tokens from.
	DefaultAccountName = "faucet"

	// DefaultDenom is the default denomination to distribute.
	DefaultDenom = "uatom"

	// DefaultAmount specifies the default amount to transfer to an account
	// on each request.
	DefaultAmount = 10000000

	// DefaultMaxAmount specifies the maximum amount that can be transferred to an
	// account in all times.
	DefaultMaxAmount = 100000000

	// DefaultRefreshWindow specifies the time after which the max amount limit
	// is refreshed for an account [1 year].
	DefaultRefreshWindow = time.Hour * 24 * 365
)

// Faucet represents a faucet.
type Faucet struct {
	// runner used to interact with blockchain's binary to transfer tokens.
	runner chaincmdrunner.Runner

	// chainID is the chain id of the chain that faucet is operating for.
	chainID string

	// accountName to transfer tokens from.
	accountName string

	// accountMnemonic is the mnemonic of the account.
	accountMnemonic string

	// coinType registered coin type number for HD derivation (BIP-0044).
	coinType string

	// accountNumber registered account number for HD derivation (BIP-0044).
	accountNumber string

	// addressIndex registered address index for HD derivation (BIP-0044).
	addressIndex string

	// coins keeps a list of coins that can be distributed by the faucet.
	coins sdk.Coins

	// coinsMax is a denom-max pair.
	// it holds the maximum amounts of coins that can be sent to a single account.
	coinsMax map[string]sdkmath.Int

	// fee to pay along with the transaction
	feeAmount sdk.Coin

	limitRefreshWindow time.Duration

	// openAPIData holds template data customizations for serving OpenAPI page & spec.
	openAPIData openAPIData

	// version holds the cosmos-sdk version.
	version cosmosver.Version

	// indexerDisabled tells whether the indexing is disabled on the node.
	indexerDisabled bool
}

// Option configures the faucetOptions.
type Option func(*Faucet)

// Account provides the account information to transfer tokens from.
// when mnemonic isn't provided, account assumed to be exists in the keyring.
func Account(name, mnemonic, coinType, accountNumber, addressIndex string) Option {
	return func(f *Faucet) {
		f.accountName = name
		f.accountMnemonic = mnemonic
		f.coinType = coinType
		f.accountNumber = accountNumber
		f.addressIndex = addressIndex
	}
}

// Coin adds a new coin to coins list to distribute by the faucet.
// the first coin added to the list considered as the default coin during transfer requests.
//
// amount is the amount of the coin can be distributed per request.
// maxAmount is the maximum amount of the coin that can be sent to a single account.
// denom is denomination of the coin to be distributed by the faucet.
func Coin(amount, maxAmount sdkmath.Int, denom string) Option {
	return func(f *Faucet) {
		f.coins = append(f.coins, sdk.NewCoin(denom, amount))
		f.coinsMax[denom] = maxAmount
	}
}

// RefreshWindow adds the duration to refresh the transfer limit to the faucet.
func RefreshWindow(refreshWindow time.Duration) Option {
	return func(f *Faucet) {
		f.limitRefreshWindow = refreshWindow
	}
}

// ChainID adds chain id to faucet. faucet will automatically fetch when it isn't provided.
func ChainID(id string) Option {
	return func(f *Faucet) {
		f.chainID = id
	}
}

// FeeAmount sets a fee that it will be paid during the transaction.
func FeeAmount(amount sdkmath.Int, denom string) Option {
	return func(f *Faucet) {
		f.feeAmount = sdk.NewCoin(denom, amount)
	}
}

// OpenAPI configures how to serve Open API page and spec.
func OpenAPI(apiAddress string) Option {
	return func(f *Faucet) {
		f.openAPIData.APIAddress = apiAddress
	}
}

// Version configures the cosmos-sdk version.
func Version(version cosmosver.Version) Option {
	return func(f *Faucet) {
		f.version = version
	}
}

// IndexerDisabled tells whether the indexing is disabled on the node.
// Without indexing, the faucet won't be able to check the limits for each account, nor verify the transaction status.
func IndexerDisabled() Option {
	return func(f *Faucet) {
		f.indexerDisabled = true
	}
}

// New creates a new faucet with ccr (to access and use blockchain's CLI) and given options.
func New(ctx context.Context, ccr chaincmdrunner.Runner, options ...Option) (Faucet, error) {
	f := Faucet{
		runner:      ccr,
		accountName: DefaultAccountName,
		coinsMax:    make(map[string]sdkmath.Int),
		openAPIData: openAPIData{"Blockchain", "http://localhost:1317"},
	}

	for _, apply := range options {
		apply(&f)
	}

	if len(f.coins) == 0 {
		Coin(sdkmath.NewInt(DefaultAmount), sdkmath.NewInt(DefaultMaxAmount), DefaultDenom)(&f)
	}
	f.coins = f.coins.Sort()

	if f.limitRefreshWindow == 0 {
		RefreshWindow(DefaultRefreshWindow)(&f)
	}

	// import the account if mnemonic is provided.
	if f.accountMnemonic != "" {
		_, err := f.runner.AddAccount(
			ctx,
			f.accountName,
			f.accountMnemonic,
			f.coinType,
			f.accountNumber,
			f.addressIndex,
		)
		if err != nil && !errors.Is(err, chaincmdrunner.ErrAccountAlreadyExists) {
			return Faucet{}, err
		}
	}

	if f.chainID == "" {
		status, err := f.runner.Status(ctx)
		if err != nil {
			return Faucet{}, err
		}

		f.chainID = status.ChainID
		f.openAPIData.ChainID = status.ChainID
	}

	return f, nil
}
