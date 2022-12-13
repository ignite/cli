package chain

import (
	"context"
	"fmt"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	chaincmdrunner "github.com/ignite/cli/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

var (
	// ErrFaucetIsNotEnabled is returned when faucet is not enabled in the config.yml.
	ErrFaucetIsNotEnabled = errors.New("faucet is not enabled in the config.yml")

	// ErrFaucetAccountDoesNotExist returned when specified faucet account in the config.yml does not exist.
	ErrFaucetAccountDoesNotExist = errors.New("specified account (faucet.name) does not exist")
)

var envAPIAddress = os.Getenv("API_ADDRESS")

// Faucet returns the faucet for the chain or an error if the faucet
// configuration is wrong or not configured (not enabled) at all.
func (c *Chain) Faucet(ctx context.Context) (cosmosfaucet.Faucet, error) {
	id, err := c.ID()
	if err != nil {
		return cosmosfaucet.Faucet{}, err
	}

	conf, err := c.Config()
	if err != nil {
		return cosmosfaucet.Faucet{}, err
	}

	commands, err := c.Commands(ctx)
	if err != nil {
		return cosmosfaucet.Faucet{}, err
	}

	// validate if the faucet initialization in the config.yml is correct.
	if conf.Faucet.Name == nil {
		return cosmosfaucet.Faucet{}, ErrFaucetIsNotEnabled
	}

	if _, err := commands.ShowAccount(ctx, *conf.Faucet.Name); err != nil {
		if errors.Is(err, chaincmdrunner.ErrAccountDoesNotExist) {
			return cosmosfaucet.Faucet{}, ErrFaucetAccountDoesNotExist
		}
		return cosmosfaucet.Faucet{}, err
	}

	// construct faucet options.
	validator, err := chainconfig.FirstValidator(conf)
	if err != nil {
		return cosmosfaucet.Faucet{}, err
	}

	servers, err := validator.GetServers()
	if err != nil {
		return cosmosfaucet.Faucet{}, err
	}

	apiAddress := servers.API.Address
	if envAPIAddress != "" {
		apiAddress = envAPIAddress
	}

	apiAddress, err = xurl.HTTP(apiAddress)
	if err != nil {
		return cosmosfaucet.Faucet{}, fmt.Errorf("invalid host api address format: %w", err)
	}

	faucetOptions := []cosmosfaucet.Option{
		cosmosfaucet.Account(*conf.Faucet.Name, "", ""),
		cosmosfaucet.ChainID(id),
		cosmosfaucet.OpenAPI(apiAddress),
	}

	// parse coins to pass to the faucet as coins.
	for _, coin := range conf.Faucet.Coins {
		parsedCoin, err := sdk.ParseCoinNormalized(coin)
		if err != nil {
			return cosmosfaucet.Faucet{}, fmt.Errorf("%w: %s", err, coin)
		}

		var amountMax uint64

		// find out the max amount for this coin.
		for _, coinMax := range conf.Faucet.CoinsMax {
			parsedMax, err := sdk.ParseCoinNormalized(coinMax)
			if err != nil {
				return cosmosfaucet.Faucet{}, fmt.Errorf("%w: %s", err, coin)
			}
			if parsedMax.Denom == parsedCoin.Denom {
				amountMax = parsedMax.Amount.Uint64()
				break
			}
		}

		faucetOptions = append(faucetOptions, cosmosfaucet.Coin(parsedCoin.Amount.Uint64(), amountMax, parsedCoin.Denom))
	}

	if conf.Faucet.RateLimitWindow != "" {
		rateLimitWindow, err := time.ParseDuration(conf.Faucet.RateLimitWindow)
		if err != nil {
			return cosmosfaucet.Faucet{}, fmt.Errorf("%w: %s", err, conf.Faucet.RateLimitWindow)
		}

		faucetOptions = append(faucetOptions, cosmosfaucet.RefreshWindow(rateLimitWindow))
	}

	// init the faucet with options and return.
	return cosmosfaucet.New(ctx, commands, faucetOptions...)
}
