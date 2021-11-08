package singleton

import (
	"fmt"
	"math/rand"
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

		// Create fields
		sampleFields := fmt.Sprintf("%[1]v: sample.AccAddress(),\n", opts.MsgSigner.UpperCamel)
		for _, field := range opts.Fields {
			sampleFields += field.GenesisArgs(rand.Intn(100) + 1)
		}

		templateGs := `%[2]v: &types.%[2]v{
		%[3]v},
		%[1]v`
		replacementGs := fmt.Sprintf(
			templateGs,
			typed.PlaceholderSimappGenesisState,
			opts.TypeName.UpperCamel,
			sampleFields,
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
