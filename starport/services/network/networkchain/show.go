package networkchain

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

// Peers return the peer addresses from the config.toml of the chain
func (c Chain) Peers() (string, error) {
	// set persistent peers
	configPath, err := c.chain.ConfigTOMLPath()
	if err != nil {
		return "", err
	}
	configToml, err := toml.LoadFile(configPath)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	persistentPeers := configToml.Get("p2p.persistent_peers")
	p2pAddresses, ok := persistentPeers.(string)
	if !ok {
		return "", fmt.Errorf("invalid p2p.persistent_peers paramever %v", persistentPeers)
	}
	return p2pAddresses, nil
}
