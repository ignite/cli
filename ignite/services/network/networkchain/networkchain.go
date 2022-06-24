package networkchain

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	sperrors "github.com/ignite/cli/ignite/errors"
	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/checksum"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/pkg/gitpod"
	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// Chain represents a network blockchain and lets you interact with its source code and binary.
type Chain struct {
	id       string
	launchID uint64

	path string
	home string

	url         string
	hash        string
	genesisURL  string
	genesisHash string
	launchTime  int64

	keyringBackend chaincmd.KeyringBackend

	isInitialized bool

	ref plumbing.ReferenceName

	chain *chain.Chain
	ev    events.Bus
	ar    cosmosaccount.Registry
}

// SourceOption sets the source for blockchain.
type SourceOption func(*Chain)

// Option sets other initialization options.
type Option func(*Chain)

// SourceRemote sets the default branch on a remote as source for the blockchain.
func SourceRemote(url string) SourceOption {
	return func(c *Chain) {
		c.url = url
	}
}

// SourceRemoteBranch sets the branch on a remote as source for the blockchain.
func SourceRemoteBranch(url, branch string) SourceOption {
	return func(c *Chain) {
		c.url = url
		c.ref = plumbing.NewBranchReferenceName(branch)
	}
}

// SourceRemoteTag sets the tag on a remote as source for the blockchain.
func SourceRemoteTag(url, tag string) SourceOption {
	return func(c *Chain) {
		c.url = url
		c.ref = plumbing.NewTagReferenceName(tag)
	}
}

// SourceRemoteHash uses a remote hash as source for the blockchain.
func SourceRemoteHash(url, hash string) SourceOption {
	return func(c *Chain) {
		c.url = url
		c.hash = hash
	}
}

// SourceLaunch returns a source option for initializing a chain from a launch
func SourceLaunch(launch networktypes.ChainLaunch) SourceOption {
	return func(c *Chain) {
		c.id = launch.ChainID
		c.launchID = launch.ID
		c.url = launch.SourceURL
		c.hash = launch.SourceHash
		c.genesisURL = launch.GenesisURL
		c.genesisHash = launch.GenesisHash
		c.home = ChainHome(launch.ID)
		c.launchTime = launch.LaunchTime
	}
}

// WithHome provides a specific home path for the blockchain for the initialization.
func WithHome(path string) Option {
	return func(c *Chain) {
		c.home = path
	}
}

// WithKeyringBackend provides the keyring backend to use to initialize the blockchain
func WithKeyringBackend(keyringBackend chaincmd.KeyringBackend) Option {
	return func(c *Chain) {
		c.keyringBackend = keyringBackend
	}
}

// WithGenesisFromURL provides a genesis url for the initial genesis of the chain blockchain
func WithGenesisFromURL(genesisURL string) Option {
	return func(c *Chain) {
		c.genesisURL = genesisURL
	}
}

// CollectEvents collects events from the chain.
func CollectEvents(ev events.Bus) Option {
	return func(c *Chain) {
		c.ev = ev
	}
}

// New initializes a network blockchain from source and options.
func New(ctx context.Context, ar cosmosaccount.Registry, source SourceOption, options ...Option) (*Chain, error) {
	c := &Chain{
		ar: ar,
	}
	source(c)
	for _, apply := range options {
		apply(c)
	}

	c.ev.Send(events.New(events.StatusOngoing, "Fetching the source code"))

	var err error
	if c.path, c.hash, err = fetchSource(ctx, c.url, c.ref, c.hash); err != nil {
		return nil, err
	}

	c.ev.Send(events.New(events.StatusDone, "Source code fetched"))
	c.ev.Send(events.New(events.StatusOngoing, "Setting up the blockchain"))

	chainOption := []chain.Option{
		chain.ID(c.id),
		chain.HomePath(c.home),
		chain.LogLevel(chain.LogSilent),
	}

	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	if gitpod.IsOnGitpod() {
		c.keyringBackend = chaincmd.KeyringBackendTest
	}

	chainOption = append(chainOption, chain.KeyringBackend(c.keyringBackend))

	chain, err := chain.New(c.path, chainOption...)
	if err != nil {
		return nil, err
	}

	if !chain.Version.IsFamily(cosmosver.Stargate) {
		return nil, sperrors.ErrOnlyStargateSupported
	}

	c.chain = chain
	c.ev.Send(events.New(events.StatusDone, "Blockchain set up"))

	return c, nil
}

func (c Chain) ChainID() (string, error) {
	return c.chain.ChainID()
}

func (c Chain) ID() (string, error) {
	return c.chain.ID()
}

func (c Chain) Name() string {
	return c.chain.Name()
}

func (c Chain) SetHome(home string) {
	c.chain.SetHome(home)
}

func (c Chain) Home() (path string, err error) {
	return c.chain.Home()
}

func (c Chain) BinaryName() (name string, err error) {
	return c.chain.Binary()
}

func (c Chain) GenesisPath() (path string, err error) {
	return c.chain.GenesisPath()
}

func (c Chain) GentxsPath() (path string, err error) {
	return c.chain.GentxsPath()
}

func (c Chain) DefaultGentxPath() (path string, err error) {
	return c.chain.DefaultGentxPath()
}

func (c Chain) AppTOMLPath() (string, error) {
	return c.chain.AppTOMLPath()
}

func (c Chain) ConfigTOMLPath() (string, error) {
	return c.chain.ConfigTOMLPath()
}

func (c Chain) SourceURL() string {
	return c.url
}

func (c Chain) SourceHash() string {
	return c.hash
}

func (c Chain) IsHomeDirExist() (ok bool, err error) {
	home, err := c.chain.Home()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(home)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// NodeID returns the chain node id
func (c Chain) NodeID(ctx context.Context) (string, error) {
	chainCmd, err := c.chain.Commands(ctx)
	if err != nil {
		return "", err
	}

	nodeID, err := chainCmd.ShowNodeID(ctx)
	if err != nil {
		return "", err
	}
	return nodeID, nil
}

// Build builds chain sources, also checks if source was already built
func (c *Chain) Build(ctx context.Context, cacheStorage cache.Storage) (binaryName string, err error) {
	// if chain was already published and has launch id check binary cache
	if c.launchID != 0 {
		if binaryName, err = c.chain.Binary(); err != nil {
			return "", err
		}
		binaryChecksum, err := checksum.Binary(binaryName)
		if err != nil && !errors.Is(err, exec.ErrNotFound) {
			return "", err
		}
		binaryMatch, err := checkBinaryCacheForLaunchID(c.launchID, binaryChecksum, c.hash)
		if err != nil {
			return "", err
		}
		if binaryMatch {
			return binaryName, nil
		}
	}

	c.ev.Send(events.New(events.StatusOngoing, "Building the chain's binary"))

	// build binary
	if binaryName, err = c.chain.Build(ctx, cacheStorage, ""); err != nil {
		return "", err
	}

	c.ev.Send(events.New(events.StatusDone, "Chain's binary built"))

	// cache built binary for launch id
	if c.launchID != 0 {
		if err := c.CacheBinary(c.launchID); err != nil {
			return "", nil
		}
	}

	return binaryName, nil
}

// CacheBinary caches last built chain binary associated with launch id
func (c *Chain) CacheBinary(launchID uint64) error {
	binaryName, err := c.chain.Binary()
	if err != nil {
		return err
	}
	binaryChecksum, err := checksum.Binary(binaryName)

	if err != nil {
		return err
	}
	return cacheBinaryForLaunchID(launchID, binaryChecksum, c.hash)
}

// fetchSource fetches the chain source from url and returns a temporary path where source is saved
func fetchSource(
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
