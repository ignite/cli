package scaffolder

import (
	"encoding/json"
	"os"
)

const (
	// DefaultChainType is the default chain type for the chain.json
	// More value are allowed by the chain registry schema, but Ignite only scaffold Cosmos chains.
	DefaultChainType = "cosmos"
)

// Status enumerates possible status for the chain.json
type Status uint8

const (
	StatusLive = iota
	StatusUpcoming
	StatusKilled
)

func (s Status) String() string {
	switch s {
	case StatusLive:
		return "live"
	case StatusUpcoming:
		return "upcoming"
	case StatusKilled:
		return "killed"
	default:
		return "unknown"
	}
}

// NetworkType enumerates possible network types for the chain.json
type NetworkType uint8

const (
	NetworkMainnet = iota
	NetworkTestnet
	NetworkDevnet
)

func (n NetworkType) String() string {
	switch n {
	case NetworkMainnet:
		return "mainnet"
	case NetworkTestnet:
		return "testnet"
	case NetworkDevnet:
		return "devnet"
	default:
		return "unknown"
	}
}

// https://raw.githubusercontent.com/cosmos/chain-registry/master/chain.schema.json
// https://github.com/cosmos/chain-registry?tab=readme-ov-file#sample
type chainJSON struct {
	ChainName    string      `json:"chain_name"`
	Status       Status      `json:"status"`
	NetworkType  NetworkType `json:"network_type"`
	Website      string      `json:"website"`
	PrettyName   string      `json:"pretty_name"`
	ChainType    string      `json:"chain_type"`
	ChainID      string      `json:"chain_id"`
	Bech32Prefix string      `json:"bech32_prefix"`
	DaemonName   string      `json:"daemon_name"`
	NodeHome     string      `json:"node_home"`
	KeyAlgos     []string    `json:"key_algos"`
	Slip44       int         `json:"slip44"`
	Fees         struct {
		FeeTokens []struct {
			Denom            string  `json:"denom"`
			FixedMinGasPrice float64 `json:"fixed_min_gas_price"`
			LowGasPrice      float64 `json:"low_gas_price"`
			AverageGasPrice  float64 `json:"average_gas_price"`
			HighGasPrice     float64 `json:"high_gas_price"`
		} `json:"fee_tokens"`
	} `json:"fees"`
	Staking struct {
		StakingTokens []struct {
			Denom string `json:"denom"`
		} `json:"staking_tokens"`
		LockDuration struct {
			Time string `json:"time"`
		} `json:"lock_duration"`
	} `json:"staking"`
	Codebase struct {
		GitRepo string `json:"git_repo"`
		Genesis struct {
			Name       string `json:"name"`
			GenesisURL string `json:"genesis_url"`
		} `json:"genesis"`
		RecommendedVersion string   `json:"recommended_version"`
		CompatibleVersions []string `json:"compatible_versions"`
		Consensus          struct {
			Type    string `json:"type"`
			Version string `json:"version"`
			Repo    string `json:"repo"`
			Tag     string `json:"tag"`
		} `json:"consensus"`
		Binaries struct {
			LinuxAmd64 string `json:"linux/amd64"`
			LinuxArm64 string `json:"linux/arm64"`
		} `json:"binaries"`
		Language struct {
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"language"`
		Sdk struct {
			Type    string `json:"type"`
			Repo    string `json:"repo"`
			Version string `json:"version"`
			Tag     string `json:"tag"`
		} `json:"sdk"`
		Ibc struct {
			Type       string   `json:"type"`
			Version    string   `json:"version"`
			IcsEnabled []string `json:"ics_enabled"`
		} `json:"ibc"`
		Cosmwasm struct {
			Version string `json:"version"`
			Repo    string `json:"repo"`
			Tag     string `json:"tag"`
			Enabled bool   `json:"enabled"`
		} `json:"cosmwasm"`
	} `json:"codebase"`
	Images []struct {
		ImageSync struct {
			ChainName string `json:"chain_name"`
			BaseDenom string `json:"base_denom"`
		} `json:"image_sync"`
		Svg   string `json:"svg"`
		Png   string `json:"png"`
		Theme struct {
			PrimaryColorHex string `json:"primary_color_hex"`
		} `json:"theme"`
	} `json:"images"`
	LogoURIs struct {
		Png string `json:"png"`
		Svg string `json:"svg"`
	} `json:"logo_URIs"`
	Description string `json:"description"`
	Peers       struct {
		Seeds []struct {
			ID       string `json:"id"`
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"seeds"`
		PersistentPeers []struct {
			ID       string `json:"id"`
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"persistent_peers"`
	} `json:"peers"`
	Apis struct {
		RPC []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"rpc"`
		Rest []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"rest"`
		Grpc []struct {
			Address  string `json:"address"`
			Provider string `json:"provider"`
		} `json:"grpc"`
	} `json:"apis"`
	Explorers []struct {
		Kind        string `json:"kind"`
		URL         string `json:"url"`
		TxPage      string `json:"tx_page,omitempty"`
		AccountPage string `json:"account_page,omitempty"`
	} `json:"explorers"`
	Keywords []string `json:"keywords"`
}

// SaveJSON saves the chainJSON to the given out directory.
func (c chainJSON) SaveJSON(out string) error {
	bz, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(out, bz, 0666)
}

// https://raw.githubusercontent.com/cosmos/chain-registry/master/assetlist.schema.json
// https://github.com/cosmos/chain-registry?tab=readme-ov-file#assetlists
type assetListJSON struct {
	ChainName string `json:"chain_name"`
	Assets    []struct {
		Description         string `json:"description"`
		ExtendedDescription string `json:"extended_description,omitempty"`
		DenomUnits          []struct {
			Denom    string `json:"denom"`
			Exponent int    `json:"exponent"`
		} `json:"denom_units"`
		Base     string `json:"base"`
		Name     string `json:"name"`
		Display  string `json:"display"`
		Symbol   string `json:"symbol"`
		LogoURIs struct {
			Png string `json:"png"`
			Svg string `json:"svg"`
		} `json:"logo_URIs"`
		CoingeckoID string `json:"coingecko_id,omitempty"`
		Images      []struct {
			Png   string `json:"png"`
			Svg   string `json:"svg"`
			Theme struct {
				PrimaryColorHex string `json:"primary_color_hex"`
			} `json:"theme"`
		} `json:"images"`
		Socials struct {
			Website string `json:"website"`
			Twitter string `json:"twitter"`
		} `json:"socials,omitempty"`
		TypeAsset string `json:"type_asset"`
		Traces    []struct {
			Type         string `json:"type"`
			Counterparty struct {
				ChainName string `json:"chain_name"`
				BaseDenom string `json:"base_denom"`
				ChannelID string `json:"channel_id"`
			} `json:"counterparty"`
			Chain struct {
				ChannelID string `json:"channel_id"`
				Path      string `json:"path"`
			} `json:"chain"`
		} `json:"traces,omitempty"`
	} `json:"assets"`
}

// SaveJSON saves the assetList to the given out directory.
func (c assetListJSON) SaveJSON(out string) error {
	bz, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(out, bz, 0666)
}

// CreateChainRegistryFiles creates a the chain registry files in the scaffolded chains.
func (s Scaffolder) CreateChainRegistryFiles() error {
	return nil
}
