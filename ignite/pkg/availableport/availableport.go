package availableport

import (
	"fmt"
	"math/rand"
	"net"
)

// Find finds n number of unused ports.
// it is not guaranteed that these ports will not be allocated to
// another program in the time of calling Find().
func Find(n int) (ports []int, err error) {
	min := 44000
	max := 55000

	for i := 0; i < n; i++ {
		for {
			port := rand.Intn(max-min+1) + min

			conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
			// if there is an error, this might mean that no one is listening from this port
			// which is what we need.
			if err == nil {
				conn.Close()
				continue
			}
			ports = append(ports, port)
			break
		}
	}
	return ports, nil
}
