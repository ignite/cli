package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/services/chain/conf"
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

	var steps step.Steps

	// cleanup persistent data from previous `serve`.
	steps.Add(step.New(
		step.PreExec(func() error {
			for _, path := range c.plugin.StoragePaths() {
				if err := os.RemoveAll(path); err != nil {
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
				c.app.D(),
				"init",
				"mynode",
				"--chain-id", chainID,
			),
			// overwrite configuration changes from Starport's config.yml to
			// over app's sdk configs.
			step.PostExec(func(err error) error {
				if err != nil {
					return err
				}

				// make sure that chain id given during chain.New() has the most priority.
				if conf.Genesis != nil {
					conf.Genesis["chain_id"] = chainID
				}

				appconfigs := []struct {
					ec      confile.EncodingCreator
					path    string
					changes map[string]interface{}
				}{
					{confile.DefaultJSONEncodingCreator, c.GenesisPath(), conf.Genesis},
					{confile.DefaultTOMLEncodingCreator, c.AppTOMLPath(), conf.Init.App},
					{confile.DefaultTOMLEncodingCreator, c.ConfigTOMLPath(), conf.Init.Config},
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
				return nil
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

func (c *Chain) setupSteps(_ context.Context, _ conf.Config) (steps step.Steps, err error) {
	if err := c.checkIBCRelayerSupport(); err == nil {
		steps.Add(step.New(
			step.PreExec(func() error {
				if err := xos.RemoveAllUnderHome(".relayer"); err != nil {
					return err
				}
				info, err := c.RelayerInfo()
				if err != nil {
					return err
				}
				fmt.Fprintf(c.stdLog(logStarport).out, "✨ Relayer info: %s\n", info)
				return nil
			}),
		))
	}

	chainID, err := c.ID()
	if err != nil {
		return nil, err
	}

	for _, execOption := range c.plugin.ConfigCommands(chainID) {
		execOption := execOption
		steps.Add(step.New(step.NewOptions().
			Add(execOption).
			Add(c.stdSteps(logAppcli)...)...,
		))
	}

	return steps, nil
}

// CreateAccount creates an account on chain.
// cmnemonic is returned when account is created but not restored.
func (c *Chain) CreateAccount(ctx context.Context, name, mnemonic string, isSilent bool) (Account, error) {
	acc := Account{
		Name: name,
	}

	var steps step.Steps

	if mnemonic != "" {
		steps.Add(
			step.New(
				step.NewOptions().
					Add(c.plugin.ImportUserCommand(name, mnemonic)...)...,
			),
		)
	} else {
		generatedMnemonic := &bytes.Buffer{}
		steps.Add(
			step.New(
				step.NewOptions().
					Add(c.plugin.AddUserCommand(name)...).
					Add(
						step.PostExec(func(exitErr error) error {
							if exitErr != nil {
								return errors.Wrapf(exitErr, "cannot create %s account", name)
							}
							if err := json.NewDecoder(generatedMnemonic).Decode(&acc); err != nil {
								return errors.Wrap(err, "cannot decode mnemonic")
							}
							if !isSilent {
								fmt.Fprintf(c.stdLog(logStarport).out, "🙂 Created an account. Password (mnemonic): %[1]v\n", acc.Mnemonic)
							}
							return nil
						}),
					).
					Add(c.stdSteps(logAppcli)...).
					// Stargate pipes from stdout, Launchpad pipes from stderr.
					Add(step.Stderr(generatedMnemonic), step.Stdout(generatedMnemonic))...,
			),
		)
	}

	key := &bytes.Buffer{}

	steps.Add(
		step.New(step.NewOptions().
			Add(
				c.plugin.ShowAccountCommand(name),
				step.PostExec(func(err error) error {
					if err != nil {
						return err
					}
					acc.Address = strings.TrimSpace(key.String())
					return nil
				}),
			).
			Add(c.stdSteps(logAppcli)...).
			Add(step.Stdout(key))...,
		),
	)

	if err := cmdrunner.New(c.cmdOptions()...).Run(ctx, steps...); err != nil {
		return Account{}, err
	}

	return acc, nil
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

var gentxRe = regexp.MustCompile(`(?m)"(.+?)"`)

// Gentx generates a gentx for v.
func (c *Chain) Gentx(ctx context.Context, v Validator) (gentxPath string, err error) {
	chainID, err := c.ID()
	if err != nil {
		return "", err
	}

	gentxPathMessage := &bytes.Buffer{}
	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(
			c.plugin.GentxCommand(chainID, v),
			step.Stderr(io.MultiWriter(gentxPathMessage, c.stdLog(logAppd).err)),
			step.Stdout(io.MultiWriter(gentxPathMessage, c.stdLog(logAppd).out)),
		)); err != nil {
		return "", err
	}
	return gentxRe.FindStringSubmatch(gentxPathMessage.String())[1], nil
}

// Account represents an account in the chain.
type Account struct {
	Name     string
	Address  string
	Mnemonic string `json:"mnemonic"`
	Coins    string
}

// AddGenesisAccount add a genesis account in the chain.
func (c *Chain) AddGenesisAccount(ctx context.Context, account Account) error {
	errb := &bytes.Buffer{}

	return cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(step.NewOptions().
			Add(
				step.Exec(
					c.app.D(),
					"add-genesis-account",
					account.Address,
					account.Coins,
					"--home",
					c.Home(),
				),
				step.Stderr(errb),
			)...,
		))
}

// CollectGentx collects gentxs on chain.
func (c *Chain) CollectGentx(ctx context.Context) error {
	var errb bytes.Buffer

	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(
			step.Exec(
				c.app.D(),
				"collect-gentxs",
				"--home",
				c.Home(),
			),
			step.Stderr(io.MultiWriter(c.stdLog(logAppd).err, &errb)),
			step.Stdout(c.stdLog(logAppd).out),
		)); err != nil {
		return errors.Wrap(err, errb.String())
	}

	return nil
}

// ShowNodeID shows node's id.
func (c *Chain) ShowNodeID(ctx context.Context) (string, error) {
	var key bytes.Buffer
	var errb bytes.Buffer

	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx,
			step.New(
				step.Exec(
					c.app.D(),
					"tendermint",
					"show-node-id",
				),
				step.Stdout(&key),
				step.Stderr(&errb),
			),
		); err != nil {
		return "", errors.Wrap(err, errb.String())
	}

	return strings.TrimSpace(key.String()), nil
}
