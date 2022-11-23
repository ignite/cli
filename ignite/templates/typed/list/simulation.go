package list

import (
	"fmt"
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

		// Create a list of two different indexes and fields to use as sample
		msgField := fmt.Sprintf("%s: sample.AccAddress(),\n", opts.MsgSigner.UpperCamel)

		// simulation genesis state
		templateGs := `	%[2]vList: []types.%[2]v{
		{
			Id: 0,
			%[3]v
		},
		{
			Id: 1,
			%[3]v
		},
	},
	%[2]vCount: 2,
	%[1]v`
		replacementGs := fmt.Sprintf(
			templateGs,
			typed.PlaceholderSimappGenesisState,
			opts.TypeName.UpperCamel,
			msgField,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderSimappGenesisState, replacementGs)

		content = typed.ModuleSimulationMsgModify(
			replacer,
			content,
			opts.ModuleName,
			opts.TypeName,
			"Create", "Update", "Delete",
		)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
