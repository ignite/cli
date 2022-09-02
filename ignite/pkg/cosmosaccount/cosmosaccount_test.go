package cosmosaccount_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
)

const testAccountName = "myTestAccount"

func TestRegistry(t *testing.T) {
	tmpDir := t.TempDir()
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(tmpDir))
	require.NoError(t, err)

	account, mnemonic, err := registry.Create(testAccountName)
	require.NoError(t, err)
	require.Equal(t, testAccountName, account.Name)
	require.NotEmpty(t, account.Record.PubKey.Value)

	getAccount, err := registry.GetByName(testAccountName)
	require.NoError(t, err)
	require.Equal(t, getAccount, account)

	sdkaddr, _ := account.Record.GetAddress()
	addr := sdkaddr.String()
	getAccount, err = registry.GetByAddress(addr)
	require.NoError(t, err)
	require.Equal(t, getAccount.Record.PubKey, account.Record.PubKey)
	require.Equal(t, getAccount.Name, testAccountName)
	require.Equal(t, getAccount.Name, account.Name)
	require.Equal(t, getAccount.Name, account.Record.Name)

	addr, err = account.Address("cosmos")
	require.NoError(t, err)
	getAccount, err = registry.GetByAddress(addr)
	require.NoError(t, err)
	require.Equal(t, getAccount.Record.PubKey, account.Record.PubKey)
	require.Equal(t, getAccount.Name, testAccountName)
	require.Equal(t, getAccount.Name, account.Name)
	require.Equal(t, getAccount.Name, account.Record.Name)

	secondTmpDir := t.TempDir()
	secondRegistry, err := cosmosaccount.New(cosmosaccount.WithHome(secondTmpDir))
	require.NoError(t, err)

	importedAccount, err := secondRegistry.Import(testAccountName, mnemonic, "")
	require.NoError(t, err)
	require.Equal(t, testAccountName, importedAccount.Name)
	require.Equal(t, importedAccount.Record.PubKey, account.Record.PubKey)

	_, _, err = registry.Create("another one")
	require.NoError(t, err)
	list, err := registry.List()
	require.NoError(t, err)
	require.Equal(t, 2, len(list))

	err = registry.DeleteByName(testAccountName)
	require.NoError(t, err)
	afterDeleteList, err := registry.List()
	require.NoError(t, err)
	require.Equal(t, 1, len(afterDeleteList))

	_, err = registry.GetByName(testAccountName)
	var expectedErr *cosmosaccount.AccountDoesNotExistError
	require.ErrorAs(t, err, &expectedErr)

	_, err = registry.GetByAddress(addr)
	require.ErrorAs(t, err, &expectedErr)
}
