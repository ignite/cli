package networkchain

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

// Peers return the peer addresses from the config.toml of the chain
func (c Chain) Peers() ([]string, error) {
	// set persistent peers
	configPath, err := c.chain.ConfigTOMLPath()
	if err != nil {
		return nil, err
	}
	configToml, err := toml.LoadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	persistentPeers := configToml.Get("p2p.persistent_peers")
	p2pAddresses, ok := persistentPeers.(string)
	if !ok {
		return nil, fmt.Errorf("invalid p2p.persistent_peers paramever %v", persistentPeers)
	}
	return strings.Split(p2pAddresses, ","), nil
}
