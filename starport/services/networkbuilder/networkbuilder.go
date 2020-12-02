package networkbuilder

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"golang.org/x/sync/errgroup"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/availableport"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/chain"
)

// Builder is network builder.
type Builder struct {
	ev        events.Bus
	spnclient *spn.Client
}

type Option func(*Builder)

// CollectEvents collects events from Builder.
func CollectEvents(ev events.Bus) Option {
	return func(b *Builder) {
		b.ev = ev
	}
}

// New creates a Builder.
func New(spnclient *spn.Client, options ...Option) (*Builder, error) {
	b := &Builder{
		spnclient: spnclient,
	}
	for _, opt := range options {
		opt(b)
	}
	return b, nil
}

// InitBlockchainFromChainID initializes blockchain from chain id.
func (b *Builder) InitBlockchainFromChainID(ctx context.Context, chainID string, mustNotInitializedBefore bool) (*Blockchain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return nil, err
	}
	chain, err := b.spnclient.ChainGet(ctx, account.Name, chainID)
	if err != nil {
		return nil, err
	}
	return b.InitBlockchainFromURL(ctx, chain.URL, chain.Hash, mustNotInitializedBefore)
}

// InitBlockchainFromURL initializes blockchain from a remote git repo.
func (b *Builder) InitBlockchainFromURL(ctx context.Context, url, hash string, mustNotInitializedBefore bool) (*Blockchain, error) {
	appPath, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	b.ev.Send(events.New(events.StatusOngoing, "Pulling the blockchain"))

	// clone the repo first and then checkout to the correct version (hash).
	repo, err := git.PlainCloneContext(ctx, appPath, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	h, err := repo.ResolveRevision(plumbing.Revision(hash))
	if err != nil {
		return nil, err
	}
	wt.Checkout(&git.CheckoutOptions{
		Hash: *h,
	})

	b.ev.Send(events.New(events.StatusDone, "Pulled the blockchain"))

	return newBlockchain(ctx, b, appPath, url, hash, mustNotInitializedBefore)
}

// InitBlockchainFromPath initializes blockchain from a local git repo.
//
// It uses the HEAD(latest commit in currently checked out branch) as the source code of blockchain.
//
// TODO: It requires that there will be no unstaged changes in the code and HEAD is synced with the upstream
// branch (if there is one).
func (b *Builder) InitBlockchainFromPath(ctx context.Context, appPath string, mustNotInitializedBefore bool) (*Blockchain, error) {
	repo, err := git.PlainOpen(appPath)
	if err != nil {
		return nil, err
	}

	// check if there are un-committed changes.
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := wt.Status()
	if err != nil {
		return nil, err
	}
	if !status.IsClean() {
		return nil, errors.New("please either revert or commit your changes")
	}

	// find out remote's url.
	// TODO use the associated upstream branch's remote.
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}
	if len(remotes) == 0 {
		return nil, errors.New("please push your blockchain first")
	}
	remote := remotes[0]
	rc := remote.Config()
	if len(rc.URLs) == 0 {
		return nil, errors.New("cannot find remote's url")
	}
	url := rc.URLs[0]

	// find the hash pointing to HEAD.
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	hash := ref.Hash()
	if err != nil {
		return nil, err
	}

	return newBlockchain(ctx, b, appPath, url, hash.String(), mustNotInitializedBefore)
}

// StartChain downloads the final version version of Genesis on the first start or fails if Genesis
// has not finalized yet.
// After overwriting the downloaded Genesis on top of app's home dir, it starts blockchain by
// executing the start command on its appd binary with optionally provided flags.
func (b *Builder) StartChain(ctx context.Context, chainID string, flags []string) error {
	app := chain.App{
		Name: chainID,
	}
	c, err := chain.New(app, true, chain.LogSilent)
	if err != nil {
		return err
	}

	account, err := b.AccountInUse()
	if err != nil {
		return err
	}
	launchInformation, err := b.spnclient.ChainGet(ctx, account.Name, chainID)
	if err != nil {
		return err
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// overwrite genesis with initial genesis.
	appHome := filepath.Join(homedir, app.ND())
	os.Rename(initialGenesisPath(appHome), genesisPath(appHome))

	// make sure that Genesis' genesis_time is set to chain's creation time on SPN.
	cf := confile.New(confile.DefaultJSONEncodingCreator, genesisPath(appHome))
	var genesis map[string]interface{}
	if err := cf.Load(&genesis); err != nil {
		return err
	}
	genesis["genesis_time"] = launchInformation.CreatedAt.UTC().Format(time.RFC3339)
	if err := cf.Save(genesis); err != nil {
		return err
	}

	// add the genesis accounts
	for _, account := range launchInformation.GenesisAccounts {
		if err = c.AddGenesisAccount(ctx, chain.Account{
			Address: account.Address.String(),
			Coins:   account.Coins.String(),
		}); err != nil {
			return err
		}
	}

	// reset gentx directory
	dir, err := ioutil.ReadDir(filepath.Join(homedir, app.ND(), "config/gentx"))
	if err != nil {
		return err
	}
	for _, d := range dir {
		if err := os.RemoveAll(filepath.Join(homedir, app.ND(), "config/gentx", d.Name())); err != nil {
			return err
		}
	}

	// add and collect the gentxs
	for i, gentx := range launchInformation.GenTxs {
		// Save the gentx in the gentx directory
		gentxPath := filepath.Join(homedir, app.ND(), fmt.Sprintf("config/gentx/gentx%v.json", i))
		if err = ioutil.WriteFile(gentxPath, gentx, 0666); err != nil {
			return err
		}
	}
	if err = c.CollectGentx(ctx); err != nil {
		return err
	}

	// prep peer configs.
	p2pAddresses := launchInformation.Peers
	chiselAddreses := make(map[string]int) // server addr-local p2p port pair.

	if xchisel.IsEnabled() {
		for i, peer := range launchInformation.Peers {
			ports, err := availableport.Find(1)
			if err != nil {
				return err
			}

			localPort := ports[0]
			sp := strings.Split(peer, "@")
			nodeID := sp[0]
			serverAddr := sp[1]

			p2pAddresses[i] = fmt.Sprintf("%s@127.0.0.1:%d", nodeID, localPort)
			chiselAddreses[serverAddr] = localPort
		}
	}

	// save the finalized version of config.toml with peers.
	configTomlPath := filepath.Join(homedir, app.ND(), "config/config.toml")
	configToml, err := toml.LoadFile(configTomlPath)
	if err != nil {
		return err
	}
	configToml.Set("p2p.persistent_peers", strings.Join(p2pAddresses, ","))
	configTomlFile, err := os.OpenFile(configTomlPath, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer configTomlFile.Close()
	if _, err = configToml.WriteTo(configTomlFile); err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	// run the start command of the chain.
	g.Go(func() error {
		return cmdrunner.New().Run(ctx, step.New(
			step.Exec(
				app.D(),
				append([]string{"start"}, flags...)...,
			),
			step.Stdout(os.Stdout),
			step.Stderr(os.Stderr),
		))
	})

	if xchisel.IsEnabled() {
		// start Chisel server.
		g.Go(func() error {
			return xchisel.StartServer(ctx, xchisel.DefaultServerPort)
		})

		// start Chisel clients for all other validators.
		for serverAddr, localPort := range chiselAddreses {
			serverAddr, localPort := serverAddr, localPort
			g.Go(func() error {
				return xchisel.StartClient(ctx, serverAddr, fmt.Sprintf("%d", localPort), "26656")
			})
		}
	}

	return g.Wait()
}
