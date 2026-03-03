package cosmosbuf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGenOptionsDefaults(t *testing.T) {
	opts := newGenOptions()
	require.Empty(t, opts.flags)
	require.Empty(t, opts.excluded)
	require.False(t, opts.fileByFile)
	require.False(t, opts.includeImports)
	require.False(t, opts.includeWKT)
	require.Empty(t, opts.moduleName)
}

func TestGenOptions(t *testing.T) {
	opts := newGenOptions()

	WithFlag("foo", "bar")(&opts)
	ExcludeFiles("*.proto")(&opts)
	IncludeImports()(&opts)
	FileByFile()(&opts)
	WithModuleName("ignite.chain")(&opts)

	require.Equal(t, "bar", opts.flags["foo"])
	require.Len(t, opts.excluded, 1)
	require.True(t, opts.includeImports)
	require.True(t, opts.fileByFile)
	require.Equal(t, "ignite.chain", opts.moduleName)

	IncludeWKT()(&opts)
	require.True(t, opts.includeImports)
	require.True(t, opts.includeWKT)
}

func TestCommandString(t *testing.T) {
	require.Equal(t, "generate", CMDGenerate.String())
}

func TestCommandReturnsErrorForInvalidCommand(t *testing.T) {
	_, err := Buf{}.command(Command("invalid"), nil)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidCommand)
}

func TestCommandBuildsExpectedArguments(t *testing.T) {
	flags := map[string]string{
		"template": "buf.gen.yaml",
		"output":   "out",
	}

	got, err := Buf{}.command(CMDGenerate, flags, "proto")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(got), 4)
	require.Equal(t, []string{"go", "tool", "github.com/bufbuild/buf/cmd/buf", "generate"}, got[:4])
	require.Contains(t, got, "proto")

	joined := strings.Join(got, " ")
	require.Contains(t, joined, "--template=buf.gen.yaml")
	require.Contains(t, joined, "--output=out")
}
