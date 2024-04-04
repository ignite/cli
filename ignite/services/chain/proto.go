package chain

import (
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
	"github.com/ignite/cli/v29/ignite/templates/app"
)

const defaultProtoFolder = "proto/"

func CheckBufProtoPath(appPath, protoPath string) (bool, error) {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return false, err
	}

	return workFile.HasProtoPath(protoPath), nil
}

func AddBufProtoPath(appPath, protoPath string) error {
	workFile, err := cosmosbuf.ParseBufWork(appPath)
	if err != nil {
		return err
	}

	return workFile.AddProtoPath(protoPath)
}

func CheckBufFiles(appPath, protoPath string) (bool, error) {
	files, err := app.BufFiles()
	if err != nil {
		return false, nil
	}
	for _, bufFile := range files {
		bufFile, ok := strings.CutPrefix(bufFile, defaultProtoFolder)
		if ok {
			bufFile = filepath.Join(protoPath, bufFile)
		}
		if !xos.FileExists(filepath.Join(appPath, bufFile)) {
			return false, nil
		}
	}
	return true, nil
}

func BoxBufFiles(runner *xgenny.Runner, appPath, protoPath string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(appPath, protoPath)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return runner.RunAndApply(g)
}
