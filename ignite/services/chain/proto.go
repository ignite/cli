package chain

import (
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
	"github.com/ignite/cli/v29/ignite/templates/app"
)

// oldBufWorkFile represents the v1 buf work file, may this check should be remove in v30
const oldBufWorkFile = "buf.work.yaml"

// CheckBufProtoDir check if the proto path exist into the directory list in the buf.work.yaml file.
func CheckBufProtoDir(appPath, protoDir string) (bool, []string, error) {
	bufCfg, err := cosmosbuf.ParseBufConfig(appPath)
	if err != nil {
		return false, nil, err
	}

	missing, err := bufCfg.MissingDirectories()
	if err != nil {
		return false, nil, err
	}

	return bufCfg.HasProtoDir(protoDir), missing, nil
}

// AddBufProtoDir add the proto path into the directory list in the buf.work.yaml file.
func AddBufProtoDir(appPath, protoDir string) error {
	workFile, err := cosmosbuf.ParseBufConfig(appPath)
	if err != nil {
		return err
	}

	return workFile.AddProtoDir(protoDir)
}

// RemoveBufProtoDirs add the proto path into the directory list in the buf.work.yaml file.
func RemoveBufProtoDirs(appPath string, protoDirs ...string) error {
	workFile, err := cosmosbuf.ParseBufConfig(appPath)
	if err != nil {
		return err
	}

	return workFile.RemoveProtoDirs(protoDirs...)
}

// CheckBufFiles check if the buf files exist, and if needs a migration to v2.
func CheckBufFiles(appPath, protoDir string) (bool, bool, error) {
	files, err := app.BufFiles()
	if err != nil {
		return false, false, nil
	}
	// if the buf.work.yaml exist, we only need the migration
	if xos.FileExists(filepath.Join(appPath, oldBufWorkFile)) {
		return true, true, nil
	}
	for _, bufFile := range files {
		bufFile, ok := app.CutTemplatePrefix(bufFile)
		if ok {
			bufFile = filepath.Join(protoDir, bufFile)
		}
		if !xos.FileExists(filepath.Join(appPath, bufFile)) {
			return false, false, nil
		}
	}
	return true, false, nil
}

// BoxBufFiles box all buf files.
func BoxBufFiles(runner *xgenny.Runner, appPath, protoDir string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(appPath, protoDir)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return runner.RunAndApply(g)
}
