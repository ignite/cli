package list

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/templates/typed"
)

func moduleSimulationModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "module_simulation.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// simulation genesis state
		templateGs := `	%[2]vList: []types.%[2]v{
		{
			Id: 0,
		},
		{
			Id: 1,
		},
	},
	%[2]vCount: 2,
	%[1]v`
		replacementGs := fmt.Sprintf(templateGs, typed.PlaceholderSimapGenesisState, opts.TypeName.UpperCamel)
		content := replacer.Replace(f.String(), typed.PlaceholderSimapGenesisState, replacementGs)

		content = typed.ModuleSimulationMsgModify(
			replacer,
			content,
			opts.TypeName,
			"Create", "Update", "Delete",
		)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
