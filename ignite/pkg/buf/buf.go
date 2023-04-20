package buf

import (
	"context"
	"errors"
	"fmt"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/xexec"
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure
	Buf struct {
		path string
	}
)

const (
	binaryName         = "buf"
	flagTemplate       = "template"
	flagOutput         = "output"
	flagIncludeImports = "include-imports"
	flagErrorFormat    = "error-format"
	flagLogFormat      = "log-format"
	fmtJSON            = "json"

	// CMDGenerate generate command
	CMDGenerate Command = "generate"
)

var (
	commands = map[Command]struct{}{
		CMDGenerate: {},
	}

	// ErrInvalidCommand error invalid command name
	ErrInvalidCommand = errors.New("invalid command name")
)

// New creates a new Buf based on the installed binary
func New() (Buf, error) {
	path, err := xexec.ResolveAbsPath(binaryName)
	if err != nil {
		return Buf{}, err
	}
	return Buf{
		path: path,
	}, nil
}

// String returns the command name
func (c Command) String() string {
	return string(c)
}

// Generate runs the buf Generate command
func (b Buf) Generate(ctx context.Context, protoDir, output, template string) error {
	cmd, err := b.generateCommand(
		CMDGenerate,
		map[string]string{
			flagTemplate:       template,
			flagOutput:         output,
			flagIncludeImports: "true",
			flagErrorFormat:    fmtJSON,
			flagLogFormat:      fmtJSON,
		},
	)
	if err != nil {
		return err
	}
	return b.runCommand(ctx, protoDir, cmd...)
}

// runCommand run the buf CLI command
func (b Buf) runCommand(ctx context.Context, workDir string, cmd ...string) error {
	execOpts := []exec.Option{
		exec.StepOption(step.Workdir(workDir)),
		exec.IncludeStdLogsToError(),
	}

	return exec.Exec(ctx, cmd, execOpts...)
}

// generateCommand generate the buf CLI command
func (b Buf) generateCommand(
	c Command,
	flags map[string]string,
	args ...string,
) ([]string, error) {
	if _, ok := commands[c]; !ok {
		return nil, fmt.Errorf("%v: %s", ErrInvalidCommand, c)
	}

	command := []string{
		b.path,
		c.String(),
	}
	command = append(command, args...)

	for flag, value := range flags {
		command = append(command,
			fmt.Sprintf("--%s=%s", flag, value),
		)
	}
	return command, nil
}
