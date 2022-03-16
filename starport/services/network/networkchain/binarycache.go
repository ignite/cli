package networkchain

import (
	"strconv"

	"github.com/tendermint/starport/starport/chainconfig"
	"github.com/tendermint/starport/starport/pkg/checksum"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/xfilepath"
)

const (
	BinaryCacheDirectory = "binary-cache"
	BinaryCacheFilename  = "list.yml"
)

type List struct {
	CachedBinary map[string]string `json:"cached_binary" yaml:"cached_binary"`
}

// CacheBinaryForLaunchID caches hash sha256(sha256(binary) + sourcehash) for launch id
func CacheBinaryForLaunchID(launchID uint64, binaryHash, sourceHash string) error {
	cachePath, err := getBinaryCacheFilepath()
	if err != nil {
		return err
	}
	cacheList := List{CachedBinary: map[string]string{}}
	err = confile.New(confile.DefaultYAMLEncodingCreator, cachePath).Load(&cacheList)
	if err != nil {
		return err
	}
	cacheList.CachedBinary[strconv.Itoa(int(launchID))] = checksum.SHA256Checksum([]byte(binaryHash), []byte(sourceHash))

	return confile.New(confile.DefaultYAMLEncodingCreator, cachePath).Save(cacheList)
}

// CheckBinaryCacheForLaunchID checks if binary for the given launch was already built
func CheckBinaryCacheForLaunchID(launchID uint64, binaryHash, sourceHash string) (bool, error) {
	cachePath, err := getBinaryCacheFilepath()
	if err != nil {
		return false, err
	}
	cacheList := List{CachedBinary: map[string]string{}}
	err = confile.New(confile.DefaultYAMLEncodingCreator, cachePath).Load(&cacheList)
	if err != nil {
		return false, err
	}
	return cacheList.CachedBinary[strconv.Itoa(int(launchID))] == checksum.SHA256Checksum([]byte(binaryHash), []byte(sourceHash)), nil
}

func getBinaryCacheFilepath() (string, error) {
	return xfilepath.Join(
		chainconfig.ConfigDirPath,
		xfilepath.Path("spn"),
		xfilepath.Path(BinaryCacheDirectory),
		xfilepath.Path(BinaryCacheFilename),
	)()
}
