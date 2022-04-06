package genesis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	gojson "github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/tendermint/starport/starport/pkg/tarball"
)

const (
	keySeparator    = "."
	genesisFilename = "genesis.json"
	paramStakeDenom = "app_state.staking.params.bond_denom"
	paramChainID    = "chain_id"
	paramAccounts   = "app_state.auth.accounts"
)

type (
	// Genesis represents the genesis reader
	Genesis struct {
		file        readWriteSeeker
		tarballPath string
		updates     map[string][]byte
	}
	// UpdateGenesisOption configures genesis update.
	UpdateGenesisOption func(map[string][]byte)

	writeTruncate interface {
		Truncate(size int64) error
	}
	readWriteSeeker interface {
		io.ReadWriteSeeker
		Close() error
		Sync() error
	}
	accounts []struct {
		Address string `json:"address"`
	}
)

var (
	ErrParamNotFound    = errors.New("parameter not found")
	ErrInvalidValueType = errors.New("invalid value type")
)

// New creates a new Genesis
func New(file readWriteSeeker) *Genesis {
	return &Genesis{
		updates: make(map[string][]byte),
		file:    file,
	}
}

// FromPath parse Genesis object from path
func FromPath(genesisPath string) (*Genesis, error) {
	file, err := os.OpenFile(genesisPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open genesis file")
	}
	return New(file), nil
}

// FromURL fetches the genesis from the given URL and returns its content.
func FromURL(ctx context.Context, url, genesisPath string) (*Genesis, error) {
	// TODO create a cache system to avoid download genesis with the same hash again

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(genesisPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create the genesis file")
	}

	tarballPath, err := tarball.ExtractFile(resp.Body, file, genesisFilename)
	if err != nil {
		return nil, err
	}
	return &Genesis{
		updates:     make(map[string][]byte),
		file:        file,
		tarballPath: tarballPath,
	}, nil
}

// StakeDenom returns the stake denom from the genesis
func (g *Genesis) StakeDenom() (denom string, err error) {
	_, err = g.Param(paramStakeDenom, &denom)
	return
}

// ChainID returns the chain id from the genesis
func (g *Genesis) ChainID() (chainID string, err error) {
	_, err = g.Param(paramChainID, &chainID)
	return
}

// Accounts returns the auth accounts from the genesis
func (g *Genesis) Accounts() ([]string, error) {
	var accs accounts
	_, err := g.Param(paramAccounts, &accs)
	accountList := make([]string, len(accs))
	for i, acc := range accs {
		accountList[i] = acc.Address
	}
	return accountList, err
}

// Param return the param and the position into byte slice from the file reader
func (g *Genesis) Param(key string, param interface{}) (int64, error) {
	// TODO find a better way to reset the reader
	if _, err := g.file.Seek(0, 0); err != nil {
		return 0, err
	}
	dec := gojson.NewDecoder(g.file)
	keys := strings.Split(key, keySeparator)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		name, ok := t.(string)
		if !ok {
			continue
		}
		if name == keys[0] {
			if len(keys) > 1 {
				keys = keys[1:]
				continue
			}
			err := dec.Decode(&param)
			if err == nil {
				return dec.InputOffset(), nil
			}

			t, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				return 0, err
			}
			switch t := t.(type) {
			case int:
				param = strconv.Itoa(t)
			case string:
				param = t
			default:
				return 0, ErrInvalidValueType
			}
			return dec.InputOffset(), nil
		}
	}
	return 0, ErrParamNotFound
}

// WithKeyValue update a genesis value object by key
func WithKeyValue(key string, value string) UpdateGenesisOption {
	return func(update map[string][]byte) {
		update[key] = []byte(`"` + value + `"`)
	}
}

// WithTime update a time value
func WithTime(key string, t int64) UpdateGenesisOption {
	return func(update map[string][]byte) {
		formatted := time.Unix(t, 0).UTC().Format(time.RFC3339Nano)
		update[key] = []byte(`"` + formatted + `"`)
	}
}

// WithKeyIntValue update a genesis int value object by key
func WithKeyIntValue(key string, value int) UpdateGenesisOption {
	return func(update map[string][]byte) {
		update[key] = []byte{byte(value)}
	}
}

// Update updates the genesis file with the new parameters by key
func (g *Genesis) Update(opts ...UpdateGenesisOption) error {
	for _, opt := range opts {
		opt(g.updates)
	}
	// TODO find a better way to reset the reader
	if _, err := g.file.Seek(0, 0); err != nil {
		return err
	}
	_, err := io.Copy(g, g.file)
	return err
}

// Write implement the write method for io.Writer interface
func (g *Genesis) Write(p []byte) (int, error) {
	var err error
	length := len(p)
	for key, value := range g.updates {
		p, err = jsonparser.Set(p, value, strings.Split(key, keySeparator)...)
		if err != nil {
			return 0, err
		}
		delete(g.updates, key)
	}

	if length > len(p) {
		err = truncate(g.file, len(p))
		if err != nil {
			return 0, err
		}
	}

	// TODO find a better way to reset the writer
	if _, err := g.file.Seek(0, 0); err != nil {
		return 0, err
	}
	n, err := g.file.Write(p)
	if err != nil {
		return n, err
	}

	if n != len(p) {
		return n, io.ErrShortWrite
	}

	// FIXME passing the new byte slice length throws an error
	// because the reader has less byte length than the writer
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/io/io.go;l=432
	return length, nil
}

// truncate remove the current file content
func truncate(rws io.WriteSeeker, size int) error {
	t, ok := rws.(writeTruncate)
	if !ok {
		return errors.New("truncate: unable to truncate")
	}
	return t.Truncate(int64(size))
}

// Close the file
func (g *Genesis) Close() error {
	return g.file.Close()
}

// Sync save the file
func (g *Genesis) Sync() error {
	return g.file.Sync()
}

// TarballPath returns the tarball path
func (g *Genesis) TarballPath() string {
	return g.tarballPath
}

// Save the genesis file
func (g *Genesis) Save(path string) error {
	if g.file != nil {
		return g.Sync()
	}
	reader, err := FromPath(path)
	if err != nil {
		return err
	}
	g.file = reader.file
	return nil
}

// Hash returns the hash of the file
func (g *Genesis) Hash() (string, error) {
	// TODO find a better way to reset the writer
	if _, err := g.file.Seek(0, 0); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, g.file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// String returns the file string
func (g *Genesis) String() (string, error) {
	// TODO find a better way to reset the writer
	if _, err := g.file.Seek(0, 0); err != nil {
		return "", err
	}
	data, err := io.ReadAll(g.file)
	return string(data), err
}
