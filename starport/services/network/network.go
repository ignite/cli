package network

import (
	"context"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/events"
)

const (
	SPNAddressPrefix = "spn"
)

// Builder is network builder.
type Builder struct {
	ev      events.Bus
	cosmos  cosmosclient.Client
	account cosmosaccount.Account
}

type Option func(*Builder)

// CollectEvents collects events from Builder.
func CollectEvents(ev events.Bus) Option {
	return func(b *Builder) {
		b.ev = ev
	}
}

// New creates a Builder.
func New(cosmos cosmosclient.Client, account cosmosaccount.Account, options ...Option) (*Builder, error) {
	b := &Builder{
		cosmos:  cosmos,
		account: account,
	}
	for _, opt := range options {
		opt(b)
	}
	return b, nil
}

// initOptions holds blockchain initialization options.
type initOptions struct {
	chainID                  string
	url                      string
	ref                      plumbing.ReferenceName
	hash                     string
	mustNotInitializedBefore bool
	homePath                 string
	keyringBackend           chaincmd.KeyringBackend
}

// SourceOption sets the source for blockchain.
type SourceOption func(*initOptions)

// InitOption sets other initialization options.
type InitOption func(*initOptions)

// SourceChainID makes source determined by the chain's id.
func SourceChainID(chainID string) SourceOption {
	return func(o *initOptions) {
		o.chainID = chainID
	}
}

// SourceRemote sets the default branch on a remote as source for the blockchain.
func SourceRemote(url string) SourceOption {
	return func(o *initOptions) {
		o.url = url
	}
}

// SourceRemoteBranch sets the branch on a remote as source for the blockchain.
func SourceRemoteBranch(url, branch string) SourceOption {
	return func(o *initOptions) {
		o.url = url
		o.ref = plumbing.NewBranchReferenceName(branch)
	}
}

// SourceRemoteTag sets the tag on a remote as source for the blockchain.
func SourceRemoteTag(url, tag string) SourceOption {
	return func(o *initOptions) {
		o.url = url
		o.ref = plumbing.NewTagReferenceName(tag)
	}
}

// SourceRemoteHash uses a remote hash as source for the blockchain.
func SourceRemoteHash(url, hash string) SourceOption {
	return func(o *initOptions) {
		o.url = url
		o.hash = hash
	}
}

// MustNotInitializedBefore makes the initialization process fail if data dir for
// the blockchain already exists.
func MustNotInitializedBefore() InitOption {
	return func(o *initOptions) {
		o.mustNotInitializedBefore = true
	}
}

// InitializationHomePath provides a specific home path for the blockchain for the initialization
func InitializationHomePath(homePath string) InitOption {
	return func(o *initOptions) {
		o.homePath = homePath
	}
}

// InitializationKeyringBackend provides the keyring backend to use to initialize the blockchain
func InitializationKeyringBackend(keyringBackend chaincmd.KeyringBackend) InitOption {
	return func(o *initOptions) {
		o.keyringBackend = keyringBackend
	}
}

// Blockchain initializes a blockchain from source and options.
func (b *Builder) Blockchain(ctx context.Context, source SourceOption, options ...InitOption) (*Blockchain, error) {
	var o initOptions
	for _, apply := range options {
		apply(&o)
	}
	source(&o)

	b.ev.Send(events.New(events.StatusOngoing, "Fetching the source code"))

	path, url, hash, err := b.fetch(ctx, o)
	if err != nil {
		return nil, err
	}

	b.ev.Send(events.New(events.StatusDone, "Source code fetched"))

	bc := &Blockchain{
		appPath: path,
		url:     url,
		hash:    hash,
		builder: b,
	}
	return bc, bc.setup(o.chainID, o.homePath, o.keyringBackend)
}

func (b *Builder) fetch(ctx context.Context, o initOptions) (path, url, hash string, err error) {
	// determine final source configuration.
	url = o.url
	ref := o.ref
	ohash := o.hash

	var repo *git.Repository

	if path, err = os.MkdirTemp("", ""); err != nil {
		return "", "", "", err
	}

	// ensure the path for chain source exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", "", "", err
	}

	// prepare clone options.
	gitoptions := &git.CloneOptions{
		URL: url,
	}

	// clone the ref when specificied. this is used by chain coordinators on create.
	if ref != "" {
		gitoptions.ReferenceName = ref
		gitoptions.SingleBranch = true
	}
	if repo, err = git.PlainCloneContext(ctx, path, false, gitoptions); err != nil {
		return "", "", "", err
	}

	if ohash != "" {
		// checkout to a certain hash when specified. this is used by validators to make sure to use
		// the locked version of the blockchain.
		wt, err := repo.Worktree()
		if err != nil {
			return "", "", "", err
		}
		h, err := repo.ResolveRevision(plumbing.Revision(ohash))
		if err != nil {
			return "", "", "", err
		}
		githash := *h
		if err := wt.Checkout(&git.CheckoutOptions{
			Hash: githash,
		}); err != nil {
			return "", "", "", err
		}
	} else {
		ref, err := repo.Head()
		if err != nil {
			return "", "", "", err
		}
		hash = ref.Hash().String()
	}

	return path, url, hash, nil
}

func (b *Builder) fetchChainLaunches(ctx context.Context) ([]launchtypes.Chain, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).ChainAll(ctx, &launchtypes.QueryAllChainRequest{
		
	})
	if err != nil {
		return launchtypes.Chain{}, err
	}
	return res.Chain, err
}