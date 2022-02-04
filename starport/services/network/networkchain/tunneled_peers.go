package networkchain

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const TunneledPeersFile = "tunneled_peers.json"

type TunneledPeerConfig struct {
	TunneledPeers []TunneledPeer `json:"tunneled_peers"`
}

type TunneledPeer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	NodeID  string `json:"node_id"`
}

func GetTunneledPeersConfig(path string) (TunneledPeerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return TunneledPeerConfig{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return TunneledPeerConfig{}, err
	}

	var result TunneledPeerConfig
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return TunneledPeerConfig{}, err
	}

	return result, nil
}

func SetTunneledPeersConfig(config TunneledPeerConfig, path string) error {
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

func (c *Chain) TunneledPeersConfigPath() (string, error) {
	home, err := c.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "config", TunneledPeersFile), nil
}
