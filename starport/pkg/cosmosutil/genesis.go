package cosmosutil

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/tendermint/starport/starport/pkg/tarball"
)

const (
	genesisFilename = "genesis.json"
)

type (
	// Genesis represents a more readable version of the stargate genesis file
	Genesis struct {
		Accounts   []string
		StakeDenom string
	}
	// ChainGenesis represents the stargate genesis file
	ChainGenesis struct {
		ChainID     string `json:"chain_id"`
		GenesisTime string `json:"genesis_time"`
		AppState    struct {
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
	// GenReader represents the genesis reader/writer
	GenReader struct {
		URL         string
		TarballPath string
		FilePath    string
		io.ReadWriter
	}
	// UpdateGenesisOption configures genesis update.
	UpdateGenesisOption func(*ChainGenesis)
)

// GenesisReaderFromPath parse GenReader object from path
func GenesisReaderFromPath(genesisPath string) (*GenReader, error) {
	genesisFile, err := os.Open(genesisPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open genesis file")
	}
	return &GenReader{
		FilePath:   genesisPath,
		ReadWriter: genesisFile,
	}, nil
}

// GenesisFromURL fetches the genesis from the given URL and returns its content.
func GenesisFromURL(ctx context.Context, url string) (genReader *GenReader, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tarballPath, err := RetrieveGenesis(resp.Body, genReader)
	genReader.URL = url
	genReader.TarballPath = tarballPath
	return genReader, err
}

// RetrieveGenesis checks if the genesis file is a tarball
// If is a tarball, extract and found the genesis file.
// If isn't a tarball, returns the input genesis again.
func RetrieveGenesis(input io.Reader, genesis io.Writer) (genesisPath string, err error) {
	err = tarball.IsTarball(input)
	switch {
	case err == tarball.ErrInvalidGzipFile:
		_, err = io.Copy(genesis, input)
		return
	case err != nil:
		return
	default:
		return tarball.ExtractFile(input, genesis, genesisFilename)
	}
}

// CheckGenesisContainsAddress returns true if the address exist into the genesis file
func CheckGenesisContainsAddress(genesisPath, addr string) (bool, error) {
	_, err := os.Stat(genesisPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	genesis, err := GenesisReaderFromPath(genesisPath)
	if err != nil {
		return false, err
	}
	return genesis.HasAccount(addr), nil
}

// HasAccount check if account exist into the genesis account
func (g Genesis) HasAccount(address string) bool {
	for _, account := range g.Accounts {
		if account == address {
			return true
		}
	}
	return false
}

func applyChanges(g *ChainGenesis, options []UpdateGenesisOption) {
	for _, applyOption := range options {
		applyOption(g)
	}
}

// WithChainID update a genesis chaind id
func WithChainID(chainID string) UpdateGenesisOption {
	return func(g *ChainGenesis) {
		g.ChainID = chainID
	}
}

// WithGenesisTime update a genesis time
func WithGenesisTime(genesisTime int64) UpdateGenesisOption {
	return func(g *ChainGenesis) {
		g.GenesisTime = time.Unix(genesisTime, 0).UTC().Format(time.RFC3339Nano)
	}
}

// UpdateGenesis update the genesis file with options
func (g *GenReader) UpdateGenesis(options ...UpdateGenesisOption) (genesis ChainGenesis, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(g).Decode(genesis)
	if err != nil {
		return
	}
	applyChanges(&genesis, options)
	return genesis, json.NewEncoder(g).Encode(genesis)
}

// HasAccount check if account exist into the genesis account
func (g *GenReader) HasAccount(address string) bool {
	genesis, err := g.Genesis()
	if err != nil {
		return false
	}
	for _, account := range genesis.Accounts {
		if account == address {
			return true
		}
	}
	return false
}

// StakeDenom returns the chain genesis stake denom
func (g *GenReader) StakeDenom() (string, error) {
	var genesis ChainGenesis
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return genesis.AppState.Staking.Params.BondDenom, json.NewDecoder(g).Decode(genesis)
}

// ChainGenesis returns the chain genesis form the reader
func (g *GenReader) ChainGenesis() (genesis ChainGenesis, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(g).Decode(genesis)
	return
}

// Genesis returns the genesis wrapper form the reader
func (g *GenReader) Genesis() (Genesis, error) {
	chainGenesis, err := g.ChainGenesis()
	if err != nil {
		return Genesis{}, err
	}
	accounts := make([]string, len(chainGenesis.AppState.Auth.Accounts))
	for i, acc := range chainGenesis.AppState.Auth.Accounts {
		accounts[i] = acc.Address
	}
	return Genesis{
		StakeDenom: chainGenesis.AppState.Staking.Params.BondDenom,
		Accounts:   accounts,
	}, nil
}

func (g *GenReader) Hash() (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, g); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (g *GenReader) String() (string, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, g); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Save saves the genesis writer to the file
func (g *GenReader) Save() error {
	if g.FilePath == "" {
		return errors.New("genesis path is empty")
	}
	genesisFile, err := os.Create(g.FilePath)
	if err != nil {
		return errors.Wrapf(err, "cannot create the genesis file %s", g.FilePath)
	}
	defer genesisFile.Close()
	_, err = genesisFile.ReadFrom(g)
	return err
}
