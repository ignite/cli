package xnet

import (
	"fmt"
	"net"
	"strconv"
)

// LocalhostIPv4Address returns a localhost IPv4 address with a port
// that represents the localhost IP address listening on that port.
func LocalhostIPv4Address(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}

// AnyIPv4Address returns an IPv4 meta address "0.0.0.0" with a port
// that represents any IP address listening on that port.
func AnyIPv4Address(port int) string {
	return fmt.Sprintf("0.0.0.0:%d", port)
}

// IncreasePort increases a port number by 1.
// This can be useful to generate port ranges or consecutive
// port numbers for the same address.
func IncreasePort(addr string) (string, error) {
	return IncreasePortBy(addr, 1)
}

// IncreasePortBy increases a port number by a factor of "inc".
// This can be useful to generate port ranges or consecutive
// port numbers for the same address.
func IncreasePortBy(addr string, inc uint64) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	v, err := strconv.ParseUint(port, 10, 0)
	if err != nil {
		return "", err
	}

	port = strconv.FormatUint(v+inc, 10)

	return net.JoinHostPort(host, port), nil
}

// MustIncreasePortBy calls IncreasePortBy and panics on error.
func MustIncreasePortBy(addr string, inc uint64) string {
	s, err := IncreasePortBy(addr, inc)
	if err != nil {
		panic(err)
	}

	return s
}
