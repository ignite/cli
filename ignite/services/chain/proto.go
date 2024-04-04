package chain

import (
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
	"github.com/ignite/cli/v29/ignite/templates/app"
)

const defaultProtoFolder = "proto/"

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

func BoxBufFiles(runner *xgenny.Runner, appPath string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(appPath)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return runner.RunAndApply(g)
}
