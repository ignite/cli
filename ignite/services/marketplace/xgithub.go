package marketplace

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	igniteTopic = "ignite-plugin"
)

type (
	Client struct {
		ts     oauth2.TokenSource
		tc     *http.Client
		client *github.Client
	}

	Query struct {
		Qualifier string
		Value     string
	}
)

func (c *Client) RepoQuery(ctx context.Context, opts *github.SearchOptions, queries ...Query) ([]github.Repository, error) {
	query, err := CreateQuery(queries)
	if err != nil {
		return nil, err
	}

	result, _, err := c.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	return result.Repositories, nil
}

func NewClient(ctx context.Context, accessToken string) *Client {
	var (
		ts = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc     = oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	)

	return &Client{
		ts:     ts,
		tc:     tc,
		client: client,
	}
}

func CreateQuery(queries []Query) (string, error) {
	var (
		qMap = make(map[string]string)
		q    = ""
	)

	for _, query := range queries {
		if _, ok := qMap[query.Qualifier]; ok {
			return "", fmt.Errorf("duplicate qualifier: %s", query.Qualifier)
		}

		qMap[query.Qualifier] = query.Value
		q += fmt.Sprintf("%s:%s ", query.Qualifier, query.Value)
	}

	return q, nil
}
