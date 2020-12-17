package networkbuilder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dariubs/percent"
	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
	"golang.org/x/sync/errgroup"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tendermint/starport/starport/pkg/availableport"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/pkg/lineprefixer"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/pkg/tendermintrpc"
	"github.com/tendermint/starport/starport/pkg/xchisel"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	tendermintrpcAddr = "http://localhost:26657"
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
	branch                   string
	hash                     string
	path                     string
	mustNotInitializedBefore bool
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

// SourceRemoteBranch sets a remote branch as source for the blockchain.
func SourceRemoteBranch(url, branch string) SourceOption {
	return func(o *initOptions) {
		o.url = url
		o.branch = branch
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

// Init initializes blockchain from by source option and init options.
func (b *Builder) Init(ctx context.Context, chainID string, source SourceOption, options ...InitOption) (*Blockchain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return nil, err
	}

	// set options.
	o := &initOptions{}
	source(o)
	for _, option := range options {
		option(o)
	}

	// determine final source configuration.
	var (
		url    = o.url
		hash   = o.hash
		path   = o.path
		branch = o.branch
	)

	if o.isChainIDSource {
		chain, err := b.spnclient.ShowChain(ctx, account.Name, chainID)
		if err != nil {
			return nil, err
		}
		url = chain.URL
		hash = chain.Hash
	}

	// pull the chain.
	b.ev.Send(events.New(events.StatusOngoing, "Pulling the blockchain"))

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
		// use a tempdir to clone the source code inside.
		if path, err = ioutil.TempDir("", ""); err != nil {
			return nil, err
		}

		// prepare clone options.
		gitoptions := &git.CloneOptions{
			URL: url,
		}

		// clone the branch when specificied. this is used by chain coordinators on create.
		// when branch isn't provided, default branch(HEAD) is used.
		// (only branch or hash can be set at the same time).
		if branch != "" {
			gitoptions.ReferenceName = plumbing.NewBranchReferenceName(branch)
			gitoptions.SingleBranch = true
		}
		if repo, err = git.PlainCloneContext(ctx, path, false, gitoptions); err != nil {
			return nil, err
		}

		if hash != "" {
			// checkout to a certain hash when specified. this is used by validators to make sure to the locked
			// version of the blockchain.
			// (only branch or hash can be set at the same time).
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

	b.ev.Send(events.New(events.StatusDone, "Pulled the blockchain"))

	if hash == "" {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		githash = ref.Hash()
	}

	return newBlockchain(ctx, b, chainID, path, url, githash.String(), o.mustNotInitializedBefore)
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
func (b *Builder) StartChain(ctx context.Context, chainID string, flags []string) error {
	chainInfo, err := b.ShowChain(ctx, chainID)
	if err != nil {
		return err
	}

	launchInfo, err := b.LaunchInformation(ctx, chainID)
	if err != nil {
		return err
	}

	// find out the app's name form url.
	u, err := url.Parse(chainInfo.URL)
	if err != nil {
		return err
	}
	importPath := path.Join(u.Host, u.Path)
	path, err := gomodulepath.Parse(importPath)
	if err != nil {
		return err
	}

	app := chain.App{
		ChainID: chainID,
		Name:    path.Root,
		Version: cosmosver.Stargate,
	}
	chainCmd, err := chain.New(app, true, chain.LogSilent)
	if err != nil {
		return err
	}

	if len(launchInfo.GenTxs) == 0 {
		return errors.New("There are no approved validators yet")
	}

	// generate the genesis file for the chain to start
	if err := generateGenesis(ctx, chainInfo, launchInfo, chainCmd); err != nil {
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
	configTomlPath := filepath.Join(chainCmd.Home(), "config/config.toml")
	configToml, err := toml.LoadFile(configTomlPath)
	if err != nil {
		return err
	}
	configToml.Set("p2p.persistent_peers", strings.Join(p2pAddresses, ","))
	configToml.Set("p2p.allow_duplicate_ip", true)
	configTomlFile, err := os.OpenFile(configTomlPath, os.O_RDWR|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	defer configTomlFile.Close()
	if _, err = configToml.WriteTo(configTomlFile); err != nil {
		return err
	}

	// peerCountPrefixer adds peer count prefix to each log line.
	peerCountPrefixer := func(w io.Writer) io.Writer {
		tc := tendermintrpc.New(tendermintrpcAddr)

		return lineprefixer.NewWriter(w, func() string {
			netInfo, err := tc.GetNetInfo(ctx)
			if err != nil {
				return ""
			}
			count := netInfo.ConnectedPeers + 1 // +1 is itself.
			prefix := fmt.Sprintf("%d (%v%%) peers online ", count, math.Trunc(percent.PercentOf(count, len(p2pAddresses))))
			return color.New(color.FgYellow).SprintFunc()(prefix)
		})
	}

	g, ctx := errgroup.WithContext(ctx)

	// run the start command of the chain.
	g.Go(func() error {
		return cmdrunner.New().Run(ctx, step.New(
			step.Exec(
				app.D(),
				append([]string{"start"}, flags...)...,
			),
			step.Stdout(peerCountPrefixer(os.Stdout)),
			step.Stderr(peerCountPrefixer(os.Stderr)),
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

// generateGenesis generate the genesis from the launch information in the specified app home
func generateGenesis(ctx context.Context, chainInfo spn.Chain, launchInfo spn.LaunchInformation, chainCmd *chain.Chain) error {
	// overwrite genesis with initial genesis.
	initialGenesis, err := ioutil.ReadFile(initialGenesisPath(chainCmd.Home()))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(genesisPath(chainCmd.Home()), initialGenesis, 755)
	if err != nil {
		return err
	}

	// make sure that Genesis' genesis_time is set to chain's creation time on SPN.
	cf := confile.New(confile.DefaultJSONEncodingCreator, genesisPath(chainCmd.Home()))
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
			Address: account.Address.String(),
			Coins:   account.Coins.String(),
		}

		if err := chainCmd.AddGenesisAccount(ctx, genesisAccount); err != nil {
			return err
		}
	}

	// reset gentx directory
	os.Mkdir(filepath.Join(chainCmd.Home(), "config/gentx"), os.ModePerm)
	dir, err := ioutil.ReadDir(filepath.Join(chainCmd.Home(), "config/gentx"))
	if err != nil {
		return err
	}

	// remove all the current gentxs
	for _, d := range dir {
		if err := os.RemoveAll(filepath.Join(chainCmd.Home(), "config/gentx", d.Name())); err != nil {
			return err
		}
	}

	// add and collect the gentxs
	for i, gentx := range launchInfo.GenTxs {
		// Save the gentx in the gentx directory
		gentxPath := filepath.Join(chainCmd.Home(), fmt.Sprintf("config/gentx/gentx%v.json", i))
		if err = ioutil.WriteFile(gentxPath, gentx, 0666); err != nil {
			return err
		}
	}
	if len(launchInfo.GenTxs) > 0 {
		if err = chainCmd.CollectGentx(ctx); err != nil {
			return err
		}
	}

	return nil
}
