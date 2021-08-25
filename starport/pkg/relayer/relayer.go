package relayer

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/ctxticker"
	tsrelayer "github.com/tendermint/starport/starport/pkg/nodetime/programs/ts-relayer"
	relayerconf "github.com/tendermint/starport/starport/pkg/relayer/config"
	"github.com/tendermint/starport/starport/pkg/xurl"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"golang.org/x/sync/errgroup"
)

const (
	ibcSetupGas   int64 = 2256000
	relayDuration       = time.Second * 5
)

type Relayer struct {
	ca cosmosaccount.Registry
}

func New(ca cosmosaccount.Registry) Relayer {
	r := Relayer{
		ca: ca,
	}

	return r
}

// Link links all chains that has a path to each other.
// paths are optional and acts as a filter to only link some chains.
// calling Link multiple times for the same paths does not have any side effects.
func (r Relayer) Link(ctx context.Context, pathIDs ...string) error {
	conf, err := relayerconf.Get()
	if err != nil {
		return err
	}

	for _, id := range pathIDs {
		path, err := conf.PathByID(id)
		if err != nil {
			return err
		}

		if path.Src.ChannelID != "" { // already linked.
			continue
		}

		if path, err = r.call(ctx, conf, path, "link"); err != nil {
			return err
		}

		if err := conf.UpdatePath(path); err != nil {
			return err
		}
		if err := relayerconf.Save(conf); err != nil {
			return err
		}
	}

	return nil
}

// Start relays packets for linked paths until ctx is canceled.
func (r Relayer) Start(ctx context.Context, pathIDs ...string) error {
	conf, err := relayerconf.Get()
	if err != nil {
		return err
	}

	wg, ctx := errgroup.WithContext(ctx)
	var m sync.Mutex // protects relayerconf.Path.

	start := func(id string) error {
		path, err := conf.PathByID(id)
		if err != nil {
			return err
		}

		if path, err = r.call(ctx, conf, path, "start"); err != nil {
			return err
		}

		m.Lock()
		defer m.Unlock()

		conf, err := relayerconf.Get()
		if err != nil {
			return err
		}

		if err := conf.UpdatePath(path); err != nil {
			return err
		}

		return relayerconf.Save(conf)
	}

	for _, id := range pathIDs {
		id := id

		wg.Go(func() error {
			return ctxticker.DoNow(ctx, relayDuration, func() error { return start(id) })
		})
	}

	return wg.Wait()
}

func (r Relayer) call(ctx context.Context, conf relayerconf.Config, path relayerconf.Path, action string) (
	relayerconf.Path, error) {
	srcChain, srcKey, err := r.prepare(ctx, conf, path.Src.ChainID)
	if err != nil {
		return relayerconf.Path{}, err
	}

	dstChain, dstKey, err := r.prepare(ctx, conf, path.Dst.ChainID)
	if err != nil {
		return relayerconf.Path{}, err
	}

	var reply relayerconf.Path

	err = tsrelayer.Call(ctx, action, []interface{}{
		path,
		srcChain,
		dstChain,
		srcKey,
		dstKey,
	}, &reply)

	return reply, err
}

func (r Relayer) prepare(ctx context.Context, conf relayerconf.Config, chainID string) (
	chain relayerconf.Chain, privKey string, err error) {
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

	for _, coin := range coins {
		if gasPrice.Denom != coin.Denom {
			continue
		}

		if gasPrice.Amount.Int64()*ibcSetupGas > coin.Amount.Int64() {
			err = fmt.Errorf("account %q on %q chain does not have enough balances", chain.Account, chain.ID)
			return relayerconf.Chain{}, "", err
		}
	}

	key, err := r.ca.ExportHex(chain.Account, "")
	if err != nil {
		return relayerconf.Chain{}, "", err
	}

	return chain, key, nil
}

func (r Relayer) balance(ctx context.Context, rpcAddress, account, addressPrefix string) (sdk.Coins, error) {
	context, err := clientCtx(rpcAddress)
	if err != nil {
		return nil, err
	}

	acc, err := r.ca.GetByName(account)
	if err != nil {
		return nil, err
	}

	addr, err := sdk.AccAddressFromBech32(acc.Address(addressPrefix))
	if err != nil {
		return nil, err
	}

	queryClient := banktypes.NewQueryClient(context)
	res, err := queryClient.AllBalances(ctx, banktypes.NewQueryAllBalancesRequest(addr, &query.PageRequest{}))
	if err != nil {
		return nil, err
	}

	return res.Balances, nil
}

// Path represents a path between two chains.
type Path struct {
	// ID is id of the path.
	ID string `json:"id"`

	// IsLinked indicates that chains of these paths are linked or not.
	IsLinked bool `json:"isLinked"`

	// Src end of the path.
	Src PathEnd `json:"src"`

	// Dst end of the path.
	Dst PathEnd `json:"dst"`
}

// PathEnd represents the chain at one side of a Path.
type PathEnd struct {
	ChannelID string `json:"channelID"`
	ChainID   string `json:"chainID"`
	PortID    string `json:"portID"`
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

func rpcClient(rpcAddress string) (*rpchttp.HTTP, error) {
	rpcAddress = fixRPCAddress(rpcAddress)
	return rpchttp.New(rpcAddress, "/websocket")
}

func clientCtx(rpcAddress string) (client.Context, error) {
	rpcClient, err := rpcClient(rpcAddress)
	if err != nil {
		return client.Context{}, err
	}
	cc := cosmosclient.NewContext(rpcClient, io.Discard, "", "")
	return cc, nil
}
