package xgit

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

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
