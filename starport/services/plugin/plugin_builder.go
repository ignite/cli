package plugin

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// for clone
	"github.com/go-git/go-git/v5"

	// for build
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"

	"github.com/tendermint/starport/starport/chainconfig"
)

//
// Builder handles download build process for new plugin.
//

const (
	pluginDir = "plugins"
)

// Builder provides interfaces to build plugins.
type Builder interface {
	Build(config chainconfig.Plugin) error
}

type builder struct {
	pluginSpec *starportplugin
}

func (b *builder) Build(config chainconfig.Plugin) error {
	path, err := b.download(config.Name, config.RepositoryURL)
	if err != nil {
		log.Println("Download ", err)
		return err
	}

	err = b.build(path, config.Name)
	return err
}

func (b *builder) download(pluginName, url string) (string, error) {
	log.Printf("Clone %s\n", url)

	starportHome, err := chainconfig.ConfigDirPath()
	if err != nil {
		log.Println(err)
		return "", err
	}

	pathTokens := strings.Split(url, "/")
	repoName := pathTokens[len(pathTokens)-1]

	repoPath := fmt.Sprintf("%s/%s/%s", starportHome, pluginDir, repoName)

	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Depth:    1,
	})

	if err != nil {
		log.Println("clone", err)
		return "", err
	}

	/*
		repository, err := git.PlainOpen(pluginHome)
		if err != nil {
			log.Println("open", err)
			return err
		}

		worktree, err := repository.Worktree()
		if err != nil {
			log.Println("open", err)
			return err
		}

		err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			log.Println("pull", err)
			return err
		}

		reference, err := repository.Head()
		if err != nil {
			log.Println("Head ", err)
			return err
		}

		_, err = repository.CommitObject(reference.Hash())
		if err != nil {
			log.Println("hash ", err)
			return err
		}
	*/

	_ = b.pluginSpec

	return fmt.Sprintf("%s/%s", repoPath, pluginName), nil
}

func (b *builder) build(path, name string) error {
	log.Printf("Build plugin %s\n", name)

	errb := &bytes.Buffer{}

	err := cmdrunner.
		New(cmdrunner.DefaultStderr(errb), cmdrunner.DefaultWorkdir(path)).
		Run(context.Background(), step.New(
			step.Exec(
				"go",
				"build",
				"-buildmode=plugin",
				"-o",
				fmt.Sprintf("%s.so", name),
			),
		),
		)

	// TODO: Check mandatory functions.

	return err
}

// NewBuilder creates new plugin builder.
func NewBuilder() (Builder, error) {
	return &builder{}, nil
}
