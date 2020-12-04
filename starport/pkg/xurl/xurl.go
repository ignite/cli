package xurl

import "strings"

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
