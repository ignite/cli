package xrelayer

import (
	"context"
	"fmt"
	"sync"
	"time"

	relayercmd "github.com/cosmos/relayer/cmd"
	"github.com/cosmos/relayer/relayer"
	"github.com/tendermint/starport/starport/pkg/ctxticker"
	"github.com/tendermint/starport/starport/pkg/looseerrgroup"
	"golang.org/x/sync/errgroup"
)

const (
	linkingTimeout    = time.Second * 10
	linkingRetryCount = 2
	txRelayFrequency  = time.Second
	maxTxSize         = 2 * relayercmd.MB
	maxMsgLength      = 5
)

// Start relays tx packeR for paths indefinitely until ctx is canceled.
func Start(ctx context.Context, paths ...string) error {
	// start relays all packets for path waiting in the queue.
	start := func(id string) error {
		conf, err := config(ctx, true)
		if err != nil {
			return err
		}

		rpath, err := conf.Paths.Get(id)
		if err != nil {
			return err
		}

		strategy, err := rpath.GetStrategy()
		if err != nil {
			return err
		}

		if naive, ok := strategy.(*relayer.NaiveStrategy); ok {
			naive.MaxTxSize = maxTxSize
			naive.MaxMsgLength = maxMsgLength
			strategy = naive
		}

		chains, src, dst, err := getChainsByPath(conf, id)
		if err != nil {
			return err
		}

		sh, err := relayer.NewSyncHeaders(chains[src], chains[dst])
		if err != nil {
			return err
		}

		sp, err := strategy.UnrelayedSequences(chains[src], chains[dst], sh)
		if err != nil {
			return err
		}

		// nothing to relay.
		if len(sp.Src) == 0 && len(sp.Dst) == 0 {
			return nil
		}

		return strategy.RelayPackets(chains[src], chains[dst], sp, sh)
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, id := range paths {
		id := id

		g.Go(func() error {
			return ctxticker.DoNow(ctx, txRelayFrequency, func() error {
				return start(id)
			})
		})
	}

	return g.Wait()
}

// Link links all chains that has a path to each other.
// paths are optional and acts as a filter to only link pointing chains.
// calling Start multiple times for the same chains does not have any side effects.
func Link(ctx context.Context, paths ...string) (linkedPaths, alreadyLinkedPaths []string, err error) {
	conf, err := config(ctx, false)
	if err != nil {
		return nil, nil, err
	}

	var m sync.Mutex

	link := func(id string) (err error) {
		defer func() {
			if err != nil {
				return
			}
			m.Lock()
			err = cfile.Save(conf)
			m.Unlock()
		}()

		// make sure path is not already linked.
		path, err := GetPath(ctx, id)
		if err != nil {
			return err
		}

		// mark the path as linked now or already linked.
		m.Lock()

		if path.IsLinked {
			alreadyLinkedPaths = append(alreadyLinkedPaths, id)

			m.Unlock()
			return nil
		}

		linkedPaths = append(linkedPaths, id)
		m.Unlock()

		// start linking the path.
		chains, src, dst, err := getChainsByPath(conf, id)
		if err != nil {
			return err
		}

		if _, err := chains[src].CreateClients(chains[dst]); err != nil {
			return err
		}
		if _, err := chains[src].CreateOpenConnections(chains[dst], linkingRetryCount, linkingTimeout); err != nil {
			return err
		}
		if _, err := chains[src].CreateOpenChannels(chains[dst], linkingRetryCount, linkingTimeout); err != nil {
			return err
		}

		return nil
	}

	g := &errgroup.Group{}

	// link non linked paths.
	for _, id := range paths {
		id := id

		g.Go(func() error {
			if err := link(id); err != nil {
				return fmt.Errorf("could not link chains for %q path: %s", id, err.Error())
			}
			return nil
		})
	}

	// cosmos/relayer does not support cancelation so we emulate it here.
	err = looseerrgroup.Wait(ctx, g)

	return linkedPaths, alreadyLinkedPaths, err
}

// Path represents a path between two chains.
type Path struct {
	// ID is id of the chain.
	ID string

	// IsLinked indicates that chains of these paths are linked or not.
	IsLinked bool

	// Path is underlying relayer path.
	Path *relayer.Path
}

// GetPath returns a path by its id.
func GetPath(ctx context.Context, id string) (Path, error) {
	conf, err := config(ctx, false)
	if err != nil {
		return Path{}, err
	}

	path, err := conf.Paths.Get(id)
	if err != nil {
		return Path{}, err
	}

	// find out if path is linked.
	chains, src, dst, err := getChainsByPath(conf, id)
	if err != nil {
		return Path{}, err
	}
	status := path.QueryPathStatus(chains[src], chains[dst]).Status
	isLinked := status.Clients && status.Connection && status.Channel

	return Path{
		ID:       id,
		IsLinked: isLinked,
		Path:     path,
	}, nil
}

// ListPaths list all the paths in relayer's database.
func ListPaths(ctx context.Context) ([]Path, error) {
	conf, err := config(ctx, false)
	if err != nil {
		return nil, err
	}

	var paths []Path

	for id := range conf.Paths {
		path, err := GetPath(ctx, id)
		if err != nil {
			return nil, err
		}

		paths = append(paths, path)
	}

	return paths, nil
}

func getChainsByPath(conf relayercmd.Config, path string) (map[string]*relayer.Chain, string, string, error) {
	pth, err := conf.Paths.Get(path)
	if err != nil {
		return nil, "", "", err
	}

	src, dst := pth.Src.ChainID, pth.Dst.ChainID
	chains, err := conf.Chains.Gets(src, dst)
	if err != nil {
		return nil, "", "", err
	}

	if err = chains[src].SetPath(pth.Src); err != nil {
		return nil, "", "", err
	}
	if err = chains[dst].SetPath(pth.Dst); err != nil {
		return nil, "", "", err
	}

	return chains, src, dst, nil
}
