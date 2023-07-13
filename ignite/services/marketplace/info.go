package marketplace

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func InfoPlugin(ctx context.Context, client *Client, repo string) error {
	result, err := client.RepoQuery(ctx, &github.SearchOptions{}, Query{Qualifier: "repo", Value: repo})
	if err != nil {
		return err
	}

	fmt.Printf(`
%s
  - Name: %s
  - Owner: %s
  - Description: %s
  - Stars: %d
  - URL: %s
`, repo, result[0].GetName(), result[0].GetOwner().GetLogin(), result[0].GetDescription(), result[0].GetStargazersCount(), result[0].GetHTMLURL())

	return nil
}
