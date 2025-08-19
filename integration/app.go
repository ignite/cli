package envtest

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	v1 "github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/availableport"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/pkg/goenv"
	"github.com/ignite/cli/v29/ignite/pkg/xurl"
	"github.com/ignite/cli/v29/ignite/templates/field"
)

const ServeTimeout = time.Minute * 15

const (
	defaultConfigFileName = "config.yml"
	defaultTestTimeout    = 30 * time.Minute // Go's default is 10m
)

type (
	// Hosts contains the "hostname:port" addresses for different service hosts.
	Hosts struct {
		RPC     string
		P2P     string
		Prof    string
		GRPC    string
		GRPCWeb string
		API     string
		Faucet  string
	}

	App struct {
		namespace   string
		name        string
		path        string
		configPath  string
		homePath    string
		testTimeout time.Duration

		env Env

		scaffolded []scaffold
	}

	scaffold struct {
		fields   field.Fields
		index    field.Field
		response field.Fields
		params   field.Fields
		module   string
		name     string
		typeName string
	}
)

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

// ScaffoldApp scaffolds an app to a unique appPath and returns it.
func (e Env) ScaffoldApp(namespace string, flags ...string) App {
	root := e.TmpDir()

	e.Exec("scaffold an app",
		step.NewSteps(step.New(
			step.Exec(
				IgniteApp,
				append([]string{
					"scaffold",
					"chain",
					namespace,
				}, flags...)...,
			),
			step.Workdir(root),
		)),
	)

	var (
		appDirName    = path.Base(namespace)
		appSourcePath = filepath.Join(root, appDirName)
		appHomePath   = e.AppHome(appDirName)
	)

	e.t.Cleanup(func() { os.RemoveAll(appHomePath) })

	return e.App(namespace, appSourcePath, AppHomePath(appHomePath))
}

func (e Env) App(namespace, appPath string, options ...AppOption) App {
	app := App{
		env:         e,
		path:        appPath,
		testTimeout: defaultTestTimeout,
		scaffolded:  make([]scaffold, 0),
		namespace:   namespace,
		name:        path.Base(namespace),
	}

	for _, apply := range options {
		apply(&app)
	}

	if app.configPath == "" {
		app.configPath = filepath.Join(appPath, defaultConfigFileName)
	}

	return app
}

func (a *App) SourcePath() string {
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
func (a *App) Binary() string {
	return path.Base(a.path) + "d"
}

// Serve serves an application lives under path with options where msg describes the
// execution from the serving action.
// unless calling with Must(), Serve() will not exit test runtime on failure.
func (a *App) Serve(msg string, options ...ExecOption) (ok bool) {
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
func (a *App) Simulate(numBlocks, blockSize int) {
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
func (a *App) EnsureSteady() {
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
func (a *App) EnableFaucet(coins, coinsMax []string) (faucetAddr string) {
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
func (a *App) RandomizeServerPorts() Hosts {
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
func (a *App) UseRandomHomeDir() (homeDirPath string) {
	dir := a.env.TmpDir()

	a.EditConfig(func(c *chainconfig.Config) {
		c.Validators[0].Home = dir
	})

	return dir
}

func (a *App) Config() chainconfig.Config {
	bz, err := os.ReadFile(a.configPath)
	require.NoError(a.env.t, err)

	var conf chainconfig.Config
	err = yaml.Unmarshal(bz, &conf)
	require.NoError(a.env.t, err)
	return conf
}

func (a *App) EditConfig(apply func(*chainconfig.Config)) {
	conf := a.Config()
	apply(&conf)

	bz, err := yaml.Marshal(conf)
	require.NoError(a.env.t, err)
	err = os.WriteFile(a.configPath, bz, 0o600)
	require.NoError(a.env.t, err)
}

// GenerateTSClient runs the command to generate the Typescript client code.
func (a *App) GenerateTSClient() bool {
	return a.env.Exec("generate typescript client", step.NewSteps(
		step.New(
			step.Exec(IgniteApp, "g", "ts-client", "--yes", "--clear-cache"),
			step.Workdir(a.path),
		),
	))
}

// MustServe serves the application and ensures success, failing the test if serving fails.
// It uses the provided context to allow cancellation.
func (a *App) MustServe(ctx context.Context) {
	a.env.Must(a.Serve("should serve chain", ExecCtx(ctx)))
}

// Scaffold scaffolds a new module or component in the app and optionally
// validates if it should fail.
// - msg: description of the scaffolding operation.
// - shouldFail: whether the scaffolding is expected to fail.
// - typeName: the type of the scaffold (e.g., "map", "message").
// - args: additional arguments for the scaffold command.
func (a *App) Scaffold(msg string, shouldFail bool, typeName string, args ...string) {
	a.generate(msg, "scaffold", shouldFail, append([]string{typeName}, args...)...)

	if !shouldFail {
		a.addScaffoldCmd(typeName, args...)
	}
}

// Generate executes a code generation command in the app and optionally
// validates if it should fail.
// - msg: description of the generation operation.
// - shouldFail: whether the generation is expected to fail.
// - args: arguments for the generation command.
func (a *App) Generate(msg string, shouldFail bool, args ...string) {
	a.generate(msg, "generate", shouldFail, args...)
}

// generate is a helper method to execute a scaffolding or generation command with the specified options.
// - msg: description of the operation.
// - command: the command to execute (e.g., "scaffold", "generate").
// - shouldFail: whether the command is expected to fail.
// - args: arguments for the command.
func (a *App) generate(msg, command string, shouldFail bool, args ...string) {
	opts := make([]ExecOption, 0)
	if shouldFail {
		opts = append(opts, ExecShouldError())
	}

	args = append([]string{command}, args...)
	a.env.Must(a.env.Exec(msg,
		step.NewSteps(step.New(
			step.Exec(IgniteApp, append(args, "--yes")...),
			step.Workdir(a.SourcePath()),
		)),
		opts...,
	))
}

// addScaffoldCmd processes the scaffold arguments and adds the scaffolded command metadata to the app.
// - typeName: the type of the scaffold (e.g., "map", "message").
// - args: arguments for the scaffold command.
func (a *App) addScaffoldCmd(typeName string, args ...string) {
	module := ""
	index := ""
	response := ""
	params := ""
	name := typeName

	// in the case of scaffolding commands that do no take arguments
	// we can skip the argument parsing
	if len(args) > 0 {
		name = args[0]
		args = args[1:]
	}

	filteredArgs := make([]string, 0)

	// remove the flags from the args
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			break
		}
		filteredArgs = append(filteredArgs, arg)
	}

	// parse the arg flags
	for i, arg := range args {
		// skip tests if the type doesn't need a message
		if arg == "--no-message" {
			return
		}
		if i+1 >= len(args) {
			break
		}
		switch arg {
		case "--module":
			module = args[i+1]
		case "--index":
			index = args[i+1]
		case "--params":
			params = args[i+1]
		case "-r", "--response":
			response = args[i+1]
		}
	}

	argsFields, err := field.ParseFields(filteredArgs, func(string) error { return nil })
	require.NoError(a.env.t, err)

	s := scaffold{
		fields:   argsFields,
		module:   module,
		typeName: typeName,
		name:     name,
	}

	// Handle field specifics based on scaffold type
	switch typeName {
	case "map":
		if index == "" {
			index = "index:string"
		}
		indexFields, err := field.ParseFields(strings.Split(index, ","), func(string) error { return nil })
		require.NoError(a.env.t, err)
		require.Len(a.env.t, indexFields, 1)
		s.index = indexFields[0]
	case "query", "message":
		if response == "" {
			break
		}
		responseFields, err := field.ParseFields(strings.Split(response, ","), func(string) error { return nil })
		require.NoError(a.env.t, err)
		require.Greater(a.env.t, len(responseFields), 0)
		s.response = responseFields
	case "module":
		s.module = name
		if params == "" {
			break
		}
		paramsFields, err := field.ParseFields(strings.Split(params, ","), func(string) error { return nil })
		require.NoError(a.env.t, err)
		require.Greater(a.env.t, len(paramsFields), 0)
		s.params = paramsFields
	case "params":
		s.params = argsFields
	}

	a.scaffolded = append(a.scaffolded, s)
}

// WaitChainUp waits the chain is up.
func (a *App) WaitChainUp(ctx context.Context, chainAPI string) {
	// check the chains is up
	env := a.env
	stepsCheckChains := step.NewSteps(
		step.New(
			step.Exec(
				a.Binary(),
				"config",
				"output", "json",
			),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, chainAPI)
			}),
		),
	)
	env.Exec(fmt.Sprintf("waiting the chain (%s) is up", chainAPI), stepsCheckChains, ExecRetry())
}
