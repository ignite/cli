package cosmosbuf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/v28/ignite/pkg/xexec"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
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
		path        string
		sdkProtoDir string
		cache       *protoanalysis.Cache
	}
<<<<<<< HEAD
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
=======

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
>>>>>>> 8e0937d9 (feat: remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta` (#4090))
	}
}

// WithFlag provides flag options for the buf generate command.
func WithFlag(flag, value string) GenOption {
	return func(o *genOptions) {
		o.flags[flag] = value
	}
}

<<<<<<< HEAD
=======
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

>>>>>>> 8e0937d9 (feat: remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta` (#4090))
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

// Update updates module dependencies.
// By default updates all dependencies unless one or more dependencies are specified.
<<<<<<< HEAD
func (b Buf) Update(ctx context.Context, modDir string, dependencies ...string) error {
	var flags map[string]string
	if dependencies != nil {
		flags = map[string]string{
			flagOnly: strings.Join(dependencies, ","),
		}
	}

	cmd, err := b.generateCommand(CMDMod, flags, "update", modDir)
=======
func (b Buf) Update(ctx context.Context, modDir string) error {
	files, err := xos.FindFilesExtension(modDir, xos.ProtoFile)
>>>>>>> 8e0937d9 (feat: remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta` (#4090))
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
<<<<<<< HEAD
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
=======
	files, err := xos.FindFilesExtension(protoDir, xos.ProtoFile)
>>>>>>> 8e0937d9 (feat: remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta` (#4090))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.Errorf("%w: %s", ErrProtoFilesNotFound, protoDir)
	}

	flags := map[string]string{
		flagOutput: output,
	}
<<<<<<< HEAD

	cmd, err := b.generateCommand(CMDExport, flags, protoDir)
=======
	cmd, err := b.command(CMDExport, flags, protoDir)
>>>>>>> 8e0937d9 (feat: remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta` (#4090))
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
<<<<<<< HEAD
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
=======
	opts := newGenOptions()
	for _, apply := range options {
		apply(&opts)
>>>>>>> 8e0937d9 (feat: remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta` (#4090))
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
	for k, v := range opts.flags {
		flags[k] = v
	}
	if opts.includeImports {
		flags[flagIncludeImports] = "true"
	}
	if opts.includeWKT {
		flags[flagIncludeWellKnownTypes] = "true"
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
