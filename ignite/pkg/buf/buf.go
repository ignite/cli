package buf

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/ignite/pkg/xexec"
	"github.com/ignite/cli/ignite/pkg/xos"
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
	cosmosSDKModulePath = "github.com/cosmos/cosmos-sdk"
	binaryName          = "buf"
	flagTemplate        = "template"
	flagOutput          = "output"
	flagErrorFormat     = "error-format"
	flagLogFormat       = "log-format"
	fmtJSON             = "json"

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
func (b Buf) Generate(ctx context.Context, protoDir, output, template string) (err error) {
	var (
		cmds  = make([][]string, 0)
		flags = map[string]string{
			flagTemplate:    template,
			flagOutput:      output,
			flagErrorFormat: fmtJSON,
			flagLogFormat:   fmtJSON,
		}
	)

	if strings.Contains(protoDir, cosmosSDKModulePath) {
		protoDir, err = prepareSDK(protoDir)
		if err != nil {
			return err
		}
		// defer os.RemoveAll(protoDir)

		cmd, err := b.generateCommand(
			CMDGenerate,
			flags,
			protoDir,
		)
		if err != nil {
			return err
		}
		cmds = append(cmds, cmd)

	} else {
		pkgs, err := protoanalysis.Parse(ctx, b.cache, protoDir)
		if err != nil {
			return err
		}

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				cmd, err := b.generateCommand(
					CMDGenerate,
					flags,
					file.Path,
				)
				if err != nil {
					return err
				}
				cmds = append(cmds, cmd)
			}
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, cmd := range cmds {
		g.Go(func() error {
			cmd := cmd
			return b.runCommand(ctx, cmd...)
		})
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

func findSDKPath(protoDir string) (string, error) {
	paths := strings.Split(protoDir, "@")
	if len(paths) < 2 {
		return "", fmt.Errorf("invalid sdk mod dir: %s", protoDir)
	}
	version := strings.Split(paths[1], "/")[0]
	return fmt.Sprintf("%s@%s/proto", paths[0], version), nil
}

func prepareSDK(protoDir string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "proto-sdk")
	srcPath, err := findSDKPath(protoDir)
	if err != nil {
		return "", err
	}
	return tmpDir, xos.CopyFolder(srcPath, tmpDir)
}
