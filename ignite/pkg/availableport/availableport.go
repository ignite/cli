package availableport

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// Find finds n number of unused ports.
// it is not guaranteed that these ports will not be allocated to
// another program in the time of calling Find().
func Find(n int) (ports []int, err error) {
	min := 44000
	max := 55000

	// If the number of ports required is bigger than the range, this stops it
	if n > (max - min) {
		return nil, fmt.Errorf("Invalid amount of ports requested: limit is %d", min-max)
	}

	// Marker to point if a port is already added in the list
	registered := make(map[int]bool)
	for i := 0; i < n; i++ {
		for {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			port := r.Intn(max-min+1) + min

			conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
			// if there is an error, this might mean that no one is listening from this port
			// which is what we need.
			if err == nil {
				conn.Close()
				continue
			}
			// if the port is already registered we skip it to the next one
			// otherwise it's added to the ports list and pointed in our map
			if registered[port] {
				continue
			}
			ports = append(ports, port)
			registered[port] = true
			break
		}
	}
	return ports, nil
}
