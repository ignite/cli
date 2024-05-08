package cosmosbuf

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/exec"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/xexec"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

type (
	// Command represents a high level command under buf.
	Command string

	// Buf represents the buf application structure.
	Buf struct {
		path         string
		protoCache   *protoanalysis.Cache
		storageCache cache.Cache[[]byte]
	}
)

const (
	binaryName      = "buf"
	flagTemplate    = "template"
	flagOutput      = "output"
	flagErrorFormat = "error-format"
	flagLogFormat   = "log-format"
	flagExcludePath = "exclude-path"
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

// New creates a new Buf based on the installed binary.
func New(cacheStorage cache.Storage) (Buf, error) {
	path, err := xexec.ResolveAbsPath(binaryName)
	if err != nil {
		return Buf{}, err
	}

	return Buf{
		path:         path,
		storageCache: cache.New[[]byte](cacheStorage, specCacheNamespace),
		protoCache:   protoanalysis.NewCache(),
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
	excluded ...string,
) (err error) {
	protoFiles, err := xos.FindFiles(protoDir, xos.ProtoFile)
	if err != nil || len(protoFiles) == 0 {
		return err
	}

	flags := map[string]string{
		flagTemplate:    template,
		flagOutput:      output,
		flagErrorFormat: fmtJSON,
		flagLogFormat:   fmtJSON,
		flagExcludePath: join(excluded...),
	}

	cmd, err := b.generateCommand(CMDGenerate, flags, protoDir)
	if err != nil {
		return err
	}

	return b.runCommand(ctx, cmd...)
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

func join(elems ...string) (value string) {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.Quote(elems[0])
	}

	var b strings.Builder
	b.WriteString(strconv.Quote(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(",")
		b.WriteString(strconv.Quote(s))
	}
	return b.String()
}
