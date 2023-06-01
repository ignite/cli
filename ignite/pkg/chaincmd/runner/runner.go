// Package chaincmdrunner provides high level access to a blockchain's commands.
package chaincmdrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"regexp"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/truncatedbuffer"
)

// Runner provides high level access to a blockchain's commands.
type Runner struct {
	chainCmd       chaincmd.ChainCmd
	stdout, stderr io.Writer
}

// Option configures Runner.
type Option func(r *Runner)

// Stdout sets stdout for executed commands.
func Stdout(w io.Writer) Option {
	return func(runner *Runner) {
		runner.stdout = w
	}
}

// Stderr sets stderr for executed commands.
func Stderr(w io.Writer) Option {
	return func(runner *Runner) {
		runner.stderr = w
	}
}

// New creates a new Runner with cc and options.
func New(ctx context.Context, chainCmd chaincmd.ChainCmd, options ...Option) (Runner, error) {
	runner := Runner{
		chainCmd: chainCmd,
		stdout:   io.Discard,
		stderr:   io.Discard,
	}

	applyOptions(&runner, options)

	// auto detect the chain id and get it applied to chaincmd if auto
	// detection is enabled.
	if chainCmd.IsAutoChainIDDetectionEnabled() {
		status, err := runner.Status(ctx)
		if err != nil {
			return Runner{}, err
		}

		runner.chainCmd = runner.chainCmd.Copy(chaincmd.WithChainID(status.ChainID))
	}

	return runner, nil
}

func applyOptions(r *Runner, options []Option) {
	for _, apply := range options {
		apply(r)
	}
}

// Copy makes a copy of runner by overwriting its options with given options.
func (r Runner) Copy(options ...Option) Runner {
	applyOptions(&r, options)

	return r
}

// Cmd returns underlying chain cmd.
func (r Runner) Cmd() chaincmd.ChainCmd {
	return r.chainCmd
}

type runOptions struct {
	// wrappedStdErrMaxLen determines the maximum length of the wrapped error logs
	// this option is used for long-running command to prevent the buffer containing stderr getting too big
	// 0 can be used for no maximum length
	wrappedStdErrMaxLen int

	// stdout and stderr used to collect a copy of command's outputs.
	stdout, stderr io.Writer

	// stdin defines input for the command
	stdin io.Reader
}

// run executes a command.
func (r Runner) run(ctx context.Context, runOptions runOptions, stepOptions ...step.Option) error {
	// we use a truncated buffer to prevent memory leak
	// this is because app currently send logs to StdErr
	// therefore if the app successfully starts, the written logs can become extensive
	errb := truncatedbuffer.NewTruncatedBuffer(runOptions.wrappedStdErrMaxLen)

	stdout := r.stdout
	if runOptions.stdout != nil {
		stdout = io.MultiWriter(stdout, runOptions.stdout)
	}

	stderr := r.stderr
	if runOptions.stderr != nil {
		stderr = io.MultiWriter(stderr, runOptions.stderr)
	}

	stderr = io.MultiWriter(stderr, errb)

	runnerOptions := []cmdrunner.Option{
		cmdrunner.DefaultStdout(stdout),
		cmdrunner.DefaultStderr(stderr),
	}

	if runOptions.stdin != nil {
		runnerOptions = append(runnerOptions, cmdrunner.DefaultStdin(runOptions.stdin))
	}

	err := cmdrunner.
		New(runnerOptions...).
		Run(ctx, step.New(stepOptions...))

	return errors.Wrap(err, errb.GetBuffer().String())
}

func newBuffer() *buffer {
	return &buffer{
		Buffer: new(bytes.Buffer),
	}
}

// buffer is a bytes.Buffer with additional features.
type buffer struct {
	*bytes.Buffer
}

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// TODO remove this after updating cosmos-sdk to v0.47.3
// https://github.com/cosmos/gogoproto/issues/66#issuecomment-1544699195
func (b *buffer) Bytes() []byte {
	re := regexp.MustCompile(`(?m)^(WARNING:)[\s\S]*?\n`)
	replaced := re.ReplaceAll(b.Buffer.Bytes(), nil)
	return bytes.TrimSpace(replaced)
}

// String returns the contents of the unread portion of the buffer
// as a string.
// TODO remove this after updating cosmos-sdk to v0.47.3
// https://github.com/cosmos/gogoproto/issues/66#issuecomment-1544699195
func (b *buffer) String() string {
	return string(b.Bytes())
}

// JSONEnsuredBytes ensures that encoding format for returned bytes is always
// JSON even if the written data is originally encoded in YAML.
func (b *buffer) JSONEnsuredBytes() ([]byte, error) {
	bz := b.Bytes()

	var out interface{}

	if err := yaml.Unmarshal(bz, &out); err == nil {
		return yaml.YAMLToJSON(bz)
	}

	return bz, nil
}

type txResult struct {
	Code   int    `json:"code"`
	RawLog string `json:"raw_log"`
	TxHash string `json:"txhash"`
}

func decodeTxResult(b *buffer) (txResult, error) {
	var r txResult

	data, err := b.JSONEnsuredBytes()
	if err != nil {
		return r, err
	}

	return r, json.Unmarshal(data, &r)
}
