package chain

import (
	"context"
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

func New(app App, logLevel LogLevel) (*Chain, error) {
	s := &Chain{
		app:            app,
		logLevel:       logLevel,
		serveRefresher: make(chan struct{}, 1),
		stdout:         ioutil.Discard,
		stderr:         ioutil.Discard,
	}

	if logLevel > LogSilent {
		s.stdout = os.Stdout
		s.stderr = os.Stderr
	}

	var err error

	if _, err := s.config(); err != nil {
		return nil, errors.New("could not locate a config.yml in your chain. please follow the link for how-to: https://github.com/tendermint/starport/blob/develop/docs/1%20Introduction/4%20Configuration.md")
	}

	s.version, err = s.appVersion()
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, err
	}

	s.plugin, err = s.pickPlugin()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Chain) appVersion() (v version, err error) {
	repo, err := git.PlainOpen(s.app.Path)
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

// ID returns the chain's id.
func (s *Chain) ID() string {
	return s.app.Name
}

// GenesisPath returns the genesis.json path of chain.
func (c *Chain) GenesisPath() (string, error) {
	return c.plugin.GenesisPath()
}

// rpcPublicAddress points to the public address of Tendermint RPC, this is shared by
// other chains for relayer related actions.
func (s *Chain) rpcPublicAddress() (string, error) {
	rpcAddress := os.Getenv("RPC_ADDRESS")
	if rpcAddress == "" {
		conf, err := s.config()
		if err != nil {
			return "", err
		}
		rpcAddress = conf.Servers.RPCAddr
	}
	return rpcAddress, nil
}

func (s *Chain) config() (conf.Config, error) {
	var paths []string
	for _, name := range conf.FileNames {
		paths = append(paths, filepath.Join(s.app.Path, name))
	}
	confFile, err := xos.OpenFirst(paths...)
	if err != nil {
		return conf.Config{}, errors.Wrap(err, "config file cannot be found")
	}
	defer confFile.Close()
	return conf.Parse(confFile)
}
