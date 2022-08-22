package envtest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/availableport"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/xurl"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const ServeTimeout = time.Minute * 15

const defaultConfigFileName = "config.yml"

type App struct {
	path       string
	configPath string
	homePath   string

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
		env:  e,
		path: path,
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
	}

	if a.homePath != "" {
		serveCommand = append(serveCommand, "--home", a.homePath)
	}
	if a.configPath != "" {
		serveCommand = append(serveCommand, "--config", a.configPath)
	}

	return a.env.Exec(msg,
		step.NewSteps(step.New(
			step.Exec(IgniteApp, serveCommand...),
			step.Workdir(a.path),
		)),
		options...,
	)
}

// Simulate runs the simulation test for the app
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
			step.Exec(gocmd.Name(), "test", "./..."),
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

	a.EditConfig(func(conf *chainconfig.Config) {
		conf.Faucet.Port = port[0]
		conf.Faucet.Coins = coins
		conf.Faucet.CoinsMax = coinsMax
	})

	addr, err := xurl.HTTP(fmt.Sprintf("0.0.0.0:%d", port[0]))
	require.NoError(a.env.t, err)

	return addr
}

// RandomizeServerPorts randomizes server ports for the app at path, updates
// its config.yml and returns new values.
func (a App) RandomizeServerPorts() chainconfig.Host {
	// generate random server ports and servers list.
	ports, err := availableport.Find(6)
	require.NoError(a.env.t, err)

	genAddr := func(port int) string {
		return fmt.Sprintf("localhost:%d", port)
	}

	servers := chainconfig.Host{
		RPC:     genAddr(ports[0]),
		P2P:     genAddr(ports[1]),
		Prof:    genAddr(ports[2]),
		GRPC:    genAddr(ports[3]),
		GRPCWeb: genAddr(ports[4]),
		API:     genAddr(ports[5]),
	}

	a.EditConfig(func(conf *chainconfig.Config) {
		conf.Host = servers
	})

	return servers
}

// UseRandomHomeDir sets in the blockchain config files generated temporary directories for home directories
// Returns the random home directory
func (a App) UseRandomHomeDir() (homeDirPath string) {
	dir := a.env.TmpDir()

	a.EditConfig(func(conf *chainconfig.Config) {
		conf.Init.Home = dir
	})

	return dir
}

func (a App) EditConfig(apply func(*chainconfig.Config)) {
	f, err := os.OpenFile(a.configPath, os.O_RDWR|os.O_CREATE, 0o755)
	require.NoError(a.env.t, err)
	defer f.Close()

	var conf chainconfig.Config
	require.NoError(a.env.t, yaml.NewDecoder(f).Decode(&conf))

	apply(&conf)

	require.NoError(a.env.t, f.Truncate(0))
	_, err = f.Seek(0, 0)
	require.NoError(a.env.t, err)
	require.NoError(a.env.t, yaml.NewEncoder(f).Encode(conf))
}
