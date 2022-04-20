package giturl

import (
	"errors"
	"net/url"
	"strings"
)

// GitURL represents a Git url.
type GitURL struct {
	// Host is a Git host.
	Host string

	// User is a user or an org.
	User string

	// Repo is a repo name.
	Repo string
}

// UserAndRepo returns the combined string representation of user and repo.
func (g GitURL) UserAndRepo() string {
	return strings.Join([]string{g.User, g.Repo}, "/")
}

// Parse parses a Git url u.
func Parse(u string) (GitURL, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return GitURL{}, err
	}

	sp := strings.Split(ur.Path, "/")
	if len(sp) < 3 {
		return GitURL{}, errors.New("invalid url")
	}

	return GitURL{
		Host: ur.Host,
		User: sp[1],
		Repo: sp[2],
	}, nil
}
