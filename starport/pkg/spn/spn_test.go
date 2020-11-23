package spn

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/stretchr/testify/require"
)

func TestAccountCreate(t *testing.T) {
	c, err := New("", "", "", Keyring(keyring.BackendMemory))
	require.NoError(t, err, "init client")

	account, err := c.AccountCreate("spn", "")
	require.NoError(t, err, "create an account")

	_, err = types.AccAddressFromBech32(account.Address)
	require.NoError(t, err, "created account's address should be valid")
	require.True(t, bip39.IsMnemonicValid(account.Mnemonic), "created account's mnemonic should be valid")
	require.Equal(t, "spn", account.Name, "account's name should be correct")

	_, err = c.AccountCreate("spn", "")
	require.Error(t, err, "should not create an account with the same name")
}

func TestAccountGet(t *testing.T) {
	c, err := New("", "", "", Keyring(keyring.BackendMemory))
	require.NoError(t, err, "init client")

	accountcreate, err := c.AccountCreate("spn", "")
	require.NoError(t, err, "create an account")

	accountget, err := c.AccountGet("spn")
	require.NoError(t, err, "should get the account")

	require.Equal(t, accountcreate.Address, accountget.Address,
		"created account should be same with the retrieved one")
}

func TestAccountExportAndImport(t *testing.T) {
	c, err := New("", "", "", Keyring(keyring.BackendMemory))
	require.NoError(t, err, "init client")

	account, err := c.AccountCreate("spn", "")
	require.NoError(t, err, "create an account")

	privateKey, err := c.AccountExport("spn", "very-secure-password")
	require.NoError(t, err, "should export the account")

	cother, err := New("", "", "", Keyring(keyring.BackendMemory))
	require.NoError(t, err, "init a new client with empty keyring")

	err = cother.AccountImport("spn", privateKey, "very-secure-password")
	require.NoError(t, err, "should import the key to the other keyring")

	accountimported, err := cother.AccountGet("spn")
	require.NoError(t, err, "should get the imported account")

	require.Equal(t, account.Address, accountimported.Address,
		"original account should be same with the imported one")
}
