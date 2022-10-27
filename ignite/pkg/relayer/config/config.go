package relayerconf

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/confile"
)

const SupportVersion = "2"

var configPath = os.ExpandEnv("$HOME/.ignite/relayer/config.yml")

var (
	ErrChainCannotBeFound = errors.New("chain cannot be found")
	ErrPathCannotBeFound  = errors.New("path cannot be found")
)

type Config struct {
	Version string  `json:"version" yaml:"version"`
	Chains  []Chain `json:"chains" yaml:"chains,omitempty"`
	Paths   []Path  `json:"paths" yaml:"paths,omitempty"`
}

func (c Config) ChainByID(id string) (Chain, error) {
	for _, chain := range c.Chains {
		if chain.ID == id {
			return chain, nil
		}
	}
	return Chain{}, errors.Wrap(ErrChainCannotBeFound, id)
}

func (c Config) PathByID(id string) (Path, error) {
	for _, path := range c.Paths {
		if path.ID == id {
			return path, nil
		}
	}
	return Path{}, errors.Wrap(ErrPathCannotBeFound, id)
}

func (c Config) UpdatePath(path Path) error {
	for i, p := range c.Paths {
		if p.ID == path.ID {
			c.Paths[i] = path
			return nil
		}
	}
	return errors.Wrap(ErrPathCannotBeFound, path.ID)
}

type Chain struct {
	ID            string `json:"id" yaml:"id"`
	Account       string `json:"account" yaml:"account"`
	AddressPrefix string `json:"address_prefix" yaml:"address_prefix"`
	RPCAddress    string `json:"rpc_address" yaml:"rpc_address"`
	GasPrice      string `json:"gas_price" yaml:"gas_price,omitempty"`
	GasLimit      int64  `json:"gas_limit" yaml:"gas_limit,omitempty"`
	ClientID      string `json:"client_id" yaml:"client_id,omitempty"`
}

type Path struct {
	ID       string  `json:"id" yaml:"id"`
	Ordering string  `json:"ordering" yaml:"ordering,omitempty"`
	Src      PathEnd `json:"src" yaml:"src"`
	Dst      PathEnd `json:"dst" yaml:"dst"`
}

type PathEnd struct {
	ChainID      string `json:"chain_id" yaml:"chain_id"`
	ConnectionID string `json:"connection_id" yaml:"connection_id,omitempty"`
	ChannelID    string `json:"channel_id" yaml:"channel_id,omitempty"`
	PortID       string `json:"port_id" yaml:"port_id"`
	Version      string `json:"version" yaml:"version,omitempty"`
	PacketHeight int64  `json:"packet_height" yaml:"packet_height,omitempty"`
	AckHeight    int64  `json:"ack_height" yaml:"ack_height,omitempty"`
}

func Get() (Config, error) {
	c := Config{}
	if err := confile.New(confile.DefaultYAMLEncodingCreator, configPath).Load(&c); err != nil {
		return c, err
	}
	if !reflect.DeepEqual(c, Config{}) && c.Version != SupportVersion {
		return c, fmt.Errorf("your relayer setup is outdated. run 'rm %s' and configure relayer again", configPath)
	}
	return c, nil
}

func Save(c Config) error {
	c.Version = SupportVersion
	return confile.New(confile.DefaultYAMLEncodingCreator, configPath).Save(c)
}

func Delete() error {
	return os.RemoveAll(configPath)
}
