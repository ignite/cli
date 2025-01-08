package scaffolder

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/services/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	chainFilename     = "chain.json"
	assetListFilename = "assetlist.json"
)

// https://raw.githubusercontent.com/cosmos/chain-registry/master/chain.schema.json
type chainJSON struct {
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
	Fees         struct {
		FeeTokens []FeeToken `json:"fee_tokens"`
	} `json:"fees"`
	Staking     Staking  `json:"staking"`
	Codebase    Codebase `json:"codebase"`
	Description string   `json:"description"`
}

type Staking struct {
	StakingTokens []StakingToken `json:"staking_tokens"`
}

type StakingToken struct {
	Denom string `json:"denom"`
}

type Codebase struct {
	GitRepo            string            `json:"git_repo"`
	Genesis            CodebaseGenesis   `json:"genesis"`
	RecommendedVersion string            `json:"recommended_version"`
	CompatibleVersions []string          `json:"compatible_versions"`
	Consensus          CodebaseConsensus `json:"consensus"`
	Sdk                CodebaseSDK       `json:"sdk"`
	Ibc                CodebaseIBC       `json:"ibc,omitempty"`
	Cosmwasm           CodebaseCosmwam   `json:"cosmwasm,omitempty"`
}

type CodebaseGenesis struct {
	GenesisURL string `json:"genesis_url"`
}

type CodebaseConsensus struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type CodebaseSDK struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type CodebaseIBC struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type CodebaseCosmwam struct {
	Version string `json:"version,omitempty"`
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

// AddChainRegistryFiles generates the chain registry files in the scaffolded chains.
func (s Scaffolder) AddChainRegistryFiles(chain *chain.Chain, cfg *chainconfig.Config) error {
	binaryName, err := chain.Binary()
	if err != nil {
		return errors.Wrap(err, "failed to get binary name")
	}

	chainHome, err := chain.DefaultHome()
	if err != nil {
		return errors.Wrap(err, "failed to get default home directory")
	}

	chainID, err := chain.ID()
	if err != nil {
		return errors.Wrap(err, "failed to get chain ID")
	}

	chainGitURL, _ /* do not fail on non-existing git repo */ := xgit.RespositoryURL(chain.AppPath())

	var (
		consensus CodebaseConsensus
		cosmwasm  CodebaseCosmwam
		ibc       CodebaseIBC
	)

	consensusVersion, err := getVersionOfFromGoMod(chain, "github.com/cometbft/cometbft")
	if err == nil {
		consensus = CodebaseConsensus{
			Type:    "cometbft",
			Version: consensusVersion,
		}
	}

	cosmwasmVersion, err := getVersionOfFromGoMod(chain, "github.com/CosmWasm/wasmd")
	if err == nil {
		cosmwasm = CodebaseCosmwam{
			Version: cosmwasmVersion,
			Enabled: true,
		}
	}

	ibcVersion, err := getVersionOfFromGoMod(chain, "github.com/cosmos/ibc-go")
	if err == nil {
		ibc = CodebaseIBC{
			Type:    "go",
			Version: ibcVersion,
		}
	}

	// get validators from config and parse their coins
	// we can assume it holds the base denom
	defaultDenom := "stake"
	if len(cfg.Validators) > 0 {
		coin, err := sdk.ParseCoinNormalized(cfg.Validators[0].Bonded)
		if err == nil {
			defaultDenom = coin.Denom
		}
	}

	chainJson := chainJSON{
		ChainName:    chain.Name(),
		PrettyName:   chain.Name(),
		ChainType:    DefaultChainType,
		Status:       DefaultChainStatus,
		NetworkType:  DefaultNetworkType,
		Website:      "https://example.com",
		ChainID:      chainID,
		Bech32Prefix: "",
		DaemonName:   binaryName,
		NodeHome:     chainHome,
		KeyAlgos:     []string{"secp256k1"},
		Slip44:       118,
		Fees: Fees{
			FeeTokens: []FeeToken{
				{
					Denom:            defaultDenom,
					FixedMinGasPrice: 0.025,
					LowGasPrice:      0.01,
					AverageGasPrice:  0.025,
					HighGasPrice:     0.03,
				},
			},
		},
		Staking: Staking{
			StakingTokens: []StakingToken{
				{
					Denom: defaultDenom,
				},
			},
		},
		Codebase: Codebase{
			GitRepo:            chainGitURL,
			RecommendedVersion: "v1.0.0",
			CompatibleVersions: []string{"v1.0.0"},
			Sdk: CodebaseSDK{
				Type:    "cosmos",
				Version: chain.Version.String(),
			},
			Consensus: consensus,
			Ibc:       ibc,
			Cosmwasm:  cosmwasm,
		},
	}

	assetListJson := assetListJSON{
		ChainName: chainJson.ChainName,
		Assets:    nil,
	}

	if err := chainJson.SaveJSON(chainFilename); err != nil {
		return err
	}

	if err := assetListJson.SaveJSON(assetListFilename); err != nil {
		return err
	}

	return nil
}

func getVersionOfFromGoMod(chain *chain.Chain, pkg string) (string, error) {
	chainPath := chain.AppPath()

	// get the version from the go.mod file
	file, err := os.Open(filepath.Join(chainPath, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pkg) {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[len(parts)-1], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("consensus version not found in go.mod")
}
