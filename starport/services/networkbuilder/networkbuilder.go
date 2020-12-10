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

// InitBlockchainFromChainID initializes blockchain from chain id.
func (b *Builder) InitBlockchainFromChainID(ctx context.Context, chainID string, mustNotInitializedBefore bool) (*Blockchain, error) {
	account, err := b.AccountInUse()
	if err != nil {
		return nil, err
	}
	chain, err := b.spnclient.ShowChain(ctx, account.Name, chainID)
	if err != nil {
		return nil, err
	}
	return b.InitBlockchainFromURL(ctx, chainID, chain.URL, chain.Hash, mustNotInitializedBefore)
}

// InitBlockchainFromURL initializes blockchain from a remote git repo.
func (b *Builder) InitBlockchainFromURL(ctx context.Context, chainID, url, rev string, mustNotInitializedBefore bool) (*Blockchain, error) {
	appPath, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	b.ev.Send(events.New(events.StatusOngoing, "Pulling the blockchain"))

	// clone the repo.
	repo, err := git.PlainCloneContext(ctx, appPath, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}

	var hash plumbing.Hash

	// checkout to the revision if provided, otherwise default branch is used.
	if rev != "" {
		wt, err := repo.Worktree()
		if err != nil {
			return nil, err
		}
		h, err := repo.ResolveRevision(plumbing.Revision(rev))
		if err != nil {
			return nil, err
		}
		hash = *h
		wt.Checkout(&git.CheckoutOptions{
			Hash: hash,
		})
	} else {
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		hash = ref.Hash()
	}

	b.ev.Send(events.New(events.StatusDone, "Pulled the blockchain"))

	return newBlockchain(ctx, b, chainID, appPath, url, hash.String(), mustNotInitializedBefore)
}

// InitBlockchainFromPath initializes blockchain from a local git repo.
//
// It uses the HEAD(latest commit in currently checked out branch) as the source code of blockchain.
//
// TODO: It requires that there will be no unstaged changes in the code and HEAD is synced with the upstream
// branch (if there is one).
func (b *Builder) InitBlockchainFromPath(ctx context.Context, chainID string, appPath string,
	mustNotInitializedBefore bool) (*Blockchain, error) {
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

	return newBlockchain(ctx, b, chainID, appPath, url, hash.String(), mustNotInitializedBefore)
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

	// get the app home
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	appHome := filepath.Join(homedir, app.ND())

	// generate the genesis file for the chain to start
	if err := generateGenesis(ctx, appHome, chainInfo, launchInfo, chainCmd); err != nil {
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
	configTomlPath := filepath.Join(appHome, "config/config.toml")
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
func generateGenesis(ctx context.Context, appHome string, chainInfo spn.Chain, launchInfo spn.LaunchInformation, chainCmd *chain.Chain) error {
	// overwrite genesis with initial genesis.
	initialGenesis, err := ioutil.ReadFile(initialGenesisPath(appHome))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(genesisPath(appHome), initialGenesis, 755)
	if err != nil {
		return err
	}

	// make sure that Genesis' genesis_time is set to chain's creation time on SPN.
	cf := confile.New(confile.DefaultJSONEncodingCreator, genesisPath(appHome))
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

		if err := chainCmd.AddGenesisAccount(ctx, genesisAccount, appHome); err != nil {
			return err
		}
	}

	// reset gentx directory
	dir, err := ioutil.ReadDir(filepath.Join(appHome, "config/gentx"))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// create the gentx folder if it doesn't exist
			if err := os.Mkdir(filepath.Join(appHome, "config/gentx"), os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// remove all the current gentxs
		for _, d := range dir {
			if err := os.RemoveAll(filepath.Join(appHome, "config/gentx", d.Name())); err != nil {
				return err
			}
		}
	}

	// add and collect the gentxs
	for i, gentx := range launchInfo.GenTxs {
		// Save the gentx in the gentx directory
		gentxPath := filepath.Join(appHome, fmt.Sprintf("config/gentx/gentx%v.json", i))
		if err = ioutil.WriteFile(gentxPath, gentx, 0666); err != nil {
			return err
		}
	}
	if len(launchInfo.GenTxs) > 0 {
		if err = chainCmd.CollectGentx(ctx, appHome); err != nil {
			return err
		}
	}

	return nil
}
