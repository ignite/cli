package scaffolder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chainconfig "github.com/ignite/cli/v29/ignite/config/chain"
	"github.com/ignite/cli/v29/ignite/pkg/chainregistry"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgit"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	chainFilename     = "chain.json"
	assetListFilename = "assetlist.json"
)

// CreateChainRegistryFiles generates the chain registry files in the scaffolded chains.
func (s Scaffolder) CreateChainRegistryFiles(chain *chain.Chain, cfg *chainconfig.Config) error {
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

	chainGitURL, _ /* do not fail on non-existing git repo */ := xgit.RepositoryURL(chain.AppPath())

	var (
		consensus chainregistry.CodebaseInfo
		ibc       chainregistry.CodebaseInfo
		cosmwasm  chainregistry.CodebaseInfoEnabled
	)

	consensusVersion, err := getVersionOfFromGoMod(chain, "github.com/cometbft/cometbft")
	if err == nil {
		consensus = chainregistry.CodebaseInfo{
			Type:    "cometbft",
			Version: consensusVersion,
		}
	}

	cosmwasmVersion, err := getVersionOfFromGoMod(chain, "github.com/CosmWasm/wasmd")
	if err == nil {
		cosmwasm = chainregistry.CodebaseInfoEnabled{
			Version: cosmwasmVersion,
			Enabled: true,
		}
	}

	ibcVersion, err := getVersionOfFromGoMod(chain, "github.com/cosmos/ibc-go")
	if err == nil {
		ibc = chainregistry.CodebaseInfo{
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

	bech32Prefix, err := chain.Bech32Prefix()
	if err != nil {
		return errors.Wrap(err, "failed to get bech32 prefix")
	}

	coinType, err := chain.CoinType()
	if err != nil {
		return errors.Wrap(err, "failed to get coin type")
	}

	chainData := chainregistry.Chain{
		ChainName:    chain.Name(),
		PrettyName:   chain.Name(),
		ChainType:    chainregistry.ChainTypeCosmos,
		Status:       chainregistry.ChainStatusUpcoming,
		NetworkType:  chainregistry.NetworkTypeDevnet,
		Website:      fmt.Sprintf("https://%s.zone", chain.Name()),
		ChainID:      chainID,
		Bech32Prefix: bech32Prefix,
		DaemonName:   binaryName,
		NodeHome:     chainHome,
		KeyAlgos:     []chainregistry.KeyAlgos{chainregistry.KeyAlgoSecp256k1},
		Slip44:       coinType,
		Fees: chainregistry.Fees{
			FeeTokens: []chainregistry.FeeToken{
				{
					Denom:            defaultDenom,
					FixedMinGasPrice: 0.025,
					LowGasPrice:      0.01,
					AverageGasPrice:  0.025,
					HighGasPrice:     0.03,
				},
			},
		},
		Staking: chainregistry.Staking{
			StakingTokens: []chainregistry.StakingToken{
				{
					Denom: defaultDenom,
				},
			},
		},
		Codebase: chainregistry.Codebase{
			GitRepo:            chainGitURL,
			RecommendedVersion: "v1.0.0",
			CompatibleVersions: []string{"v1.0.0"},
			Sdk: chainregistry.CodebaseInfo{
				Type:    "cosmos",
				Version: chain.Version.String(),
			},
			Consensus: consensus,
			Ibc:       ibc,
			Cosmwasm:  cosmwasm,
		},
		APIs: chainregistry.APIs{
			RPC: []chainregistry.APIProvider{
				{
					Address:  "http://localhost:26657",
					Provider: "localhost",
				},
			},
			Rest: []chainregistry.APIProvider{
				{
					Address:  "http://localhost:1317",
					Provider: "localhost",
				},
			},
			Grpc: []chainregistry.APIProvider{
				{
					Address:  "localhost:9090",
					Provider: "localhost",
				},
			},
		},
	}

	assetListData := chainregistry.AssetList{
		ChainName: chainData.ChainName,
		Assets: []chainregistry.Asset{
			{
				Description: fmt.Sprintf("The native token of the %s chain", chainData.ChainName),
				DenomUnits: []chainregistry.DenomUnit{
					{
						Denom:    defaultDenom,
						Exponent: 0,
					},
				},
				Base:   defaultDenom,
				Name:   chainData.ChainName,
				Symbol: strings.ToUpper(defaultDenom),
				LogoURIs: chainregistry.LogoURIs{
					Png: "https://ignite.com/favicon.ico",
					Svg: "https://ignite.com/favicon.ico",
				},
				TypeAsset: "sdk.coin",
				Socials: chainregistry.Socials{
					Website: "https://ignite.com",
					Twitter: "https://x.com/ignite",
				},
			},
		},
	}

	if err := chainData.SaveJSON(chainFilename); err != nil {
		return err
	}

	if err := assetListData.SaveJSON(assetListFilename); err != nil {
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
