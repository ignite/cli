package chain

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
	"github.com/ignite/cli/v29/ignite/templates/app"
)

// CheckBufProtoPath check if the proto path exist into the directory list in the buf.work.yaml file.
func CheckBufProtoPath(appPath, protoPath string) (bool, []string, error) {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return false, nil, err
	}

	missing, err := workFile.MissingDirectories()
	if err != nil {
		return false, nil, err
	}

	return workFile.HasProtoPath(protoPath), missing, nil
}

// AddBufProtoPath add the proto path into the directory list in the buf.work.yaml file.
func AddBufProtoPath(appPath, protoPath string) error {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return err
	}

	return workFile.AddProtoPath(protoPath)
}

// RemoveBufProtoPath add the proto path into the directory list in the buf.work.yaml file.
func RemoveBufProtoPath(appPath string, protoPaths ...string) error {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return err
	}

	return workFile.RemoveProtoPaths(protoPaths...)
}

// CheckBufFiles check if the buf files exist.
func CheckBufFiles(appPath, protoPath string) (bool, error) {
	files, err := app.BufFiles()
	if err != nil {
		return false, nil
	}
	for _, bufFile := range files {
		bufFile, ok := strings.CutPrefix(bufFile, fmt.Sprintf("%s/", defaults.ProtoPath))
		if ok {
			bufFile = filepath.Join(protoPath, bufFile)
		}
		if !xos.FileExists(filepath.Join(appPath, bufFile)) {
			return false, nil
		}
	}
	return true, nil
}

// BoxBufFiles box all buf files.
func BoxBufFiles(runner *xgenny.Runner, appPath, protoPath string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(appPath, protoPath)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return runner.RunAndApply(g)
}
