package networkbuilder

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/spn"
)

var (
	starportConfDir = os.ExpandEnv("$HOME/.starport")
	confPath        = filepath.Join(starportConfDir, "networkbuilder")
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
	chain, err := b.spnclient.ShowChain(account.Name, chainID)
	if err != nil {
		return nil, err
	}
	return b.InitBlockchainFromURL(ctx, chain.URL, chain.Hash)
}

// InitBlockchainFromURL initializes blockchain from a remote git repo.
func (b *Builder) InitBlockchainFromURL(ctx context.Context, url, hash string) (*Blockchain, error) {
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
