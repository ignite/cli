package chainregistry

import (
	"encoding/json"
	"os"
)

const (
	// DefaultChainType is the default chain type for the chain.json
	// More value are allowed by the chain registry schema, but Ignite only scaffolds Cosmos chains.
	DefaultChainType = "cosmos"

	// DefaultChainStatus is the default chain status for the chain.json
	// More value are allowed by the chain registry schema, but Ignite only scaffolds upcoming chains.
	DefaultChainStatus = "upcoming"

	// DefaultNetworkType is the default network type for the chain.json
	// More value are allowed by the chain registry schema, but Ignite only scaffolds devnet chains.
	DefaultNetworkType = "devnet"
)

// Chain represents the chain.json file from the chain registry.
// https://raw.githubusercontent.com/cosmos/chain-registry/master/chain.schema.json
type Chain struct {
	ChainName    string   `json:"chain_name"`
	Status       string   `json:"status"`
	NetworkType  string   `json:"network_type"`
	Website      string   `json:"website"`
	PrettyName   string   `json:"pretty_name"`
	ChainType    string   `json:"chain_type"`
	ChainID      string   `json:"chain_id"`
	Bech32Prefix string   `json:"bech32_prefix"`
	DaemonName   string   `json:"daemon_name"`
	NodeHome     string   `json:"node_home"`
	KeyAlgos     []string `json:"key_algos"`
	Slip44       int      `json:"slip44"`
	Fees         Fees     `json:"fees"`
	Staking      Staking  `json:"staking"`
	Codebase     Codebase `json:"codebase"`
	Description  string   `json:"description"`
	Apis         Apis     `json:"apis"`
}

type Staking struct {
	StakingTokens []StakingToken `json:"staking_tokens"`
}

type StakingToken struct {
	Denom string `json:"denom"`
}

type Codebase struct {
	GitRepo            string              `json:"git_repo"`
	Genesis            CodebaseGenesis     `json:"genesis"`
	RecommendedVersion string              `json:"recommended_version"`
	CompatibleVersions []string            `json:"compatible_versions"`
	Consensus          CodebaseInfo        `json:"consensus"`
	Sdk                CodebaseInfo        `json:"sdk"`
	Ibc                CodebaseInfo        `json:"ibc,omitempty"`
	Cosmwasm           CodebaseInfoEnabled `json:"cosmwasm,omitempty"`
}

type CodebaseGenesis struct {
	GenesisURL string `json:"genesis_url"`
}

type CodebaseInfo struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Repo    string `json:"repo,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type CodebaseInfoEnabled struct {
	Version string `json:"version,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Enabled bool   `json:"enabled"`
}

type Fees struct {
	FeeTokens []FeeToken `json:"fee_tokens"`
}

type FeeToken struct {
	Denom            string  `json:"denom"`
	FixedMinGasPrice float64 `json:"fixed_min_gas_price"`
	LowGasPrice      float64 `json:"low_gas_price"`
	AverageGasPrice  float64 `json:"average_gas_price"`
	HighGasPrice     float64 `json:"high_gas_price"`
}

type Apis struct {
	RPC  []ApiProvider `json:"rpc"`
	Rest []ApiProvider `json:"rest"`
	Grpc []ApiProvider `json:"grpc"`
}

type ApiProvider struct {
	Address  string `json:"address"`
	Provider string `json:"provider"`
}

// SaveJSON saves the chainJSON to the given out directory.
func (c Chain) SaveJSON(out string) error {
	bz, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(out, bz, 0o600)
}
