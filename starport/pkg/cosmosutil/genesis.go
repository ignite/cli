package cosmosutil

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/tendermint/starport/starport/pkg/tarball"
)

const (
	genesisFilename  = "genesis.json"
	genesisTimeField = "genesis_time"
	chainIDField     = "chain_id"
)

type (
	// Genesis represents a more readable version of the stargate genesis file
	Genesis struct {
		Accounts   []string
		StakeDenom string
	}
	// ChainGenesis represents the stargate genesis file
	ChainGenesis struct {
		ChainID  string `json:"chain_id"`
		AppState struct {
			Auth struct {
				Accounts []struct {
					Address string `json:"address"`
				} `json:"accounts"`
			} `json:"auth"`
			Staking struct {
				Params struct {
					BondDenom string `json:"bond_denom"`
				} `json:"params"`
			} `json:"staking"`
		} `json:"app_state"`
	}
)

// HasAccount check if account exist into the genesis account
func (g Genesis) HasAccount(address string) bool {
	for _, account := range g.Accounts {
		if account == address {
			return true
		}
	}
	return false
}

// ParseGenesisFromPath parse ChainGenesis object from a genesis file
func ParseGenesisFromPath(genesisPath string) (Genesis, error) {
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return Genesis{}, errors.Wrap(err, "cannot open genesis file")
	}
	return ParseGenesis(genesisFile)
}

// ParseChainGenesis parse ChainGenesis object from a byte slice
func ParseChainGenesis(genesisFile []byte) (chainGenesis ChainGenesis, err error) {
	if err := json.Unmarshal(genesisFile, &chainGenesis); err != nil {
		return chainGenesis, errors.New("cannot unmarshal the chain genesis file: " + err.Error())
	}
	return chainGenesis, err
}

// ParseGenesis parse ChainGenesis object from a byte slice into a Genesis object
func ParseGenesis(genesisFile []byte) (Genesis, error) {
	chainGenesis, err := ParseChainGenesis(genesisFile)
	if err != nil {
		return Genesis{}, errors.New("cannot unmarshal the genesis file: " + err.Error())
	}
	genesis := Genesis{StakeDenom: chainGenesis.AppState.Staking.Params.BondDenom}
	for _, acc := range chainGenesis.AppState.Auth.Accounts {
		genesis.Accounts = append(genesis.Accounts, acc.Address)
	}
	return genesis, nil
}

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(genesisPath, addr string) (bool, error) {
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	genesis, err := ParseGenesisFromPath(genesisPath)
	if err != nil {
		return false, err
	}
	return genesis.HasAccount(addr), nil
}

// SetGenesisTime sets the genesis time inside a genesis file
func SetGenesisTime(genesisPath string, genesisTime int64) error {
	// fetch and parse genesis
	genesisBytes, err := os.ReadFile(genesisPath)
	if err != nil {
		return err
	}

	var genesis map[string]interface{}
	if err := json.Unmarshal(genesisBytes, &genesis); err != nil {
		return err
	}

	// check the genesis time with the RFC3339 standard format
	formattedTime := time.Unix(genesisTime, 0).UTC().Format(time.RFC3339Nano)

	// modify and save the new genesis
	genesis[genesisTimeField] = &formattedTime
	genesisBytes, err = json.Marshal(genesis)
	if err != nil {
		return err
	}
	return os.WriteFile(genesisPath, genesisBytes, 0644)
}

// TODO refactor
func SetChainID(genesisPath, chainID string) error {
	genesisBytes, err := os.ReadFile(genesisPath)
	if err != nil {
		return err
	}

	var genesis map[string]interface{}
	if err := json.Unmarshal(genesisBytes, &genesis); err != nil {
		return err
	}

	genesis[chainIDField] = chainID
	genesisBytes, err = json.Marshal(genesis)
	if err != nil {
		return err
	}
	return os.WriteFile(genesisPath, genesisBytes, 0644)
}

// GenesisAndHashFromURL fetches the genesis from the given url and returns its content along with the sha256 hash.
func GenesisAndHashFromURL(ctx context.Context, url string) (genesis []byte, hash string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	genesis, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(genesis)); err != nil {
		return nil, "", err
	}

	hexHash := hex.EncodeToString(h.Sum(nil))

	return genesis, hexHash, nil
}

// GenesisFromTarball checks if the genesis file is a tarball
// If is a tarball, extract and found the genesis file.
// If isn't a tarball, returns the input genesis again.
func GenesisFromTarball(genesis []byte) (out []byte, isTarball bool, err error) {
	err = tarball.IsTarball(genesis)
	switch {
	case err == tarball.ErrInvalidGzipFile:
		return genesis, false, nil
	case err != nil:
		return genesis, false, err
	default:
		genesis, err = tarball.ReadFile(genesis, genesisFilename)
		return genesis, true, err
	}
}
