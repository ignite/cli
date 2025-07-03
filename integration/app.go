package envtest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	chainconfig "github.com/ignite/cli/v28/ignite/config/chain"
	v1 "github.com/ignite/cli/v28/ignite/config/chain/v1"
	"github.com/ignite/cli/v28/ignite/pkg/availableport"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/goenv"
	"github.com/ignite/cli/v28/ignite/pkg/xurl"
)

const ServeTimeout = time.Minute * 15

const (
	defaultConfigFileName = "config.yml"
	defaultTestTimeout    = 30 * time.Minute // Go's default is 10m
)

// Hosts contains the "hostname:port" addresses for different service hosts.
type Hosts struct {
	RPC     string
	P2P     string
	Prof    string
	GRPC    string
	GRPCWeb string
	API     string
	Faucet  string
}

type App struct {
	path        string
	configPath  string
	homePath    string
	testTimeout time.Duration

	env Env
}

type AppOption func(*App)

func AppConfigPath(path string) AppOption {
	return func(o *App) {
		o.configPath = path
	}
}

func AppHomePath(path string) AppOption {
	return func(o *App) {
		o.homePath = path
	}
}

func AppTestTimeout(d time.Duration) AppOption {
	return func(o *App) {
		o.testTimeout = d
	}
}

// Scaffold scaffolds an app to a unique appPath and returns it.
func (e Env) Scaffold(name string, flags ...string) App {
	root := e.TmpDir()

	e.Exec("scaffold an app",
		step.NewSteps(step.New(
			step.Exec(
				IgniteApp,
				append([]string{
					"scaffold",
					"chain",
					name,
				}, flags...)...,
			),
			step.Workdir(root),
		)),
	)

	var (
		appDirName    = path.Base(name)
		appSourcePath = filepath.Join(root, appDirName)
		appHomePath   = e.AppHome(appDirName)
	)

	e.t.Cleanup(func() { os.RemoveAll(appHomePath) })

	return e.App(appSourcePath, AppHomePath(appHomePath))
}

func (e Env) App(path string, options ...AppOption) App {
	app := App{
		env:         e,
		path:        path,
		testTimeout: defaultTestTimeout,
	}

	for _, apply := range options {
		apply(&app)
	}

	if app.configPath == "" {
		app.configPath = filepath.Join(path, defaultConfigFileName)
	}

	return app
}

func (a App) SourcePath() string {
	return a.path
}

func (a *App) SetHomePath(homePath string) {
	a.homePath = homePath
}

func (a *App) SetConfigPath(path string) {
	a.configPath = path
}

// Binary returns the binary name of the app. Can be executed directly w/o any
// path after app.Serve is called, since it should be in the $PATH.
func (a App) Binary() string {
	return path.Base(a.path) + "d"
}

// Serve serves an application lives under path with options where msg describes the
// execution from the serving action.
// unless calling with Must(), Serve() will not exit test runtime on failure.
func (a App) Serve(msg string, options ...ExecOption) (ok bool) {
	serveCommand := []string{
		"chain",
		"serve",
		"-v",
		"--quit-on-fail",
	}

	if a.homePath != "" {
		serveCommand = append(serveCommand, "--home", a.homePath)
	}
	if a.configPath != "" {
		serveCommand = append(serveCommand, "--config", a.configPath)
	}
	a.env.t.Cleanup(func() {
		// Serve install the app binary in GOBIN, let's clean that.
		appBinary := path.Join(goenv.Bin(), a.Binary())
		os.Remove(appBinary)
	})

	return a.env.Exec(msg,
		step.NewSteps(step.New(
			step.Exec(IgniteApp, serveCommand...),
			step.Workdir(a.path),
		)),
		options...,
	)
}

// Simulate runs the simulation test for the app.
func (a App) Simulate(numBlocks, blockSize int) {
	a.env.Exec("running the simulation tests",
		step.NewSteps(step.New(
			step.Exec(
				IgniteApp, // TODO
				"chain",
				"simulate",
				"--numBlocks",
				strconv.Itoa(numBlocks),
				"--blockSize",
				strconv.Itoa(blockSize),
			),
			step.Workdir(a.path),
		)),
	)
}

// EnsureSteady ensures that app living at the path can compile and its tests are passing.
func (a App) EnsureSteady() {
	_, statErr := os.Stat(a.configPath)

	require.False(a.env.t, os.IsNotExist(statErr), "config.yml cannot be found")

	a.env.Exec("make sure app is steady",
		step.NewSteps(step.New(
			step.Exec(gocmd.Name(), "test", "-timeout", a.testTimeout.String(), "./..."),
			step.Workdir(a.path),
		)),
	)
}

// EnableFaucet enables faucet by finding a random port for the app faucet and update config.yml
// with this port and provided coins options.
func (a App) EnableFaucet(coins, coinsMax []string) (faucetAddr string) {
	// find a random available port
	port, err := availableport.Find(1)
	require.NoError(a.env.t, err)

	a.EditConfig(func(c *chainconfig.Config) {
		c.Faucet.Port = port[0]
		c.Faucet.Coins = coins
		c.Faucet.CoinsMax = coinsMax
	})

	addr, err := xurl.HTTP(fmt.Sprintf("0.0.0.0:%d", port[0]))
	require.NoError(a.env.t, err)

	return addr
}

// RandomizeServerPorts randomizes server ports for the app at path, updates
// its config.yml and returns new values.
func (a App) RandomizeServerPorts() Hosts {
	// generate random server ports
	ports, err := availableport.Find(7)
	require.NoError(a.env.t, err)

	genAddr := func(port uint) string {
		return fmt.Sprintf("127.0.0.1:%d", port)
	}

	hosts := Hosts{
		RPC:     genAddr(ports[0]),
		P2P:     genAddr(ports[1]),
		Prof:    genAddr(ports[2]),
		GRPC:    genAddr(ports[3]),
		GRPCWeb: genAddr(ports[4]),
		API:     genAddr(ports[5]),
		Faucet:  genAddr(ports[6]),
	}

	a.EditConfig(func(c *chainconfig.Config) {
		c.Faucet.Host = hosts.Faucet

		s := v1.Servers{}
		s.GRPC.Address = hosts.GRPC
		s.GRPCWeb.Address = hosts.GRPCWeb
		s.API.Address = hosts.API
		s.P2P.Address = hosts.P2P
		s.RPC.Address = hosts.RPC
		s.RPC.PProfAddress = hosts.Prof

		v := &c.Validators[0]
		require.NoError(a.env.t, v.SetServers(s))
	})

	return hosts
}

// UseRandomHomeDir sets in the blockchain config files generated temporary directories for home directories.
// Returns the random home directory.
func (a App) UseRandomHomeDir() (homeDirPath string) {
	dir := a.env.TmpDir()

	a.EditConfig(func(c *chainconfig.Config) {
		c.Validators[0].Home = dir
	})

	return dir
}

func (a App) Config() chainconfig.Config {
	bz, err := os.ReadFile(a.configPath)
	require.NoError(a.env.t, err)

	var conf chainconfig.Config
	err = yaml.Unmarshal(bz, &conf)
	require.NoError(a.env.t, err)
	return conf
}

func (a App) EditConfig(apply func(*chainconfig.Config)) {
	conf := a.Config()
	apply(&conf)

	bz, err := yaml.Marshal(conf)
	require.NoError(a.env.t, err)
	err = os.WriteFile(a.configPath, bz, 0o644)
	require.NoError(a.env.t, err)
}

// GenerateTSClient runs the command to generate the Typescript client code.
func (a App) GenerateTSClient() bool {
	return a.env.Exec("generate typescript client", step.NewSteps(
		step.New(
			step.Exec(IgniteApp, "g", "ts-client", "--yes", "--clear-cache"),
			step.Workdir(a.path),
		),
	))
}
