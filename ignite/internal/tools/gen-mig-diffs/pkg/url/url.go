package url

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// URL represents a Git URL in any supported protocol.
type URL struct {
	// Protocol is the protocol of the endpoint (e.g. ssh, https).
	Protocol string
	// Host is the host.
	Host string
	// Path is the repository path.
	Path string
}

var (
	isSchemeRegExp   = regexp.MustCompile(`^[^:]+://`)
	scpLikeUrlRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\].*)$`)
)

// New creates a new URL object.
func New(endpoint string) (URL, error) {
	if e, ok := parseSCPLike(endpoint); ok {
		return e, nil
	}
	return parseURL(endpoint)
}

func (u URL) Compare(cp URL) error {
	switch {
	case u.Host != cp.Host:
		return errors.Errorf("host mismatch for %s != %s", u.Host, cp.Host)
	case u.Path != cp.Path:
		return errors.Errorf("path mismatch for %s != %s", u.Path, cp.Path)
	default:
		return nil
	}
}

func (u URL) String() string {
	if u.Protocol == "ssh" {
		return fmt.Sprintf("git@%s:%s.git", u.Host, u.Path)
	}
	return fmt.Sprintf("%s://%s/%s.git", u.Protocol, u.Host, u.Path)
}

// parseSCPLike returns an URL object from SCP git URL.
func parseSCPLike(endpoint string) (URL, bool) {
	if matchesScheme(endpoint) || !matchesScpLike(endpoint) {
		return URL{}, false
	}

	_, host, _, path := findScpLikeComponents(endpoint)

	return URL{
		Protocol: "ssh",
		Host:     host,
		Path:     strings.TrimSuffix(path, ".git"),
	}, true
}

// parseURL returns an URL object from an endpoint.
func parseURL(endpoint string) (URL, error) {
	if !matchesScheme(endpoint) || matchesScpLike(endpoint) {
		return URL{}, errors.Errorf("invalid endpoint scheme: %s", endpoint)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return URL{}, err
	}

	if !u.IsAbs() {
		return URL{}, errors.Errorf("invalid endpoint URL: %s", endpoint)
	}

	return URL{
		Protocol: u.Scheme,
		Host:     u.Hostname(),
		Path:     getPath(u),
	}, nil
}

// matchesScheme returns true if the given string matches a URL-like format scheme.
func matchesScheme(url string) bool {
	return isSchemeRegExp.MatchString(url)
}

// matchesScpLike returns true if the given string matches an SCP-like format scheme.
func matchesScpLike(url string) bool {
	return scpLikeUrlRegExp.MatchString(url)
}

// findScpLikeComponents returns the user, host, port and path of the given SCP-like URL.
func findScpLikeComponents(url string) (user, host, port, path string) {
	m := scpLikeUrlRegExp.FindStringSubmatch(url)
	user = m[1]
	host = m[2]
	port = m[3]
	path = m[4]
	return m[1], m[2], m[3], m[4]
}

// getPath returns the path from an *url.URL.
func getPath(u *url.URL) string {
	res := u.Path
	if u.RawQuery != "" {
		res += "?" + u.RawQuery
	}
	if u.Fragment != "" {
		res += "#" + u.Fragment
	}

	res = strings.Trim(res, "/")
	res = strings.TrimSuffix(res, ".git")
	return strings.Split(res, ":")[0]
}
