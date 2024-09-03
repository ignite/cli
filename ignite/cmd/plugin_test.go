package ignitecmd

import (
	"context"
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

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	"github.com/ignite/cli/v29/ignite/services/plugin/mocks"
)

func buildRootCmd(ctx context.Context) *cobra.Command {
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
	rootCmd.SetContext(ctx)
	return rootCmd
}

func assertFlags(t *testing.T, expectedFlags plugin.Flags, execCmd *plugin.ExecutedCommand) {
	t.Helper()
	var (
		have     []string
		expected []string
	)

	t.Helper()

	flags, err := execCmd.NewFlags()
	assert.NoError(t, err)

	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == "help" {
			// ignore help flag
			return
		}

		have = append(have, f.Name)
	})

	for _, f := range expectedFlags {
		expected = append(expected, f.Name)
	}

	assert.Equal(t, expected, have)
}

func TestLinkPluginCmds(t *testing.T) {
	t.Skip("passes locally and with act, but fails in CI")

	var (
		args         = []string{"arg1", "arg2"}
		pluginParams = map[string]string{"key": "val"}
		// define a plugin with command flags
		pluginWithFlags = &plugin.Command{
			Use: "flaggy",
			Flags: plugin.Flags{
				{Name: "flag1", Type: plugin.FlagTypeString},
				{Name: "flag2", Type: plugin.FlagTypeInt, DefaultValue: "0", Value: "0"},
			},
		}
	)

	// helper to assert pluginInterface.Execute() calls
	expectExecute := func(t *testing.T, _ context.Context, p *mocks.PluginInterface, cmd *plugin.Command) {
		t.Helper()
		p.EXPECT().
			Execute(
				mock.Anything,
				mock.MatchedBy(func(execCmd *plugin.ExecutedCommand) bool {
					fmt.Println(cmd.Use == execCmd.Use, cmd.Use, execCmd.Use)
					return cmd.Use == execCmd.Use
				}),
				mock.Anything,
			).
			Run(func(_ context.Context, execCmd *plugin.ExecutedCommand, _ plugin.ClientAPI) {
				// Assert execCmd is populated correctly
				assert.True(t, strings.HasSuffix(execCmd.Path, cmd.Use), "wrong path %s", execCmd.Path)
				assert.Equal(t, args, execCmd.Args)
				assertFlags(t, cmd.Flags, execCmd)
				assert.Equal(t, pluginParams, execCmd.With)
			}).
			Return(nil)
	}

	tests := []struct {
		name            string
		setup           func(*testing.T, context.Context, *mocks.PluginInterface)
		expectedDumpCmd string
		expectedError   string
	}{
		{
			name: "ok: link foo at root",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				cmd := &plugin.Command{
					Use: "foo",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Commands: []*plugin.Command{cmd}}, nil)
				expectExecute(t, ctx, p, cmd)
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
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				cmd := &plugin.Command{
					Use:               "foo",
					PlaceCommandUnder: "ignite scaffold",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Commands: []*plugin.Command{cmd}}, nil)
				expectExecute(t, ctx, p, cmd)
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
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				cmd := &plugin.Command{
					Use:               "foo",
					PlaceCommandUnder: "scaffold",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Commands: []*plugin.Command{cmd}}, nil)
				expectExecute(t, ctx, p, cmd)
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
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Commands: []*plugin.Command{
							{
								Use:               "foo",
								PlaceCommandUnder: "ignite scaffold chain",
							},
						},
					},
						nil,
					)
			},
			expectedError: `can't attach app command "foo" to runnable command "ignite scaffold chain"`,
		},
		{
			name: "fail: link to unknown command",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Commands: []*plugin.Command{
							{
								Use:               "foo",
								PlaceCommandUnder: "ignite unknown",
							},
						},
					},
						nil,
					)
			},
			expectedError: `unable to find command path "ignite unknown" for app "foo"`,
		},
		{
			name: "fail: plugin name exists in legacy commands",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Commands: []*plugin.Command{
							{
								Use: "scaffold",
							},
						},
					},
						nil,
					)
			},
			expectedError: `app command "scaffold" already exists in Ignite's commands`,
		},
		{
			name: "fail: plugin name with args exists in legacy commands",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Commands: []*plugin.Command{
							{
								Use: "scaffold [args]",
							},
						},
					},
						nil,
					)
			},
			expectedError: `app command "scaffold" already exists in Ignite's commands`,
		},
		{
			name: "fail: plugin name exists in legacy sub commands",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Commands: []*plugin.Command{
							{
								Use:               "chain",
								PlaceCommandUnder: "scaffold",
							},
						},
					},
						nil,
					)
			},
			expectedError: `app command "chain" already exists in Ignite's commands`,
		},
		{
			name: "ok: link multiple at root",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				fooCmd := &plugin.Command{
					Use: "foo",
				}
				barCmd := &plugin.Command{
					Use: "bar",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Commands: []*plugin.Command{
							fooCmd, barCmd, pluginWithFlags,
						},
					}, nil)
				expectExecute(t, ctx, p, fooCmd)
				expectExecute(t, ctx, p, barCmd)
				expectExecute(t, ctx, p, pluginWithFlags)
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
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				cmd := &plugin.Command{
					Use: "foo",
					Commands: []*plugin.Command{
						{Use: "bar"},
						{Use: "baz"},
						pluginWithFlags,
					},
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Commands: []*plugin.Command{cmd}}, nil)
				// cmd is not executed because it's not runnable, only sub-commands
				// are executed.
				expectExecute(t, ctx, p, cmd.Commands[0])
				expectExecute(t, ctx, p, cmd.Commands[1])
				expectExecute(t, ctx, p, cmd.Commands[2])
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
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				cmd := &plugin.Command{
					Use: "foo",
					Commands: []*plugin.Command{
						{Use: "bar", Commands: []*plugin.Command{{Use: "baz"}}},
						{Use: "qux", Commands: []*plugin.Command{{Use: "quux"}, {Use: "corge"}}},
					},
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Commands: []*plugin.Command{cmd}}, nil)
				expectExecute(t, ctx, p, cmd.Commands[0].Commands[0])
				expectExecute(t, ctx, p, cmd.Commands[1].Commands[0])
				expectExecute(t, ctx, p, cmd.Commands[1].Commands[1])
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			require := require.New(t)
			assert := assert.New(t)
			pi := mocks.NewPluginInterface(t)
			p := &plugin.Plugin{
				Plugin: pluginsconfig.Plugin{
					Path: "foo",
					With: pluginParams,
				},
				Interface: pi,
			}
			rootCmd := buildRootCmd(ctx)
			tt.setup(t, ctx, pi)

			_ = linkPlugins(ctx, rootCmd, []*plugin.Plugin{p})

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
	t.Skip("passes locally and with act, but fails in CI")

	var (
		args         = []string{"arg1", "arg2"}
		pluginParams = map[string]string{"key": "val"}
		ctx          = context.Background()

		// helper to assert pluginInterface.ExecuteHook*() calls in expected order
		// (pre, then post, then cleanup)
		expectExecuteHook = func(t *testing.T, p *mocks.PluginInterface, expectedFlags plugin.Flags, hooks ...*plugin.Hook) {
			t.Helper()
			matcher := func(hook *plugin.Hook) any {
				return mock.MatchedBy(func(execHook *plugin.ExecutedHook) bool {
					return hook.Name == execHook.Hook.Name &&
						hook.PlaceHookOn == execHook.Hook.PlaceHookOn
				})
			}
			asserter := func(hook *plugin.Hook) func(_ context.Context, hook *plugin.ExecutedHook, _ plugin.ClientAPI) {
				return func(_ context.Context, execHook *plugin.ExecutedHook, _ plugin.ClientAPI) {
					assert.True(t, strings.HasSuffix(execHook.ExecutedCommand.Path, hook.PlaceHookOn), "wrong path %q want %q", execHook.ExecutedCommand.Path, hook.PlaceHookOn)
					assert.Equal(t, args, execHook.ExecutedCommand.Args)
					assertFlags(t, expectedFlags, execHook.ExecutedCommand)
					assert.Equal(t, pluginParams, execHook.ExecutedCommand.With)
				}
			}
			var lastPre *mock.Call
			for _, hook := range hooks {
				pre := p.EXPECT().
					ExecuteHookPre(ctx, matcher(hook), mock.Anything).
					Run(asserter(hook)).
					Return(nil).
					Call
				if lastPre != nil {
					pre.NotBefore(lastPre)
				}
				lastPre = pre
			}
			for _, hook := range hooks {
				post := p.EXPECT().
					ExecuteHookPost(ctx, matcher(hook), mock.Anything).
					Run(asserter(hook)).
					Return(nil).
					Call
				cleanup := p.EXPECT().
					ExecuteHookCleanUp(ctx, matcher(hook), mock.Anything).
					Run(asserter(hook)).
					Return(nil).
					Call
				post.NotBefore(lastPre)
				cleanup.NotBefore(post)
			}
		}
	)
	tests := []struct {
		name          string
		expectedError string
		setup         func(*testing.T, context.Context, *mocks.PluginInterface)
	}{
		{
			name: "fail: command not runnable",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Hooks: []*plugin.Hook{
							{
								Name:        "test-hook",
								PlaceHookOn: "ignite scaffold",
							},
						},
					},
						nil,
					)
			},
			expectedError: `can't attach app hook "test-hook" to non executable command "ignite scaffold"`,
		},
		{
			name: "fail: command doesn't exists",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{
						Hooks: []*plugin.Hook{
							{
								Name:        "test-hook",
								PlaceHookOn: "ignite chain",
							},
						},
					},
						nil,
					)
			},
			expectedError: `unable to find command path "ignite chain" for app hook "test-hook"`,
		},
		{
			name: "ok: single hook",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				hook := &plugin.Hook{
					Name:        "test-hook",
					PlaceHookOn: "scaffold chain",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Hooks: []*plugin.Hook{hook}}, nil)
				expectExecuteHook(t, p, plugin.Flags{{Name: "path"}}, hook)
			},
		},
		{
			name: "ok: multiple hooks on same command",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				hook1 := &plugin.Hook{
					Name:        "test-hook-1",
					PlaceHookOn: "scaffold chain",
				}
				hook2 := &plugin.Hook{
					Name:        "test-hook-2",
					PlaceHookOn: "scaffold chain",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Hooks: []*plugin.Hook{hook1, hook2}}, nil)
				expectExecuteHook(t, p, plugin.Flags{{Name: "path"}}, hook1, hook2)
			},
		},
		{
			name: "ok: multiple hooks on different commands",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				hookChain1 := &plugin.Hook{
					Name:        "test-hook-1",
					PlaceHookOn: "scaffold chain",
				}
				hookChain2 := &plugin.Hook{
					Name:        "test-hook-2",
					PlaceHookOn: "scaffold chain",
				}
				hookModule := &plugin.Hook{
					Name:        "test-hook-3",
					PlaceHookOn: "scaffold module",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Hooks: []*plugin.Hook{hookChain1, hookChain2, hookModule}}, nil)
				expectExecuteHook(t, p, plugin.Flags{{Name: "path"}}, hookChain1, hookChain2)
				expectExecuteHook(t, p, nil, hookModule)
			},
		},
		{
			name: "ok: duplicate hook names on same command",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				hooks := []*plugin.Hook{
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
					{
						Name:        "test-hook",
						PlaceHookOn: "ignite scaffold chain",
					},
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Hooks: hooks}, nil)
				expectExecuteHook(t, p, plugin.Flags{{Name: "path"}}, hooks...)
			},
		},
		{
			name: "ok: duplicate hook names on different commands",
			setup: func(t *testing.T, ctx context.Context, p *mocks.PluginInterface) {
				t.Helper()
				hookChain := &plugin.Hook{
					Name:        "test-hook",
					PlaceHookOn: "ignite scaffold chain",
				}
				hookModule := &plugin.Hook{
					Name:        "test-hook",
					PlaceHookOn: "ignite scaffold module",
				}
				p.EXPECT().
					Manifest(ctx).
					Return(&plugin.Manifest{Hooks: []*plugin.Hook{hookChain, hookModule}}, nil)
				expectExecuteHook(t, p, plugin.Flags{{Name: "path"}}, hookChain)
				expectExecuteHook(t, p, nil, hookModule)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			require := require.New(t)
			pi := mocks.NewPluginInterface(t)
			p := &plugin.Plugin{
				Plugin: pluginsconfig.Plugin{
					Path: "foo",
					With: pluginParams,
				},
				Interface: pi,
			}
			rootCmd := buildRootCmd(ctx)
			tt.setup(t, ctx, pi)

			_ = linkPlugins(ctx, rootCmd, []*plugin.Plugin{p})

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
	t.Helper()
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
