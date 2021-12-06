package plugin

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/chainconfig"
)

func TestNewPlugin(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name           string
		pluginCfg      chainconfig.Plugin
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
			pluginCfg: chainconfig.Plugin{Path: "/xxx/yyy/plugin"},
			expectedPlugin: Plugin{
				Error: errors.Errorf(`local plugin path "/xxx/yyy/plugin" not found`),
			},
		},
		{
			name:      "fail: local plugin is not a dir",
			pluginCfg: chainconfig.Plugin{Path: path.Join(wd, "testdata/fakebin")},
			expectedPlugin: Plugin{
				Error: errors.Errorf(fmt.Sprintf("local plugin path %q is not a dir", path.Join(wd, "testdata/fakebin"))),
			},
		},
		{
			name:      "ok: local plugin",
			pluginCfg: chainconfig.Plugin{Path: path.Join(wd, "testdata")},
			expectedPlugin: Plugin{
				srcPath:    path.Join(wd, "testdata"),
				binaryName: "testdata",
			},
		},
		{
			name:      "fail: remote plugin with only domain",
			pluginCfg: chainconfig.Plugin{Path: "github.com"},
			expectedPlugin: Plugin{
				Error: errors.Errorf(`plugin path "github.com" is not a valid repository URL`),
			},
		},
		{
			name:      "fail: remote plugin with incomplete URL",
			pluginCfg: chainconfig.Plugin{Path: "github.com/starport"},
			expectedPlugin: Plugin{
				Error: errors.Errorf(`plugin path "github.com/starport" is not a valid repository URL`),
			},
		},
		{
			name:      "ok: remote plugin",
			pluginCfg: chainconfig.Plugin{Path: "github.com/starport/plugin"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/starport/plugin",
				cloneURL:   "https://github.com/starport/plugin",
				cloneDir:   ".starport/plugins/github.com/starport/plugin",
				reference:  "",
				srcPath:    ".starport/plugins/github.com/starport/plugin",
				binaryName: "plugin",
			},
		},
		{
			name:      "ok: remote plugin with @ref",
			pluginCfg: chainconfig.Plugin{Path: "github.com/starport/plugin@develop"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/starport/plugin@develop",
				cloneURL:   "https://github.com/starport/plugin",
				cloneDir:   ".starport/plugins/github.com/starport/plugin@develop",
				reference:  "develop",
				srcPath:    ".starport/plugins/github.com/starport/plugin@develop",
				binaryName: "plugin",
			},
		},
		{
			name:      "ok: remote plugin with subpath",
			pluginCfg: chainconfig.Plugin{Path: "github.com/starport/plugin/plugin1"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/starport/plugin",
				cloneURL:   "https://github.com/starport/plugin",
				cloneDir:   ".starport/plugins/github.com/starport/plugin",
				reference:  "",
				srcPath:    ".starport/plugins/github.com/starport/plugin/plugin1",
				binaryName: "plugin1",
			},
		},
		{
			name:      "ok: remote plugin with subpath and @ref",
			pluginCfg: chainconfig.Plugin{Path: "github.com/starport/plugin/plugin1@develop"},
			expectedPlugin: Plugin{
				repoPath:   "github.com/starport/plugin@develop",
				cloneURL:   "https://github.com/starport/plugin",
				cloneDir:   ".starport/plugins/github.com/starport/plugin@develop",
				reference:  "develop",
				srcPath:    ".starport/plugins/github.com/starport/plugin@develop/plugin1",
				binaryName: "plugin1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectedPlugin.Plugin = tt.pluginCfg

			p := newPlugin(".starport/plugins", tt.pluginCfg)

			assertPlugin(t, tt.expectedPlugin, *p)
		})
	}
}

func TestPluginLoad(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// Use a common temp dir for all the cases to facilitate cleaning.
	tmpDir := path.Join(os.TempDir(), "starport_"+t.Name())
	err = os.MkdirAll(tmpDir, 0700)
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	mkdirTmp := func(t *testing.T, dir string) string {
		tmp, err := os.MkdirTemp(tmpDir, dir)
		require.NoError(t, err)
		return tmp
	}

	// Helper to make a local git repository with gofile committed.
	// Returns the repo directory and the git.Repository
	makeGitRepo := func(t *testing.T, name string) (string, *git.Repository) {
		require := require.New(t)
		repoDir := mkdirTmp(t, "plugin_repo")
		err = Scaffold(repoDir, "github.com/starport/"+name)
		require.NoError(err)
		repo, err := git.PlainInit(repoDir, false)
		require.NoError(err)
		w, err := repo.Worktree()
		require.NoError(err)
		w.Add(".")
		w.Commit("msg", &git.CommitOptions{})
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
			expectedError: `no Go files in ` + wd + `/testdata`,
		},
		{
			name: "ok: from local",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir := mkdirTmp(t, "plugin_local")
				err = Scaffold(repoDir, "github.com/foo/bar")
				require.NoError(t, err)
				return Plugin{
					srcPath:    path.Join(repoDir, "bar"),
					binaryName: "bar",
				}
			},
		},
		{
			name: "ok: from git repo",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, _ := makeGitRepo(t, "remote")
				cloneDir := mkdirTmp(t, "clone_dir")

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
				cloneDir := mkdirTmp(t, "clone_dir")

				return Plugin{
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

				cloneDir := mkdirTmp(t, "clone_dir")

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

				cloneDir := mkdirTmp(t, "clone_dir")

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
			name: "fail: git ref not found",
			buildPlugin: func(t *testing.T) Plugin {
				repoDir, _ := makeGitRepo(t, "remote-no-ref")

				cloneDir := mkdirTmp(t, "clone_dir")

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
			p := tt.buildPlugin(t)

			p.load()

			if tt.expectedError != "" {
				require.Error(t, p.Error, "expected error %q", tt.expectedError)
				require.Regexp(t, tt.expectedError, p.Error.Error())
				return
			}
			require.NoError(t, p.Error)
			require.NotNil(t, p.Interface)
			assert.Equal(t, p.binaryName, p.Interface.Commands()[0].Use)
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
				Plugin: chainconfig.Plugin{Path: "/local"},
			},
		},
		{
			name:   "dont clean plugin with errors",
			plugin: &Plugin{Error: errors.New("oups")},
		},
		{
			name: "ok",
			plugin: &Plugin{
				cloneURL: "https://github.com/starport/plugin",
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
