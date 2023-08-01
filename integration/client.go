package envtest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/stretchr/testify/require"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
)

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
func (a App) RunClientTests(options ...ClientOption) bool {
	npm, err := exec.LookPath("npm")
	require.NoError(a.env.t, err, "npm binary not found")

	// The root dir for the tests must be an absolute path.
	// It is used as the start search point to find test files.
	rootDir, err := os.Getwd()
	require.NoError(a.env.t, err)

	// The filename of this module is required to be able to define the location
	// of the TS client test runner package to be used as working directory when
	// running the tests.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		a.env.t.Fatal("failed to read file name")
	}

	opts := clientOptions{
		env: map[string]string{
			// Absolute path to the blockchain app directory
			"TEST_CHAIN_PATH": a.path,
			// Absolute path to the TS client directory
			"TEST_TSCLIENT_DIR": filepath.Join(a.path, chainconfig.DefaultTSClientPath),
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
	ok = a.env.Exec("install client dependencies", step.NewSteps(
		step.New(
			step.Workdir(filepath.Join(a.path, chainconfig.DefaultTSClientPath)),
			step.Stdout(&output),
			step.Exec(npm, "install"),
			step.PostExec(func(err error) error {
				// Print the npm output when there is an error
				if err != nil {
					a.env.t.Log("\n", output.String())
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
	return a.env.Exec("run client tests", step.NewSteps(
		// Make sure the test runner dependencies are installed
		step.New(
			step.Workdir(runnerDir),
			step.Stdout(&output),
			step.Exec(npm, "install"),
			step.PostExec(func(err error) error {
				// Print the npm output when there is an error
				if err != nil {
					a.env.t.Log("\n", output.String())
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
				a.env.t.Log("\n", output.String())

				return err
			}),
		),
	))
}
