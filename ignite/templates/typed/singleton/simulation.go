package singleton

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/templates/typed"
)

func moduleSimulationModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module_simulation.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content := typed.ModuleSimulationMsgModify(
			replacer,
			f.String(),
			opts.ModuleName,
			opts.TypeName,
			"Create", "Update", "Delete",
		)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
