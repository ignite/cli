package v1_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	v1 "github.com/ignite/cli/v29/ignite/services/plugin/grpc/v1"
)

func TestManifestImportCobraCommand(t *testing.T) {
	manifest := &v1.Manifest{
		Name: "hey",
		Commands: []*v1.Command{
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

	expectedManifest := &v1.Manifest{
		Name: "hey",
		Commands: []*v1.Command{
			{Use: "existing"},
			{
				Use:               "new",
				Aliases:           []string{"n"},
				Short:             "short",
				Long:              "long",
				Hidden:            true,
				PlaceCommandUnder: "under",
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
					},
					{
						Name:         "persistent",
						DefaultValue: "hello",
						Value:        "hello",
						Usage:        "a persistent string",
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
						Flags: []*v1.Flag{
							{
								Name:         "subbool",
								Shorthand:    "b",
								DefaultValue: "true",
								Value:        "true",
								Usage:        "a bool",
								Type:         v1.Flag_TYPE_FLAG_BOOL,
							},
							{
								Name:         "substring",
								DefaultValue: "hello",
								Value:        "hello",
								Usage:        "a string",
								Type:         v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED,
							},
						},
						Commands: []*v1.Command{{Use: "subsub"}},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedManifest, manifest)
}
