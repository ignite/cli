package xrelayer

import (
	"context"
	"fmt"
	"sync"
	"time"

	relayercmd "github.com/cosmos/relayer/cmd"
	"github.com/cosmos/relayer/relayer"
	"golang.org/x/sync/errgroup"
)

const (
	linkingTimeout    = time.Second * 10
	linkingRetryCount = 2
)

// Start links all chains that has a path to each other.
// paths are optional and acts as a filter to only link pointing chains.
// calling Start multiple times for the same chains does not have any side effects.
func Start(ctx context.Context, paths ...string) (linkedPaths, alreadyLinkedPaths []string, err error) {
	conf, err := config(ctx)
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

	g, ctx := errgroup.WithContext(ctx)

	for _, id := range paths {
		// link non linked paths.
		id := id

		g.Go(func() (err error) {
			if err = link(id); err != nil {
				return &CouldNotLinkPathError{id, err}
			}
			return nil
		})
	}

	// cosmos/relayer does not support cancelation so we emulate it here.
	doneC := make(chan error)

	go func() { doneC <- g.Wait() }()

	select {
	case <-ctx.Done():
		err = ctx.Err()

	case err = <-doneC:
	}

	return linkedPaths, alreadyLinkedPaths, err
}

// Path represents a path between two chains.
type Path struct {
	// ID is id of the chain.
	ID string

	// IsLinked indicates that chains of these paths are linked or not.
	IsLinked bool
}

// GetPath returns a path by its id.
func GetPath(ctx context.Context, id string) (Path, error) {
	conf, err := config(ctx)
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
	}, nil
}

// ListPaths list all the paths in relayer's database.
func ListPaths(ctx context.Context) ([]Path, error) {
	conf, err := config(ctx)
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

type relayerErr struct {
	message string
}

func newRelayerError(message string) *relayerErr {
	return &relayerErr{message}
}

func (e *relayerErr) Error() string {
	return fmt.Sprintf("relayer error: %s", e.message)
}

type CouldNotLinkPathError struct {
	PathID string

	err error
}

func (e *CouldNotLinkPathError) Unwrap() error { return e.err }

func (e *CouldNotLinkPathError) Error() string {
	return fmt.Sprintf("couldn not link chains for %q path: %s\n", e.PathID, e.err.Error())
}
