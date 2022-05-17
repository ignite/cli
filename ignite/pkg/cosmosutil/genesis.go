package cosmosutil

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

const (
	FieldGenesisTime                 = "genesis_time"
	FieldChainID                     = "chain_id"
	FieldConsumerChainID             = "app_state.monitoringp.params.consumerChainID"
	FieldLastBlockHeight             = "app_state.monitoringp.params.lastBlockHeight"
	FieldConsensusTimestamp          = "app_state.monitoringp.params.consumerConsensusState.timestamp"
	FieldConsensusNextValidatorsHash = "app_state.monitoringp.params.consumerConsensusState.nextValidatorsHash"
	FieldConsensusRootHash           = "app_state.monitoringp.params.consumerConsensusState.root.hash"
	FieldConsumerUnbondingPeriod     = "app_state.monitoringp.params.consumerUnbondingPeriod"
	FieldConsumerRevisionHeight      = "app_state.monitoringp.params.consumerRevisionHeight"
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

	// fields to update from genesis
	fields map[string]string
	// GenesisField configures the genesis key value fields.
	GenesisField func(fields)
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

// WithKeyValue sets key and value field to genesis file
func WithKeyValue(key, value string) GenesisField {
	return func(f fields) {
		f[key] = value
	}
}

// WithKeyValueInt sets key and int64 value field to genesis file
func WithKeyValueInt(key string, value int64) GenesisField {
	return func(f fields) {
		f[key] = strconv.FormatInt(value, 10)
	}
}

// WithKeyValueUint sets key and uint64 value field to genesis file
func WithKeyValueUint(key string, value uint64) GenesisField {
	return func(f fields) {
		f[key] = strconv.FormatUint(value, 10)
	}
}

// WithKeyValueTimestamp sets key and timestamp value field to genesis file
func WithKeyValueTimestamp(key string, value int64) GenesisField {
	return func(f fields) {
		f[key] = time.Unix(value, 0).UTC().Format(time.RFC3339Nano)
	}
}

// WithKeyValueBoolean sets key and boolean value field to genesis file
func WithKeyValueBoolean(key string, value bool) GenesisField {
	return func(f fields) {
		if value {
			f[key] = "true"
		} else {
			f[key] = "false"
		}
	}
}

func UpdateGenesis(genesisPath string, options ...GenesisField) error {
	f := fields{}
	for _, applyField := range options {
		applyField(f)
	}

	genesisBytes, err := os.ReadFile(genesisPath)
	if err != nil {
		return err
	}

	for key, value := range f {
		genesisBytes, err = jsonparser.Set(
			genesisBytes,
			[]byte(fmt.Sprintf(`"%s"`, value)),
			strings.Split(key, ".")...,
		)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(genesisPath, genesisBytes, 0644)
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
