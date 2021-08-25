package networkbuilder

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tendermint/starport/starport/pkg/gitpod"
	"github.com/tendermint/starport/starport/pkg/xfilepath"

	"github.com/tendermint/starport/starport/services"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/dariubs/percent"
	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
	"golang.org/x/sync/errgroup"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/availableport"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/ctxticker"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/pkg/tendermintrpc"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	tendermintrpcAddr = "http://localhost:26657"
)

var (
	// spnChainSourcePath returns the path used for the chain source used to build spn chains
	spnChainSourcePath = xfilepath.Join(
		services.StarportConfPath,
		xfilepath.Path("spn-chains"),
	)

	spnChainHomesDir = xfilepath.JoinFromHome(
		xfilepath.Path(".spn-chain-homes"),
	)
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

// initOptions holds blockchain initialization options.
type initOptions struct {
	isChainIDSource          bool
	url                      string
	ref                      plumbing.ReferenceName
	hash                     string
	path                     string
	mustNotInitializedBefore bool
	homePath                 string
	keyringBackend           chaincmd.KeyringBackend
}

// newInitOptions initializes initOptions
func newInitOptions(chainID string, options ...InitOption) (initOpts initOptions, err error) {
	initOpts.homePath, err = xfilepath.Join(
		spnChainHomesDir,
		xfilepath.Path(chainID),
	)()
	if err != nil {
		return initOpts, err
	}

	// set custom options
	for _, option := range options {
		option(&initOpts)
	}

	return initOpts, nil
}

// SourceOption sets the source for blockchain.
type SourceOption func(*initOptions)

// InitOption sets other initialization options.
type InitOption func(*initOptions)

// SourceChainID makes source determined by the chain's id.
func SourceChainID() SourceOption {
	return func(o *initOptions) {
		o.isChainIDSource = true
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

// SourceLocal uses a local git repo as source for the blockchain.
func SourceLocal(path string) SourceOption {
	return func(o *initOptions) {
		o.path = path
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

// Init initializes blockchain from by source option and init options.
func (b *Builder) Init(ctx context.Context, chainID string, source SourceOption, options ...InitOption) (*Blockchain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return nil, err
	}

	o, err := newInitOptions(chainID, options...)
	if err != nil {
		return nil, err
	}
	source(&o)

	// determine final source configuration.
	var (
		url  = o.url
		hash = o.hash
		path = o.path
		ref  = o.ref
	)

	if o.isChainIDSource {
		chain, err := b.spnclient.ShowChain(ctx, account.Name, chainID)
		if err != nil {
			return nil, err
		}

		// verify chain information
		if err := b.VerifyChain(ctx, chain); err != nil {
			return nil, err
		}

		url = chain.URL
		hash = chain.Hash
	}

	// pull the chain.
	b.ev.Send(events.New(events.StatusOngoing, "Fetching the source code"))

	var (
		repo    *git.Repository
		githash plumbing.Hash
	)

	switch {
	// clone git repo from local filesystem. this option only used by chain coordinators.
	case path != "":
		if repo, err = git.PlainOpen(path); err != nil {
			return nil, err
		}
		if url, err = b.ensureRemoteSynced(repo); err != nil {
			return nil, err
		}

	// otherwise clone from the remote. this option can be used by chain coordinators
	// as well as validators.
	default:
		sourcePath, err := spnChainSourcePath()
		if err != nil {
			return nil, err
		}

		// ensure the path for chain source exists
		if err := os.MkdirAll(sourcePath, 0755); err != nil {
			return nil, err
		}

		path = filepath.Join(sourcePath, chainID)
		if _, err := os.Stat(path); err == nil {
			// if the directory already exists, we overwrite it to ensure we have the last version
			if err := os.RemoveAll(path); err != nil {
				return nil, err
			}
		} else if !os.IsNotExist(err) {
			return nil, err
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
			return nil, err
		}

		if hash != "" {
			// checkout to a certain hash when specified. this is used by validators to make sure to use
			// the locked version of the blockchain.
			wt, err := repo.Worktree()
			if err != nil {
				return nil, err
			}
			h, err := repo.ResolveRevision(plumbing.Revision(hash))
			if err != nil {
				return nil, err
			}
			githash = *h
			if err := wt.Checkout(&git.CheckoutOptions{
				Hash: githash,
			}); err != nil {
				return nil, err
			}
		}
	}

	b.ev.Send(events.New(events.StatusDone, "Source code fetched"))

	if hash == "" {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		githash = ref.Hash()
	}

	return newBlockchain(
		ctx,
		b,
		chainID,
		path,
		url,
		githash.String(),
		o.homePath,
		o.keyringBackend,
		o.mustNotInitializedBefore,
	)
}

// ensureRemoteSynced ensures that current worktree in the repository has no unstaged
// changes and synced up with the remote.
// it returns the url of repo or an error related to unstaged changes.
func (b *Builder) ensureRemoteSynced(repo *git.Repository) (url string, err error) {
	// check if there are un-committed changes.
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}
	status, err := wt.Status()
	if err != nil {
		return "", err
	}
	if !status.IsClean() {
		return "", errors.New("please either revert or commit your changes")
	}

	// find out remote's url.
	// TODO use the associated upstream branch's remote.
	remotes, err := repo.Remotes()
	if err != nil {
		return "", err
	}
	if len(remotes) == 0 {
		return "", errors.New("please push your blockchain first")
	}
	remote := remotes[0]
	rc := remote.Config()
	if len(rc.URLs) == 0 {
		return "", errors.New("cannot find remote's url")
	}
	return rc.URLs[0], nil
}

// StartChain downloads the final version version of Genesis on the first start or fails if Genesis
// has not finalized yet.
// After overwriting the downloaded Genesis on top of app's home dir, it starts blockchain by
// executing the start command on its appd binary with optionally provided flags.
func (b *Builder) StartChain(ctx context.Context, chainID string, flags []string, options ...InitOption) error {
	o, err := newInitOptions(chainID, options...)
	if err != nil {
		return err
	}

	chainInfo, err := b.ShowChain(ctx, chainID)
	if err != nil {
		return err
	}

	launchInfo, err := b.LaunchInformation(ctx, chainID)
	if err != nil {
		return err
	}

	chainOption := []chain.Option{
		chain.LogLevel(chain.LogSilent),
	}

	// Custom home paths
	if o.homePath != "" {
		chainOption = append(chainOption, chain.HomePath(o.homePath))
	}

	// use test keyring backend on Gitpod in order to prevent prompting for keyring
	// password. This happens because Gitpod uses containers.
	if gitpod.IsOnGitpod() {
		chainOption = append(chainOption, chain.KeyringBackend(chaincmd.KeyringBackendTest))
	}

	sourcePath, err := spnChainSourcePath()
	if err != nil {
		return err
	}

	appPath := filepath.Join(sourcePath, chainID)
	chainHandler, err := chain.New(ctx, appPath, chainOption...)
	if err != nil {
		return err
	}

	commands, err := chainHandler.Commands(ctx)
	if err != nil {
		return err
	}

	if len(launchInfo.GenTxs) == 0 {
		return errors.New("there are no approved validators yet")
	}

	// generate the genesis file for the chain to start
	if err := generateGenesis(ctx, chainInfo, launchInfo, chainHandler); err != nil {
		return err
	}

	// prep peer configs.
	p2pAddresses := launchInfo.Peers
	chiselAddreses := make(map[string]int) // server addr-local p2p port pair.
	ports, err := availableport.Find(len(launchInfo.Peers))
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2) // make sure that ports are released by the OS before being used.

	if xchisel.IsEnabled() {
		for i, peer := range launchInfo.Peers {
			localPort := ports[i]
			sp := strings.Split(peer, "@")
			nodeID := sp[0]
			serverAddr := sp[1]

			p2pAddresses[i] = fmt.Sprintf("%s@127.0.0.1:%d", nodeID, localPort)
			chiselAddreses[serverAddr] = localPort
		}
	}

	// save the finalized version of config.toml with peers.
	home, err := chainHandler.Home()
	if err != nil {
		return err
	}
	configTomlPath := filepath.Join(home, "config/config.toml")
	configToml, err := toml.LoadFile(configTomlPath)
	if err != nil {
		return err
	}
	configToml.Set("p2p.persistent_peers", strings.Join(p2pAddresses, ","))
	configToml.Set("p2p.allow_duplicate_ip", true)
	configTomlFile, err := os.OpenFile(configTomlPath, os.O_RDWR|os.O_TRUNC, 0644)
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
		return commands.
			Copy(
				chaincmdrunner.Stdout(os.Stdout),
				chaincmdrunner.Stderr(os.Stderr)).
			Start(ctx, flags...)
	})

	// log connected peers info.
	g.Go(func() error {
		tc := tendermintrpc.New(tendermintrpcAddr)

		return ctxticker.DoNow(ctx, time.Second*5, func() error {
			netInfo, err := tc.GetNetInfo(ctx)
			if err == nil {
				count := netInfo.ConnectedPeers + 1 // +1 is itself.
				color.New(color.FgYellow).Printf("%d (%v%%) PEERS ONLINE\n", count, math.Trunc(percent.PercentOf(count, len(p2pAddresses))))
			}
			return nil
		})
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

// GenerateTemporaryGenesis generates the genesis from the launch information in a given temporary directory and return the genesis path
func (b *Builder) GenerateGenesisWithHome(
	ctx context.Context,
	chainID string,
	launchInfo spn.LaunchInformation,
	homeDir string,
) (string, error) {
	chainInfo, err := b.ShowChain(ctx, chainID)
	if err != nil {
		return "", err
	}

	sourcePath, err := spnChainSourcePath()
	if err != nil {
		return "", err
	}

	appPath := filepath.Join(sourcePath, chainID)
	chainHandler, err := chain.New(ctx, appPath,
		chain.HomePath(homeDir),
		chain.LogLevel(chain.LogSilent),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	)
	if err != nil {
		return "", err
	}

	// Run the commands to generate genesis
	if err := generateGenesis(ctx, chainInfo, launchInfo, chainHandler); err != nil {
		return "", err
	}

	return chainHandler.GenesisPath()
}

// generateGenesis generate the genesis from the launch information in the specified app home
func generateGenesis(ctx context.Context, chainInfo spn.Chain, launchInfo spn.LaunchInformation, chainHandler *chain.Chain) error {
	commands, err := chainHandler.Commands(ctx)
	if err != nil {
		return err
	}

	home, err := chainHandler.Home()
	if err != nil {
		return err
	}

	var initialGenesis []byte

	if chainInfo.GenesisURL != "" {
		var hash string
		if initialGenesis, hash, err = genesisAndHashFromURL(ctx, chainInfo.GenesisURL); err != nil {
			return err
		}
		if hash != chainInfo.GenesisHash {
			return errors.New("hash mismatch for the downloaded genesis")
		}
	} else if initialGenesis, err = ioutil.ReadFile(initialGenesisPath(home)); err != nil {
		return err
	}

	// overwrite genesis with initial genesis.
	if err := ioutil.WriteFile(genesisPath(home), initialGenesis, 0755); err != nil {
		return err
	}

	// make sure that Genesis' genesis_time is set to chain's creation time on SPN.
	cf := confile.New(confile.DefaultJSONEncodingCreator, genesisPath(home))
	var genesis map[string]interface{}
	if err := cf.Load(&genesis); err != nil {
		return err
	}
	genesis["genesis_time"] = chainInfo.CreatedAt.UTC().Format(time.RFC3339)
	if err := cf.Save(genesis); err != nil {
		return err
	}

	// add the genesis accounts
	for _, account := range launchInfo.GenesisAccounts {
		genesisAccount := chain.Account{
			Address: account.Address,
			Coins:   account.Coins.String(),
		}

		if err := commands.AddGenesisAccount(ctx, genesisAccount.Address, genesisAccount.Coins); err != nil {
			return err
		}
	}

	// reset gentx directory
	os.Mkdir(filepath.Join(home, "config/gentx"), os.ModePerm)
	dir, err := ioutil.ReadDir(filepath.Join(home, "config/gentx"))
	if err != nil {
		return err
	}

	// remove all the current gentxs
	for _, d := range dir {
		if err := os.RemoveAll(filepath.Join(home, "config/gentx", d.Name())); err != nil {
			return err
		}
	}

	// add and collect the gentxs
	for i, gentx := range launchInfo.GenTxs {
		// Save the gentx in the gentx directory
		gentxPath := filepath.Join(home, fmt.Sprintf("config/gentx/gentx%v.json", i))
		if err = ioutil.WriteFile(gentxPath, gentx, 0666); err != nil {
			return err
		}
	}
	if len(launchInfo.GenTxs) > 0 {
		if err = commands.CollectGentxs(ctx); err != nil {
			return err
		}
	}

	return nil
}
