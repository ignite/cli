package cosmosbuf

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gobwas/glob"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v29/ignite/pkg/dircache"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/goenv"
	"github.com/ignite/cli/v29/ignite/pkg/xexec"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

const (
	binaryName                = "buf"
	flagTemplate              = "template"
	flagOutput                = "output"
	flagErrorFormat           = "error-format"
	flagLogFormat             = "log-format"
	flagIncludeImports        = "include-imports"
	flagIncludeWellKnownTypes = "include-wkt"
	flagPath                  = "path"
	fmtJSON                   = "json"

	// CMDGenerate generate command.
	CMDGenerate Command = "generate"
	CMDExport   Command = "export"
	CMDDep      Command = "dep"

	specCacheNamespace = "generate.buf"
)

var (
	commands = map[Command]struct{}{
		CMDGenerate: {},
		CMDExport:   {},
		CMDDep:      {},
	}

	// ErrInvalidCommand indicates an invalid command name.
	ErrInvalidCommand = errors.New("invalid command name")

	// ErrProtoFilesNotFound indicates that no ".proto" files were found.
	ErrProtoFilesNotFound = errors.New("no proto files found")
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure.
	Buf struct {
		path  string
		cache dircache.Cache
	}

	// genOptions used to configure code generation.
	genOptions struct {
		excluded       []glob.Glob
		flags          map[string]string
		fileByFile     bool
		includeImports bool
		includeWKT     bool
	}

	// GenOption configures code generation.
	GenOption func(*genOptions)
)

func newGenOptions() genOptions {
	return genOptions{
		flags:          make(map[string]string),
		excluded:       make([]glob.Glob, 0),
		fileByFile:     false,
		includeWKT:     false,
		includeImports: false,
	}
}

// WithFlag provides flag options for the buf generate command.
func WithFlag(flag, value string) GenOption {
	return func(o *genOptions) {
		o.flags[flag] = value
	}
}

// ExcludeFiles exclude file names from the generate command using glob.
func ExcludeFiles(patterns ...string) GenOption {
	return func(o *genOptions) {
		for _, pattern := range patterns {
			o.excluded = append(o.excluded, glob.MustCompile(pattern))
		}
	}
}

// IncludeImports also generate all imports except for Well-Known Types.
func IncludeImports() GenOption {
	return func(o *genOptions) {
		o.includeImports = true
	}
}

// IncludeWKT also generate Well-Known Types.
// Cannot be set without IncludeImports.
func IncludeWKT() GenOption {
	return func(o *genOptions) {
		o.includeImports = true
		o.includeWKT = true
	}
}

// FileByFile runs the generate command for each proto file.
func FileByFile() GenOption {
	return func(o *genOptions) {
		o.fileByFile = true
	}
}

// New creates a new Buf based on the installed binary.
func New(cacheStorage cache.Storage, goModPath string) (Buf, error) {
	path, err := xexec.ResolveAbsPath(filepath.Join(goenv.Bin(), binaryName))
	if err != nil {
		return Buf{}, err
	}

	bufCacheDir := filepath.Join("buf", goModPath)
	c, err := dircache.New(cacheStorage, bufCacheDir, specCacheNamespace)
	if err != nil {
		return Buf{}, err
	}

	return Buf{
		path:  path,
		cache: c,
	}, nil
}

// String returns the command name.
func (c Command) String() string {
	return string(c)
}

// Update updates module dependencies.
// By default updates all dependencies unless one or more dependencies are specified.
func (b Buf) Update(ctx context.Context, modDir string) error {
	files, err := xos.FindFilesExtension(modDir, xos.ProtoFile)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.Errorf("%w: %s", ErrProtoFilesNotFound, modDir)
	}

	cmd, err := b.command(CMDDep, nil, "update", modDir)
	if err != nil {
		return err
	}
	return b.runCommand(ctx, cmd...)
}

// Export runs the buf Export command for the files in the proto directory.
func (b Buf) Export(ctx context.Context, protoDir, output string) error {
	files, err := xos.FindFilesExtension(protoDir, xos.ProtoFile)
	if err != nil {
		return err
	}
	if len(files) == 0 {
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
	opts := newGenOptions()
	for _, apply := range options {
		apply(&opts)
	}

	// find all proto files into the path.
	foundFiles, err := xos.FindFilesExtension(protoPath, xos.ProtoFile)
	if err != nil || len(foundFiles) == 0 {
		return err
	}

	// check if already exist a cache for the template.
	key, err := b.cache.CopyTo(protoPath, output, template)
	if err != nil && !errors.Is(err, dircache.ErrCacheNotFound) {
		return err
	} else if err == nil {
		return nil
	}

	// remove excluded and cached files.
	protoFiles := make([]string, 0)
	for _, file := range foundFiles {
		okExclude := false
		for _, g := range opts.excluded {
			if g.Match(file) {
				okExclude = true
				break
			}
		}
		if !okExclude {
			protoFiles = append(protoFiles, file)
		}
	}
	if len(protoFiles) == 0 {
		return nil
	}

	flags := map[string]string{
		flagTemplate:    template,
		flagOutput:      output,
		flagErrorFormat: fmtJSON,
		flagLogFormat:   fmtJSON,
	}
	for k, v := range opts.flags {
		flags[k] = v
	}
	if opts.includeImports {
		flags[flagIncludeImports] = "true"
	}
	if opts.includeWKT {
		flags[flagIncludeWellKnownTypes] = "true"
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

	return b.cache.Save(output, key)
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
