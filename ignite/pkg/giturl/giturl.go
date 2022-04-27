package giturl

import (
	"errors"
	"net/url"
	"strings"
)

var (
	// ErrInvalidURL is returned when a URL doesn't follow the format domain.com/username/reponame.
	ErrInvalidURL = errors.New("invalid url")
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

// Parse a Git URL with format "domain.com/username/reponame".
// The URL scheme is optional.
func Parse(gitURL string) (GitURL, error) {
	u, err := url.Parse(gitURL)
	if err != nil {
		return GitURL{}, err
	}

	p := strings.Split(strings.TrimLeft(u.Path, "/"), "/")
	g := GitURL{}

	if u.Host != "" {
		if len(p) < 2 {
			return GitURL{}, ErrInvalidURL
		}

		g.Host = u.Host
		g.User = p[0]
		g.Repo = p[1]
	} else {
		// URL parses the domain name as part of the path when the git URL has no scheme
		// so the first path element is assumed to be a domain name when it contains a "."
		// TODO: should we use a regexp or the simplistic check is enough?
		if len(p) < 3 || !strings.Contains(p[0], ".") {
			return GitURL{}, ErrInvalidURL
		}

		g.Host = p[0]
		g.User = p[1]
		g.Repo = p[2]
	}

	return g, nil
}
