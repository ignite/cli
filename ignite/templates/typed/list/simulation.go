package list

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

func moduleSimulationModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "module/simulation.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Create a list of two different indexes and fields to use as sample
		msgField := fmt.Sprintf("%s: sample.AccAddress(),\n", opts.MsgSigner.UpperCamel)

		// simulation genesis state
		content, err := xast.ModifyFunction(
			f.String(),
			"GenerateGenesisState",
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
				fmt.Sprintf(
					"[]types.%[2]v{{ Id: 0, %[3]v }, { Id: 1, %[3]v }}",
					opts.TypeName.UpperCamel,
					opts.TypeName.PascalCase,
					msgField,
				),
			),
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vCount", opts.TypeName.UpperCamel),
				"2",
			),
		)
		if err != nil {
			return err
		}

		content, err = typed.ModuleSimulationMsgModify(
			content,
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
