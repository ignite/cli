package cache_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cache"
)

type TestStruct struct {
	Num int
}

func TestCreateStorage(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	_, err := cache.NewStorage(filepath.Join(tmpDir1, "test.db"))
	require.NoError(t, err)

	_, err = cache.NewStorage(filepath.Join(tmpDir2, "test.db"))
	require.NoError(t, err)
}

func TestStoreString(t *testing.T) {
	tmpDir := t.TempDir()
	cacheStorage, err := cache.NewStorage(filepath.Join(tmpDir, "testdbfile.db"))
	require.NoError(t, err)

	strNamespace := cache.New[string](cacheStorage, "myNameSpace")

	err = strNamespace.Put("myKey", "myValue")
	require.NoError(t, err)

	val, err := strNamespace.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, "myValue", val)

	strNamespaceAgain := cache.New[string](cacheStorage, "myNameSpace")

	valAgain, err := strNamespaceAgain.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, "myValue", valAgain)
}

func TestStoreObjects(t *testing.T) {
	tmpDir := t.TempDir()
	cacheStorage, err := cache.NewStorage(filepath.Join(tmpDir, "testdbfile.db"))
	require.NoError(t, err)

	structCache := cache.New[TestStruct](cacheStorage, "mySimpleNamespace")

	err = structCache.Put("myKey", TestStruct{
		Num: 42,
	})
	require.NoError(t, err)

	val, err := structCache.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, val, TestStruct{
		Num: 42,
	})

	arrayNamespace := cache.New[[]TestStruct](cacheStorage, "myArrayNamespace")

	err = arrayNamespace.Put("myKey", []TestStruct{
		{
			Num: 42,
		},
		{
			Num: 420,
		},
	})
	require.NoError(t, err)

	val2, err := arrayNamespace.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, 2, len(val2))
	require.Equal(t, 42, (val2)[0].Num)
	require.Equal(t, 420, (val2)[1].Num)

	empty, err := arrayNamespace.Get("doesNotExists")
	require.Equal(t, cache.ErrorNotFound, err)
	require.Nil(t, empty)
}

func TestConflicts(t *testing.T) {
	tmpDir := t.TempDir()
	tmpDir2 := t.TempDir()
	cacheStorage1, err := cache.NewStorage(filepath.Join(tmpDir, "testdbfile.db"))
	require.NoError(t, err)
	cacheStorage2, err := cache.NewStorage(filepath.Join(tmpDir2, "testdbfile.db"))
	require.NoError(t, err)

	sameStorageDifferentNamespaceCache1 := cache.New[int](cacheStorage1, "ns1")

	sameStorageDifferentNamespaceCache2 := cache.New[int](cacheStorage1, "ns2")

	differentStorageSameNamespace := cache.New[int](cacheStorage2, "ns1")

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
	val1, err := sameStorageDifferentNamespaceCache1.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, 42, val1)

	val2, err := sameStorageDifferentNamespaceCache2.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, 1337, val2)

	val3, err := differentStorageSameNamespace.Get("myKey")
	require.NoError(t, err)
	require.Equal(t, 9001, val3)
}

func TestDeleteKey(t *testing.T) {
	tmpDir := t.TempDir()
	cacheStorage, err := cache.NewStorage(filepath.Join(tmpDir, "testdbfile.db"))
	require.NoError(t, err)

	strNamespace := cache.New[string](cacheStorage, "myNameSpace")
	err = strNamespace.Put("myKey", "someValue")
	require.NoError(t, err)

	err = strNamespace.Delete("myKey")
	require.NoError(t, err)

	_, err = strNamespace.Get("myKey")
	require.Equal(t, cache.ErrorNotFound, err)
}

func TestClearStorage(t *testing.T) {
	tmpDir := t.TempDir()
	cacheStorage, err := cache.NewStorage(filepath.Join(tmpDir, "testdbfile.db"))
	require.NoError(t, err)

	strNamespace := cache.New[string](cacheStorage, "myNameSpace")

	err = strNamespace.Put("myKey", "myValue")
	require.NoError(t, err)

	err = cacheStorage.Clear()
	require.NoError(t, err)

	_, err = strNamespace.Get("myKey")
	require.Equal(t, cache.ErrorNotFound, err)
}

func TestKey(t *testing.T) {
	singleKey := cache.Key("test1")
	require.Equal(t, "test1", singleKey)

	multiKey := cache.Key("test1", "test2", "test3")
	require.Equal(t, "test1test2test3", multiKey)
}
