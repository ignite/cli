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
	options := []availableport.Options{
		availableport.WithMinPort(1),
		availableport.WithMaxPort(1),
	}
	ports, err = availableport.Find(10, options...)
	require.Equal(t, fmt.Errorf("invalid amount of ports requested: limit is 0"), err)
	require.Equal(t, 0, len(ports))

	//	// Case max < min
	options = []availableport.Options{
		availableport.WithMinPort(5),
		availableport.WithMaxPort(1),
	}
	ports, err = availableport.Find(10, options...)
	require.Equal(t, fmt.Errorf("invalid ports range: max < min (1 < 5)"), err)
	require.Equal(t, 0, len(ports))

	// Case max < min min restriction given
	options = []availableport.Options{
		availableport.WithMinPort(55001),
		availableport.WithMaxPort(1),
	}
	ports, err = availableport.Find(10, options...)
	require.Equal(t, fmt.Errorf("invalid ports range: max < min (1 < 55001)"), err)
	require.Equal(t, 0, len(ports))

	//	 Case max < min max restriction given
	options = []availableport.Options{
		availableport.WithMaxPort(43999),
	}
	ports, err = availableport.Find(10, options...)
	require.Equal(t, fmt.Errorf("invalid ports range: max < min (43999 < 44000)"), err)
	require.Equal(t, 0, len(ports))

	//	Case randomizer given

	options = []availableport.Options{
		availableport.WithRandomizer(rand.New(rand.NewSource(2023))),
		availableport.WithMinPort(100),
		availableport.WithMaxPort(200),
	}

	ports, err = availableport.Find(10, options...)
	require.Equal(t, 10, len(ports))
	require.Equal(t, []uint([]uint{0xc3, 0x81, 0x96, 0x6c, 0x78, 0xa9, 0xa6, 0x79, 0x83, 0xa0}), ports)
}
