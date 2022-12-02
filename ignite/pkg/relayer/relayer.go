package relayer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/ctxticker"
	tsrelayer "github.com/ignite/cli/ignite/pkg/nodetime/programs/ts-relayer"
	relayerconf "github.com/ignite/cli/ignite/pkg/relayer/config"
	"github.com/ignite/cli/ignite/pkg/xurl"
)

const (
	algoSecp256k1       = "secp256k1"
	ibcSetupGas   int64 = 2256000
	relayDuration       = time.Second * 5
)

// ErrLinkedPath indicates that an IBC path is already liked.
var ErrLinkedPath = errors.New("path already linked")

// Relayer is an IBC relayer.
type Relayer struct {
	ca cosmosaccount.Registry
}

// New creates a new IBC relayer and uses ca to access accounts.
func New(ca cosmosaccount.Registry) Relayer {
	return Relayer{
		ca: ca,
	}
}

// LinkPaths links all chains that has a path from config file to each other.
// paths are optional and acts as a filter to only link some chains.
// calling Link multiple times for the same paths does not have any side effects.
func (r Relayer) LinkPaths(
	ctx context.Context,
	pathIDs ...string,
) error {
	conf, err := relayerconf.Get()
	if err != nil {
		return err
	}

	for _, id := range pathIDs {
		conf, err = r.Link(ctx, conf, id)
		if err != nil {
			// Continue with next path when current one is already linked
			if errors.Is(err, ErrLinkedPath) {
				continue
			}
			return err
		}
		if err := relayerconf.Save(conf); err != nil {
			return err
		}
	}
	return nil
}

// Link links chain path to each other.
func (r Relayer) Link(
	ctx context.Context,
	conf relayerconf.Config,
	pathID string,
) (relayerconf.Config, error) {
	path, err := conf.PathByID(pathID)
	if err != nil {
		return conf, err
	}

	if path.Src.ChannelID != "" {
		return conf, fmt.Errorf("%w: %s", ErrLinkedPath, path.ID)
	}

	if path, err = r.call(ctx, conf, path, "link"); err != nil {
		return conf, err
	}

	return conf, conf.UpdatePath(path)
}

// StartPaths relays packets for linked paths from config file until ctx is canceled.
func (r Relayer) StartPaths(ctx context.Context, pathIDs ...string) error {
	conf, err := relayerconf.Get()
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	var m sync.Mutex // protects relayerconf.Path.
	for _, id := range pathIDs {
		id := id
		g.Go(func() error {
			return r.Start(ctx, conf, id, func(path relayerconf.Config) error {
				m.Lock()
				defer m.Unlock()
				return relayerconf.Save(conf)
			})
		})
	}
	return g.Wait()
}

// Start relays packets for linked path until ctx is canceled.
func (r Relayer) Start(
	ctx context.Context,
	conf relayerconf.Config,
	pathID string,
	postExecute func(path relayerconf.Config) error,
) error {
	return ctxticker.DoNow(ctx, relayDuration, func() error {
		path, err := conf.PathByID(pathID)
		if err != nil {
			return err
		}
		path, err = r.call(ctx, conf, path, "start")
		if err != nil {
			return err
		}
		if err := conf.UpdatePath(path); err != nil {
			return err
		}
		if postExecute != nil {
			return postExecute(conf)
		}
		return nil
	})
}

func (r Relayer) call(
	ctx context.Context,
	conf relayerconf.Config,
	path relayerconf.Path,
	action string,
) (
	reply relayerconf.Path, err error,
) {
	srcChain, srcKey, err := r.prepare(ctx, conf, path.Src.ChainID)
	if err != nil {
		return relayerconf.Path{}, err
	}

	dstChain, dstKey, err := r.prepare(ctx, conf, path.Dst.ChainID)
	if err != nil {
		return relayerconf.Path{}, err
	}

	args := []interface{}{
		path,
		srcChain,
		dstChain,
		srcKey,
		dstKey,
	}
	return reply, tsrelayer.Call(ctx, action, args, &reply)
}

func (r Relayer) prepare(ctx context.Context, conf relayerconf.Config, chainID string) (
	chain relayerconf.Chain, privKey string, err error,
) {
	chain, err = conf.ChainByID(chainID)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	coins, err := r.balance(ctx, chain.RPCAddress, chain.Account, chain.AddressPrefix)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	gasPrice, err := sdk.ParseCoinNormalized(chain.GasPrice)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	account, err := r.ca.GetByName(chain.Account)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	addr, err := account.Address(chain.AddressPrefix)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	errMissingBalance := fmt.Errorf(`account "%s(%s)" on %q chain does not have enough balances`,
		addr,
		chain.Account,
		chain.ID,
	)

	if len(coins) == 0 {
		return relayerconf.Chain{}, "", errMissingBalance
	}

	for _, coin := range coins {
		if gasPrice.Denom != coin.Denom {
			continue
		}

		if gasPrice.Amount.Int64()*ibcSetupGas > coin.Amount.Int64() {
			return relayerconf.Chain{}, "", errMissingBalance
		}
	}

	// Get the key in ASCII armored format
	passphrase := ""
	key, err := r.ca.Export(chain.Account, passphrase)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	// Unarmor the key to be able to read it as bytes
	priv, algo, err := crypto.UnarmorDecryptPrivKey(key, passphrase)
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	// Check the algorithm because the TS relayer expects a secp256k1 private key
	if algo != algoSecp256k1 {
		return relayerconf.Chain{}, "", fmt.Errorf("private key algorithm must be secp256k1 instead of %s", algo)
	}

	return chain, hex.EncodeToString(priv.Bytes()), nil
}

func (r Relayer) balance(ctx context.Context, rpcAddress, account, addressPrefix string) (sdk.Coins, error) {
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddress))
	if err != nil {
		return nil, err
	}

	acc, err := r.ca.GetByName(account)
	if err != nil {
		return nil, err
	}

	addr, err := acc.Address(addressPrefix)
	if err != nil {
		return nil, err
	}

	queryClient := banktypes.NewQueryClient(client.Context())
	res, err := queryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: addr})
	if err != nil {
		return nil, err
	}

	return res.Balances, nil
}

// GetPath returns a path by its id.
func (r Relayer) GetPath(_ context.Context, id string) (relayerconf.Path, error) {
	conf, err := relayerconf.Get()
	if err != nil {
		return relayerconf.Path{}, err
	}

	return conf.PathByID(id)
}

// ListPaths list all the paths.
func (r Relayer) ListPaths(_ context.Context) ([]relayerconf.Path, error) {
	conf, err := relayerconf.Get()
	if err != nil {
		return nil, err
	}

	return conf.Paths, nil
}

func fixRPCAddress(rpcAddress string) string {
	return strings.TrimSuffix(xurl.HTTPEnsurePort(rpcAddress), "/")
}
