package chainregistry

import (
	"encoding/json"
	"os"
)

// AssetList represents the assetlist.json file from the chain registry.
// https://raw.githubusercontent.com/cosmos/chain-registry/master/assetlist.schema.json
// https://github.com/cosmos/chain-registry?tab=readme-ov-file#assetlists
type AssetList struct {
	ChainName string  `json:"chain_name"`
	Assets    []Asset `json:"assets"`
}

type Asset struct {
	Description string      `json:"description"`
	DenomUnits  []DenomUnit `json:"denom_units"`
	Base        string      `json:"base"`
	Name        string      `json:"name"`
	Display     string      `json:"display"`
	Symbol      string      `json:"symbol"`
	LogoURIs    LogoURIs    `json:"logo_URIs"`
	CoingeckoID string      `json:"coingecko_id,omitempty"`
	Socials     Socials     `json:"socials,omitempty"`
	TypeAsset   string      `json:"type_asset"`
}

type DenomUnit struct {
	Denom    string `json:"denom"`
	Exponent int    `json:"exponent"`
}

type LogoURIs struct {
	Png string `json:"png"`
	Svg string `json:"svg"`
}

type Socials struct {
	Website string `json:"website"`
	Twitter string `json:"twitter"`
}

// SaveJSON saves the assetlist.json to the given out directory.
func (c AssetList) SaveJSON(out string) error {
	bz, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(out, bz, 0o600)
}
