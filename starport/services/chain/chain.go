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
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

var (
	appBackendSourceWatchPaths = []string{
		"app",
		"cmd",
		"x",
		"proto",
		"third_party",
	}

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
	serveCancel    context.CancelFunc
	serveRefresher chan struct{}
	served         bool

	// protoBuiltAtLeastOnce indicates that app's proto generation at least made once.
	protoBuiltAtLeastOnce bool

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

	// isThirdPartyModuleCodegen indicates if proto code generation should be made
	// for 3rd party modules. SDK modules are also considered as a 3rd party.
	isThirdPartyModuleCodegenEnabled bool

	// path of a custom config file
	ConfigFile string
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

// KeyringBackend specifies the keyring backend to use for the chain command
func KeyringBackend(keyringBackend chaincmd.KeyringBackend) Option {
	return func(c *Chain) {
		c.options.keyringBackend = keyringBackend
	}
}

// ConfigFile specifies a custom config file to use
func ConfigFile(configFile string) Option {
	return func(c *Chain) {
		c.options.ConfigFile = configFile
	}
}

// EnableThirdPartyModuleCodegen enables code generation for third party modules,
// including the SDK.
func EnableThirdPartyModuleCodegen() Option {
	return func(c *Chain) {
		c.options.isThirdPartyModuleCodegenEnabled = true
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
		rpcAddress = conf.Host.RPC
	}
	return rpcAddress, nil
}

// StoragePaths returns the home and the cli home (for Launchpad blockchain)
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

// ConfigPath returns the config path of the chain
// Empty string means that the chain has no defined config
func (c *Chain) ConfigPath() string {
	if c.options.ConfigFile != "" {
		return c.options.ConfigFile
	}
	path, err := conf.LocateDefault(c.app.Path)
	if err != nil {
		return ""
	}
	return path
}

// Config returns the config of the chain
func (c *Chain) Config() (conf.Config, error) {
	configPath := c.ConfigPath()
	if configPath == "" {
		return conf.DefaultConf, nil
	}
	return conf.ParseFile(configPath)
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

// Binary returns the name of app's default (appd) binary.
func (c *Chain) Binary() (string, error) {
	conf, err := c.Config()
	if err != nil {
		return "", err
	}

	if conf.Build.Binary != "" {
		return conf.Build.Binary, nil
	}

	return c.app.D(), nil
}

// BinaryCLI returns the name of app's secondary (appcli) binary.
func (c *Chain) BinaryCLI() string {
	return c.app.CLI()
}

// Binaries returns the list of binaries available for the chain.
func (c *Chain) Binaries() ([]string, error) {
	binary, err := c.Binary()
	if err != nil {
		return nil, err
	}

	binaries := []string{
		binary,
	}

	if c.Version.Major().Is(cosmosver.Launchpad) {
		binaries = append(binaries, c.app.CLI())
	}

	return binaries, nil
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
	// Return home dir for Stargate app
	if c.SDKVersion().Is(cosmosver.Stargate) {
		return c.Home()
	}

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
func (c *Chain) Commands(ctx context.Context) (chaincmdrunner.Runner, error) {
	id, err := c.ID()
	if err != nil {
		return chaincmdrunner.Runner{}, err
	}

	home, err := c.Home()
	if err != nil {
		return chaincmdrunner.Runner{}, err
	}

	binary, err := c.Binary()
	if err != nil {
		return chaincmdrunner.Runner{}, err
	}

	config, err := c.Config()
	if err != nil {
		return chaincmdrunner.Runner{}, err
	}

	chainCommandOptions := []chaincmd.Option{
		chaincmd.WithChainID(id),
		chaincmd.WithHome(home),
		chaincmd.WithVersion(c.Version),
		chaincmd.WithNodeAddress(xurl.TCP(config.Host.RPC)),
	}

	if c.plugin.Version() == cosmosver.Launchpad {
		cliHome, err := c.CLIHome()
		if err != nil {
			return chaincmdrunner.Runner{}, err
		}

		chainCommandOptions = append(chainCommandOptions,
			chaincmd.WithLaunchpadCLI(c.BinaryCLI()),
			chaincmd.WithLaunchpadCLIHome(cliHome),
		)
	}

	// use keyring backend if specified
	if c.options.keyringBackend != chaincmd.KeyringBackendUnspecified {
		chainCommandOptions = append(chainCommandOptions, chaincmd.WithKeyringBackend(c.options.keyringBackend))
	} else {
		// check if keyring backend is specified in config
		if config.Init.KeyringBackend != "" {
			configKeyringBackend, err := chaincmd.KeyringBackendFromString(config.Init.KeyringBackend)
			if err != nil {
				return chaincmdrunner.Runner{}, err
			}
			chainCommandOptions = append(chainCommandOptions, chaincmd.WithKeyringBackend(configKeyringBackend))
		} else {
			// default keyring backend used is OS
			chainCommandOptions = append(chainCommandOptions, chaincmd.WithKeyringBackend(chaincmd.KeyringBackendOS))
		}
	}

	cc := chaincmd.New(binary, chainCommandOptions...)

	ccroptions := []chaincmdrunner.Option{}
	if c.logLevel == LogVerbose {
		ccroptions = append(ccroptions,
			chaincmdrunner.Stdout(os.Stdout),
			chaincmdrunner.Stderr(os.Stderr),
			chaincmdrunner.DaemonLogPrefix(c.genPrefix(logAppd)),
			chaincmdrunner.CLILogPrefix(c.genPrefix(logAppcli)),
		)
	}

	return chaincmdrunner.New(ctx, cc, ccroptions...)
}
