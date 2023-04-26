package xgit

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	commitMsg       = "Initialized with Ignite CLI"
	defaultOpenOpts = git.PlainOpenOptions{DetectDotGit: true}
	devXAuthor      = &object.Signature{
		Name:  "Developer Experience team at Ignite",
		Email: "hello@ignite.com",
		When:  time.Now(),
	}
)

// InitAndCommit creates a git repo in path if path isn't already inside a git
// repository, then commits path content.
func InitAndCommit(path string) error {
	repo, err := git.PlainOpenWithOptions(path, &defaultOpenOpts)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryNotExists) {
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

// AreChangesCommitted returns true if dir is a clean git repository with no
// pending changes. It returns also true if dir is NOT a git repository.
func AreChangesCommitted(dir string) (bool, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return false, err
	}

	repository, err := git.PlainOpen(dir)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
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

// Clone clones a git repository represented by urlRef, into dir.
// urlRef is the URL of the repository, with an optional ref, suffixed to the
// URL with a `@`. Ref can be a tag, a branch or a hash.
// Valid examples of urlRef: github.com/org/repo, github.com/org/repo@v1,
// github.com/org/repo@develop, github.com/org/repo@ab88cdf.
func Clone(ctx context.Context, urlRef, dir string) error {
	// Ensure dir is empty if it exists (if it doesn't exist, the call to
	// git.PlainCloneContext below will create it).
	files, _ := os.ReadDir(dir)
	if len(files) > 0 {
		return fmt.Errorf("clone: target directory %q is not empty", dir)
	}
	// Split urlRef
	var (
		parts = strings.Split(urlRef, "@")
		url   = parts[0]
		ref   string
	)
	if len(parts) > 1 {
		ref = parts[1]
	}
	// First clone the repo
	repo, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return err
	}
	if ref == "" {
		// if ref is not provided, job is done
		return nil
	}
	// Reference provided, try to resolve
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	var h *plumbing.Hash
	for _, ref := range []string{ref, "origin/" + ref} {
		h, err = repo.ResolveRevision(plumbing.Revision(ref))
		if err == nil {
			break
		}
	}
	if err != nil {
		// Ref not found, clean up dir and return error
		os.RemoveAll(dir)
		return err
	}
	return wt.Checkout(&git.CheckoutOptions{
		Hash: *h,
	})
}

// IsRepository checks if a path contains a Git repository.
func IsRepository(path string) (bool, error) {
	if _, err := git.PlainOpenWithOptions(path, &defaultOpenOpts); err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
