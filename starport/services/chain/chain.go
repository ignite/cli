package chain

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/services/chain/conf"
	secretconf "github.com/tendermint/starport/starport/services/chain/conf/secret"
)

var (
	appBackendWatchPaths = append([]string{
		"app",
		"cmd",
		"x",
		"proto",
		"third_party",
		secretconf.SecretFile,
	}, conf.FileNames...)

	vuePath = "vue"

	errorColor = color.Red.Render
	infoColor  = color.Yellow.Render
)

type version struct {
	tag  string
	hash string
}

type LogLvl int

const (
	LogSilent LogLvl = iota
	LogRegular
	LogVerbose
)

// Chain provides programatic access and tools for a Cosmos SDK blockchain.
type Chain struct {
	// app holds info about blockchain app.
	app App

	options chainOptions

	Version cosmosver.Version

	plugin         Plugin
	sourceVersion  version
	logLevel       LogLvl
	cmd            chaincmdrunner.Runner
	serveCancel    context.CancelFunc
	serveRefresher chan struct{}
	stdout, stderr io.Writer
}

// chainOptions holds user given options that overwrites chain's defaults.
type chainOptions struct {
	// chainID is the chain's id.
	chainID string

	// homePath of the chain's config dir.
	homePath string
}

// Option configures Chain.
type Option func(*Chain)

// LogLevel sets logging level.
func LogLevel(level LogLvl) Option {
	return func(c *Chain) {
		c.logLevel = level
	}
}

// ID replaces chain's id with given id.
func ID(id string) Option {
	return func(c *Chain) {
		c.options.chainID = id
	}
}

// HomePath replaces chain's configuration home path with given path.
func HomePath(path string) Option {
	return func(c *Chain) {
		c.options.homePath = path
	}
}

// New initializes a new Chain with options that its source lives at path.
func New(path string, options ...Option) (*Chain, error) {
	app, err := NewAppAt(path)
	if err != nil {
		return nil, err
	}

	c := &Chain{
		app:            app,
		logLevel:       LogSilent,
		serveRefresher: make(chan struct{}, 1),
		stdout:         ioutil.Discard,
		stderr:         ioutil.Discard,
	}

	for _, apply := range options {
		apply(c)
	}

	if c.logLevel == LogVerbose {
		c.stdout = os.Stdout
		c.stderr = os.Stderr
	}

	c.sourceVersion, err = c.appVersion()
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, err
	}

	c.Version, err = cosmosver.Detect(c.app.Path)
	if err != nil {
		return nil, err
	}

	// initialize the plugin depending on the version of the chain
	c.plugin, err = c.pickPlugin()
	if err != nil {
		return nil, err
	}

	// initialize the chain commands
	id, err := c.ID()
	if err != nil {
		return nil, err
	}

	ccoptions := []chaincmd.Option{
		chaincmd.WithChainID(id),
		chaincmd.WithHome(c.Home()),
		chaincmd.WithKeyringBackend(chaincmd.KeyringBackendTest),
		chaincmd.WithVersion(c.Version),
	}
	if c.plugin.Version() == cosmosver.Launchpad {
		ccoptions = append(ccoptions,
			chaincmd.WithSecondaryCLI(c.app.CLI()),
			//chaincmd.WithSecondaryCLIHome(),
		)
	}

	cc := chaincmd.New(c.app.D(), ccoptions...)

	ccroptions := []chaincmdrunner.Option{}
	if c.logLevel == LogVerbose {
		ccroptions = append(ccroptions,
			chaincmdrunner.Stdout(os.Stdout),
			chaincmdrunner.Stderr(os.Stderr),
			chaincmdrunner.DaemonLogPrefix(c.genPrefix(logAppd)),
			chaincmdrunner.CLILogPrefix(c.genPrefix(logAppcli)),
		)
	}
	c.cmd = chaincmdrunner.New(cc, ccroptions...)

	return c, nil
}

func (c *Chain) appVersion() (v version, err error) {
	repo, err := git.PlainOpen(c.app.Path)
	if err != nil {
		return version{}, err
	}
	iter, err := repo.Tags()
	if err != nil {
		return version{}, err
	}
	ref, err := iter.Next()
	if err != nil {
		return version{}, nil
	}
	v.tag = strings.TrimPrefix(ref.Name().Short(), "v")
	v.hash = ref.Hash().String()
	return v, nil
}

// SDKVersion returns the version of SDK used to build the blockchain.
func (c *Chain) SDKVersion() cosmosver.MajorVersion {
	return c.plugin.Version()
}

// RPCPublicAddress points to the public address of Tendermint RPC, this is shared by
// other chains for relayer related actions.
func (c *Chain) RPCPublicAddress() (string, error) {
	rpcAddress := os.Getenv("RPC_ADDRESS")
	if rpcAddress == "" {
		conf, err := c.Config()
		if err != nil {
			return "", err
		}
		rpcAddress = conf.Servers.RPCAddr
	}
	return rpcAddress, nil
}

func (c *Chain) StoragePaths() []string {
	return c.plugin.StoragePaths()
}

func (c *Chain) Config() (conf.Config, error) {
	path, err := conf.Locate(c.app.Path)
	if err != nil {
		return conf.DefaultConf, nil
	}
	return conf.ParseFile(path)
}

// ID returns the chain's id.
func (c *Chain) ID() (string, error) {
	// chainID in App has the most priority.
	if c.options.chainID != "" {
		return c.options.chainID, nil
	}

	// otherwise uses defined in config.yml
	chainConfig, err := c.Config()
	if err != nil {
		return "", err
	}
	genid, ok := chainConfig.Genesis["chain_id"]
	if ok {
		return genid.(string), nil
	}

	// use app name by default.
	return c.app.N(), nil
}

// Home returns the blockchain node's home dir.
func (c *Chain) Home() string {
	if c.options.homePath != "" {
		return c.options.homePath
	}
	return c.DefaultHome()
}

// DefaultHome returns the blockchain node's default home dir when not specified.
func (c *Chain) DefaultHome() string {
	return c.plugin.Home()
}

// GenesisPath returns genesis.json path of the app.
func (c *Chain) GenesisPath() string {
	return fmt.Sprintf("%s/config/genesis.json", c.Home())
}

// AppTOMLPath returns app.toml path of the app.
func (c *Chain) AppTOMLPath() string {
	return fmt.Sprintf("%s/config/app.toml", c.Home())
}

// ConfigTOMLPath returns config.toml path of the app.
func (c *Chain) ConfigTOMLPath() string {
	return fmt.Sprintf("%s/config/config.toml", c.Home())
}

// Commands returns the runner execute commands on the chain's binary
func (c *Chain) Commands() chaincmdrunner.Runner {
	return c.cmd
}
