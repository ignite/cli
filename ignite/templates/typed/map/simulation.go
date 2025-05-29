package maptype

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

		// Create a list of two different index/fields to use as sample
		sampleIndexes := make([]string, 2)
		for i := 0; i < 2; i++ {
			sampleIndexes[i] = fmt.Sprintf("%s: sample.AccAddress(),\n", opts.MsgSigner.UpperCamel)
			sampleIndexes[i] += opts.Index.GenesisArgs(i)
		}

		// simulation genesis state
		content, err := xast.ModifyFunction(
			f.String(),
			"GenerateGenesisState",
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vMap", opts.TypeName.UpperCamel),
				fmt.Sprintf(
					"[]types.%[1]v{{ %[2]v }, { %[3]v }}",
					opts.TypeName.PascalCase,
					sampleIndexes[0],
					sampleIndexes[1],
				),
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
