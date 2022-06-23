package ignitecmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/numbers"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networkchain"
)

// NewNetworkRequestVerify verify the request and simulate the chain.
func NewNetworkRequestVerify() *cobra.Command {
	c := &cobra.Command{
		Use:   "verify [launch-id] [number<,...>]",
		Short: "Verify the request and simulate the chain genesis from them",
		RunE:  networkRequestVerifyHandler,
		Args:  cobra.ExactArgs(2),
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	return c
}

func networkRequestVerifyHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// get the list of request ids
	ids, err := numbers.ParseList(args[1])
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	// verify the requests
	if err := verifyRequest(cmd.Context(), cacheStorage, nb, launchID, ids...); err != nil {
		session.Printf("%s Request(s) %s not valid\n", icons.NotOK, numbers.List(ids, "#"))
		return err
	}

	return session.Printf("%s Request(s) %s verified\n", icons.OK, numbers.List(ids, "#"))
}

// verifyRequest initialize the chain from the launch ID in a temporary directory
// and simulate the launch of the chain from genesis with the request IDs
func verifyRequest(
	ctx context.Context,
	cacheStorage cache.Storage,
	nb NetworkBuilder,
	launchID uint64,
	requestIDs ...uint64,
) error {
	n, err := nb.Network()
	if err != nil {
		return err
	}

	// initialize the chain with a temporary dir
	chainLaunch, err := n.ChainLaunch(ctx, launchID)
	if err != nil {
		return err
	}

	homeDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(homeDir)

	c, err := nb.Chain(
		networkchain.SourceLaunch(chainLaunch),
		networkchain.WithHome(homeDir),
		networkchain.WithKeyringBackend(chaincmd.KeyringBackendTest),
	)
	if err != nil {
		return err
	}

	// fetch the current genesis information and the requests for the chain for simulation
	genesisInformation, err := n.GenesisInformation(ctx, launchID)
	if err != nil {
		return err
	}

	requests, err := n.RequestFromIDs(ctx, launchID, requestIDs...)
	if err != nil {
		return err
	}

	return c.SimulateRequests(
		ctx,
		cacheStorage,
		genesisInformation,
		requests,
	)
}
