package xgit

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	commitMsg  = "Initialized with Ignite CLI"
	devXAuthor = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

// InitAndCommit creates a git repo in path if path isn't already inside a git
// repository, then commits path content.
func InitAndCommit(path string) error {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return fmt.Errorf("open git repo %s: %w", path, err)
		}
		// not a git repo, creates a new one
		repo, err = git.PlainInit(path, false)
		if err != nil {
			return fmt.Errorf("init git repo %s: %w", path, err)
		}
	}
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("worktree %s: %w", path, err)
	}
	// wt.Add(path) takes only relative path, we need to turn path relative to
	// repo path.
	repoPath := wt.Filesystem.Root()
	path, err = filepath.Rel(repoPath, path)
	if err != nil {
		return fmt.Errorf("find relative path %s %s: %w", repoPath, path, err)
	}
	if _, err := wt.Add(path); err != nil {
		return fmt.Errorf("git add %s: %w", path, err)
	}
	_, err = wt.Commit(commitMsg, &git.CommitOptions{
		All:    true,
		Author: devXAuthor,
	})
	if err != nil {
		return fmt.Errorf("git commit %s: %w", path, err)
	}
	return nil
}

func AreChangesCommitted(appPath string) (bool, error) {
	appPath, err := filepath.Abs(appPath)
	if err != nil {
		return false, err
	}

	repository, err := git.PlainOpen(appPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return true, nil
		}
		return false, err
	}

	w, err := repository.Worktree()
	if err != nil {
		return false, err
	}

	ws, err := w.Status()
	if err != nil {
		return false, err
	}
	return ws.IsClean(), nil
}
