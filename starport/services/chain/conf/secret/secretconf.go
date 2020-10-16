package secretconf

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/go-bip39"
	"github.com/tendermint/starport/starport/services/chain/conf"
	"gopkg.in/yaml.v2"
)

const (
	SecretFile             = "secret.yml"
	SelfRelayerAccountName = "relayer"
)

var (
	selfRelayerAccountDefaultCoins = []string{"800token"}
)

type Config struct {
	// Accounts of the local chain.
	Accounts []conf.Account `yaml:"accounts"`

	// Relayer configuration.
	Relayer Relayer `yaml:"relayer"`
}

func (c *Config) SelfRelayerAccount(name string) (account conf.Account, found bool) {
	for _, a := range c.Accounts {
		if a.Name == name {
			return a, true
		}
	}
	return conf.Account{}, false
}

func (c *Config) SetSelfRelayerAccount(accName string) error {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return err
	}
	c.Accounts = append(c.Accounts, conf.Account{
		Name:     accName,
		Coins:    selfRelayerAccountDefaultCoins,
		Mnemonic: mnemonic,
	})
	return nil
}

func (c *Config) UpsertRelayerAccount(acc conf.Account) {
	var found bool
	for i, account := range c.Relayer.Accounts {
		if account.Name == acc.Name {
			found = true
			c.Relayer.Accounts[i] = acc
			break
		}
	}
	if !found {
		c.Relayer.Accounts = append(c.Relayer.Accounts, acc)
	}
}

// Account holds the options related to setting up Cosmos wallets.
type Relayer struct {
	// Accounts of remote chains.
	Accounts []conf.Account `yaml:"accounts"`
}

// Parse parses config.yml into Config.
func Parse(r io.Reader) (*Config, error) {
	var conf Config
	return &conf, yaml.NewDecoder(r).Decode(&conf)
}

func Open(path string) (*Config, error) {
	file, err := os.Open(filepath.Join(path, SecretFile))
	if err != nil {
		return &Config{}, nil
	}
	defer file.Close()
	return Parse(file)
}

func Save(path string, conf *Config) error {
	file, err := os.OpenFile(filepath.Join(path, SecretFile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(conf)
}
