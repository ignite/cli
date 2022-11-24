package plugin_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/ignite/cli/ignite/services/plugin"
)

func TestManifestImportCobraCommand(t *testing.T) {
	manifest := plugin.Manifest{
		Name: "hey",
		Commands: []plugin.Command{
			{Use: "existing"},
		},
	}
	cmd := &cobra.Command{
		Use:     "new",
		Aliases: []string{"n"},
		Short:   "short",
		Long:    "long",
	}
	cmd.Flags().BoolP("bool", "b", true, "a bool")
	cmd.Flags().String("string", "hello", "a string")
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

	expectedManifest := plugin.Manifest{
		Name: "hey",
		Commands: []plugin.Command{
			{Use: "existing"},
			{
				Use:               "new",
				Aliases:           []string{"n"},
				Short:             "short",
				Long:              "long",
				PlaceCommandUnder: "under",
				Flags: []plugin.Flag{
					{
						Name:      "bool",
						Shorthand: "b",
						DefValue:  "true",
						Value:     "true",
						Usage:     "a bool",
						Type:      plugin.FlagTypeBool,
					},
					{
						Name:     "string",
						DefValue: "hello",
						Value:    "hello",
						Usage:    "a string",
						Type:     plugin.FlagTypeString,
					},
				},
				Commands: []plugin.Command{
					{
						Use:     "sub",
						Aliases: []string{"s"},
						Short:   "sub short",
						Long:    "sub long",
						Flags: []plugin.Flag{
							{
								Name:      "subbool",
								Shorthand: "b",
								DefValue:  "true",
								Value:     "true",
								Usage:     "a bool",
								Type:      plugin.FlagTypeBool,
							},
							{
								Name:     "substring",
								DefValue: "hello",
								Value:    "hello",
								Usage:    "a string",
								Type:     plugin.FlagTypeString,
							},
						},
						Commands: []plugin.Command{{Use: "subsub"}},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedManifest, manifest)
}
