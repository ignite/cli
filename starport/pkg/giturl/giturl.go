package giturl

import (
	"net/url"
	"strings"
)

// UserAndRepo returns only the user and repo portion of the git url u.
func UserAndRepo(u string) string {
	ur, err := url.Parse(u)
	if err != nil {
		return u
	}

	sp := strings.Split(ur.Path, "/")
	if len(sp) < 3 {
		return u
	}

	return strings.Join(sp[1:3], "/")
}
