package networkbuilder

import (
	"context"
	"errors"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/services/chain"
)

// Builder is network builder.
type Builder struct {
	ev        events.Bus
	spnclient spn.Client
}

type Option func(*Builder)

// CollectEvents collects events from Builder.
func CollectEvents(ev events.Bus) Option {
	return func(b *Builder) {
		b.ev = ev
	}
}

// New creates a Builder.
func New(spnclient spn.Client, options ...Option) (*Builder, error) {
	b := &Builder{
		spnclient: spnclient,
	}
	for _, opt := range options {
		opt(b)
	}
	return b, nil
}

// InitBlockchainFromChainID initializes blockchain from chain id.
func (b *Builder) InitBlockchainFromChainID(ctx context.Context, chainID string) (*Blockchain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return nil, err
	}
	chain, err := b.spnclient.ChainGet(ctx, account.Name, chainID)
	if err != nil {
		return nil, err
	}
	return b.InitBlockchainFromURL(ctx, chain.URL, chain.Hash)
}

// InitBlockchainFromURL initializes blockchain from a remote git repo.
func (b *Builder) InitBlockchainFromURL(ctx context.Context, url, hash string) (*Blockchain, error) {
	appPath, err := ioutil.TempDir("./tmp", "")
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

	return newBlockchain(ctx, b, appPath, url, hash)
}

// InitBlockchainFromPath initializes blockchain from a local git repo.
//
// It uses the HEAD(latest commit in currently checked out branch) as the source code of blockchain.
//
// TODO: It requires that there will be no unstaged changes in the code and HEAD is synced with the upstream
// branch (if there is one).
func (b *Builder) InitBlockchainFromPath(ctx context.Context, appPath string) (*Blockchain, error) {
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

	return newBlockchain(ctx, b, appPath, url, hash.String())
}

// StartChain downloads the final version version of Genesis on the first start or fails if Genesis
// has not finalized yet.
// After overwriting the downloaded Genesis on top of app's home dir, it starts blockchain by
// executing the start command on its appd binary with optionally provided flags.
func (b *Builder) StartChain(ctx context.Context, chainID string, flags []string) error {
	app := chain.App{
		Name: chainID,
	}
	c, err := ConfigGet()
	if err != nil {
		return err
	}

	// save the finalized version of configurations if this isn't done before
	if !c.IsChainMarkedFinalized(chainID) {
		// save the finalized version of Genesis
		account, err := b.AccountInUse()
		if err != nil {
			return err
		}
		chain, err := b.spnclient.ChainGet(ctx, account.Name, chainID)
		if err != nil {
			return err
		}
		homedir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		genesisPath := filepath.Join(homedir, app.ND(), "config/genesis.json")
		if err := ioutil.WriteFile(genesisPath, chain.Genesis, 0644); err != nil {
			return err
		}

		// save the finalized version of config.toml
		configTomlPath := filepath.Join(homedir, app.ND(), "config/config.toml")
		configToml, err := toml.LoadFile(configTomlPath)
		if err != nil {
			return err
		}
		configToml.Set("p2p.persistent_peers", strings.Join(chain.Peers, ","))
		configTomlFile, err := os.OpenFile(configTomlPath, os.O_RDWR|os.O_TRUNC, 644)
		if err != nil {
			return err
		}
		defer configTomlFile.Close()
		_, err = configToml.WriteTo(configTomlFile)
		if err != nil {
			return err
		}

		// mark starport config as finalized and save it
		c.MarkFinalized(chainID)
		if err := ConfigSave(c); err != nil {
			return err
		}
	}

	return cmdrunner.New().Run(ctx, step.New(
		step.Exec(
			app.D(),
			append([]string{"start"}, flags...)...,
		),
		step.Stdout(os.Stdout),
		step.Stderr(os.Stderr),
	))
}
