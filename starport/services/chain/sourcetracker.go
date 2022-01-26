package chain

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	moduleanalysis "github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/dirchange"
)

const (
	moduleWildcard = "*"
)

type ArtifactType uint8

const (
	ArtifactTypeBinary = iota
	ArtifactTypeDart
	ArtifactTypeOpenAPI
	ArtifactTypeJS
	ArtifactTypeModuleProto
	// ArtifactTypeState represents app state created from genesis
	ArtifactTypeState
)

// Artifact represents a concrete building process output e.g. module: gov artifactType: ArtifactTypeDart
type Artifact struct {
	module       string
	artifactType ArtifactType
}

type Artifacts map[Artifact]struct{}

func newArtifacts() Artifacts {
	return Artifacts{}
}

func (a Artifacts) addArtifacts(artifacts Artifacts) {
	for artifact := range artifacts {
		a[artifact] = struct{}{}
	}
}

// Contains checks whether Artifacts contain specific artifact of not
func (a Artifacts) Contains(module string, artifactType ArtifactType) bool {
	_, specificArtifactExists := a[Artifact{
		module:       module,
		artifactType: artifactType,
	}]

	_, wildcardArtifactExists := a[Artifact{
		module:       "",
		artifactType: artifactType,
	}]

	return specificArtifactExists || wildcardArtifactExists
}

type TrackerType uint8

const (
	GoSourceTracker = iota
	BinaryTracker
	ModuleProtoSourceTracker
	VuexOutputTracker
	ModuleProtoGoOutputTracker
	OpenApiOutputTracker
	DartOutputTracker
	ConfigTracker
	ThirdPartyProtoTracker
	GoModTracker
)

// Tracker tracks specified file changes for specified module set. Empty module slice represents all modules at once
type Tracker struct {
	paths       []string
	modules     []string
	trackerType TrackerType
}

// trackerConfig specifies which artifacts derive from different tracker types
var trackerConfig = map[TrackerType][]ArtifactType{
	GoSourceTracker:            {ArtifactTypeBinary},
	BinaryTracker:              {ArtifactTypeBinary},
	ModuleProtoSourceTracker:   {ArtifactTypeModuleProto, ArtifactTypeDart, ArtifactTypeJS, ArtifactTypeOpenAPI, ArtifactTypeBinary},
	VuexOutputTracker:          {ArtifactTypeJS},
	ModuleProtoGoOutputTracker: {ArtifactTypeModuleProto, ArtifactTypeBinary},
	OpenApiOutputTracker:       {ArtifactTypeOpenAPI},
	DartOutputTracker:          {ArtifactTypeDart},
	ConfigTracker:              {ArtifactTypeState},
	ThirdPartyProtoTracker:     {ArtifactTypeModuleProto, ArtifactTypeDart, ArtifactTypeJS, ArtifactTypeOpenAPI, ArtifactTypeBinary},
	GoModTracker:               {ArtifactTypeModuleProto, ArtifactTypeDart, ArtifactTypeJS, ArtifactTypeOpenAPI, ArtifactTypeBinary},
}

func (t Tracker) getTargetArtifacts() Artifacts {
	// using trackerConfig gets list of artifacts derived from the tracker
	return nil
}

// isChanged checks either tracked file were changed since last persist or not
func (t Tracker) isChanged(workdir, savePath string) (bool, error) {
	return dirchange.HasDirChecksumChanged(workdir, t.paths, savePath, t.checksumFilename())
}

// persistChecksum saves checksum of tracked files
func (t Tracker) persistChecksum(workdir, savePath string) error {
	return dirchange.SaveDirChecksum(workdir, t.paths, savePath, t.checksumFilename())
}

// checksumFilename returns filename for checksum file
func (t Tracker) checksumFilename() string {
	var moduleName string
	if len(t.modules) == 0 {
		moduleName = moduleWildcard
	} else {
		sort.Strings(t.modules)
		moduleName = strings.Join(t.modules, "|")
	}
	var trackerName string
	switch t.trackerType {
	case ArtifactTypeBinary:
		trackerName = "binary"
	case ArtifactTypeDart:
		trackerName = "dart"
	case ArtifactTypeOpenAPI:
		trackerName = "openapi"
	case ArtifactTypeJS:
		trackerName = "js"
	case ArtifactTypeModuleProto:
		trackerName = "proto"
	case ArtifactTypeState:
		trackerName = "appstate"
	}
	return fmt.Sprintf("%s_%s_checksum.txt", trackerName, moduleName)
}

// GetTargetArtifacts checks trackers for changes and returns list of artifacts meant to be built
func (c *Chain) GetTargetArtifacts(ctx context.Context) (Artifacts, error) {
	// get trackers
	// check trackers against changes and get artifacts
	// save current trackers hashes
	// return artifacts
	return nil, nil
}

// GetTrackers gets trackers for chain source code and artifacts
func (c Chain) GetTrackers(ctx context.Context) ([]Tracker, error) {
	conf, err := c.Config()
	if err != nil {
		return nil, err
	}

	modules, err := moduleanalysis.Discover(ctx, c.app.Path, conf.Build.Proto.Path)
	if err != nil {
		return nil, err
	}

	trackers := make([]Tracker, 0)
	for _, m := range modules {
		trackers = append(trackers, c.newProtoOutputTracker(m, c.app.Path), c.newModuleProtoSourceTracker(m))
	}
	binaryName, err := c.Binary()
	if err != nil {
		return nil, err
	}
	binaryTracker, err := c.newBinaryTracker(binaryName)
	if err != nil {
		return nil, err
	}
	return append(
		trackers,
		binaryTracker,
		c.newGoSourceTracker(),
		c.newConfigTracker(c.ConfigPath()),
		c.newDartTracker(conf.Client.Dart.Path),
		c.newOpenAPITracker(conf.Client.OpenAPI.Path),
		c.newVuexOutputTracker(conf.Client.Vuex.Path),
		c.newGoModTracker(filepath.Join(c.app.Path, "go.mod")),
		c.newThirdPartyProtoTracker(filepath.Join(c.app.Path, "third_party")),
	), nil
}

func (c Chain) getArtifactsFromTrackers(trackers []Tracker) error {
	// check trackers changes and returns artifacts
	return nil
}

// persistTrackersChecksums persists checksums for list of trackers
func (c Chain) persistTrackersChecksums(trackers []Tracker) error {
	savePath, err := c.chainSavePath()
	if err != nil {
		return err
	}
	for _, tracker := range trackers {
		err := tracker.persistChecksum(c.app.Path, savePath)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	QueryProtoOutput   = "query.pb.go"
	QueryGWProtoOutput = "query.pb.gw.go"
	TxProtoOutput      = "tx.pb.go"
)

func (c Chain) newConfigTracker(configPath string) Tracker {
	return Tracker{
		trackerType: ConfigTracker,
		paths:       []string{configPath},
	}
}

func (c Chain) newOpenAPITracker(openApiPath string) Tracker {
	if openApiPath == "" {
		openApiPath = defaultOpenAPIPath
	}
	return Tracker{
		trackerType: OpenApiOutputTracker,
		paths:       []string{openApiPath},
	}
}

func (c Chain) newDartTracker(dartOutputPath string) Tracker {
	if dartOutputPath == "" {
		dartOutputPath = defaultDartPath
	}
	return Tracker{
		trackerType: DartOutputTracker,
		paths:       []string{dartOutputPath},
	}
}

func (c Chain) newGoSourceTracker() Tracker {
	return Tracker{
		trackerType: GoSourceTracker,
		paths:       appBackendSourceWatchPaths,
	}
}

func (c Chain) newBinaryTracker(binaryName string) (Tracker, error) {
	binaryPath, err := exec.LookPath(binaryName)
	if err != nil {
		return Tracker{}, err
	}
	return Tracker{
		trackerType: BinaryTracker,
		paths:       []string{binaryPath},
	}, nil
}

func (c Chain) newProtoOutputTracker(m moduleanalysis.Module, appPath string) Tracker {
	return Tracker{
		paths: []string{
			filepath.Join(appPath, "x", m.Name, "types", QueryProtoOutput),
			filepath.Join(appPath, "x", m.Name, "types", QueryGWProtoOutput),
			filepath.Join(appPath, "x", m.Name, "types", TxProtoOutput),
		},
		modules:     []string{m.Name},
		trackerType: ModuleProtoGoOutputTracker,
	}
}

func (c Chain) newVuexOutputTracker(vuexPath string) Tracker {
	if vuexPath == "" {
		vuexPath = defaultVuexPath
	}
	return Tracker{
		paths:       []string{filepath.Join(vuexPath, "generated")},
		trackerType: VuexOutputTracker,
	}
}

func (c Chain) newModuleProtoSourceTracker(m moduleanalysis.Module) Tracker {
	return Tracker{
		paths:       []string{m.Pkg.Path},
		modules:     []string{m.Name},
		trackerType: ModuleProtoSourceTracker,
	}

}

func (c Chain) newGoModTracker(goModPath string) Tracker {
	return Tracker{
		paths:       []string{goModPath},
		trackerType: GoModTracker,
	}
}

func (c Chain) newThirdPartyProtoTracker(thirdPartyProtoPath string) Tracker {
	return Tracker{
		paths:       []string{thirdPartyProtoPath},
		trackerType: ThirdPartyProtoTracker,
	}
}
