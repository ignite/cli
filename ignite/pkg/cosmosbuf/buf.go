package cosmosbuf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

<<<<<<< HEAD
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/v28/ignite/pkg/xexec"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure.
	Buf struct {
		path        string
		sdkProtoDir string
		cache       *protoanalysis.Cache
	}
=======
	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v29/ignite/pkg/dircache"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goenv"
	"github.com/ignite/cli/v29/ignite/pkg/xexec"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
>>>>>>> cfea8dd5 (use buf build from the gobin path (#4242))
)

const (
	binaryName      = "buf"
	flagTemplate    = "template"
	flagOutput      = "output"
	flagErrorFormat = "error-format"
	flagLogFormat   = "log-format"
	flagOnly        = "only"
	fmtJSON         = "json"

	// CMDGenerate generate command.
	CMDGenerate Command = "generate"
	CMDExport   Command = "export"
	CMDMod      Command = "mod"
)

var (
	commands = map[Command]struct{}{
		CMDGenerate: {},
		CMDExport:   {},
		CMDMod:      {},
	}

	// ErrInvalidCommand indicates an invalid command name.
	ErrInvalidCommand = errors.New("invalid command name")

	// ErrProtoFilesNotFound indicates that no ".proto" files were found.
	ErrProtoFilesNotFound = errors.New("no proto files found")
)

// New creates a new Buf based on the installed binary.
<<<<<<< HEAD
func New() (Buf, error) {
	path, err := xexec.ResolveAbsPath(binaryName)
=======
func New(cacheStorage cache.Storage, goModPath string) (Buf, error) {
	path, err := xexec.ResolveAbsPath(filepath.Join(goenv.Bin(), binaryName))
>>>>>>> cfea8dd5 (use buf build from the gobin path (#4242))
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

// Update updates module dependencies.
// By default updates all dependencies unless one or more dependencies are specified.
func (b Buf) Update(ctx context.Context, modDir string, dependencies ...string) error {
	var flags map[string]string
	if dependencies != nil {
		flags = map[string]string{
			flagOnly: strings.Join(dependencies, ","),
		}
	}

	cmd, err := b.generateCommand(CMDMod, flags, "update", modDir)
	if err != nil {
		return err
	}

	return b.runCommand(ctx, cmd...)
}

// Export runs the buf Export command for the files in the proto directory.
func (b Buf) Export(ctx context.Context, protoDir, output string) error {
	// Check if the proto directory is the Cosmos SDK one
	// TODO(@julienrbrt): this whole custom handling can be deleted
	// after https://github.com/cosmos/cosmos-sdk/pull/18993 in v29.
	if strings.Contains(protoDir, cosmosver.CosmosSDKRepoName) {
		if b.sdkProtoDir == "" {
			// Copy Cosmos SDK proto path without the Buf workspace.
			// This is done because the workspace contains a reference to
			// a "orm/internal" proto folder that is not present by default
			// in the SDK repository.
			d, err := copySDKProtoDir(protoDir)
			if err != nil {
				return err
			}

			b.sdkProtoDir = d
		}

		// Split absolute path into an absolute prefix and a relative suffix
		paths := strings.Split(protoDir, "/proto")
		if len(paths) < 2 {
			return errors.Errorf("invalid Cosmos SDK mod path: %s", protoDir)
		}

		// Use the SDK copy to resolve SDK proto files
		protoDir = filepath.Join(b.sdkProtoDir, paths[1])
	}
	specs, err := xos.FindFiles(protoDir, xos.ProtoFile)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return errors.Errorf("%w: %s", ErrProtoFilesNotFound, protoDir)
	}
	flags := map[string]string{
		flagOutput: output,
	}

	cmd, err := b.generateCommand(CMDExport, flags, protoDir)
	if err != nil {
		return err
	}

	return b.runCommand(ctx, cmd...)
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

	// TODO(@julienrbrt): this whole custom handling can be deleted
	// after https://github.com/cosmos/cosmos-sdk/pull/18993 in v29.
	if strings.Contains(protoDir, cosmosver.CosmosSDKRepoName) {
		if b.sdkProtoDir == "" {
			b.sdkProtoDir, err = copySDKProtoDir(protoDir)
			if err != nil {
				return err
			}
		}
		dirs := strings.Split(protoDir, "/proto/")
		if len(dirs) < 2 {
			return errors.Errorf("invalid Cosmos SDK mod path: %s", dirs)
		}
		protoDir = filepath.Join(b.sdkProtoDir, dirs[1])
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

			specs, err := xos.FindFiles(protoDir, "proto")
			if err != nil {
				return err
			}
			if len(specs) == 0 {
				continue
			}

			cmd, err := b.generateCommand(CMDGenerate, flags, file.Path)
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

// Cleanup deletes temporary files and directories.
func (b Buf) Cleanup() error {
	if b.sdkProtoDir != "" {
		return os.RemoveAll(b.sdkProtoDir)
	}
	return nil
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
		return nil, errors.Errorf("%w: %s", ErrInvalidCommand, c)
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

// findSDKProtoPath finds the Cosmos SDK proto folder path.
func findSDKProtoPath(protoDir string) string {
	paths := strings.Split(protoDir, "@")
	if len(paths) < 2 {
		return protoDir
	}
	version := strings.Split(paths[1], "/")[0]
	return fmt.Sprintf("%s@%s/proto", paths[0], version)
}

// copySDKProtoDir copies the Cosmos SDK proto folder to a temporary directory.
// The temporary directory must be removed by the caller.
func copySDKProtoDir(protoDir string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "proto-sdk")
	if err != nil {
		return "", err
	}

	srcPath := findSDKProtoPath(protoDir)
	return tmpDir, xos.CopyFolder(srcPath, tmpDir)
}
