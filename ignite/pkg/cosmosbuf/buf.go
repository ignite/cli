package cosmosbuf

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xexec"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure.
	Buf struct {
		path         string
		sdkProtoDir  string
		storageCache cache.Cache[[]byte]
	}

	// genOptions used to configure code generation.
	genOptions struct {
		excluded   map[string]struct{}
		fileByFile bool
	}

	// GenOption configures code generation.
	GenOption func(*genOptions)
)

const (
	binaryName      = "buf"
	flagTemplate    = "template"
	flagOutput      = "output"
	flagErrorFormat = "error-format"
	flagLogFormat   = "log-format"
	flagPath        = "path"
	flagOnly        = "only"
	fmtJSON         = "json"

	// CMDGenerate generate command.
	CMDGenerate Command = "generate"
	CMDExport   Command = "export"
	CMDMod      Command = "mod"

	specCacheNamespace = "generate.buf"
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

func ExcludeFiles(files ...string) GenOption {
	return func(o *genOptions) {
		excluded := make(map[string]struct{})
		for _, file := range files {
			excluded[file] = struct{}{}
		}
		o.excluded = excluded
	}
}

func FileByFile() GenOption {
	return func(o *genOptions) {
		o.fileByFile = true
	}
}

// New creates a new Buf based on the installed binary.
func New(cacheStorage cache.Storage) (Buf, error) {
	path, err := xexec.ResolveAbsPath(binaryName)
	if err != nil {
		return Buf{}, err
	}

	return Buf{
		path:         path,
		storageCache: cache.New[[]byte](cacheStorage, specCacheNamespace),
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

	cmd, err := b.command(CMDMod, flags, "update", modDir)
	if err != nil {
		return err
	}

	return b.runCommand(ctx, cmd...)
}

// Export runs the buf Export command for the files in the proto directory.
func (b Buf) Export(ctx context.Context, protoDir, output string) error {
	specs, err := xos.FindFilesExtension(protoDir, xos.ProtoFile)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return errors.Errorf("%w: %s", ErrProtoFilesNotFound, protoDir)
	}
	flags := map[string]string{
		flagOutput: output,
	}

	cmd, err := b.command(CMDExport, flags, protoDir)
	if err != nil {
		return err
	}

	return b.runCommand(ctx, cmd...)
}

// Generate runs the buf Generate command for each file into the proto directory.
func (b Buf) Generate(
	ctx context.Context,
	protoPath,
	output,
	template string,
	options ...GenOption,
) (err error) {
	opts := genOptions{}
	for _, apply := range options {
		apply(&opts)
	}

	// find all proto files into the path
	protoFiles, err := xos.FindFilesExtension(protoPath, xos.ProtoFile)
	if err != nil || len(protoFiles) == 0 {
		return err
	}

	cached, err := b.getDirCache(protoPath, output, template)
	if err != nil {
		return err
	}

	// remove excluded and cached files
	for i, file := range protoFiles {
		if _, ok := opts.excluded[filepath.Base(file)]; ok {
			protoFiles = append(protoFiles[:i], protoFiles[i+1:]...)
		} else if _, ok := cached[file]; ok {
			protoFiles = append(protoFiles[:i], protoFiles[i+1:]...)
		}
	}

	flags := map[string]string{
		flagTemplate:    template,
		flagOutput:      output,
		flagErrorFormat: fmtJSON,
		flagLogFormat:   fmtJSON,
	}

	if !opts.fileByFile {
		cmd, err := b.command(CMDGenerate, flags, protoPath)
		if err != nil {
			return err
		}
		for _, file := range protoFiles {
			cmd = append(cmd, fmt.Sprintf("--%s=%s", flagPath, file))
		}
		if err := b.runCommand(ctx, cmd...); err != nil {
			return err
		}
	} else {
		g, ctx := errgroup.WithContext(ctx)
		for _, file := range protoFiles {
			cmd, err := b.command(CMDGenerate, flags, file)
			if err != nil {
				return err
			}

			g.Go(func() error {
				cmd := cmd
				return b.runCommand(ctx, cmd...)
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}

	return b.saveDirCache(output, template)
}

// runCommand run the buf CLI command.
func (b Buf) runCommand(ctx context.Context, cmd ...string) error {
	execOpts := []exec.Option{
		exec.IncludeStdLogsToError(),
	}
	return exec.Exec(ctx, cmd, execOpts...)
}

// command generate the buf CLI command.
func (b Buf) command(
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
