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
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/xos"
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
}

// TODO document noCheck (basically it stands to enable Chain initialization without
// need of source code)
func New(app App, noCheck bool, logLevel LogLevel) (*Chain, error) {
	c := &Chain{
		app:            app,
		logLevel:       logLevel,
		serveRefresher: make(chan struct{}, 1),
		stdout:         ioutil.Discard,
		stderr:         ioutil.Discard,
	}

	if logLevel == LogVerbose {
		c.stdout = os.Stdout
		c.stderr = os.Stderr
	}

	var err error

	if !noCheck {
		if _, err := c.Config(); err != nil {
			return nil, errors.New("could not locate a config.yml in your chain. please follow the link for how-to: https://github.com/tendermint/starport/blob/develop/docs/1%20Introduction/4%20Configuration.md")
		}

		c.version, err = c.appVersion()
		if err != nil && err != git.ErrRepositoryNotExists {
			return nil, err
		}
	}

	c.plugin, err = c.pickPlugin()
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
	var paths []string
	for _, name := range conf.FileNames {
		paths = append(paths, filepath.Join(c.app.Path, name))
	}
	confFile, err := xos.OpenFirst(paths...)
	if err != nil {
		return conf.Config{}, errors.Wrap(err, "config file cannot be found")
	}
	defer confFile.Close()
	return conf.Parse(confFile)
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
func (c *Chain) Home() string {
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
