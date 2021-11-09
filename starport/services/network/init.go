package network

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/chain"
)

// Init initializes blockchain by building the binaries and running the init command and
// create the initial genesis of the chain
func (b *Blockchain) Init(ctx context.Context) error {
	chainHome, err := b.chain.Home()
	if err != nil {
		return err
	}

	// cleanup home dir of app if exists.
	if err := os.RemoveAll(chainHome); err != nil {
		return err
	}

	// build the chain and initialize it with a new validator key
	b.builder.ev.Send(events.New(events.StatusOngoing, "Building the blockchain"))
	if _, err := b.chain.Build(ctx, ""); err != nil {
		return err
	}
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain built"))
	b.builder.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))
	if err := b.chain.Init(ctx, false); err != nil {
		return err
	}
	b.builder.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))

	// initialize and verify the genesis
	b.initGenesis(ctx)

	b.isInitialized = true

	return nil
}

// InitAccount initializes an account for the blockchain and issue a gentx in config/gentx/gentx.json
func (b *Blockchain) InitAccount(ctx context.Context, v chain.Validator, accountName string) (string, error) {
	if !b.isInitialized {
		return "", errors.New("the blockchain must be initialized to initialize an account")
	}

	// create the chain account
	chainCmd, err := b.chain.Commands(ctx)
	if err != nil {
		return "", err
	}
	address, err := b.ImportAccount(ctx, accountName)
	if err != nil {
		return "", err
	}

	// add account into the genesis
	err = chainCmd.AddGenesisAccount(ctx, address, v.StakingAmount)
	if err != nil {
		return "", err
	}

	// create the gentx
	issuedGentxPath, err := b.chain.IssueGentx(ctx, v)
	if err != nil {
		return "", err
	}

	// rename the issued gentx into gentx.json
	gentxPath := filepath.Join(filepath.Dir(issuedGentxPath), gentxFilename)
	return gentxPath, os.Rename(issuedGentxPath, gentxPath)
}

// initGenesis creates the initial genesis of the genesis depending on the initial genesis type (default, url, ...)
func (b *Blockchain) initGenesis(ctx context.Context) error {
	// if the blockchain has a genesis URL, the initial genesis is fetched from the url
	// otherwise, default genesis is used, which requires no action since the default genesis is generated from the init command
	if b.genesisURL != "" {
		genesis, hash, err := genesisAndHashFromURL(ctx, b.genesisURL)
		if err != nil {
			return err
		}

		// if the blockchain has been initialized with no genesis hash, we assign the fetched hash to it
		// otherwise we check the genesis integrity with the existing hash
		if b.genesisHash != "" {
			b.genesisHash = hash
		} else if hash != b.genesisHash {
			return fmt.Errorf("genesis from URL %s is invalid. Expected hash %s, actual hash %s", b.genesisURL, b.genesisHash, hash)
		}

		// replace the default genesis with the fetched genesis
		genesisPath, err := b.chain.GenesisPath()
		if err != nil {
			return err
		}
		if err := os.WriteFile(genesisPath, genesis, 0644); err != nil {
			return err
		}
	}

	// check the genesis is valid
	return b.checkGenesis(ctx)
}

// ImportAccount imports an account from Starport into the chain
// we first export the account into a temporary key file and import it with the chain CLI
func (b *Blockchain) ImportAccount(ctx context.Context, name string) (string, error) {
	// keys import command of chain CLI requires that the key file is encrypted with a passphrase of at least 8 characters
	// we generate a random passphrase to import the account
	passphrase := randomPassphrase()

	// export the key in a temporary file
	armored, err := b.builder.cosmos.AccountRegistry.Export(name, passphrase)
	if err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	keyFile := filepath.Join(tmpDir, name)
	if err := os.WriteFile(keyFile, []byte(armored), 0644); err != nil {
		return "", err
	}

	// import the key file into the chain
	chainCmd, err := b.chain.Commands(ctx)
	if err != nil {
		return "", err
	}
	acc, err := chainCmd.ImportAccount(ctx, name, keyFile, passphrase)
	return acc.Address, err
}

// checkGenesis checks the stored genesis is valid
func (b *Blockchain) checkGenesis(ctx context.Context) error {
	// Perform static analysis of the chain with the validate-genesis command
	commands, err := b.chain.Commands(ctx)
	if err != nil {
		return err
	}
	return commands.ValidateGenesis(ctx)

	// TODO: static analysis of the genesis with validate-genesis doesn't check the full validity of the genesis
	// example: gentxs formats are not checked
	// to perform a full validity check of the genesis we must try to start the chain with sample accounts
}

// genesisAndHashFromURL fetches the genesis from the given url and returns its content along with the sha256 hash
func genesisAndHashFromURL(ctx context.Context, url string) (genesis []byte, hash string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	genesis, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, bytes.NewReader(genesis)); err != nil {
		return nil, "", err
	}

	hexHash := hex.EncodeToString(h.Sum(nil))

	return genesis, hexHash, nil
}

// randomPassphrase generates a random passphrase of 32 characters from account export and import
func randomPassphrase() string {
	var passphraseRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	passBytes := make([]rune, 32)
	for i := range passBytes {
		passBytes[i] = passphraseRunes[rand.Intn(len(passphraseRunes))]
	}
	return string(passBytes)
}
