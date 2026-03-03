package repoversion

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

func commitFile(t *testing.T, repo *git.Repository, dir, name, content string, when time.Time) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	wt, err := repo.Worktree()
	require.NoError(t, err)
	_, err = wt.Add(name)
	require.NoError(t, err)

	hash, err := wt.Commit(name, &git.CommitOptions{
		Author:    &object.Signature{Name: "test", Email: "test@example.com", When: when},
		Committer: &object.Signature{Name: "test", Email: "test@example.com", When: when},
	})
	require.NoError(t, err)

	return hash.String()
}

func TestDetermineWithoutTag(t *testing.T) {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)

	headHash := commitFile(t, repo, dir, "a.txt", "a", time.Unix(100, 0))

	v, err := Determine(dir)
	require.NoError(t, err)
	require.Empty(t, v.Tag)
	require.Equal(t, headHash, v.Hash)
}

func TestDetermineWithTagOnHead(t *testing.T) {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)

	headHash := commitFile(t, repo, dir, "a.txt", "a", time.Unix(100, 0))
	head, err := repo.Head()
	require.NoError(t, err)
	_, err = repo.CreateTag("v1.2.3", head.Hash(), nil)
	require.NoError(t, err)

	v, err := Determine(dir)
	require.NoError(t, err)
	require.Equal(t, "1.2.3", v.Tag)
	require.Equal(t, headHash, v.Hash)
}

func TestDetermineWithOlderTagUsesSuffix(t *testing.T) {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)

	_ = commitFile(t, repo, dir, "a.txt", "a", time.Unix(100, 0))
	head, err := repo.Head()
	require.NoError(t, err)
	_, err = repo.CreateTag("v1.0.0", head.Hash(), nil)
	require.NoError(t, err)

	headHash := commitFile(t, repo, dir, "b.txt", "b", time.Unix(200, 0))

	v, err := Determine(dir)
	require.NoError(t, err)
	require.Equal(t, "1.0.0-"+headHash[:8], v.Tag)
	require.Equal(t, headHash, v.Hash)
}
