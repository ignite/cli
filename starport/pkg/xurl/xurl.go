package xurl

import (
	"fmt"
	"net/url"
	"strings"
)

// TCP unsures that s url contains TCP protocol identifier.
func TCP(s string) string {
	if strings.HasPrefix(s, "tcp") {
		return s
	}
	return "tcp://" + Address(s)
}

// HTTP unsures that s url contains HTTP protocol identifier.
func HTTP(s string) string {
	if strings.HasPrefix(s, "http") {
		return s
	}
	return "http://" + Address(s)
}

// WS unsures that s url contains WS protocol identifier.
func WS(s string) string {
	if strings.HasPrefix(s, "ws") {
		return s
	}
	return "ws://" + Address(s)
}

// HTTPEnsurePort ensures that url has a port number suits with the connection type.
func HTTPEnsurePort(s string) string {
	u, err := url.Parse(s)
	if err != nil || u.Port() != "" {
		return s
	}

	port := "80"

	if u.Scheme == "https" {
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
