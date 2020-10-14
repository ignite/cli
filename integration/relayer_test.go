package integration_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/iowait"
	"github.com/tendermint/starport/starport/pkg/randstr"
	starportconf "github.com/tendermint/starport/starport/services/serve/conf"
	"golang.org/x/sync/errgroup"
)

func TestRelayerWithMultipleChains(t *testing.T) {
	t.Parallel()

	relayerWithMultipleChains(t, 3)
}

func relayerWithMultipleChains(t *testing.T, chainCount int) {
	type Chain struct {
		name, path    string
		servers       starportconf.Servers
		logs, connstr *bytes.Buffer
	}

	var (
		env    = newEnv(t)
		chains []*Chain
	)

	// init & scaffold chains.
	newChain := func() *Chain {
		var (
			name    = randstr.Runes(15)
			path    = env.Scaffold(name, Stargate)
			servers = env.RandomizeServerPorts(path)
		)
		return &Chain{
			name:    name,
			path:    path,
			servers: servers,
			logs:    &bytes.Buffer{},
			connstr: &bytes.Buffer{},
		}
	}
	for i := 0; i < chainCount; i++ {
		chains = append(chains, newChain())
	}

	// serve chains.
	ctx, serveCancel := context.WithCancel(env.Ctx())
	defer serveCancel()
	g, ctx := errgroup.WithContext(ctx)
	for _, chain := range chains {
		chain := chain
		g.Go(func() error {
			ok := env.Serve(fmt.Sprintf("should serve app %q", chain.name),
				chain.path,
				ExecCtx(ctx),
				ExecStdout(chain.logs),
			)
			if !ok {
				return errors.New("cannot serve")
			}
			return nil
		})
	}
	defer func() {
		// wait untill all chains stop serving.
		// a chain will stop serving either by a failure or cancelation.
		// failure is not expected. so, test will exit with error in case of a failure in any of the served chains.
		if err := g.Wait(); err != nil {
			t.FailNow()
		}
	}()

	// wait for chains to be properly served. we could have skip this but having this step
	// is useful to test if chains will restart and detect chains added by `starport chain add`.
	for _, chain := range chains {
		require.NoError(t, env.IsAppServed(ctx, chain.servers), "some apps cannot get online in time")
	}

	// retrieve each chain's relayer connection string.
	for _, chain := range chains {
		env.Must(env.Exec(fmt.Sprintf("should get base64 relayer connection string from chain %q", chain.name),
			step.New(
				step.Exec(
					"starport",
					"chain",
					"me",
				),
				step.Workdir(chain.path),
				step.Stdout(chain.connstr),
			),
			ExecCtx(ctx),
		))
	}

	// cross connect all chains with each other.
	for _, srcchain := range chains {
		for _, dstchain := range chains {
			if srcchain == dstchain {
				continue
			}
			env.Must(env.Exec(fmt.Sprintf("adding chain %q to chain %q", dstchain.name, srcchain.name),
				step.New(
					step.Exec(
						"starport",
						"chain",
						"add",
						dstchain.connstr.String(),
					),
					step.Workdir(srcchain.path),
				),
				ExecCtx(ctx),
			))
		}
	}

	// check each chain's logs to see if they report that they're connected with
	// other chains successfully. we should expect len(chains) x (len(chains) - 1)
	// connections at total. but things aren't stable yet so, for we expect at least len(chains)
	// connections.
	var readers []io.Reader
	for _, chain := range chains {
		readers = append(readers, chain.logs)
	}
	r := io.MultiReader(readers...)
	require.NoError(t, iowait.Untill(r, "linked", len(chains)))
	serveCancel()
}
