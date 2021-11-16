package network

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"github.com/tendermint/starport/starport/pkg/events"
)

const (
	SPNAddressPrefix = "spn"
	gentxFilename    = "gentx.json"
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
	launchID                 uint64
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

// SourceLaunchID makes source determined by the launch id
func SourceLaunchID(launchID uint64) SourceOption {
	return func(o *initOptions) {
		o.launchID = launchID
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

// AccountRegistry returns the account registry used by the network builder
func (b Builder) AccountRegistry() cosmosaccount.Registry {
	return b.cosmos.AccountRegistry
}

// Blockchain initializes a blockchain from source and options.
func (b *Builder) Blockchain(ctx context.Context, source SourceOption, options ...InitOption) (*Blockchain, error) {
	var o initOptions
	for _, apply := range options {
		apply(&o)
	}
	source(&o)

	var (
		chainID     string
		genesisURL  string
		genesisHash string
		home        = o.homePath
		url         = o.url
		ref         = o.ref
		hash        = o.hash
	)

	// if a launch id is provided, chain information are fetched from Starport Network
	if o.launchID > 0 {
		b.ev.Send(events.New(events.StatusOngoing, "Fetching chain information"))
		chainLaunch, err := b.fetchChainLaunch(ctx, o.launchID)
		if err != nil {
			return nil, err
		}
		b.ev.Send(events.New(events.StatusOngoing, "Chain information fetched"))
		url = chainLaunch.SourceURL
		hash = chainLaunch.SourceHash
		chainID = chainLaunch.GenesisChainID

		// Check if custom genesis URL is provided
		if customGenesisURL := chainLaunch.InitialGenesis.GetGenesisURL(); customGenesisURL != nil {
			genesisURL = customGenesisURL.Url
			genesisHash = customGenesisURL.Hash
		}

		// If no custom home is provided, a default home determined from the launch ID is used
		if home == "" {
			home, err = ChainHome(o.launchID)
			if err != nil {
				return nil, err
			}
		}
	}

	b.ev.Send(events.New(events.StatusOngoing, "Fetching the source code"))
	path, hash, err := b.fetchSource(ctx, url, ref, hash)
	if err != nil {
		return nil, err
	}
	b.ev.Send(events.New(events.StatusDone, "Source code fetched"))

	bc := &Blockchain{
		appPath:     path,
		url:         url,
		hash:        hash,
		builder:     b,
		genesisURL:  genesisURL,
		genesisHash: genesisHash,
	}
	return bc, bc.setup(chainID, home, o.keyringBackend)
}

// fetchChainLaunch fetches the chain launch from Starport Network from a launch id
func (b *Builder) fetchChainLaunch(ctx context.Context, launchID uint64) (launchtypes.Chain, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).Chain(ctx, &launchtypes.QueryGetChainRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return launchtypes.Chain{}, err
	}
	return res.Chain, err
}

// fetchSource fetches the chain source from url and returns a temporary path where source is saved
func (b *Builder) fetchSource(
	ctx context.Context,
	url string,
	ref plumbing.ReferenceName,
	customHash string,
) (path, hash string, err error) {
	var repo *git.Repository

	if path, err = os.MkdirTemp("", ""); err != nil {
		return "", "", err
	}

	// ensure the path for chain source exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", "", err
	}

	// prepare clone options.
	gitoptions := &git.CloneOptions{
		URL: url,
	}

	// clone the ref when specified, this is used by chain coordinators on create.
	if ref != "" {
		gitoptions.ReferenceName = ref
		gitoptions.SingleBranch = true
	}
	if repo, err = git.PlainCloneContext(ctx, path, false, gitoptions); err != nil {
		return "", "", err
	}

	if customHash != "" {
		hash = customHash

		// checkout to a certain hash when specified. this is used by validators to make sure to use
		// the locked version of the blockchain.
		wt, err := repo.Worktree()
		if err != nil {
			return "", "", err
		}
		h, err := repo.ResolveRevision(plumbing.Revision(customHash))
		if err != nil {
			return "", "", err
		}
		githash := *h
		if err := wt.Checkout(&git.CheckoutOptions{
			Hash: githash,
		}); err != nil {
			return "", "", err
		}
	} else {
		// when no specific hash is provided. HEAD is fetched
		ref, err := repo.Head()
		if err != nil {
			return "", "", err
		}
		hash = ref.Hash().String()
	}

	return path, hash, nil
}
