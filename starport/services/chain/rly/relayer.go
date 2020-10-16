package rly

// types are copy pasted from https://github.com/ovrclk/relayer because this package
// is not go get-able for now. once that's fixed, this pkg will be removed.
import (
	"github.com/google/uuid"
)

type GlobalConfig struct {
	Timeout       string `yaml:"timeout" json:"timeout"`
	LiteCacheSize int    `yaml:"lite-cache-size" json:"lite-cache-size"`
}

type Chains []*Chain

type Chain struct {
	Key            string  `yaml:"key" json:"key"`
	ChainID        string  `yaml:"chain-id" json:"chain-id"`
	RPCAddr        string  `yaml:"rpc-addr" json:"rpc-addr"`
	AccountPrefix  string  `yaml:"account-prefix" json:"account-prefix"`
	Gas            uint64  `yaml:"gas,omitempty" json:"gas,omitempty"`
	GasAdjustment  float64 `yaml:"gas-adjustment,omitempty" json:"gas-adjustment,omitempty"`
	GasPrices      string  `yaml:"gas-prices,omitempty" json:"gas-prices,omitempty"`
	DefaultDenom   string  `yaml:"default-denom,omitempty" json:"default-denom,omitempty"`
	Memo           string  `yaml:"memo,omitempty" json:"memo,omitempty"`
	TrustingPeriod string  `yaml:"trusting-period" json:"trusting-period"`
}

type Paths map[string]*Path

type Path struct {
	Src      *PathEnd     `yaml:"src" json:"src"`
	Dst      *PathEnd     `yaml:"dst" json:"dst"`
	Strategy *StrategyCfg `yaml:"strategy" json:"strategy"`
}

type StrategyCfg struct {
	Type string `json:"type" yaml:"type"`
}

type PathEnd struct {
	ChainID      string `yaml:"chain-id,omitempty" json:"chain-id,omitempty"`
	ClientID     string `yaml:"client-id,omitempty" json:"client-id,omitempty"`
	ConnectionID string `yaml:"connection-id,omitempty" json:"connection-id,omitempty"`
	ChannelID    string `yaml:"channel-id,omitempty" json:"channel-id,omitempty"`
	PortID       string `yaml:"port-id,omitempty" json:"port-id,omitempty"`
	Order        string `yaml:"order,omitempty" json:"order,omitempty"`
	Version      string `yaml:"version,omitempty" json:"version,omitempty"`
}

type Config struct {
	Global GlobalConfig `yaml:"global" json:"global"`
	Chains Chains       `yaml:"chains" json:"chains"`
	Paths  Paths        `yaml:"paths" json:"paths"`
}

func NewChain(id, addr string) *Chain {
	return &Chain{
		Key:            "testkey",
		ChainID:        id,
		RPCAddr:        addr,
		AccountPrefix:  "cosmos",
		GasAdjustment:  1.5,
		TrustingPeriod: "336h",
	}
}

func NewPath(src, dst *PathEnd) *Path {
	return &Path{
		Src:      src,
		Dst:      dst,
		Strategy: &StrategyCfg{"naive"},
	}
}

func NewPathEnd(sid, did string) *PathEnd {
	return &PathEnd{
		ChainID:      sid,
		ClientID:     uuid.New().String(),
		ConnectionID: uuid.New().String(),
		ChannelID:    uuid.New().String(),
		PortID:       "transfer",
		Order:        "unordered",
		Version:      "ics20-1",
	}
}
