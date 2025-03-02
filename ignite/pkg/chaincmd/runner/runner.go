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
// This method is purposely verbose to trim gibberish output.
func (b *buffer) JSONEnsuredBytes() ([]byte, error) {
	bz := b.Bytes()
	content := strings.TrimSpace(string(bz))

	// Early detection - check first non-whitespace character
	if len(content) > 0 {
		firstChar := content[0]

		// Quick check for JSON format (starts with { or [)
		if firstChar == '{' || firstChar == '[' {
			// Attempt to validate and extract clean JSON
			return cleanAndValidateJSON(bz)
		}

		// Quick check for YAML format (common indicators)
		if firstChar == '-' || strings.HasPrefix(content, "---") ||
			strings.Contains(content, ":\n") || strings.Contains(content, ": ") {
			// Likely YAML, convert to JSON directly
			var out any
			if err := yaml.Unmarshal(bz, &out); err == nil {
				return yaml.YAMLToJSON(bz)
			}
		}
	}

	// If format wasn't immediately obvious, try the more thorough approach
	return fallbackFormatDetection(bz)
}

// cleanAndValidateJSON attempts to extract valid JSON from potentially messy output.
func cleanAndValidateJSON(bz []byte) ([]byte, error) {
	// Find the first JSON opening character
	startIndex := strings.IndexAny(string(bz), "{[")
	if startIndex < 0 {
		return bz, nil // No JSON structure found
	}

	// Determine matching closing character
	opening := bz[startIndex]
	var closing byte
	if opening == '{' {
		closing = '}'
	} else {
		closing = ']'
	}

	endIndex := findMatchingCloseBracket(bz[startIndex:], opening, closing)
	if endIndex < 0 {
		// no proper closing found, try last instance
		endIndex = bytes.LastIndexByte(bz, closing)
		if endIndex <= startIndex {
			return bz[startIndex:], nil // Return from start to end if no closing found
		}
	} else {
		endIndex += startIndex
	}

	// validate JSON
	jsonData := bz[startIndex : endIndex+1]
	var jsonTest any
	if err := json.Unmarshal(jsonData, &jsonTest); err == nil {
		return jsonData, nil
	}

	// if validation failed, return from start to end
	return bz[startIndex:], nil
}

// findMatchingCloseBracket returns the accounting for nested structures.
func findMatchingCloseBracket(data []byte, openChar, closeChar byte) int {
	depth := 0
	for i, b := range data {
		if b == openChar {
			depth++
		} else if b == closeChar {
			depth--
			if depth == 0 {
				return i // Found matching closing bracket
			}
		}
	}
	return -1 // No matching bracket found
}

// fallbackFormatDetection tries different approaches to detect and convert format.
func fallbackFormatDetection(bz []byte) ([]byte, error) {
	// first try to find and extract JSON
	startIndex := strings.IndexAny(string(bz), "{[")
	if startIndex >= 0 {
		result, err := cleanAndValidateJSON(bz)
		if err == nil {
			return result, nil
		}

		// if extraction failed but we found a start, return from there
		return bz[startIndex:], nil
	}

	// fallback to yaml parsing
	var out any
	if err := yaml.Unmarshal(bz, &out); err == nil {
		return yaml.YAMLToJSON(bz)
	}

	// nothing worked, return original
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
