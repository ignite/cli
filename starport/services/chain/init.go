package chain

import (
	"context"
	"fmt"
	"os"

	"github.com/imdario/mergo"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/xos"
)

const (
	moniker = "mynode"
)

// Init initializes chain.
func (c *Chain) Init(ctx context.Context) error {
	chainID, err := c.ID()
	if err != nil {
		return err
	}

	conf, err := c.Config()
	if err != nil {
		return err
	}

	// cleanup persistent data from previous `serve`.
	for _, path := range c.plugin.StoragePaths() {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	// init node.
	if err := c.cmd.Init(ctx, moniker); err != nil {
		return err
	}

	// overwrite configuration changes from Starport's config.yml to
	// over app's sdk configs.

	// make sure that chain id given during chain.New() has the most priority.
	if conf.Genesis != nil {
		conf.Genesis["chain_id"] = chainID
	}

	// Initilize app config
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}
	appTOMLPath, err := c.AppTOMLPath()
	if err != nil {
		return err
	}
	configTOMLPath, err := c.ConfigTOMLPath()
	if err != nil {
		return err
	}

	appconfigs := []struct {
		ec      confile.EncodingCreator
		path    string
		changes map[string]interface{}
	}{
		{confile.DefaultJSONEncodingCreator, genesisPath, conf.Genesis},
		{confile.DefaultTOMLEncodingCreator, appTOMLPath, conf.Init.App},
		{confile.DefaultTOMLEncodingCreator, configTOMLPath, conf.Init.Config},
	}

	for _, ac := range appconfigs {
		cf := confile.New(ac.ec, ac.path)
		var conf map[string]interface{}
		if err := cf.Load(&conf); err != nil {
			return err
		}
		if err := mergo.Merge(&conf, ac.changes, mergo.WithOverride); err != nil {
			return err
		}
		if err := cf.Save(conf); err != nil {
			return err
		}
	}

	// run post init handler
	return c.plugin.PostInit(conf)
}

func (c *Chain) configure(ctx context.Context) error {
	// setup IBC Relayer.
	if err := c.checkIBCRelayerSupport(); err == nil {
		if err := xos.RemoveAllUnderHome(".relayer"); err != nil {
			return err
		}
		info, err := c.RelayerInfo()
		if err != nil {
			return err
		}
		fmt.Fprintf(c.stdLog(logStarport).out, "✨ Relayer info: %s\n", info)
		return nil
	}

	// configure blockchain.
	chainID, err := c.ID()
	if err != nil {
		return err
	}

	return c.plugin.Configure(ctx, c.Commands(), chainID)
}

type Validator struct {
	Name                    string
	Moniker                 string
	StakingAmount           string
	CommissionRate          string
	CommissionMaxRate       string
	CommissionMaxChangeRate string
	MinSelfDelegation       string
	GasPrices               string
}

// Account represents an account in the chain.
type Account struct {
	Name     string
	Address  string
	Mnemonic string `json:"mnemonic"`
	Coins    string
}
