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
	"github.com/tendermint/starport/starport/pkg/xos"
	"github.com/tendermint/starport/starport/services/chain/conf"
)

// Init initializes chain.
func (c *Chain) Init(ctx context.Context) error {
	conf, err := c.Config()
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

	return steps, nil
}

// CreateAccount creates an account on chain.
// cmnemonic is returned when account is created but not restored.
func (s *Chain) CreateAccount(ctx context.Context, name, mnemonic string, coins []string, isSilent bool) (cmnemonic string, err error) {
	var steps step.Steps
	var user struct {
		Mnemonic string `json:"mnemonic"`
	}

	if mnemonic != "" {
		steps.Add(
			step.New(
				step.NewOptions().
					Add(s.plugin.ImportUserCommand(name, mnemonic)...)...,
			),
		)
	} else {
		generatedMnemonic := &bytes.Buffer{}
		steps.Add(
			step.New(
				step.NewOptions().
					Add(s.plugin.AddUserCommand(name)...).
					Add(
						step.PostExec(func(exitErr error) error {
							if exitErr != nil {
								return errors.Wrapf(exitErr, "cannot create %s account", name)
							}
							if err := json.NewDecoder(generatedMnemonic).Decode(&user); err != nil {
								return errors.Wrap(err, "cannot decode mnemonic")
							}
							if !isSilent {
								fmt.Fprintf(s.stdLog(logStarport).out, "🙂 Created an account. Password (mnemonic): %[1]v\n", user.Mnemonic)
							}
							return nil
						}),
					).
					Add(s.stdSteps(logAppcli)...).
					// Stargate pipes from stdout, Launchpad pipes from stderr.
					Add(step.Stderr(generatedMnemonic), step.Stdout(generatedMnemonic))...,
			),
		)
	}

	key := &bytes.Buffer{}
	steps.Add(
		step.New(step.NewOptions().
			Add(
				s.plugin.ShowAccountCommand(name),
				step.PostExec(func(err error) error {
					if err != nil {
						return err
					}
					coins := strings.Join(coins, ",")
					key := strings.TrimSpace(key.String())
					return cmdrunner.
						New().
						Run(ctx, step.New(step.NewOptions().
							Add(step.Exec(
								s.app.d(),
								"add-genesis-account",
								key,
								coins,
							)).
							Add(s.stdSteps(logAppd)...)...,
						))
				}),
			).
			Add(s.stdSteps(logAppcli)...).
			Add(step.Stdout(key))...,
		),
	)

	return user.Mnemonic, cmdrunner.New(s.cmdOptions()...).Run(ctx, steps...)
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
	gentxPathMessage := &bytes.Buffer{}
	if err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(
			c.plugin.GentxCommand(v),
			step.Stderr(io.MultiWriter(gentxPathMessage, c.stdLog(logAppd).err)),
			step.Stdout(io.MultiWriter(gentxPathMessage, c.stdLog(logAppd).out)),
		)); err != nil {
		return "", err
	}
	return gentxRe.FindStringSubmatch(gentxPathMessage.String())[1], nil
}

// CollectGentx collects gentxs on chain.
func (c *Chain) CollectGentx(ctx context.Context) error {
	return cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(step.NewOptions().
			Add(step.Exec(
				c.app.d(),
				"collect-gentxs",
			)).
			Add(c.stdSteps(logAppd)...)...,
		))
}

// ShowNodeID shows node's id.
func (c *Chain) ShowNodeID(ctx context.Context) (string, error) {
	key := &bytes.Buffer{}
	err := cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx,
			step.New(
				step.Exec(
					c.app.d(),
					"tendermint",
					"show-node-id",
				),
				step.Stdout(key),
			),
		)
	return strings.TrimSpace(key.String()), err
}
