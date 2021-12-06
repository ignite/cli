package starportcmd

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/services/plugin"
)

// pluginInterface implements plugin.Interface for testing purpose.
type pluginInterface struct {
	commands []plugin.Command
}

func (p pluginInterface) Commands() []plugin.Command {
	return p.commands
}

func (pluginInterface) Execute(plugin.Command, []string) error {
	return nil
}

func TestLinkPluginCmds(t *testing.T) {
	buildRootCmd := func() *cobra.Command {
		var (
			rootCmd = &cobra.Command{
				Use: "starport",
			}
			scaffoldCmd = &cobra.Command{
				Use: "scaffold",
			}
			scaffoldChainCmd = &cobra.Command{
				Use: "chain",
				Run: func(*cobra.Command, []string) {},
			}
		)
		scaffoldCmd.AddCommand(scaffoldChainCmd)
		rootCmd.AddCommand(scaffoldCmd)
		return rootCmd
	}
	tests := []struct {
		name            string
		pluginInterface pluginInterface
		expectedDumpCmd string
		expectedError   string
	}{
		{
			name: "ok: link foo at root",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "foo",
					},
				},
			},
			expectedDumpCmd: `
starport
  foo*
  scaffold
    chain*
`,
		},
		{
			name: "ok: link foo at subcommand",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use:               "foo",
						PlaceCommandUnder: "starport scaffold",
					},
				},
			},
			expectedDumpCmd: `
starport
  scaffold
    chain*
    foo*
`,
		},
		{
			name: "ok: link foo at subcommand with incomplete PlaceCommandUnder",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use:               "foo",
						PlaceCommandUnder: "scaffold",
					},
				},
			},
			expectedDumpCmd: `
starport
  scaffold
    chain*
    foo*
`,
		},
		{
			name: "fail: link to runnable command",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use:               "foo",
						PlaceCommandUnder: "starport scaffold chain",
					},
				},
			},
			expectedError: `can't attach plugin command "foo" to runnable command "starport scaffold chain"`,
		},
		{
			name: "fail: link to unknown command",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use:               "foo",
						PlaceCommandUnder: "starport unknown",
					},
				},
			},
			expectedError: `unable to find commandPath "starport unknown" for plugin "foo"`,
		},
		{
			name: "fail: plugin name exists in legacy commands",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "scaffold",
					},
				},
			},
			expectedError: `plugin command "scaffold" already exists in starport's commands`,
		},
		{
			name: "fail: plugin name exists in legacy sub commands",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use:               "chain",
						PlaceCommandUnder: "scaffold",
					},
				},
			},
			expectedError: `plugin command "chain" already exists in starport's commands`,
		},
		{
			name: "ok: link foo and bar at root",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "foo",
					},
					{
						Use: "bar",
					},
				},
			},
			expectedDumpCmd: `
starport
  bar*
  foo*
  scaffold
    chain*
`,
		},
		{
			name: "ok: link with subcommands",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "foo",
						Commands: []plugin.Command{
							{Use: "bar"},
							{Use: "baz"},
						},
					},
				},
			},
			expectedDumpCmd: `
starport
  foo
    bar*
    baz*
  scaffold
    chain*
`,
		},
		{
			name: "ok: link with multiple subcommands",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "foo",
						Commands: []plugin.Command{
							{Use: "bar", Commands: []plugin.Command{{Use: "baz"}}},
							{Use: "qux", Commands: []plugin.Command{{Use: "quux"}, {Use: "corge"}}},
						},
					},
				},
			},
			expectedDumpCmd: `
starport
  foo
    bar
      baz*
    qux
      corge*
      quux*
  scaffold
    chain*
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			p := &plugin.Plugin{
				Plugin:    chainconfig.Plugin{Name: "foo"},
				Interface: tt.pluginInterface,
			}
			rootCmd := buildRootCmd()

			linkPluginCmds(rootCmd, p)

			if tt.expectedError != "" {
				require.Error(p.Error)
				require.EqualError(p.Error, tt.expectedError)
				return
			}
			require.NoError(p.Error)
			var s strings.Builder
			s.WriteString("\n")
			dumpCmd(rootCmd, &s, 0)
			assert.Equal(tt.expectedDumpCmd, s.String())
		})
	}
}

// dumpCmd helps in comparing cobra.Command by writing their Use and Commands.
// Runnable commands are marked with a *.
func dumpCmd(c *cobra.Command, w io.Writer, ntabs int) {
	fmt.Fprintf(w, "%s%s", strings.Repeat("  ", ntabs), c.Use)
	ntabs++
	if c.Runnable() {
		fmt.Fprintf(w, "*")
	}
	fmt.Fprintf(w, "\n")
	for _, cc := range c.Commands() {
		dumpCmd(cc, w, ntabs)
	}
}
