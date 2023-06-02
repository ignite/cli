package grpc_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/services/plugin/grpc"
	v1 "github.com/ignite/cli/ignite/services/plugin/grpc/v1"
)

func TestCommandToCobraCommand(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		pcmd    = grpc.Command{
			Use:     "new",
			Aliases: []string{"n"},
			Short:   "short",
			Long:    "long",
			Hidden:  true,
			Flags: []*grpc.Flag{
				{
					Name:         "bool",
					Shorthand:    "b",
					DefaultValue: "true",
					Value:        "true",
					Usage:        "a bool",
					Type:         v1.FlagType_FLAG_TYPE_BOOL,
				},
				{
					Name:         "string",
					DefaultValue: "hello",
					Value:        "hello",
					Usage:        "a string",
					Type:         v1.FlagType_FLAG_TYPE_STRING,
					Persistent:   true,
				},
			},
			Commands: []*grpc.Command{
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

func TestManifestImportCobraCommand(t *testing.T) {
	manifest := &grpc.Manifest{
		Name: "hey",
		Commands: []*grpc.Command{
			{Use: "existing"},
		},
	}
	cmd := &cobra.Command{
		Use:     "new",
		Aliases: []string{"n"},
		Short:   "short",
		Long:    "long",
		Hidden:  true,
	}
	cmd.Flags().BoolP("bool", "b", true, "a bool")
	cmd.Flags().String("string", "hello", "a string")
	cmd.PersistentFlags().String("persistent", "hello", "a persistent string")
	subcmd := &cobra.Command{
		Use:     "sub",
		Aliases: []string{"s"},
		Short:   "sub short",
		Long:    "sub long",
	}
	subcmd.Flags().BoolP("subbool", "b", true, "a bool")
	subcmd.Flags().String("substring", "hello", "a string")
	subcmd.AddCommand(&cobra.Command{
		Use: "subsub",
	})
	cmd.AddCommand(subcmd)

	manifest.ImportCobraCommand(cmd, "under")

	expectedManifest := &grpc.Manifest{
		Name: "hey",
		Commands: []*grpc.Command{
			{Use: "existing"},
			{
				Use:               "new",
				Aliases:           []string{"n"},
				Short:             "short",
				Long:              "long",
				Hidden:            true,
				PlaceCommandUnder: "under",
				Flags: []*grpc.Flag{
					{
						Name:         "bool",
						Shorthand:    "b",
						DefaultValue: "true",
						Value:        "true",
						Usage:        "a bool",
						Type:         v1.FlagType_FLAG_TYPE_BOOL,
					},
					{
						Name:         "string",
						DefaultValue: "hello",
						Value:        "hello",
						Usage:        "a string",
						Type:         v1.FlagType_FLAG_TYPE_STRING,
					},
					{
						Name:         "persistent",
						DefaultValue: "hello",
						Value:        "hello",
						Usage:        "a persistent string",
						Type:         v1.FlagType_FLAG_TYPE_STRING,
						Persistent:   true,
					},
				},
				Commands: []*grpc.Command{
					{
						Use:     "sub",
						Aliases: []string{"s"},
						Short:   "sub short",
						Long:    "sub long",
						Flags: []*grpc.Flag{
							{
								Name:         "subbool",
								Shorthand:    "b",
								DefaultValue: "true",
								Value:        "true",
								Usage:        "a bool",
								Type:         v1.FlagType_FLAG_TYPE_BOOL,
							},
							{
								Name:         "substring",
								DefaultValue: "hello",
								Value:        "hello",
								Usage:        "a string",
								Type:         v1.FlagType_FLAG_TYPE_STRING,
							},
						},
						Commands: []*grpc.Command{{Use: "subsub"}},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedManifest, manifest)
}
