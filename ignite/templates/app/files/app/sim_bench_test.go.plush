//go:build sims

package app

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simsx"

	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
)

// Profile with:
// `go test -benchmem -run=^$ ./app -bench ^BenchmarkFullAppSimulation$ -Commit=true -cpuprofile cpu.out`
func BenchmarkFullAppSimulation(b *testing.B) {
	b.ReportAllocs()

	config := simcli.NewConfigFromFlags()
	config.ChainID = simsx.SimAppChainID

	simsx.RunWithSeed(b, config, New, setupStateFactory, 1, nil)
}
