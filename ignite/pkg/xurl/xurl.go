package xurl

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	schemeTCP   = "tcp"
	schemeHTTP  = "http"
	schemeHTTPS = "https"
	schemeWS    = "ws"
)

// TCP unsures that s url contains TCP protocol identifier.
func TCP(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	if u.Scheme != schemeTCP {
		u.Scheme = schemeTCP
	}
	return u.String(), nil
}

// HTTP unsures that s url contains HTTP protocol identifier.
func HTTP(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	if u.Scheme != schemeHTTP {
		u.Scheme = schemeHTTP
	}
	return u.String(), nil
}

// HTTPS unsures that s url contains HTTPS protocol identifier.
func HTTPS(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	if u.Scheme != schemeHTTPS {
		u.Scheme = schemeHTTPS
	}
	return u.String(), nil
}

// WS unsures that s url contains WS protocol identifier.
func WS(s string) (string, error) {
	u, err := parseURL(s)
	if err != nil {
		return "", err
	}

	if u.Scheme != schemeWS {
		u.Scheme = schemeWS
	}
	return u.String(), nil
}

// HTTPEnsurePort ensures that url has a port number suits with the connection type.
func HTTPEnsurePort(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.Port() != "" {
		return s
	}

	port := "80"

	if u.Scheme == schemeHTTPS {
		port = "443"
	}

	u.Host = fmt.Sprintf("%s:%s", u.Hostname(), port)

	return u.String()
}

// CleanPath cleans path from the url.
func CleanPath(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	u.Path = ""

	return u.String()
}

// Address unsures that address contains localhost as host if non specified.
func Address(address string) string {
	if strings.HasPrefix(address, ":") {
		return "localhost" + address
	}
	return address
}

// IsLocalPath checks if given address is a local fs path or a URL.
func IsLocalPath(address string) bool {
	for _, pattern := range []string{
		"http://",
		"https://",
		"git@",
	} {
		if strings.HasPrefix(address, pattern) {
			return false
		}
	}
	return true
}

func IsHTTP(address string) bool {
	return strings.HasPrefix(address, "http")
}

func parseURL(s string) (*url.URL, error) {
	if s == "" {
		return nil, errors.New("url is empty")
	}

	return url.Parse(Address(s))
}
