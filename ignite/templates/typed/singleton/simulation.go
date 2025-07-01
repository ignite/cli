package singleton

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/templates/typed"
)

func moduleSimulationModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "module/simulation.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := typed.ModuleSimulationMsgModify(
			f.String(),
			opts.ModulePath,
			opts.ModuleName,
			opts.TypeName,
			opts.MsgSigner,
			"Create", "Update", "Delete",
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
