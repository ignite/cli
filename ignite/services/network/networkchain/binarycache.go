package networkchain

import (
	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/pkg/checksum"
	"github.com/ignite/cli/ignite/pkg/confile"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

const (
	SPNCacheDirectory    = "spn"
	BinaryCacheDirectory = "binary-cache"
	BinaryCacheFilename  = "checksums.yml"
)

type BinaryCacheList struct {
	CachedBinaries []Binary `yaml:"cached_binaries"`
}

// Binary associates launch id with build hash where build hash is sha256(binary, source)
type Binary struct {
	LaunchID  uint64
	BuildHash string
}

func (l *BinaryCacheList) Set(launchID uint64, buildHash string) {
	for i, binary := range l.CachedBinaries {
		if binary.LaunchID == launchID {
			l.CachedBinaries[i].BuildHash = buildHash
			return
		}
	}
	l.CachedBinaries = append(l.CachedBinaries, Binary{
		LaunchID:  launchID,
		BuildHash: buildHash,
	})
}

func (l *BinaryCacheList) Get(launchID uint64) (string, bool) {
	for _, binary := range l.CachedBinaries {
		if binary.LaunchID == launchID {
			return binary.BuildHash, true
		}
	}
	return "", false
}

// cacheBinaryForLaunchID caches hash sha256(sha256(binary) + sourcehash) for launch id
func cacheBinaryForLaunchID(launchID uint64, binaryHash, sourceHash string) error {
	cachePath, err := getBinaryCacheFilepath()
	if err != nil {
		return err
	}
	var cacheList = BinaryCacheList{}
	err = confile.New(confile.DefaultYAMLEncodingCreator, cachePath).Load(&cacheList)
	if err != nil {
		return err
	}
	cacheList.Set(launchID, checksum.Strings(binaryHash, sourceHash))

	return confile.New(confile.DefaultYAMLEncodingCreator, cachePath).Save(cacheList)
}

// checkBinaryCacheForLaunchID checks if binary for the given launch was already built
func checkBinaryCacheForLaunchID(launchID uint64, binaryHash, sourceHash string) (bool, error) {
	cachePath, err := getBinaryCacheFilepath()
	if err != nil {
		return false, err
	}
	var cacheList = BinaryCacheList{}
	err = confile.New(confile.DefaultYAMLEncodingCreator, cachePath).Load(&cacheList)
	if err != nil {
		return false, err
	}
	buildHash, ok := cacheList.Get(launchID)
	return ok && buildHash == checksum.Strings(binaryHash, sourceHash), nil
}

func getBinaryCacheFilepath() (string, error) {
	return xfilepath.Join(
		chainconfig.ConfigDirPath,
		xfilepath.Path(SPNCacheDirectory),
		xfilepath.Path(BinaryCacheDirectory),
		xfilepath.Path(BinaryCacheFilename),
	)()
}
