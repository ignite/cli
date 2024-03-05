package chain

import (
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
	"github.com/ignite/cli/v28/ignite/templates/app"
)

var bufFiles = []string{
	"buf.work.yaml",
	"proto/buf.gen.gogo.yaml",
	"proto/buf.gen.pulsar.yaml",
	"proto/buf.gen.swagger.yaml",
	"proto/buf.gen.ts.yaml",
	"proto/buf.lock",
	"proto/buf.yaml",
}

func CheckBufFiles(appPath string) bool {
	for _, bufFile := range bufFiles {
		if !xos.FileExists(filepath.Join(appPath, bufFile)) {
			return false
		}
	}
	return true
}

func BoxBufFiles(runner *xgenny.Runner, appPath string) error {
	g, err := app.NewBufGenerator(appPath)
	if err != nil {
		return err
	}
	return runner.Run(g)
}
