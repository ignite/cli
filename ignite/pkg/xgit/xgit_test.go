package xgit_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/randstr"
	"github.com/ignite/cli/ignite/pkg/xgit"
)

func TestInitAndCommit(t *testing.T) {
	tests := []struct {
		name                  string
		dirFunc               func(*testing.T) string
		expectDotGitFolder    bool
		expectedNumCommits    int
		expectedFilesInCommit []string
	}{
		{
			name: "dir is not inside an existing repo",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(path.Join(dir, "foo"), []byte("hello"), 0o755)
				require.NoError(t, err)
				return dir
			},
			expectDotGitFolder:    true,
			expectedNumCommits:    1,
			expectedFilesInCommit: []string{"foo"},
		},
		{
			name: "dir is inside an existing repo",
			// In this repo, there's no existing commit but a standalone uncommitted
			// foo file that shouldn't be included in the xgit.InitAndCommit's commit.
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "foo"), []byte("hello"), 0o755)
				require.NoError(t, err)
				dirInsideRepo := path.Join(dir, "bar")
				err = os.Mkdir(dirInsideRepo, 0o0755)
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dirInsideRepo, "baz"), []byte("hello"), 0o755)
				require.NoError(t, err)
				return dirInsideRepo
			},
			expectDotGitFolder:    false,
			expectedNumCommits:    1,
			expectedFilesInCommit: []string{"bar/baz"},
		},
		{
			name: "dir is an existing repo",
			dirFunc: func(t *testing.T) string {
				// In this repo, there's one existing commit, and an uncommitted baz file
				// that must be included in the xgit.InitAndCommit's commit.
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "foo"), []byte("hello"), 0o755)
				require.NoError(t, err)
				repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{})
				require.NoError(t, err)
				wt, err := repo.Worktree()
				require.NoError(t, err)
				_, err = wt.Add(".")
				require.NoError(t, err)
				_, err = wt.Commit("First commit", &git.CommitOptions{
					Author: &object.Signature{},
				})
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "bar"), []byte("hello"), 0o755)
				require.NoError(t, err)
				return dir
			},
			expectDotGitFolder:    true,
			expectedNumCommits:    2,
			expectedFilesInCommit: []string{"bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.dirFunc(t)

			err := xgit.InitAndCommit(dir)

			require.NoError(t, err)
			_, err = os.Stat(path.Join(dir, ".git"))
			require.Equal(t, tt.expectDotGitFolder, !os.IsNotExist(err))
			// Assert repository commits. For that we need to open the repo and
			// iterate over existing commits.
			repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
				DetectDotGit: true,
			})
			require.NoError(t, err)
			logs, err := repo.Log(&git.LogOptions{})
			require.NoError(t, err)
			var (
				numCommits int
				lastCommit *object.Commit
			)
			err = logs.ForEach(func(c *object.Commit) error {
				if numCommits == 0 {
					lastCommit = c
				}
				numCommits++
				return nil
			})
			require.NoError(t, err)
			require.Equal(t, tt.expectedNumCommits, numCommits)
			if assert.NotNil(t, lastCommit) {
				require.Equal(t, "Initialized with Ignite CLI", lastCommit.Message)
				require.WithinDuration(t, time.Now(), lastCommit.Committer.When, 10*time.Second)
				require.Equal(t, "Developer Experience team at Ignite", lastCommit.Author.Name)
				require.Equal(t, "hello@ignite.com", lastCommit.Author.Email)
				stats, err := lastCommit.Stats()
				require.NoError(t, err)
				var files []string
				for _, s := range stats {
					files = append(files, s.Name)
				}
				require.Equal(t, tt.expectedFilesInCommit, files)
			}
		})
	}
}

func TestAreChangesCommitted(t *testing.T) {
	tests := []struct {
		name           string
		dirFunc        func(*testing.T) string
		expectedResult bool
	}{
		{
			name: "dir is not a git repo",
			dirFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectedResult: true,
		},
		{
			name: "dir is a empty git repo",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				return dir
			},
			expectedResult: true,
		},
		{
			name: "dir is a dirty empty git repo",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "foo"), []byte("hello"), 0o755)
				require.NoError(t, err)
				return dir
			},
			expectedResult: false,
		},
		{
			name: "dir is a cleaned git repo",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "foo"), []byte("hello"), 0o755)
				require.NoError(t, err)
				repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{})
				require.NoError(t, err)
				wt, err := repo.Worktree()
				require.NoError(t, err)
				_, err = wt.Add(".")
				require.NoError(t, err)
				_, err = wt.Commit("First commit", &git.CommitOptions{
					Author: &object.Signature{},
				})
				require.NoError(t, err)
				return dir
			},
			expectedResult: true,
		},
		{
			name: "dir is a dirty git repo",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "foo"), []byte("hello"), 0o755)
				require.NoError(t, err)
				repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{})
				require.NoError(t, err)
				wt, err := repo.Worktree()
				require.NoError(t, err)
				_, err = wt.Add(".")
				require.NoError(t, err)
				_, err = wt.Commit("First commit", &git.CommitOptions{
					Author: &object.Signature{},
				})
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "bar"), []byte("hello"), 0o755)
				require.NoError(t, err)
				return dir
			},
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.dirFunc(t)

			res, err := xgit.AreChangesCommitted(dir)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestClone(t *testing.T) {
	// Create a folder with content
	notEmptyDir := t.TempDir()
	err := os.WriteFile(path.Join(notEmptyDir, ".foo"), []byte("hello"), 0o755)
	require.NoError(t, err)
	// Create a local git repo for all the test cases
	repoDir := t.TempDir()
	repo, err := git.PlainInit(repoDir, false)
	require.NoError(t, err)
	err = os.WriteFile(path.Join(repoDir, "foo"), []byte("hello"), 0o755)
	require.NoError(t, err)
	// Add a first commit
	w, err := repo.Worktree()
	require.NoError(t, err)
	_, err = w.Add(".")
	require.NoError(t, err)
	commit1, err := w.Commit("commit1", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "bob",
			Email: "bob@example.com",
			When:  time.Now(),
		},
	})
	// Add a branch on commit1
	require.NoError(t, err)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("my-branch"),
		Create: true,
	})
	require.NoError(t, err)
	// Back to master
	err = w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName("master")})
	require.NoError(t, err)
	// Add a tag on commit1
	_, err = repo.CreateTag("v1", commit1, &git.CreateTagOptions{
		Tagger:  &object.Signature{Name: "me"},
		Message: "v1",
	})
	require.NoError(t, err)
	// Add a second commit
	err = os.WriteFile(path.Join(repoDir, "bar"), []byte("hello"), 0o755)
	require.NoError(t, err)
	_, err = w.Add(".")
	require.NoError(t, err)
	commit2, err := w.Commit("commit2", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "bob",
			Email: "bob@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name          string
		dir           string
		urlRef        string
		expectedError string
		expectedRef   plumbing.Hash
	}{
		{
			name:          "fail: repo doesn't exist",
			dir:           t.TempDir(),
			urlRef:        "/tmp/not/exists",
			expectedError: "repository not found",
		},
		{
			name:          "fail: target dir isn't empty",
			dir:           notEmptyDir,
			urlRef:        repoDir,
			expectedError: fmt.Sprintf(`clone: target directory "%s" is not empty`, notEmptyDir),
		},
		{
			name:        "ok: target dir doesn't exists",
			dir:         "/tmp/not/exists/" + randstr.Runes(6),
			urlRef:      repoDir,
			expectedRef: commit2,
		},
		{
			name:        "ok: no ref",
			dir:         t.TempDir(),
			urlRef:      repoDir,
			expectedRef: commit2,
		},
		{
			name:        "ok: empty ref",
			dir:         t.TempDir(),
			urlRef:      repoDir + "@",
			expectedRef: commit2,
		},
		{
			name:        "ok: with tag ref",
			dir:         t.TempDir(),
			urlRef:      repoDir + "@v1",
			expectedRef: commit1,
		},
		{
			name:        "ok: with branch ref",
			dir:         t.TempDir(),
			urlRef:      repoDir + "@my-branch",
			expectedRef: commit1,
		},
		{
			name:        "ok: with commit1 hash ref",
			dir:         t.TempDir(),
			urlRef:      repoDir + "@" + commit1.String(),
			expectedRef: commit1,
		},
		{
			name:        "ok: with commit2 hash ref",
			dir:         t.TempDir(),
			urlRef:      repoDir + "@" + commit2.String(),
			expectedRef: commit2,
		},
		{
			name:          "fail: ref not found",
			dir:           t.TempDir(),
			urlRef:        repoDir + "@what",
			expectedError: "reference not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				files, _    = os.ReadDir(tt.dir)
				dirWasEmpty = len(files) == 0
			)

			err := xgit.Clone(context.Background(), tt.urlRef, tt.dir)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				if dirWasEmpty {
					// If it was empty, ensure target dir is still clean
					files, _ := os.ReadDir(tt.dir)
					require.Empty(t, files, "target dir should be empty in case of error")
				}
				return
			}
			require.NoError(t, err)
			_, err = os.Stat(tt.dir)
			require.False(t, os.IsNotExist(err), "dir %s should exist", tt.dir)
			repo, err := git.PlainOpen(tt.dir)
			require.NoError(t, err)
			h, err := repo.Head()
			require.NoError(t, err)
			require.Equal(t, tt.expectedRef, h.Hash())
		})
	}
}

func TestIsRepository(t *testing.T) {
	tests := []struct {
		name       string
		dirFunc    func(*testing.T) string
		shouldFail bool
		expected   bool
	}{
		{
			name: "path is a repository",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				_, err := git.PlainInit(dir, false)
				require.NoError(t, err)
				return dir
			},
			expected: true,
		},
		{
			name: "path is not a repository",
			dirFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expected: false,
		},
		{
			name: "repository error",
			dirFunc: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.Chmod(dir, 0)
				require.NoError(t, err)
				return dir
			},
			shouldFail: true,
			expected:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			exists, err := xgit.IsRepository(tt.dirFunc(t))

			// Assert
			require.Equal(t, tt.expected, exists)

			if tt.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
