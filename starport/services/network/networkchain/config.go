package networkchain

import (
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
)

const (
	HTTPTunnelChisel = "chisel"
)

const SPNConfigFile = "spn.yml"

type Config struct {
	TunneledPeers []TunneledPeer `json:"tunneled_peers" yaml:"tunneled_peers"`
}

type TunneledPeer struct {
	Name      string `json:"name" yaml:"name"`
	Address   string `json:"address" yaml:"address"`
	NodeID    string `json:"node_id" yaml:"node_id"`
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
