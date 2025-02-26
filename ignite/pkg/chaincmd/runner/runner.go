// Package chaincmdrunner provides high level access to a blockchain's commands.
package chaincmdrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/truncatedbuffer"
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

// JSONEnsuredBytes ensures that encoding format for returned bytes is always
// JSON even if the written data is originally encoded in YAML.
func (b *buffer) JSONEnsuredBytes() ([]byte, error) {
	bz := b.Bytes()

	// Attempt to find valid JSON in the buffer
	startIndex := strings.IndexAny(string(bz), "{[")
	if startIndex >= 0 {
		// Check if we need to find the matching closing bracket
		opening := bz[startIndex]
		var closing byte
		if opening == '{' {
			closing = '}'
		} else {
			closing = ']'
		}

		// Look for the last matching closing bracket
		endIndex := bytes.LastIndexByte(bz, closing)
		if endIndex > startIndex {
			// Extract what appears to be valid JSON
			bz = bz[startIndex : endIndex+1]

			// Verify it's actually valid JSON
			var jsonTest interface{}
			if err := json.Unmarshal(bz, &jsonTest); err == nil {
				return bz, nil
			}
		}
	}

	// If we couldn't extract valid JSON, try parsing as YAML
	var out interface{}
	if err := yaml.Unmarshal(bz, &out); err == nil {
		return yaml.YAMLToJSON(bz)
	}

	// If neither JSON nor YAML parsing succeeded, return the original bytes
	// starting from the first opening brace if found, or the entire buffer
	if startIndex >= 0 {
		return bz[startIndex:], nil
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
