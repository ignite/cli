package availableport

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

type availablePortOptions struct {
	randomizer *rand.Rand
	minPort    uint
	maxPort    uint
}

type Options func(o *availablePortOptions)

func WithRandomizer(r *rand.Rand) Options {
	return func(o *availablePortOptions) {
		o.randomizer = r
	}
}

func WithMaxPort(maxPort uint) Options {
	return func(o *availablePortOptions) {
		o.maxPort = maxPort
	}
}

func WithMinPort(minPort uint) Options {
	return func(o *availablePortOptions) {
		o.minPort = minPort
	}
}

// Find finds n number of unused ports.
// it is not guaranteed that these ports will not be allocated to
// another program in the time of calling Find().
func Find(n uint, options ...Options) (ports []uint, err error) {
	// Defining them before so we can set a value depending on the AvailablePortOptions
	opts := availablePortOptions{
		minPort:    44000,
		maxPort:    55000,
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, apply := range options {
		apply(&opts)
	}
	// If the number of ports required is bigger than the range, this stops it
	if opts.maxPort < opts.minPort {
		return nil, errors.Errorf("invalid ports range: max < min (%d < %d)", opts.maxPort, opts.minPort)
	}

	// If the number of ports required is bigger than the range, this stops it
	if n > (opts.maxPort - opts.minPort) {
		return nil, errors.Errorf("invalid amount of ports requested: limit is %d", opts.maxPort-opts.minPort)
	}

	// Marker to point if a port is already added in the list
	registered := make(map[uint]bool)
	for len(registered) < int(n) {
		// Greater or equal to min and lower than max
		totalPorts := opts.maxPort - opts.minPort + 1
		randomPort := opts.randomizer.Intn(int(totalPorts))
		port := uint(randomPort) + opts.minPort

		conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		// if there is an error, this might mean that no one is listening from this port
		// which is what we need.
		if err == nil {
			conn.Close()
			continue
		}
		if conn != nil {
			defer conn.Close()
		}

		// if the port is already registered we skip it to the next one
		// otherwise it's added to the ports list and pointed in our map
		if registered[port] {
			continue
		}
		ports = append(ports, port)
		registered[port] = true
	}
	return ports, nil
}
