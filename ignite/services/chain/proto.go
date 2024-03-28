package chain

import (
	"path/filepath"

<<<<<<< HEAD
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
	"github.com/ignite/cli/v28/ignite/templates/app"
=======
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
	"github.com/ignite/cli/v29/ignite/templates/app"
>>>>>>> 2ad41ee3 (feat(pkg): improve xgenny dry run (#4001))
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

func BoxBufFiles(runner *xgenny.Runner, appPath string) (xgenny.SourceModification, error) {
	g, err := app.NewBufGenerator(appPath)
	if err != nil {
		return xgenny.SourceModification{}, err
	}
	return runner.RunAndApply(g)
}
