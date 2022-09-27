package envtest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/availableport"
	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/xurl"
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

type clientOptions struct {
	env                    map[string]string
	testName, testFilePath string
}

// ClientOption defines options for the TS client test runner.
type ClientOption func(*clientOptions)

// ClientEnv option defines environment values for the tests.
func ClientEnv(env map[string]string) ClientOption {
	return func(o *clientOptions) {
		for k, v := range env {
			o.env[k] = v
		}
	}
}

// ClientTestName option defines a pattern to match the test names that should be run.
func ClientTestName(pattern string) ClientOption {
	return func(o *clientOptions) {
		o.testName = pattern
	}
}

// ClientTestFile option defines the name of the file where to look for tests.
func ClientTestFile(filePath string) ClientOption {
	return func(o *clientOptions) {
		o.testFilePath = filePath
	}
}

// RunClientTests runs the Typescript client tests.
func (e Env) RunClientTests(path string, options ...ClientOption) bool {
	npm, err := exec.LookPath("npm")
	require.NoError(e.t, err, "npm binary not found")

	// The root dir for the tests must be an absolute path.
	// It is used as the start search point to find test files.
	rootDir, err := os.Getwd()
	require.NoError(e.t, err)

	// The filename of this module is required to be able to define the location
	// of the TS client test runner package to be used as working directory when
	// running the tests.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		e.t.Fatal("failed to read file name")
	}

	opts := clientOptions{
		env: map[string]string{
			"TEST_CHAIN_PATH": path,
		},
	}
	for _, o := range options {
		o(&opts)
	}

	var (
		output bytes.Buffer
		env    []string
	)

	//  Install the dependencies needed to run TS client tests
	ok = e.Exec("install client dependencies", step.NewSteps(
		step.New(
			step.Workdir(fmt.Sprintf("%s/vue", path)),
			step.Stdout(&output),
			step.Exec(npm, "install"),
			step.PostExec(func(err error) error {
				// Print the npm output when there is an error
				if err != nil {
					e.t.Log("\n", output.String())
				}

				return err
			}),
		),
	))
	if !ok {
		return false
	}

	output.Reset()

	args := []string{"run", "test", "--", "--dir", rootDir}
	if opts.testName != "" {
		args = append(args, "-t", opts.testName)
	}

	if opts.testFilePath != "" {
		args = append(args, opts.testFilePath)
	}

	for k, v := range opts.env {
		env = append(env, cmdrunner.Env(k, v))
	}

	// The tests are run from the TS client test runner package directory
	runnerDir := filepath.Join(filepath.Dir(filename), "testdata/tstestrunner")

	// TODO: Ignore stderr ? Errors are already displayed with traceback in the stdout
	return e.Exec("run client tests", step.NewSteps(
		// Make sure the test runner dependencies are installed
		step.New(
			step.Workdir(runnerDir),
			step.Stdout(&output),
			step.Exec(npm, "install"),
			step.PostExec(func(err error) error {
				// Print the npm output when there is an error
				if err != nil {
					e.t.Log("\n", output.String())
				}

				return err
			}),
		),
		// Run the TS client tests
		step.New(
			step.Workdir(runnerDir),
			step.Stdout(&output),
			step.Env(env...),
			step.PreExec(func() error {
				// Clear the output from the previous step
				output.Reset()

				return nil
			}),
			step.Exec(npm, args...),
			step.PostExec(func(err error) error {
				// Always print tests output to be available on errors or when verbose is enabled
				e.t.Log("\n", output.String())

				return err
			}),
		),
	))
}
