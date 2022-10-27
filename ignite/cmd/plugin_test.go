package ignitecmd

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/services/plugin"
)

// pluginInterface implements plugin.Interface for testing purpose.
type pluginInterface struct {
	commands []plugin.Command
	hooks    []plugin.Hook
}

func (p pluginInterface) Commands() []plugin.Command {
	return p.commands
}

func (p pluginInterface) Hooks() []plugin.Hook {
	return p.hooks
}

func (pluginInterface) Execute(plugin.Command, []string) error {
	return nil
}

func (pluginInterface) ExecuteHookPre(name string, args []string) error {
	fmt.Printf("Executing pre run behavior for %s\n", name)
	return nil
}

func (pluginInterface) ExecuteHookPost(name string, args []string) error {
	fmt.Printf("Executing post run behavior for %s\n", name)
	return nil
}

func (pluginInterface) ExecuteHookCleanUp(name string, args []string) error {
	fmt.Printf("Executing cleanup behavior for %s\n", name)
	return nil
}

func TestLinkPluginCmds(t *testing.T) {
	buildRootCmd := func() *cobra.Command {
		var (
			rootCmd = &cobra.Command{
				Use: "ignite",
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
ignite
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
						PlaceCommandUnder: "ignite scaffold",
					},
				},
			},
			expectedDumpCmd: `
ignite
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
ignite
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
						PlaceCommandUnder: "ignite scaffold chain",
					},
				},
			},
			expectedError: `can't attach plugin command "foo" to runnable command "ignite scaffold chain"`,
		},
		{
			name: "fail: link to unknown command",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use:               "foo",
						PlaceCommandUnder: "ignite unknown",
					},
				},
			},
			expectedError: `unable to find commandPath "ignite unknown" for plugin "foo"`,
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
			expectedError: `plugin command "scaffold" already exists in ignite's commands`,
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
			expectedError: `plugin command "chain" already exists in ignite's commands`,
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
ignite
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
ignite
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
ignite
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
				Plugin:    chainconfig.Plugin{Path: "foo"},
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

func TestLinkPluginHooks(t *testing.T) {
	buildRootCmd := func() *cobra.Command {
		var (
			rootCmd = &cobra.Command{
				Use: "ignite",
				RunE: func(*cobra.Command, []string) error {
					return nil
				},
			}
			scaffoldCmd = &cobra.Command{
				Use: "scaffold",
				RunE: func(*cobra.Command, []string) error {
					return nil
				},
			}
			scaffoldChainCmd = &cobra.Command{
				Use: "chain",
				RunE: func(*cobra.Command, []string) error {
					return nil
				},
			}
		)
		scaffoldCmd.AddCommand(scaffoldChainCmd)
		rootCmd.AddCommand(scaffoldCmd)
		return rootCmd
	}
	tests := []struct {
		name            string
		pluginInterface pluginInterface
		shouldError     bool
	}{
		{
			name: "error: invalid command path",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "foo",
					},
				},
				hooks: []plugin.Hook{
					{
						Name:        "test-hook-1",
						PlaceHookOn: "ignite scaffold",
					},
				},
			},
			shouldError: false,
		},
		{
			name: "ok: link foo at subcommand",
			pluginInterface: pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook-2",
						PlaceHookOn: "ignite chain",
					},
				},
			},
			shouldError: true,
		},
		{
			name: "ok: should prepend root command",
			pluginInterface: pluginInterface{
				commands: []plugin.Command{
					{
						Use: "foo",
					},
				},
				hooks: []plugin.Hook{
					{
						Name:        "test-hook-3",
						PlaceHookOn: "scaffold chain",
					},
				},
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			p := &plugin.Plugin{
				Plugin:    chainconfig.Plugin{Path: "foo"},
				Interface: tt.pluginInterface,
			}
			rootCmd := buildRootCmd()

			linkPluginHooks(rootCmd, p)
			if tt.shouldError {
				require.Error(p.Error)
			} else {
				require.NoError(p.Error)
			}

			for _, hook := range p.Interface.Hooks() {
				areHooksDefined := checkCmd(rootCmd, hook.PlaceHookOn)
				if tt.shouldError {
					assert.False(areHooksDefined, "hooks should have pre and post runs undefined")
				} else {
					assert.True(areHooksDefined, "hooks should have pre and post runs defined")
				}
			}
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

func checkCmd(c *cobra.Command, path string) bool {
	// bring helper for path prefix from plugin.go
	if !strings.HasPrefix(path, "ignite") {
		// cmdPath must start with `ignite ` before comparison with
		// cmd.CommandPath()
		path = "ignite " + path
	}

	isDefined := false
	command := findCommandByPath(c, path)
	if command == nil {
		return false
	}

	isDefined = command.PreRun != nil

	if !isDefined {
		return isDefined
	}

	isDefined = command.PostRunE != nil

	return isDefined
}
