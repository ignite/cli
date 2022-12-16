package plugin

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodule"
)

func TestNewPlugin(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name           string
		pluginCfg      pluginsconfig.Plugin
		expectedPlugin Plugin
	}{
		{
			name: "fail: empty path",
			expectedPlugin: Plugin{
				Error: errors.Errorf(`missing plugin property "path"`),
			},
		},
		{
			name:      "fail: local plugin doesnt exists",
			pluginCfg: pluginsconfig.Plugin{Path: "/xxx/yyy/plugin"},
			expectedPlugin: Plugin{
				Error: errors.Errorf(`local plugin path "/xxx/yyy/plugin" not found`),
			},
		},
		{
			name:      "fail: local plugin is not a dir",
			pluginCfg: pluginsconfig.Plugin{Path: path.Join(wd, "testdata/fakebin")},
			expectedPlugin: Plugin{
				Error: errors.Errorf(fmt.Sprintf("local plugin path %q is not a dir", path.Join(wd, "testdata/fakebin"))),
			},
		},
		{
			name:      "ok: local plugin",
			pluginCfg: pluginsconfig.Plugin{Path: path.Join(wd, "testdata")},
			expectedPlugin: Plugin{
				srcPath:    path.Join(wd, "testdata"),
				binaryName: "testdata",
			},
		},
		{
			name:      "fail: remote plugin with only domain",
			pluginCfg: pluginsconfig.Plugin{Path: "github.com"},
			expectedPlugin: Plugin{
				Error: errors.Errorf(`plugin path "github.com" is not a valid repository URL`),
			},
		},
		{
			name:      "fail: remote plugin with incomplete URL",
			pluginCfg: pluginsconfig.Plugin{Path: "github.com/ignite"},
			expectedPlugin: Plugin{
				Error: errors.Errorf(`plugin path "github.com/ignite" is not a valid repository URL`),
			},
		},
		{
			name:      "ok: remote plugin",
			pluginCfg: pluginsconfig.Plugin{Path: "github.com/ignite/plugin"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/ignite/plugin",
				cloneURL:   "https://github.com/ignite/plugin",
				cloneDir:   ".ignite/plugins/github.com/ignite/plugin",
				reference:  "",
				srcPath:    ".ignite/plugins/github.com/ignite/plugin",
				binaryName: "plugin",
			},
		},
		{
			name:      "ok: remote plugin with @ref",
			pluginCfg: pluginsconfig.Plugin{Path: "github.com/ignite/plugin@develop"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/ignite/plugin@develop",
				cloneURL:   "https://github.com/ignite/plugin",
				cloneDir:   ".ignite/plugins/github.com/ignite/plugin@develop",
				reference:  "develop",
				srcPath:    ".ignite/plugins/github.com/ignite/plugin@develop",
				binaryName: "plugin",
			},
		},
		{
			name:      "ok: remote plugin with subpath",
			pluginCfg: pluginsconfig.Plugin{Path: "github.com/ignite/plugin/plugin1"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/ignite/plugin",
				cloneURL:   "https://github.com/ignite/plugin",
				cloneDir:   ".ignite/plugins/github.com/ignite/plugin",
				reference:  "",
				srcPath:    ".ignite/plugins/github.com/ignite/plugin/plugin1",
				binaryName: "plugin1",
			},
		},
		{
			name:      "ok: remote plugin with subpath and @ref",
			pluginCfg: pluginsconfig.Plugin{Path: "github.com/ignite/plugin/plugin1@develop"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/ignite/plugin@develop",
				cloneURL:   "https://github.com/ignite/plugin",
				cloneDir:   ".ignite/plugins/github.com/ignite/plugin@develop",
				reference:  "develop",
				srcPath:    ".ignite/plugins/github.com/ignite/plugin@develop/plugin1",
				binaryName: "plugin1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectedPlugin.Plugin = tt.pluginCfg

			p := newPlugin(".ignite/plugins", tt.pluginCfg)

			assertPlugin(t, tt.expectedPlugin, *p)
		})
	}
}

func TestPluginLoad(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// scaffoldPlugin runs Scaffold and updates the go.mod so it uses the
	// current ignite/cli sources.
	scaffoldPlugin := func(t *testing.T, dir, name string) string {
		require := require.New(t)
		path, err := Scaffold(dir, name)
		require.NoError(err)
		// We want the scaffolded plugin to use the current version of ignite/cli,
		// for that we need to update the plugin go.mod and add a replace to target
		// current ignite/cli
		gomod, err := gomodule.ParseAt(path)
		require.NoError(err)
		// use GOMOD env to get current directory module path
		modpath, err := gocmd.Env(gocmd.EnvGOMOD)
		require.NoError(err)
		modpath = filepath.Dir(modpath)
		gomod.AddReplace("github.com/ignite/cli", "", modpath, "")
		// Save go.mod
		data, err := gomod.Format()
		require.NoError(err)
		err = os.WriteFile(filepath.Join(path, "go.mod"), data, 0o644)
		require.NoError(err)
		return path
	}

	// Helper to make a local git repository with gofile committed.
	// Returns the repo directory and the git.Repository
	makeGitRepo := func(t *testing.T, name string) (string, *git.Repository) {
		require := require.New(t)
		repoDir := t.TempDir()
		scaffoldPlugin(t, repoDir, "github.com/ignite/"+name)
		require.NoError(err)
		repo, err := git.PlainInit(repoDir, false)
		require.NoError(err)
		w, err := repo.Worktree()
		require.NoError(err)
		_, err = w.Add(".")
		require.NoError(err)
		_, err = w.Commit("msg", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "bob",
				Email: "bob@example.com",
				When:  time.Now(),
			},
		})
		require.NoError(err)
		return repoDir, repo
	}
	tests := []struct {
		name          string
		buildPlugin   func(t *testing.T) Plugin
		expectedError string
	}{
		{
			name: "fail: plugin is already in error",
			buildPlugin: func(t *testing.T) Plugin {
				return Plugin{
					Error: errors.New("oups"),
				}
			},
			expectedError: `oups`,
		},
		{
			name: "fail: no go files in srcPath",
			buildPlugin: func(t *testing.T) Plugin {
				return Plugin{
					srcPath:    path.Join(wd, "testdata"),
					binaryName: "testdata",
				}
			},
			expectedError: `no packages to build`,
		},
		{
			name: "ok: from local",
			buildPlugin: func(t *testing.T) Plugin {
				return Plugin{
					srcPath:    scaffoldPlugin(t, t.TempDir(), "github.com/foo/bar"),
					binaryName: "bar",
				}
			},
		},
		{
			name: "ok: from git repo",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, _ := makeGitRepo(t, "remote")
				cloneDir := t.TempDir()

				return Plugin{
					cloneURL:   repoDir,
					cloneDir:   cloneDir,
					srcPath:    path.Join(cloneDir, "remote"),
					binaryName: "remote",
				}
			},
		},
		{
			name: "fail: git repo doesnt exists",
			buildPlugin: func(t *testing.T) Plugin {
				cloneDir := t.TempDir()

				return Plugin{
					repoPath: "/xxxx/yyyy",
					cloneURL: "/xxxx/yyyy",
					cloneDir: cloneDir,
					srcPath:  path.Join(cloneDir, "plugin"),
				}
			},
			expectedError: `cloning "/xxxx/yyyy": repository not found`,
		},
		{
			name: "ok: from git repo with tag",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, repo := makeGitRepo(t, "remote-tag")
				h, err := repo.Head()
				require.NoError(t, err)
				_, err = repo.CreateTag("v1", h.Hash(), &git.CreateTagOptions{
					Tagger:  &object.Signature{Name: "me"},
					Message: "v1",
				})
				require.NoError(t, err)

				cloneDir := t.TempDir()

				return Plugin{
					cloneURL:   repoDir,
					reference:  "v1",
					cloneDir:   cloneDir,
					srcPath:    path.Join(cloneDir, "remote-tag"),
					binaryName: "remote-tag",
				}
			},
		},
		{
			name: "ok: from git repo with branch",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, repo := makeGitRepo(t, "remote-branch")
				w, err := repo.Worktree()
				require.NoError(t, err)
				err = w.Checkout(&git.CheckoutOptions{
					Branch: plumbing.NewBranchReferenceName("branch1"),
					Create: true,
				})
				require.NoError(t, err)

				cloneDir := t.TempDir()

				return Plugin{
					cloneURL:   repoDir,
					reference:  "branch1",
					cloneDir:   cloneDir,
					srcPath:    path.Join(cloneDir, "remote-branch"),
					binaryName: "remote-branch",
				}
			},
		},
		{
			name: "ok: from git repo with hash",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, repo := makeGitRepo(t, "remote-hash")
				h, err := repo.Head()
				require.NoError(t, err)

				cloneDir := t.TempDir()

				return Plugin{
					cloneURL:   repoDir,
					reference:  h.Hash().String(),
					cloneDir:   cloneDir,
					srcPath:    path.Join(cloneDir, "remote-hash"),
					binaryName: "remote-hash",
				}
			},
		},
		{
			name: "fail: git ref not found",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, _ := makeGitRepo(t, "remote-no-ref")

				cloneDir := t.TempDir()

				return Plugin{
					cloneURL:   repoDir,
					reference:  "doesnt_exists",
					cloneDir:   cloneDir,
					srcPath:    path.Join(cloneDir, "remote-no-ref"),
					binaryName: "remote-no-ref",
				}
			},
			expectedError: `cloning ".*": reference not found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			p := tt.buildPlugin(t)
			defer p.KillClient()

			p.load(context.Background())

			if tt.expectedError != "" {
				require.Error(p.Error, "expected error %q", tt.expectedError)
				require.Regexp(tt.expectedError, p.Error.Error())
				return
			}
			require.NoError(p.Error)
			require.NotNil(p.Interface)
			manifest, err := p.Interface.Manifest()
			require.NoError(err)
			assert.Equal(p.binaryName, manifest.Name)
			assert.NoError(p.Interface.Execute(ExecutedCommand{}))
			assert.NoError(p.Interface.ExecuteHookPre(ExecutedHook{}))
			assert.NoError(p.Interface.ExecuteHookPost(ExecutedHook{}))
			assert.NoError(p.Interface.ExecuteHookCleanUp(ExecutedHook{}))
		})
	}
}

func TestPluginClean(t *testing.T) {
	tests := []struct {
		name         string
		plugin       *Plugin
		expectRemove bool
	}{
		{
			name: "dont clean local plugin",
			plugin: &Plugin{
				Plugin: pluginsconfig.Plugin{Path: "/local"},
			},
		},
		{
			name:   "dont clean plugin with errors",
			plugin: &Plugin{Error: errors.New("oups")},
		},
		{
			name: "ok",
			plugin: &Plugin{
				cloneURL: "https://github.com/ignite/plugin",
			},
			expectRemove: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := os.MkdirTemp("", "cloneDir")
			require.NoError(t, err)
			tt.plugin.cloneDir = tmp

			err = tt.plugin.clean()

			require.NoError(t, err)
			if tt.expectRemove {
				_, err := os.Stat(tmp)
				assert.True(t, os.IsNotExist(err), "cloneDir not removed")
			}
		})
	}
}

func assertPlugin(t *testing.T, want, have Plugin) {
	if want.Error != nil {
		require.Error(t, have.Error)
		assert.Regexp(t, want.Error.Error(), have.Error.Error())
	} else {
		require.NoError(t, have.Error)
	}
	// Errors aren't comparable with assert.Equal, because of the different stacks
	want.Error = nil
	have.Error = nil
	assert.Equal(t, want, have)
}
