package xgit

import (
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

var (
	commitMsg  = "Initialized with Ignite CLI"
	devXAuthor = &object.Signature{
		Name:  "Developer Experience team at Tendermint",
		Email: "hello@tendermint.com",
		When:  time.Now(),
	}
)

func InitAndCommit(path string) error {
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return errors.WithStack(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := wt.Add("."); err != nil {
		return errors.WithStack(err)
	}
	_, err = wt.Commit(commitMsg, &git.CommitOptions{
		All:    true,
		Author: devXAuthor,
	})
	return errors.WithStack(err)
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
