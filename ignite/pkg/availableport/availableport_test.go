package availableport_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ignite/cli/ignite/pkg/availableport"
	"github.com/stretchr/testify/require"
)

func TestAvailablePort(t *testing.T) {

	ports, err := availableport.Find(10)
	require.Equal(t, nil, err)
	require.Equal(t, 10, len(ports))

	for idx := 0; idx < 9; idx++ {
		for jdx := idx + 1; jdx < 10; jdx++ {
			require.NotEqual(t, ports[idx], ports[jdx])
		}
	}

	// Case no ports generated
	options := availableport.OptionalParameters{
		WithMinPort: 1,
		WithMaxPort: 1,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, fmt.Errorf("invalid amount of ports requested: limit is 0"), err)
	require.Equal(t, 0, len(ports))

	// Case max < min
	options = availableport.OptionalParameters{
		WithMinPort: 5,
		WithMaxPort: 1,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, fmt.Errorf("invalid ports range: max < min (1 < 5)"), err)
	require.Equal(t, 0, len(ports))

	// Case max < min min restriction given
	options = availableport.OptionalParameters{
		WithMinPort: 55001,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, fmt.Errorf("invalid ports range: max < min (55000 < 55001)"), err)
	require.Equal(t, 0, len(ports))

	// Case max < min max restriction given
	options = availableport.OptionalParameters{
		WithMaxPort: 43999,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, fmt.Errorf("invalid ports range: max < min (43999 < 44000)"), err)
	require.Equal(t, 0, len(ports))

	// Case negative min
	options = availableport.OptionalParameters{
		WithMinPort: -10,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, fmt.Errorf("ports can't be negative (negative min port given)"), err)
	require.Equal(t, 0, len(ports))

	// Case negative max
	options = availableport.OptionalParameters{
		WithMaxPort: -10,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, fmt.Errorf("ports can't be negative (negative max port given)"), err)
	require.Equal(t, 0, len(ports))

	// Case randomizer given
	options = availableport.OptionalParameters{
		WithRandomizer: rand.New(rand.NewSource(2023)),
		WithMinPort:    100,
		WithMaxPort:    200,
	}
	ports, err = availableport.Find(10, options)
	require.Equal(t, 10, len(ports))
	require.Equal(t, []int{195, 129, 150, 108, 120, 169, 166, 121, 131, 160}, ports)
}
