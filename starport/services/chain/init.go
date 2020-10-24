package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/imdario/mergo"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/services/chain/conf"
	secretconf "github.com/tendermint/starport/starport/services/chain/conf/secret"
)

// Init initializes chain.
func (c *Chain) Init(ctx context.Context) error {
	conf, err := c.config()
	if err != nil {
		return err
	}

	var steps step.Steps

	// cleanup persistent data from previous `serve`.
	steps.Add(step.New(
		step.PreExec(func() error {
			for _, path := range c.plugin.StoragePaths() {
				if err := xos.RemoveAllUnderHome(path); err != nil {
					return err
				}
			}
			return nil
		}),
	))

	// init node.
	steps.Add(step.New(step.NewOptions().
		Add(
			step.Exec(
				c.app.d(),
				"init",
				"mynode",
				"--chain-id", c.app.n(),
			),
			step.PostExec(func(err error) error {
				// overwrite Genesis with user configs.
				if err != nil {
					return err
				}
				if conf.Genesis == nil {
					return nil
				}
				path, err := c.plugin.GenesisPath()
				if err != nil {
					return err
				}
				file, err := os.OpenFile(path, os.O_RDWR, 644)
				if err != nil {
					return err
				}
				defer file.Close()
				var genesis map[string]interface{}
				if err := json.NewDecoder(file).Decode(&genesis); err != nil {
					return err
				}
				if err := mergo.Merge(&genesis, conf.Genesis, mergo.WithOverride); err != nil {
					return err
				}
				if err := file.Truncate(0); err != nil {
					return err
				}
				if _, err := file.Seek(0, 0); err != nil {
					return err
				}
				return json.NewEncoder(file).Encode(&genesis)
			}),
			step.PostExec(func(err error) error {
				if err != nil {
					return err
				}
				return c.plugin.PostInit(conf)
			}),
		).
		Add(c.stdSteps(logAppd)...)...,
	))

	return cmdrunner.New(c.cmdOptions()...).Run(ctx, steps...)
}

func (s *Chain) setupSteps(ctx context.Context, conf conf.Config) (steps step.Steps, err error) {
	sconf, err := secretconf.Open(s.app.Path)
	if err != nil {
		return nil, err
	}

	for _, account := range conf.Accounts {
		steps.Add(s.createAccountSteps(ctx, account.Name, "", account.Coins, false)...)
	}

	for _, account := range sconf.Accounts {
		steps.Add(s.createAccountSteps(ctx, account.Name, account.Mnemonic, account.Coins, false)...)
	}

	if err := s.checkIBCRelayerSupport(); err == nil {
		steps.Add(step.New(
			step.PreExec(func() error {
				if err := xos.RemoveAllUnderHome(".relayer"); err != nil {
					return err
				}
				info, err := s.RelayerInfo()
				if err != nil {
					return err
				}
				fmt.Fprintf(s.stdLog(logStarport).out, "✨ Relayer info: %s\n", info)
				return nil
			}),
		))
	}

	for _, execOption := range s.plugin.ConfigCommands() {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(execOption).
			Add(s.stdSteps(logAppcli)...)...,
		))
	}

	steps.Add(step.New(step.NewOptions().
		Add(s.plugin.GentxCommand(conf)).
		Add(s.stdSteps(logAppd)...)...,
	))
	return steps, nil
}
