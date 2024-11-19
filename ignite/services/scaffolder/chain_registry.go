package scaffolder

// https://raw.githubusercontent.com/cosmos/chain-registry/master/chain.schema.json
// https://github.com/cosmos/chain-registry?tab=readme-ov-file#sample
type chainJSON struct {
	Schema       string   `json:"$schema"`
	ChainName    string   `json:"chain_name"`
	ChainType    string   `json:"chain_type"`
	ChainID      string   `json:"chain_id"`
	Website      string   `json:"website"`
	PrettyName   string   `json:"pretty_name"`
	Status       string   `json:"status"`
	NetworkType  string   `json:"network_type"`
	Bech32Prefix string   `json:"bech32_prefix"`
	DaemonName   string   `json:"daemon_name"`
	NodeHome     string   `json:"node_home"`
	KeyAlgos     []string `json:"key_algos"`
	Slip44       int      `json:"slip44"`
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
	} `json:"staking"`
	Codebase struct {
		GitRepo            string   `json:"git_repo"`
		RecommendedVersion string   `json:"recommended_version"`
		CompatibleVersions []string `json:"compatible_versions"`
		Consensus          struct {
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"consensus"`
		Binaries struct {
			LinuxAmd64  string `json:"linux/amd64"`
			LinuxArm64  string `json:"linux/arm64"`
			DarwinAmd64 string `json:"darwin/amd64"`
			DarwinArm64 string `json:"darwin/arm64"`
		} `json:"binaries"`
		Genesis struct {
			GenesisURL string `json:"genesis_url"`
		} `json:"genesis"`
		Versions []struct {
			Name               string   `json:"name"`
			Tag                string   `json:"tag"`
			RecommendedVersion string   `json:"recommended_version"`
			CompatibleVersions []string `json:"compatible_versions"`
			Consensus          struct {
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"consensus"`
			Height   int `json:"height"`
			Binaries struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
			} `json:"binaries,omitempty"`
			NextVersionName string `json:"next_version_name"`
			Sdk             struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Tag     string `json:"tag"`
			} `json:"sdk"`
			Ibc struct {
				Type    string `json:"type"`
				Version string `json:"version"`
			} `json:"ibc"`
			Proposal  int `json:"proposal,omitempty"`
			Binaries0 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries1 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries2 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries3 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries4 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries5 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries6 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Binaries7 struct {
				LinuxAmd64   string `json:"linux/amd64"`
				LinuxArm64   string `json:"linux/arm64"`
				DarwinAmd64  string `json:"darwin/amd64"`
				DarwinArm64  string `json:"darwin/arm64"`
				WindowsAmd64 string `json:"windows/amd64"`
				WindowsArm64 string `json:"windows/arm64"`
			} `json:"binaries,omitempty"`
			Cosmwasm struct {
				Version string `json:"version"`
				Repo    string `json:"repo"`
				Tag     string `json:"tag"`
			} `json:"cosmwasm,omitempty"`
			Binaries8 struct {
				LinuxAmd64  string `json:"linux/amd64"`
				LinuxArm64  string `json:"linux/arm64"`
				DarwinAmd64 string `json:"darwin/amd64"`
				DarwinArm64 string `json:"darwin/arm64"`
			} `json:"binaries,omitempty"`
		} `json:"versions"`
		Sdk struct {
			Type    string `json:"type"`
			Version string `json:"version"`
			Tag     string `json:"tag"`
		} `json:"sdk"`
		Ibc struct {
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"ibc"`
		Cosmwasm struct {
			Version string `json:"version"`
			Repo    string `json:"repo"`
			Tag     string `json:"tag"`
		} `json:"cosmwasm"`
	} `json:"codebase"`
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
			Provider string `json:"provider,omitempty"`
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
		Kind          string `json:"kind"`
		URL           string `json:"url"`
		TxPage        string `json:"tx_page,omitempty"`
		AccountPage   string `json:"account_page,omitempty"`
		ValidatorPage string `json:"validator_page,omitempty"`
		ProposalPage  string `json:"proposal_page,omitempty"`
		BlockPage     string `json:"block_page,omitempty"`
	} `json:"explorers"`
	Images []struct {
		Png   string `json:"png"`
		Svg   string `json:"svg"`
		Theme struct {
			PrimaryColorHex string `json:"primary_color_hex"`
		} `json:"theme"`
	} `json:"images"`
}

// https://raw.githubusercontent.com/cosmos/chain-registry/master/assetlist.schema.json
// https://github.com/cosmos/chain-registry?tab=readme-ov-file#assetlists
type assetlistJson struct {
	Schema    string `json:"$schema"`
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

// CreateChainRegistryFiles creates a the chain registry files in the scaffolded chains.
func (s Scaffolder) CreateChainRegistryFiles() error {
	return nil
}
