package availableport

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

type optionalParameters struct {
	WithRandomizer *rand.Rand
	WithMinPort    int
	WithMaxPort    int
}

// Find finds n number of unused ports.
// it is not guaranteed that these ports will not be allocated to
// another program in the time of calling Find().
func Find(n int, moreParameters ...optionalParameters) (ports []int, err error) {
	// Defining them before so we can set a value depending on the optionalParameters
	var min int
	var max int
	var r *rand.Rand

	if len(moreParameters) != 0 {
		extra := moreParameters[0]
		if extra.WithMinPort != 0 {
			if extra.WithMinPort > -1 {
				min = extra.WithMinPort
			} else {
				// This is not required since the port would become 0
				// but the user could not notice that sent a negative port
				return nil, fmt.Errorf("ports can't be negative (negative min port given)")
			}
		} else {
			min = 44000
		}
		if extra.WithMaxPort != 0 {
			if extra.WithMaxPort > -1 {
				max = extra.WithMaxPort
			} else {
				// This is not required since the port would become 0
				// but the user could not notice that sent a negative port
				return nil, fmt.Errorf("ports can't be negative (negative max port given)")
			}
		} else {
			max = 55000
		}
		if extra.WithRandomizer != nil {
			r = extra.WithRandomizer
		} else {
			r = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
	} else {
		// If we don't require special conditions, we can
		// return to the original parameters
		min = 44000
		max = 55000
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// If the number of ports required is bigger than the range, this stops it
	if max < min {
		return nil, fmt.Errorf("invalid ports range: max < min (%d < %d)", max, min)
	}

	// If the number of ports required is bigger than the range, this stops it
	if n > (max - min) {
		return nil, fmt.Errorf("invalid amount of ports requested: limit is %d", max-min)
	}

	// Marker to point if a port is already added in the list
	registered := make(map[int]bool)
	for i := 0; i < n; i++ {
		for {
			// Greater or equal to min and lower than max
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
