package chain

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"io"
	"regexp"
	"strings"
)

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

// AddGenesisAccount add a genesis account in the chain
func (c *Chain) AddGenesisAccount(ctx context.Context, account spn.GenesisAccount) error {
	return cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(step.NewOptions().
			Add(step.Exec(
				c.app.D(),
				"add-genesis-account",
				account.Address.String(),
				account.Coins.String(),
			)).
			Add(c.stdSteps(logAppd)...)...,
		))
}

// CollectGentx collects gentxs on chain.
func (c *Chain) CollectGentx(ctx context.Context) error {
	return cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx, step.New(step.NewOptions().
			Add(step.Exec(
				c.app.D(),
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
					c.app.D(),
					"tendermint",
					"show-node-id",
				),
				step.Stdout(key),
			),
		)
	return strings.TrimSpace(key.String()), err
}

// GetInitialGenesis gets the initial genesis of the chain
// Currently the only way to get the initial genesis of a chain through its CLI is by running the init command
// Therefore, we run the init command in a temporary directory and then delete it
func (c *Chain) GetInitialGenesis(ctx context.Context, chainID string) ([]byte, error) {
	// Generate the temporary dir
	tmpDir, err := ioutil.TempDir("/tmp", "prefix")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Init command
	err = cmdrunner.
		New(c.cmdOptions()...).
		Run(ctx,
			step.New(
				step.Exec(
					c.app.D(),
					"init",
					"moniker",
					"--chain-id",
					chainID,
					"--home",
					tmpDir,
				),
			),
		)
	if err != nil {
		return nil, err
	}

	// Read and return the genesis file
	return ioutil.ReadFile(filepath.Join(tmpDir, "config/genesis.json"))
}