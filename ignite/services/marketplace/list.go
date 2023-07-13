package marketplace

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"

	"github.com/ignite/cli/ignite/pkg/cliui/icons"
)

func ListPlugins(ctx context.Context, client *Client, opts *github.SearchOptions) error {
	result, err := client.RepoQuery(
		ctx,
		opts,
		Query{Qualifier: "topic", Value: "ignite-plugin"},
		// Query{Qualifier: "language", Value: "go"},
		// Query{Qualifier: "license", Value: "MIT"},
		// Query{Qualifier: "stars", Value: ">5"}, // TODO: Uncommend this line
	)
	if err != nil {
		return err
	}

	homeDir := os.Getenv("HOME")

	owners, err := allDirectoriesInDir(homeDir + "/.ignite/plugins")
	if err != nil {
		return err
	}

	plugins := make(map[string]bool)
	for owner := range owners {
		pl, err := allDirectoriesInDir(homeDir + "/.ignite/plugins/" + owner)
		if err != nil {
			return err
		}

		for p := range pl {
			plugins[owner+"/"+p] = true
		}
	}

	for _, repo := range result {
		fullName := *repo.Owner.Login + "/" + *repo.Name

		if _, ok := plugins[fullName]; ok {
			fmt.Printf("\n%s %s: ", icons.OK, fullName)
		} else {
			fmt.Printf("\n%s %s: ", icons.NotOK, fullName)
		}

		if repo.Description != nil {
			fmt.Printf("%s", *repo.Description)
		}
	}

	return nil
}

func allDirectoriesInDir(dirpath string) (map[string]bool, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	dirs := make(map[string]bool)
	for _, p := range fileInfos {
		if p.IsDir() {
			dirs[p.Name()] = true
		}
	}

	return dirs, nil
}
