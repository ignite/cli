package ignitecmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/services/plugin"
)

// pluginInterface implements plugin.Interface for testing purpose.
type pluginInterface struct {
	commands []plugin.Command
	hooks    []plugin.Hook

	// hookCalls holds trace of ExecuteHook* methods' invocation.
	hookCalls map[string][]string
	// holds arguments tied to the ExecuteHook* methods' invocation.
	hookArgs map[string]map[string][]string
}

func (p *pluginInterface) Commands() []plugin.Command {
	return p.commands
}

func (p *pluginInterface) Hooks() []plugin.Hook {
	return p.hooks
}

func (p *pluginInterface) Execute(c plugin.Command, args []string) error {
	return nil
}

func (p *pluginInterface) ExecuteHookPre(hook plugin.Hook, args []string) error {
	if p.hookCalls == nil {
		p.hookCalls = make(map[string][]string)
	}

	p.hookCalls[hook.PlaceHookOn] = append(p.hookCalls[hook.PlaceHookOn],
		fmt.Sprintf("pre-%s", hook.Name))

	if p.hookArgs == nil && len(args) > 0 {
		p.hookArgs = make(map[string]map[string][]string)
	} else if len(args) > 0 {
		p.hookArgs[hook.PlaceHookOn] = make(map[string][]string)
		p.hookArgs[hook.PlaceHookOn]["pre"] = args
	}
	return nil
}

func (p *pluginInterface) ExecuteHookPost(hook plugin.Hook, args []string) error {
	if p.hookCalls == nil {
		p.hookCalls = make(map[string][]string)
	}
	p.hookCalls[hook.PlaceHookOn] = append(p.hookCalls[hook.PlaceHookOn],
		fmt.Sprintf("post-%s", hook.Name))

	if p.hookArgs == nil && len(args) > 0 {
		p.hookArgs = make(map[string]map[string][]string)
	} else if len(args) > 0 {
		if p.hookArgs[hook.PlaceHookOn] == nil {
			return fmt.Errorf("post hook executed before pre for hook %q aborting", hook.Name)
		}
		p.hookArgs[hook.PlaceHookOn]["post"] = args
	}
	return nil
}

func (p *pluginInterface) ExecuteHookCleanUp(hook plugin.Hook, args []string) error {
	if p.hookCalls == nil {
		p.hookCalls = make(map[string][]string)
	}
	p.hookCalls[hook.PlaceHookOn] = append(p.hookCalls[hook.PlaceHookOn],
		fmt.Sprintf("cleanup-%s", hook.Name))

	if p.hookArgs == nil && len(args) > 0 {
		p.hookArgs = make(map[string]map[string][]string)
	} else if len(args) > 0 {
		if p.hookArgs[hook.PlaceHookOn] == nil {
			return fmt.Errorf("cleanup hook executed before pre for hook %q aborting", hook.Name)
		}
		p.hookArgs[hook.PlaceHookOn]["cleanup"] = args
	}
	return nil
}

func buildRootCmd() *cobra.Command {
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
		scaffoldModuleCmd = &cobra.Command{
			Use: "module",
			Run: func(*cobra.Command, []string) {},
		}
	)

	// test flag for passing to hook life cycles
	scaffoldChainCmd.Flags().AddFlag(&pflag.Flag{
		Name:      "flag",
		Shorthand: "f",
		Usage:     "test flag",
	})
	scaffoldModuleCmd.Flags().AddFlag(&pflag.Flag{
		Name:      "flag",
		Shorthand: "f",
		Usage:     "test flag",
	})

	scaffoldCmd.AddCommand(scaffoldChainCmd)
	scaffoldCmd.AddCommand(scaffoldModuleCmd)
	rootCmd.AddCommand(scaffoldCmd)
	return rootCmd
}

func TestLinkPluginCmds(t *testing.T) {
	tests := []struct {
		name            string
		pluginInterface *pluginInterface
		expectedDumpCmd string
		expectedError   string
	}{
		{
			name: "ok: link foo at root",
			pluginInterface: &pluginInterface{
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
    module*
`,
		},
		{
			name: "ok: link foo at subcommand",
			pluginInterface: &pluginInterface{
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
    module*
`,
		},
		{
			name: "ok: link foo at subcommand with incomplete PlaceCommandUnder",
			pluginInterface: &pluginInterface{
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
    module*
`,
		},
		{
			name: "fail: link to runnable command",
			pluginInterface: &pluginInterface{
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
			pluginInterface: &pluginInterface{
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
			pluginInterface: &pluginInterface{
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
			pluginInterface: &pluginInterface{
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
			pluginInterface: &pluginInterface{
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
    module*
`,
		},
		{
			name: "ok: link with subcommands",
			pluginInterface: &pluginInterface{
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
    module*
`,
		},
		{
			name: "ok: link with multiple subcommands",
			pluginInterface: &pluginInterface{
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
    module*
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

func TestLinkPluginHooks(t *testing.T) {
	tests := []struct {
		name            string
		pluginInterface *pluginInterface
		args            []string
		expectedError   string
		expectedCalls   map[string][]string
		epectedArgs     map[string]map[string][]string
	}{
		{
			name: "fail: hook plugin command",
			pluginInterface: &pluginInterface{
				commands: []plugin.Command{{
					Use: "test-plugin",
				}},
				hooks: []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite test-plugin",
					},
				},
			},
			expectedError: `unable to find commandPath "ignite test-plugin" for plugin hook "test-hook"`,
		},
		{
			name: "fail: command not runnable",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold",
					},
				},
			},
			expectedError: `can't attach plugin hook "test-hook" to non executable command "ignite scaffold"`,
		},
		{
			name: "fail: command doesn't exists",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite chain",
					},
				},
			},
			expectedError: `unable to find commandPath "ignite chain" for plugin hook "test-hook"`,
		},
		{
			name: "ok: single hook",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "scaffold chain",
					},
				},
			},
			expectedCalls: map[string][]string{
				"scaffold chain": {
					"pre-test-hook", "post-test-hook", "cleanup-test-hook",
				},
			},
		},
		{
			name: "ok: multiple hooks on same command",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook-1",
						PlaceHookOn: "scaffold chain",
					},
					{
						Name:        "test-hook-2",
						PlaceHookOn: "scaffold chain",
					},
				},
			},
			expectedCalls: map[string][]string{
				"scaffold chain": {
					"pre-test-hook-1", "pre-test-hook-2",
					"post-test-hook-1", "cleanup-test-hook-1",
					"post-test-hook-2", "cleanup-test-hook-2",
				},
			},
		},
		{
			name: "ok: multiple hooks on different commands",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook-1",
						PlaceHookOn: "scaffold chain",
					},
					{
						Name:        "test-hook-2",
						PlaceHookOn: "scaffold chain",
					},
					{
						Name:        "test-hook-3",
						PlaceHookOn: "scaffold module",
					},
				},
			},
			args: []string{"flag foo"},
			expectedCalls: map[string][]string{
				"scaffold chain": {
					"pre-test-hook-1", "pre-test-hook-2",
					"post-test-hook-1", "cleanup-test-hook-1",
					"post-test-hook-2", "cleanup-test-hook-2",
				},
				"scaffold module": {
					"pre-test-hook-3", "post-test-hook-3", "cleanup-test-hook-3",
				},
			},
			epectedArgs: map[string]map[string][]string{
				"scaffold chain": {
					"pre":     {"flag foo"},
					"post":    {"flag foo"},
					"cleanup": {"flag foo"},
				},
				"scaffold module": {
					"pre":     {"flag foo"},
					"post":    {"flag foo"},
					"cleanup": {"flag foo"},
				},
			},
		},
		{
			name: "ok: duplicate hook names on same command",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
				},
			},
			args: []string{"flag foo"},
			expectedCalls: map[string][]string{
				"ignite scaffold chain": {
					"pre-test-hook", "pre-test-hook",
					"post-test-hook", "cleanup-test-hook",
					"post-test-hook", "cleanup-test-hook",
				},
			},
			epectedArgs: map[string]map[string][]string{
				"ignite scaffold chain": {
					"pre":     {"flag foo"},
					"post":    {"flag foo"},
					"cleanup": {"flag foo"},
				},
			},
		},
		{
			name: "ok: duplicate hook names on different commands",
			pluginInterface: &pluginInterface{
				hooks: []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold module",
					},
				},
			},
			expectedCalls: map[string][]string{
				"ignite scaffold chain": {
					"pre-test-hook", "post-test-hook", "cleanup-test-hook",
				},
				"ignite scaffold module": {
					"pre-test-hook", "post-test-hook", "cleanup-test-hook",
				},
			},
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

			if tt.expectedError != "" {
				require.EqualError(p.Error, tt.expectedError)
				return
			}
			require.NoError(p.Error)
			execCmd(t, rootCmd, tt.args)
			assert.Equal(tt.expectedCalls, tt.pluginInterface.hookCalls)
			assert.Equal(tt.epectedArgs, tt.pluginInterface.hookArgs)
		})
	}
}

// execCmd executes all the runnable commands contained in c.
func execCmd(t *testing.T, c *cobra.Command, args []string) {
	if c.Runnable() {
		os.Args = strings.Fields(c.CommandPath())
		os.Args = append(os.Args, args...)
		err := c.Execute()
		require.NoError(t, err)
		return
	}
	for _, c := range c.Commands() {
		execCmd(t, c, args)
	}
}
