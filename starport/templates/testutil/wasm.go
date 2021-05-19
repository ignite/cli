package testutil

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
)

// app.NewApp modification in testutil module on Stargate when importing wasm
func testutilAppModifyStargate() genny.RunFn {
	return func(r *genny.Runner) error {
		for _, path := range []string{
			"testutil/simapp/simapp.go",
			"testutil/network/network.go",
		} {
			f, err := r.Disk.Find(path)
			if err != nil {
				return err
			}

			templateenabledProposals := `%[1]v
			app.GetEnabledProposals(), nil,`
			replacementAppArgument := fmt.Sprintf(templateenabledProposals, placeholderSgTestutilAppArgument)
			content := strings.Replace(f.String(), placeholderSgTestutilAppArgument, replacementAppArgument, 1)

			newFile := genny.NewFileS(path, content)
			if err := r.File(newFile); err != nil {
				return err
			}
		}
		return nil
	}
}

// WASMRegister register testutil modifiers that should be applied when wasm is imported.
func WASMRegister(_ *plush.Context, gen *genny.Generator) error {
	gen.RunFn(testutilAppModifyStargate())
	return nil
}
