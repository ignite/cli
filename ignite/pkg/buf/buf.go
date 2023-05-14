package buf

import (
	"context"
	"errors"
	"fmt"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/ignite/pkg/xexec"
	"golang.org/x/sync/errgroup"
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure.
	Buf struct {
		path  string
		cache protoanalysis.Cache
	}
)

const (
	binaryName      = "buf"
	flagTemplate    = "template"
	flagOutput      = "output"
	flagErrorFormat = "error-format"
	flagLogFormat   = "log-format"
	fmtJSON         = "json"

	// CMDGenerate generate command.
	CMDGenerate Command = "generate"
)

var (
	commands = map[Command]struct{}{
		CMDGenerate: {},
	}

	// ErrInvalidCommand error invalid command name.
	ErrInvalidCommand = errors.New("invalid command name")
)

// New creates a new Buf based on the installed binary.
func New() (Buf, error) {
	path, err := xexec.ResolveAbsPath(binaryName)
	if err != nil {
		return Buf{}, err
	}
	return Buf{
		path:  path,
		cache: protoanalysis.NewCache(),
	}, nil
}

// String returns the command name.
func (c Command) String() string {
	return string(c)
}

// Generate runs the buf Generate command for each file into the proto directory.
func (b Buf) Generate(ctx context.Context, protoDir, output, template string) error {
	pkgs, err := protoanalysis.Parse(ctx, b.cache, protoDir)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			cmd, err := b.generateCommand(
				CMDGenerate,
				map[string]string{
					flagTemplate:    template,
					flagOutput:      output,
					flagErrorFormat: fmtJSON,
					flagLogFormat:   fmtJSON,
				},
				file.Path,
			)
			if err != nil {
				return err
			}
			g.Go(func() error {
				cmd := cmd
				return b.runCommand(ctx, cmd...)
			})
		}
	}
	return g.Wait()
}

// runCommand run the buf CLI command.
func (b Buf) runCommand(ctx context.Context, cmd ...string) error {
	execOpts := []exec.Option{
		exec.IncludeStdLogsToError(),
	}
	return exec.Exec(ctx, cmd, execOpts...)
}

// generateCommand generate the buf CLI command.
func (b Buf) generateCommand(
	c Command,
	flags map[string]string,
	args ...string,
) ([]string, error) {
	if _, ok := commands[c]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, c)
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
