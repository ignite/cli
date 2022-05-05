package cache

import (
	"github.com/ignite-hq/cli/ignite/pkg/xfilepath"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestStruct struct {
	Num int
}

func TestCreateStorage(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()
	cacheFolder = xfilepath.Path(tmpDir1) // Overriding $HOME/.ignite folder

	_, err := NewChainStorage("myChain")
	require.NoError(t, err)

	_, err = NewStorage(tmpDir2)
	require.NoError(t, err)
}

func TestStoreString(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFolder = xfilepath.Path(tmpDir) // Overriding $HOME/.ignite folder
	cacheStorage, err := NewChainStorage("myChain")
	require.NoError(t, err)

	strNamespace := New[string](cacheStorage, "myNameSpace")

	err = strNamespace.Put("myKey", "myValue")
	require.NoError(t, err)

	val, found, err := strNamespace.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "myValue", val)

	strNamespaceAgain := New[string](cacheStorage, "myNameSpace")

	valAgain, found, err := strNamespaceAgain.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "myValue", valAgain)
}

func TestStoreObjects(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFolder = xfilepath.Path(tmpDir) // Overriding $HOME/.ignite folder
	cacheStorage, err := NewChainStorage("myChain")
	require.NoError(t, err)

	cache := New[TestStruct](cacheStorage, "mySimpleNamespace")

	err = cache.Put("myKey", TestStruct{
		Num: 42,
	})
	require.NoError(t, err)

	val, found, err := cache.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, val, TestStruct{
		Num: 42,
	})

	arrayNamespace := New[[]TestStruct](cacheStorage, "myArrayNamespace")

	err = arrayNamespace.Put("myKey", []TestStruct{
		{
			Num: 42,
		},
		{
			Num: 420,
		},
	})
	require.NoError(t, err)

	val2, found, err := arrayNamespace.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, 2, len(val2))
	require.Equal(t, 42, (val2)[0].Num)
	require.Equal(t, 420, (val2)[1].Num)

	empty, found, err := arrayNamespace.Get("doesNotExists")
	require.NoError(t, err)
	require.False(t, found)
	require.Nil(t, empty)
}

func TestConflicts(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFolder = xfilepath.Path(tmpDir) // Overriding $HOME/.ignite folder
	cacheStorage1, err := NewChainStorage("myChain")
	require.NoError(t, err)
	cacheStorage2, err := NewChainStorage("myChain2")
	require.NoError(t, err)

	sameStorageDifferentNamespaceCache1 := New[int](cacheStorage1, "ns1")

	sameStorageDifferentNamespaceCache2 := New[int](cacheStorage1, "ns2")

	differentStorageSameNamespace := New[int](cacheStorage2, "ns1")

	// Put values in caches
	err = sameStorageDifferentNamespaceCache1.Put("myKey", 41)
	require.NoError(t, err)

	err = sameStorageDifferentNamespaceCache2.Put("myKey", 1337)
	require.NoError(t, err)

	err = differentStorageSameNamespace.Put("myKey", 9001)
	require.NoError(t, err)

	// Overwrite a value
	err = sameStorageDifferentNamespaceCache1.Put("myKey", 42)
	require.NoError(t, err)

	// Check that everything comes back as expected
	val1, found, err := sameStorageDifferentNamespaceCache1.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, 42, val1)

	val2, found, err := sameStorageDifferentNamespaceCache2.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, 1337, val2)

	val3, found, err := differentStorageSameNamespace.Get("myKey")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, 9001, val3)
}

func TestDeleteKey(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFolder = xfilepath.Path(tmpDir) // Overriding $HOME/.ignite folder
	cacheStorage, err := NewChainStorage("myChain")
	require.NoError(t, err)

	strNamespace := New[string](cacheStorage, "myNameSpace")
	err = strNamespace.Put("myKey", "someValue")
	require.NoError(t, err)

	err = strNamespace.Delete("myKey")
	require.NoError(t, err)

	_, found, err := strNamespace.Get("myKey")
	require.NoError(t, err)
	require.False(t, found)
}

func TestClearStorage(t *testing.T) {
	tmpDir := t.TempDir()
	cacheFolder = xfilepath.Path(tmpDir) // Overriding $HOME/.ignite folder
	cacheStorage, err := NewChainStorage("myChain")
	require.NoError(t, err)

	strNamespace := New[string](cacheStorage, "myNameSpace")

	err = strNamespace.Put("myKey", "myValue")
	require.NoError(t, err)

	err = cacheStorage.Clear()
	require.NoError(t, err)

	_, found, err := strNamespace.Get("myKey")
	require.NoError(t, err)
	require.False(t, found)
}
