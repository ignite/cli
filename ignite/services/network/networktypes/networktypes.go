package networktypes

import (
	"github.com/ignite-hq/cli/ignite/pkg/relayer"
	relayerconf "github.com/ignite-hq/cli/ignite/pkg/relayer/config"
)

const (
	// SPNChainID name used as SPN chain id.
	SPNChainID = "spn-1"

	// SPN name used as an address prefix and as a home dir for chains to publish.
	SPN = "spn"

	// SPNDenom is the denom used for the spn chain native token
	SPNDenom = "uspn"

	spnVersion  = "monitoring-1"
	spnPortID   = "monitoringc"
	chainPortID = "monitoringp"
)

func SPNRelayerConfig(srcChain, dstChain relayer.Chain) (string, relayerconf.Config) {
	pathID := relayer.PathID(srcChain.ID, dstChain.ID)
	return pathID, relayerconf.Config{
		Version: relayerconf.SupportVersion,
		Chains:  []relayerconf.Chain{srcChain.Config(), dstChain.Config()},
		Paths: []relayerconf.Path{
			{
				ID:       pathID,
				Ordering: relayer.OrderingOrdered,
				Src: relayerconf.PathEnd{
					ChainID: srcChain.ID,
					PortID:  spnPortID,
					Version: spnVersion,
				},
				Dst: relayerconf.PathEnd{
					ChainID: dstChain.ID,
					PortID:  chainPortID,
					Version: spnVersion,
				},
			},
		},
	}
}
