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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/services/plugin"
	"github.com/ignite/cli/ignite/services/plugin/mocks"
)

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
	scaffoldChainCmd.Flags().String("path", "", "the path")
	scaffoldCmd.AddCommand(scaffoldChainCmd)
	scaffoldCmd.AddCommand(scaffoldModuleCmd)
	rootCmd.AddCommand(scaffoldCmd)
	return rootCmd
}

func TestLinkPluginCmds(t *testing.T) {
	var (
		args         = []string{"arg1", "arg2"}
		pluginParams = map[string]string{"key": "val"}
		// define a plugin with command flags
		pluginWithFlags = plugin.Command{
			Use: "flaggy",
		}
	)
	pluginWithFlags.Flags().String("flag1", "", "")
	pluginWithFlags.Flags().Int("flag2", 0, "")

	// helper to assert pluginInterface.Execute() calls
	expectExecute := func(t *testing.T, p *mocks.PluginInterface, cmd plugin.Command) {
		p.EXPECT().Execute(
			mock.MatchedBy(func(execCmd plugin.Command) bool {
				return cmd.Use == execCmd.Use
			}),
			args,
		).Run(func(execCmd plugin.Command, args []string) {
			// Assert flags are passed to Execute
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				assert.NotNil(t, execCmd.Flags().Lookup(f.Name), "flag %s not found", f.Name)
			})
			// Assert With is populated
			assert.Equal(t, pluginParams, execCmd.With)
		}).Return(nil)
	}

	tests := []struct {
		name            string
		setup           func(*testing.T, *mocks.PluginInterface)
		expectedDumpCmd string
		expectedError   string
	}{
		{
			name: "ok: link foo at root",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				cmd := plugin.Command{
					Use: "foo",
				}
				p.EXPECT().Commands().Return([]plugin.Command{cmd}, nil)
				expectExecute(t, p, cmd)
			},
			expectedDumpCmd: `
ignite
  foo*
  scaffold
    chain* --path=string
    module*
`,
		},
		{
			name: "ok: link foo at subcommand",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				cmd := plugin.Command{
					Use:               "foo",
					PlaceCommandUnder: "ignite scaffold",
				}
				p.EXPECT().Commands().Return([]plugin.Command{cmd}, nil)
				expectExecute(t, p, cmd)
			},
			expectedDumpCmd: `
ignite
  scaffold
    chain* --path=string
    foo*
    module*
`,
		},
		{
			name: "ok: link foo at subcommand with incomplete PlaceCommandUnder",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				cmd := plugin.Command{
					Use:               "foo",
					PlaceCommandUnder: "scaffold",
				}
				p.EXPECT().Commands().Return([]plugin.Command{cmd}, nil)
				expectExecute(t, p, cmd)
			},
			expectedDumpCmd: `
ignite
  scaffold
    chain* --path=string
    foo*
    module*
`,
		},
		{
			name: "fail: link to runnable command",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				p.EXPECT().Commands().Return(
					[]plugin.Command{
						{
							Use:               "foo",
							PlaceCommandUnder: "ignite scaffold chain",
						},
					},
					nil,
				)
			},
			expectedError: `can't attach plugin command "foo" to runnable command "ignite scaffold chain"`,
		},
		{
			name: "fail: link to unknown command",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				p.EXPECT().Commands().Return(
					[]plugin.Command{
						{
							Use:               "foo",
							PlaceCommandUnder: "ignite unknown",
						},
					},
					nil,
				)
			},
			expectedError: `unable to find commandPath "ignite unknown" for plugin "foo"`,
		},
		{
			name: "fail: plugin name exists in legacy commands",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				p.EXPECT().Commands().Return(
					[]plugin.Command{
						{
							Use: "scaffold",
						},
					},
					nil,
				)
			},
			expectedError: `plugin command "scaffold" already exists in ignite's commands`,
		},
		{
			name: "fail: plugin name exists in legacy sub commands",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				p.EXPECT().Commands().Return(
					[]plugin.Command{
						{
							Use:               "chain",
							PlaceCommandUnder: "scaffold",
						},
					},
					nil,
				)
			},
			expectedError: `plugin command "chain" already exists in ignite's commands`,
		},
		{
			name: "ok: link multiple at root",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				fooCmd := plugin.Command{
					Use: "foo",
				}
				barCmd := plugin.Command{
					Use: "bar",
				}
				p.EXPECT().Commands().Return([]plugin.Command{
					fooCmd, barCmd, pluginWithFlags,
				}, nil)
				expectExecute(t, p, fooCmd)
				expectExecute(t, p, barCmd)
				expectExecute(t, p, pluginWithFlags)
			},
			expectedDumpCmd: `
ignite
  bar*
  flaggy* --flag1=string --flag2=int
  foo*
  scaffold
    chain* --path=string
    module*
`,
		},
		{
			name: "ok: link with subcommands",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				cmd := plugin.Command{
					Use: "foo",
					Commands: []plugin.Command{
						{Use: "bar"},
						{Use: "baz"},
						pluginWithFlags,
					},
				}
				p.EXPECT().Commands().Return([]plugin.Command{cmd}, nil)
				// cmd is not executed because it's not runnable, only sub-commands
				// are executed.
				expectExecute(t, p, cmd.Commands[0])
				expectExecute(t, p, cmd.Commands[1])
				expectExecute(t, p, cmd.Commands[2])
			},
			expectedDumpCmd: `
ignite
  foo
    bar*
    baz*
    flaggy* --flag1=string --flag2=int
  scaffold
    chain* --path=string
    module*
`,
		},
		{
			name: "ok: link with multiple subcommands",
			setup: func(t *testing.T, p *mocks.PluginInterface) {
				cmd := plugin.Command{
					Use: "foo",
					Commands: []plugin.Command{
						{Use: "bar", Commands: []plugin.Command{{Use: "baz"}}},
						{Use: "qux", Commands: []plugin.Command{{Use: "quux"}, {Use: "corge"}}},
					},
				}
				p.EXPECT().Commands().Return([]plugin.Command{cmd}, nil)
				expectExecute(t, p, cmd.Commands[0].Commands[0])
				expectExecute(t, p, cmd.Commands[1].Commands[0])
				expectExecute(t, p, cmd.Commands[1].Commands[1])
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
    chain* --path=string
    module*
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			pi := mocks.NewPluginInterface(t)
			p := &plugin.Plugin{
				Plugin: chainconfig.Plugin{
					Path: "foo",
					With: pluginParams,
				},
				Interface: pi,
			}
			rootCmd := buildRootCmd()
			tt.setup(t, pi)

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
			execCmd(t, rootCmd, args)
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
	c.Flags().VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(w, " --%s=%s", f.Name, f.Value.Type())
	})
	fmt.Fprintf(w, "\n")
	for _, cc := range c.Commands() {
		dumpCmd(cc, w, ntabs)
	}
}

func TestLinkPluginHooks(t *testing.T) {
	var (
		args         = []string{"arg1", "arg2"}
		pluginParams = map[string]string{"key": "val"}

		// helper to assert pluginInterface.ExecuteHook*() calls in expected order
		// (pre, then post, then cleanup)
		expectExecuteHook = func(p *mocks.PluginInterface, hooks ...plugin.Hook) {
			var lastPre *mock.Call
			for _, hook := range hooks {
				// assert that hook.With has been populated before call to ExecuteHook
				hook.With = pluginParams
				pre := p.EXPECT().ExecuteHookPre(hook, args).Return(nil).Call
				if lastPre != nil {
					pre.NotBefore(lastPre)
				}
				lastPre = pre
			}
			for _, hook := range hooks {
				// assert that hook.With has been populated before call to ExecuteHook
				hook.With = pluginParams
				post := p.EXPECT().ExecuteHookPost(hook, args).Return(nil).Call
				cleanup := p.EXPECT().ExecuteHookCleanUp(hook, args).Return(nil).Call
				post.NotBefore(lastPre)
				cleanup.NotBefore(post)
			}
		}
	)
	tests := []struct {
		name          string
		expectedError string
		setup         func(*mocks.PluginInterface)
	}{
		// TODO(tb): commented because linkPluginCmds is not invoked in this test,
		// so it's not possible to assert that a hook can't be placed on a plugin
		// command.
		/*
			{
				name: "fail: hook plugin command",
				setup: func(p *mocks.PluginInterface) {
					p.EXPECT().Commands().Return([]plugin.Command{{Use: "test-plugin"}}, nil)
					p.EXPECT().Hooks().Return(
						[]plugin.Hook{
							{
								Name:        "test-hook",
								PlaceHookOn: "ignite test-plugin",
							},
						},
						nil,
					)
				},
				expectedError: `unable to find commandPath "ignite test-plugin" for plugin hook "test-hook"`,
			},
		*/
		{
			name: "fail: command not runnable",
			setup: func(p *mocks.PluginInterface) {
				p.EXPECT().Hooks().Return(
					[]plugin.Hook{
						{
							Name:        "test-hook",
							PlaceHookOn: "ignite scaffold",
						},
					},
					nil,
				)
			},
			expectedError: `can't attach plugin hook "test-hook" to non executable command "ignite scaffold"`,
		},
		{
			name: "fail: command doesn't exists",
			setup: func(p *mocks.PluginInterface) {
				p.EXPECT().Hooks().Return(
					[]plugin.Hook{
						{
							Name:        "test-hook",
							PlaceHookOn: "ignite chain",
						},
					},
					nil,
				)
			},
			expectedError: `unable to find commandPath "ignite chain" for plugin hook "test-hook"`,
		},
		{
			name: "ok: single hook",
			setup: func(p *mocks.PluginInterface) {
				hook := plugin.Hook{
					Name:        "test-hook",
					PlaceHookOn: "scaffold chain",
				}
				p.EXPECT().Hooks().Return([]plugin.Hook{hook}, nil)
				expectExecuteHook(p, hook)
			},
		},
		{
			name: "ok: multiple hooks on same command",
			setup: func(p *mocks.PluginInterface) {
				hook1 := plugin.Hook{
					Name:        "test-hook-1",
					PlaceHookOn: "scaffold chain",
				}
				hook2 := plugin.Hook{
					Name:        "test-hook-2",
					PlaceHookOn: "scaffold chain",
				}
				p.EXPECT().Hooks().Return([]plugin.Hook{hook1, hook2}, nil)
				expectExecuteHook(p, hook1, hook2)
			},
		},
		{
			name: "ok: multiple hooks on different commands",
			setup: func(p *mocks.PluginInterface) {
				hookChain1 := plugin.Hook{
					Name:        "test-hook-1",
					PlaceHookOn: "scaffold chain",
				}
				hookChain2 := plugin.Hook{
					Name:        "test-hook-2",
					PlaceHookOn: "scaffold chain",
				}
				hookModule := plugin.Hook{
					Name:        "test-hook-3",
					PlaceHookOn: "scaffold module",
				}
				p.EXPECT().Hooks().Return([]plugin.Hook{hookChain1, hookChain2, hookModule}, nil)
				expectExecuteHook(p, hookChain1, hookChain2)
				expectExecuteHook(p, hookModule)
			},
		},
		{
			name: "ok: duplicate hook names on same command",
			setup: func(p *mocks.PluginInterface) {
				hooks := []plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
				}
				p.EXPECT().Hooks().Return(hooks, nil)
				expectExecuteHook(p, hooks...)
			},
		},
		{
			name: "ok: duplicate hook names on different commands",
			setup: func(p *mocks.PluginInterface) {
				hookChain := plugin.Hook{
					Name:        "test-hook",
					PlaceHookOn: "ignite scaffold chain",
				}
				hookModule := plugin.Hook{
					Name:        "test-hook",
					PlaceHookOn: "ignite scaffold module",
				}
				p.EXPECT().Hooks().Return([]plugin.Hook{hookChain, hookModule}, nil)
				expectExecuteHook(p, hookChain)
				expectExecuteHook(p, hookModule)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			// assert := assert.New(t)
			pi := mocks.NewPluginInterface(t)
			p := &plugin.Plugin{
				Plugin: chainconfig.Plugin{
					Path: "foo",
					With: pluginParams,
				},
				Interface: pi,
			}
			rootCmd := buildRootCmd()
			tt.setup(pi)

			linkPluginHooks(rootCmd, p)

			if tt.expectedError != "" {
				require.EqualError(p.Error, tt.expectedError)
				return
			}
			require.NoError(p.Error)
			execCmd(t, rootCmd, args)
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
