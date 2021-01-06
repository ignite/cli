package chain

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tendermint/starport/starport/pkg/chaincmd"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
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

type LogLevel int

const (
	LogSilent LogLevel = iota
	LogRegular
	LogVerbose
)

type Chain struct {
	app            App
	plugin         Plugin
	version        version
	logLevel       LogLevel
	serveCancel    context.CancelFunc
	serveRefresher chan struct{}
	stdout, stderr io.Writer
	cmd            chaincmd.ChainCmd
	keyringBackend chaincmd.KeyringBackend
}

type Option func(*Chain)

// ChainWithKeyringBackend specify the keyring backend to use for the chain command
func WithKeyringBackend(keyringBackend chaincmd.KeyringBackend) Option {
	return func(c *Chain) {
		c.keyringBackend = keyringBackend
	}
}

// TODO document noCheck (basically it stands to enable Chain initialization without
// need of source code)
func New(app App, noCheck bool, logLevel LogLevel, chainOptions ...Option) (*Chain, error) {
	c := &Chain{
		app:            app,
		logLevel:       logLevel,
		serveRefresher: make(chan struct{}, 1),
		stdout:         ioutil.Discard,
		stderr:         ioutil.Discard,
		keyringBackend: chaincmd.KeyringBackendUnspecified,
	}

	// Apply the options
	for _, applyOption := range chainOptions {
		applyOption(c)
	}

	if logLevel == LogVerbose {
		c.stdout = os.Stdout
		c.stderr = os.Stderr
	}

	var err error

	// Check
	if !noCheck {
		c.version, err = c.appVersion()
		if err != nil && err != git.ErrRepositoryNotExists {
			return nil, err
		}
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
	home, err := c.Home()
	if err != nil {
		return nil, err
	}
	options := append([]chaincmd.Option{}, chaincmd.WithChainID(id), chaincmd.WithHome(home))

	// append Launchpad CLI if Launchpad version is used
	version, err := c.CosmosVersion()
	if err != nil {
		return nil, err
	}
	if version == cosmosver.Launchpad {
		cliHome, err := c.CLIHome()
		if err != nil {
			return nil, err
		}
		options = append(options, chaincmd.WithLaunchpad(c.app.CLI()), chaincmd.WithLaunchpadCLIHome(cliHome))
	}

	// append keyring backend if specified
	if c.keyringBackend != chaincmd.KeyringBackendUnspecified {
		options = append(options, chaincmd.WithKeyringBackend(c.keyringBackend))
	}

	c.cmd = chaincmd.New(
		app.D(),
		options...,
	)

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
	if c.app.ChainID != "" {
		return c.app.ChainID, nil
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
	appHome := c.app.Home()
	if appHome != "" {
		return appHome, nil
	}

	// check if home is defined in config
	config, err := c.Config()
	if err != nil {
		return "", err
	}
	if config.Init.Home != "" {
		return config.Init.Home, nil
	}

	// Return default home otherwise
	return c.DefaultHome(), nil
}

// DefaultHome returns the blockchain node's default home dir when not specified.
func (c *Chain) DefaultHome() string {
	return c.plugin.Home()
}

// CLIHome returns the blockchain node's home dir.
// This directory is the same as home for Stargate, it is a separate directory for Launchpad
func (c *Chain) CLIHome() (string, error) {
	// check if cli home is explicitly defined for the app
	cliHome := c.app.CLIHome()
	if cliHome != "" {
		return cliHome, nil
	}

	// check if cli home is defined in config
	config, err := c.Config()
	if err != nil {
		return "", err
	}
	if config.Init.CLIHome != "" {
		return config.Init.CLIHome, nil
	}

	// Return default home otherwise
	return c.plugin.CLIHome(), nil
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

// Commands returns the chaincmd object to perform command with the chain binary
func (c *Chain) Commands() chaincmd.ChainCmd {
	return c.cmd
}

func (c *Chain) CosmosVersion() (cosmosver.MajorVersion, error) {
	version := c.app.Version
	if version == "" {
		var err error
		version, err = cosmosver.Detect(c.app.Path)
		if err != nil {
			return cosmosver.Stargate, err
		}
	}

	return version, nil
}
