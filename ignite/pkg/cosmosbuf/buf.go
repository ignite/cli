package cosmosbuf

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/ignite/pkg/xexec"
	"github.com/ignite/cli/ignite/pkg/xos"
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure.
	Buf struct {
		path     string
		sdkCache string
		cache    *protoanalysis.Cache
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
func (b Buf) Generate(
	ctx context.Context,
	protoDir,
	output,
	template string,
	excludeFilename ...string,
) (err error) {
	var (
		excluded = make(map[string]struct{})
		flags    = map[string]string{
			flagTemplate:    template,
			flagOutput:      output,
			flagErrorFormat: fmtJSON,
			flagLogFormat:   fmtJSON,
		}
	)
	for _, file := range excludeFilename {
		excluded[file] = struct{}{}
	}

	// TODO find a better way to generate the cosmos-sdk files
	// the buf.work.yaml contains the `orm/internal` folder, but the `go mod`
	// can't download this folder because is unused as a dependency. We need to
	// change the workspace copying the files to another folder and generate the
	// files.
	if strings.Contains(protoDir, cosmosver.CosmosModulePath) {
		if b.sdkCache == "" {
			b.sdkCache, err = prepareSDK(protoDir)
			if err != nil {
				return err
			}
		}
		dirs := strings.Split(protoDir, "/proto/")
		if len(dirs) < 2 {
			return fmt.Errorf("invalid cosmos sdk mod path: %s", dirs)
		}
		protoDir = filepath.Join(b.sdkCache, dirs[1])
	}

	pkgs, err := protoanalysis.Parse(ctx, b.cache, protoDir)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			if _, ok := excluded[filepath.Base(file.Path)]; ok {
				continue
			}
			cmd, err := b.generateCommand(
				CMDGenerate,
				flags,
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

// findSDKProtoPath find the cosmos-sdk proto folder path.
func findSDKProtoPath(protoDir string) (string, error) {
	paths := strings.Split(protoDir, "@")
	if len(paths) < 2 {
		return "", fmt.Errorf("invalid sdk mod dir: %s", protoDir)
	}
	version := strings.Split(paths[1], "/")[0]
	return fmt.Sprintf("%s@%s/proto", paths[0], version), nil
}

// prepareSDK copy the cosmos sdk proto folder to a temporary directory
// so we can skip the buf workspace.
func prepareSDK(protoDir string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "proto-sdk")
	if err != nil {
		return "", err
	}
	srcPath, err := findSDKProtoPath(protoDir)
	if err != nil {
		return "", err
	}
	return tmpDir, xos.CopyFolder(srcPath, tmpDir)
}
