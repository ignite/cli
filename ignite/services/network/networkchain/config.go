package networkchain

import (
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/confile"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
)

const (
	HTTPTunnelChisel = "chisel"
)

const SPNConfigFile = "spn.yml"

type Config struct {
	TunneledPeers []TunneledPeer `json:"tunneled_peers" yaml:"tunneled_peers"`
}

// TunneledPeer represents http tunnel to a peer which can't be reached via regular tcp connection
type TunneledPeer struct {
	// Name represents tunnel type e.g. "chisel"
	Name string `json:"name" yaml:"name"`

	// Address represents http address of the tunnel e.g. "https://tendermint-starport-i5e75cplx02.ws-eu31.gitpod.io/"
	Address string `json:"address" yaml:"address"`

	// NodeID tendermint node id of the node behind the tunnel e.g. "e6a59e37b2761f26a21c9168f78a7f2b07c120c7"
	NodeID string `json:"node_id" yaml:"node_id"`

	// LocalPort specifies port which has to be used for local tunnel client
	LocalPort string `json:"local_port" yaml:"local_port"`
}

func GetSPNConfig(path string) (conf Config, err error) {
	err = confile.New(confile.DefaultYAMLEncodingCreator, path).Load(&conf)
	return
}

func SetSPNConfig(config Config, path string) error {
	return confile.New(confile.DefaultYAMLEncodingCreator, path).Save(config)
}

func (c *Chain) SPNConfigPath() (string, error) {
	home, err := c.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, cosmosutil.ChainConfigDir, SPNConfigFile), nil
}
