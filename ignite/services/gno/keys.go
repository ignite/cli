// Package gno provides ignite's gno.land integration: dev chain serving,
// realm/package scaffolding, deployment and key management backed by the
// gno keybase.
package gno

import (
	"fmt"
	"os"

	"github.com/gnolang/gno/gnovm/pkg/gnoenv"
	"github.com/gnolang/gno/tm2/pkg/crypto/bip39"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys/armor"
)

// HomeDir returns the default gno home directory (~/.config/gno or $GNOHOME).
func HomeDir() string { return gnoenv.HomeDir() }

// openKeybase returns the gno keybase stored under the default gno home.
func openKeybase() (keys.Keybase, error) {
	return keys.NewKeyBaseFromDir(HomeDir())
}

// KeyInfo describes a key stored in the gno keybase.
type KeyInfo struct {
	Name    string
	Address string
	Type    string
}

func toKeyInfo(info keys.Info) KeyInfo {
	return KeyInfo{
		Name:    info.GetName(),
		Address: info.GetAddress().String(),
		Type:    info.GetType().String(),
	}
}

// CreateKey creates a new key in the gno keybase. An empty passphrase creates
// an unencrypted key (dev friendly). An empty mnemonic generates a new one,
// which is returned so it can be shown once to the user.
func CreateKey(name, mnemonic, passphrase string, account, index uint32) (KeyInfo, string, error) {
	kb, err := openKeybase()
	if err != nil {
		return KeyInfo{}, "", fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	if mnemonic == "" {
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			return KeyInfo{}, "", fmt.Errorf("generating entropy: %w", err)
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return KeyInfo{}, "", fmt.Errorf("generating mnemonic: %w", err)
		}
	}

	info, err := kb.CreateAccount(name, mnemonic, "", passphrase, account, index)
	if err != nil {
		return KeyInfo{}, "", fmt.Errorf("creating account: %w", err)
	}
	return toKeyInfo(info), mnemonic, nil
}

// RecoverKey creates a key from an existing mnemonic.
func RecoverKey(name, mnemonic, passphrase string, account, index uint32) (KeyInfo, error) {
	kb, err := openKeybase()
	if err != nil {
		return KeyInfo{}, fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	info, err := kb.CreateAccount(name, mnemonic, "", passphrase, account, index)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("recovering account: %w", err)
	}
	return toKeyInfo(info), nil
}

// ShowKey returns the info of a key by name or bech32 address.
func ShowKey(nameOrAddress string) (KeyInfo, error) {
	kb, err := openKeybase()
	if err != nil {
		return KeyInfo{}, fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	info, err := kb.GetByNameOrAddress(nameOrAddress)
	if err != nil {
		return KeyInfo{}, err
	}
	return toKeyInfo(info), nil
}

// ListKeys returns all keys in the gno keybase.
func ListKeys() ([]KeyInfo, error) {
	kb, err := openKeybase()
	if err != nil {
		return nil, fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	infos, err := kb.List()
	if err != nil {
		return nil, err
	}
	list := make([]KeyInfo, 0, len(infos))
	for _, info := range infos {
		list = append(list, toKeyInfo(info))
	}
	return list, nil
}

// DeleteKey removes a key from the keybase. An empty passphrase deletes
// without prompting (matches keys created with an empty passphrase).
func DeleteKey(name, passphrase string) error {
	kb, err := openKeybase()
	if err != nil {
		return fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	return kb.Delete(name, passphrase, passphrase == "")
}

// ExportKey returns an armored private key for name, or writes it to
// outputPath when provided.
func ExportKey(name, passphrase, outputPath string) (armorOut string, err error) {
	kb, err := openKeybase()
	if err != nil {
		return "", fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	privKey, err := kb.ExportPrivKey(name, passphrase)
	if err != nil {
		return "", fmt.Errorf("exporting key: %w", err)
	}
	armorOut = armor.EncryptArmorPrivKey(privKey, passphrase)

	if outputPath == "" {
		return armorOut, nil
	}
	if err := os.WriteFile(outputPath, []byte(armorOut), 0o600); err != nil {
		return "", err
	}
	return "", nil
}

// ImportKey imports an armored private key file into the keybase.
func ImportKey(name, armorPath, passphrase string) (KeyInfo, error) {
	kb, err := openKeybase()
	if err != nil {
		return KeyInfo{}, fmt.Errorf("opening keybase: %w", err)
	}
	defer kb.CloseDB()

	armorBytes, err := os.ReadFile(armorPath)
	if err != nil {
		return KeyInfo{}, err
	}
	privKey, err := armor.UnarmorDecryptPrivKey(string(armorBytes), passphrase)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("decrypting armor: %w", err)
	}
	if err := kb.ImportPrivKey(name, privKey, passphrase); err != nil {
		return KeyInfo{}, fmt.Errorf("importing key: %w", err)
	}
	return ShowKey(name)
}
