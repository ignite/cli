package networkchain

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/randstr"
	"github.com/tendermint/starport/starport/services/chain"
)

const (
	passphraseLength = 32
)

// InitAccount initializes an account for the blockchain and issue a gentx in config/gentx/gentx.json
func (c Chain) InitAccount(ctx context.Context, v chain.Validator, accountName string) (string, error) {
	if !c.isInitialized {
		return "", errors.New("the blockchain must be initialized to initialize an account")
	}

	chainCmd, err := c.chain.Commands(ctx)
	if err != nil {
		return "", err
	}

	// create the chain account.
	address, err := c.ImportAccount(ctx, accountName)
	if err != nil {
		return "", err
	}

	// add account into the genesis
	err = chainCmd.AddGenesisAccount(ctx, address, v.StakingAmount)
	if err != nil {
		return "", err
	}

	// create the gentx.
	issuedGentxPath, err := c.chain.IssueGentx(ctx, v)
	if err != nil {
		return "", err
	}

	// rename the issued gentx into gentx.json
	gentxPath := filepath.Join(filepath.Dir(issuedGentxPath), cosmosutil.GentxFilename)
	return gentxPath, os.Rename(issuedGentxPath, gentxPath)
}

// ImportAccount imports an account from Starport into the chain.
// we first export the account into a temporary key file and import it with the chain CLI.
func (c *Chain) ImportAccount(ctx context.Context, name string) (string, error) {
	// keys import command of chain CLI requires that the key file is encrypted with a passphrase of at least 8 characters
	// we generate a random passphrase to import the account
	passphrase := randstr.Runes(passphraseLength)

	// export the key in a temporary file.
	armored, err := c.ar.Export(name, passphrase)
	if err != nil {
		return "", err
	}

	keyFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(keyFile.Name())

	if _, err := keyFile.Write([]byte(armored)); err != nil {
		return "", err
	}

	// import the key file into the chain.
	chainCmd, err := c.chain.Commands(ctx)
	if err != nil {
		return "", err
	}

	acc, err := chainCmd.ImportAccount(ctx, name, keyFile.Name(), passphrase)
	return acc.Address, err
}
