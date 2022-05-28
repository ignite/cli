package v0

import "github.com/ignite/cli/ignite/chainconfig/common"

// ConfigYaml is the user given configuration to do additional setup
// during serve.
type ConfigYaml struct {
	Accounts  []common.Account       `yaml:"accounts"`
	Validator common.Validator       `yaml:"validator"`
	Faucet    common.Faucet          `yaml:"faucet"`
	Client    common.Client          `yaml:"client"`
	Build     common.Build           `yaml:"build"`
	Init      common.Init            `yaml:"init"`
	Genesis   map[string]interface{} `yaml:"genesis"`
	Host      common.Host            `yaml:"host"`
}

// AccountByName finds account by name.
func (c ConfigYaml) AccountByName(name string) (acc common.Account, found bool) {
	for _, acc := range c.Accounts {
		if acc.Name == name {
			return acc, true
		}
	}
	return common.Account{}, false
}
