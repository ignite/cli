package chain

import (
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
	"github.com/ignite/cli/v28/ignite/templates/app"
)

// CheckBufProtoDir check if the proto path exist into the directory list in the buf.work.yaml file.
func CheckBufProtoDir(appPath, protoDir string) (bool, []string, error) {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return false, nil, err
	}

	missing, err := workFile.MissingDirectories()
	if err != nil {
		return false, nil, err
	}

	return workFile.HasProtoDir(protoDir), missing, nil
}

// AddBufProtoDir add the proto path into the directory list in the buf.work.yaml file.
func AddBufProtoDir(appPath, protoDir string) error {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return err
	}

	return workFile.AddProtoDir(protoDir)
}

// RemoveBufProtoDirs add the proto path into the directory list in the buf.work.yaml file.
func RemoveBufProtoDirs(appPath string, protoDirs ...string) error {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return err
	}

	return workFile.RemoveProtoDirs(protoDirs...)
}

// CheckBufFiles check if the buf files exist.
func CheckBufFiles(appPath, protoDir string) (bool, error) {
	files, err := app.BufFiles()
	if err != nil {
		return false, nil
	}
	for _, bufFile := range files {
		bufFile, ok := app.CutTemplatePrefix(bufFile)
		if ok {
			bufFile = filepath.Join(protoDir, bufFile)
		}
		if !xos.FileExists(filepath.Join(appPath, bufFile)) {
			return false, nil
		}
	}
	return true, nil
}

// BoxBufFiles box all buf files.
func BoxBufFiles(runner *xgenny.Runner, appPath, protoDir string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(appPath, protoDir)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return xgenny.RunWithValidation(placeholder.New(), g)
}
