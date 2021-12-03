package plugin

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/tendermint/starport/starport/chainconfig"
)

//
// Builder handles download build process for new plugin.
//

// TODO: How to divide GIT module?

const (
	pluginDir = "plugin"
)

// Builder provides interfaces to build plugins.
type Builder interface {
	Build(config chainconfig.Plugin) error
}

type builder struct {
	pluginSpec *starportplugin
}

func (b *builder) Build(config chainconfig.Plugin) error {
	// TODO:
	err := b.download(config.RepositoryURL)
	if err != nil {
		return err
	}

	// who can fill-up the config.Name?
	err = b.build(config.Name)
	return err
}

func (b *builder) download(url string) error {
	_ = pluginDir

	_, err := git.PlainClone(pluginDir, false, &git.CloneOptions{
		URL: url,
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println(err)
	}

	repository, err := git.PlainOpen(pluginDir)
	if err != nil {
		fmt.Println(err)
	}

	worktree, err := repository.Worktree()
		if err != nil {
			panic(err)
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		fmt.Println(err)
	}

	reference, err := repository.Head()
	if err != nil {
		panic(err)
	}
	commit, err := repository.CommitObject(reference.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Println(commit)
	// TODO:
	// Create `pluginDir` if not exist.
	// Clone repo from url on `pluginDir`.
	return nil
}

func (b *builder) build(name string) error {
	// TODO:
	// Build plugin.
	// go build -buildmode=plugin
	return nil
}

// NewBuilder creates new plugin builder.
// TODO: Need parameters?
func NewBuilder() (Builder, error) {
	// TODO:
	return &builder{}, nil
}
