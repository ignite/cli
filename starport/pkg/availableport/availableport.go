package availableport

import (
	"net"
)

// Find finds n number of unused ports.
// it is not guaranteed that these ports will not be allocated to
// another program in the time of calling Find().
func Find(n int) (ports []int, err error) {
	for i := 0; i < n; i++ {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, err
		}
		if err := ln.Close(); err != nil {
			return nil, err
		}
		ports = append(ports, ln.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}
