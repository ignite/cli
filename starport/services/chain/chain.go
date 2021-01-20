package chain

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	conf "github.com/tendermint/starport/starport/chainconf"
	secretconf "github.com/tendermint/starport/starport/chainconf/secret"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
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
	served         bool
	stdout, stderr io.Writer
}

// chainOptions holds user given options that overwrites chain's defaults.
type chainOptions struct {
	// chainID is the chain's id.
	chainID string

	// homePath of the chain's config dir.
	homePath string

	// cliHomePath of the chain's config dir.
	cliHomePath string

	// keyring backend used by commands if not specified in configuration
	keyringBackend chaincmd.KeyringBackend
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

// CLIHomePath replaces chain's cli configuration home path with given path.
func CLIHomePath(path string) Option {
	return func(c *Chain) {
		c.options.cliHomePath = path
	}
}

// KeyringBackend specify the keyring backend to use for the chain command
func KeyringBackend(keyringBackend chaincmd.KeyringBackend) Option {
	return func(c *Chain) {
		c.options.keyringBackend = keyringBackend
	}
}

// New initializes a new Chain with options that its source lives at path.
func New(ctx context.Context, path string, options ...Option) (*Chain, error) {
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

	// Apply the options
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
	c.plugin = c.pickPlugin()

	// initialize the chain commands
	id, err := c.ID()
	if err != nil {
		return nil, err
	}

	home, err := c.Home()
	if err != nil {
		return nil, err
	}
	ccoptions := []chaincmd.Option{
		chaincmd.WithChainID(id),
		chaincmd.WithHome(home),
		chaincmd.WithVersion(c.Version),
	}
	if c.plugin.Version() == cosmosver.Launchpad {
		cliHome, err := c.CLIHome()
		if err != nil {
			return nil, err
		}
		ccoptions = append(ccoptions,
			chaincmd.WithLaunchpadCLI(c.app.CLI()),
			chaincmd.WithLaunchpadCLIHome(cliHome),
		)
	}

	config, err := c.Config()
	if err != nil {
		return nil, err
	}

	// use keyring backend if specified
	if c.options.keyringBackend != chaincmd.KeyringBackendUnspecified {
		ccoptions = append(ccoptions, chaincmd.WithKeyringBackend(c.options.keyringBackend))
	} else {
		// check if keyring backend is specified in config
		if config.Init.KeyringBackend != "" {
			configKeyringBackend, err := chaincmd.KeyringBackendFromString(config.Init.KeyringBackend)
			if err != nil {
				return nil, err
			}
			ccoptions = append(ccoptions, chaincmd.WithKeyringBackend(configKeyringBackend))
		} else {
			// default keyring backend used is OS
			ccoptions = append(ccoptions, chaincmd.WithKeyringBackend(chaincmd.KeyringBackendOS))
		}
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
	c.cmd, err = chaincmdrunner.New(ctx, cc, ccroptions...)
	if err != nil {
		return nil, err
	}

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

func (c *Chain) StoragePaths() (paths []string, err error) {
	home, err := c.Home()
	if err != nil {
		return paths, err
	}
	paths = append(paths, home)

	cliHome, err := c.CLIHome()
	if err != nil {
		return paths, err
	}
	if cliHome != home {
		paths = append(paths, cliHome)
	}

	return paths, nil
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
func (c *Chain) Home() (string, error) {
	// check if home is explicitly defined for the app
	home := c.options.homePath
	if home == "" {
		// return default home otherwise
		var err error
		home, err = c.DefaultHome()
		if err != nil {
			return "", err
		}

	}

	// expand environment variables in home
	home = filepath.Join(os.ExpandEnv(home))

	return home, nil
}

// DefaultHome returns the blockchain node's default home dir when not specified in the app
func (c *Chain) DefaultHome() (string, error) {
	// check if home is defined in config
	config, err := c.Config()
	if err != nil {
		return "", err
	}
	if config.Init.Home != "" {
		return config.Init.Home, nil
	}

	return c.plugin.Home(), nil
}

// CLIHome returns the blockchain node's home dir.
// This directory is the same as home for Stargate, it is a separate directory for Launchpad
func (c *Chain) CLIHome() (string, error) {
	// check if cli home is explicitly defined for the app
	home := c.options.cliHomePath
	if home == "" {
		// check if cli home is defined in config
		config, err := c.Config()
		if err != nil {
			return "", err
		}
		if config.Init.CLIHome != "" {
			home = config.Init.CLIHome
		} else {
			// Use default for cli home otherwise
			home = c.plugin.CLIHome()
		}
	}

	// expand environment variables in home
	home = filepath.Join(os.ExpandEnv(home))

	// Return default home otherwise
	return home, nil
}

// GenesisPath returns genesis.json path of the app.
func (c *Chain) GenesisPath() (string, error) {
	home, err := c.Home()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/config/genesis.json", home), nil
}

// AppTOMLPath returns app.toml path of the app.
func (c *Chain) AppTOMLPath() (string, error) {
	home, err := c.Home()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/config/app.toml", home), nil
}

// ConfigTOMLPath returns config.toml path of the app.
func (c *Chain) ConfigTOMLPath() (string, error) {
	home, err := c.Home()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/config/config.toml", home), nil
}

// Commands returns the runner execute commands on the chain's binary
func (c *Chain) Commands() chaincmdrunner.Runner {
	return c.cmd
}
