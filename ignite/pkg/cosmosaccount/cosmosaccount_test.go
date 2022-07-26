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
	require.False(t, account.Info.GetAddress().Empty())

	getAccount, err := registry.GetByName(testAccountName)
	require.NoError(t, err)
	require.True(t, getAccount.Info.GetAddress().Equals(account.Info.GetAddress()))

	secondTmpDir := t.TempDir()
	secondRegistry, err := cosmosaccount.New(cosmosaccount.WithHome(secondTmpDir))
	require.NoError(t, err)

	importedAccount, err := secondRegistry.Import(testAccountName, mnemonic, "")
	require.NoError(t, err)
	require.Equal(t, testAccountName, importedAccount.Name)
	require.True(t, importedAccount.Info.GetAddress().Equals(account.Info.GetAddress()))

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
}
