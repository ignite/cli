package v1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/ignite/cli/ignite/services/plugin/grpc/v1"
)

func TestCommandToCobraCommand(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		pcmd    = v1.Command{
			Use:     "new",
			Aliases: []string{"n"},
			Short:   "short",
			Long:    "long",
			Hidden:  true,
			Flags: []*v1.Flag{
				{
					Name:         "bool",
					Shorthand:    "b",
					DefaultValue: "true",
					Value:        "true",
					Usage:        "a bool",
					Type:         v1.Flag_TYPE_FLAG_BOOL,
				},
				{
					Name:         "string",
					DefaultValue: "hello",
					Value:        "hello",
					Usage:        "a string",
					Type:         v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED,
					Persistent:   true,
				},
			},
			Commands: []*v1.Command{
				{
					Use:     "sub",
					Aliases: []string{"s"},
					Short:   "sub short",
					Long:    "sub long",
				},
			},
		}
	)

	cmd, err := pcmd.ToCobraCommand()

	require.NoError(err)
	require.NotNil(cmd)
	assert.Empty(cmd.Commands()) // subcommands aren't converted
	assert.Equal(pcmd.Use, cmd.Use)
	assert.Equal(pcmd.Short, cmd.Short)
	assert.Equal(pcmd.Long, cmd.Long)
	assert.Equal(pcmd.Aliases, cmd.Aliases)
	assert.Equal(pcmd.Hidden, cmd.Hidden)
	for _, f := range pcmd.Flags {
		if f.Persistent {
			assert.NotNil(cmd.PersistentFlags().Lookup(f.Name), "missing pflag %s", f.Name)
		} else {
			assert.NotNil(cmd.Flags().Lookup(f.Name), "missing flag %s", f.Name)
		}
	}
}
