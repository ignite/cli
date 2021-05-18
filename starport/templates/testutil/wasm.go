package testutil

import (
	"context"
	"fmt"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

// app.NewApp modification in testutil module on Stargate when importing wasm
func testutilAppModifyStargate(ctx context.Context) genny.RunFn {
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
			content := placeholder.Replace(ctx, f.String(), placeholderSgTestutilAppArgument, replacementAppArgument)

			newFile := genny.NewFileS(path, content)
			if err := r.File(newFile); err != nil {
				return err
			}
		}
		return nil
	}
}

// WASMRegister register testutil modifiers that should be applied when wasm is imported.
func WASMRegister(ctx context.Context, _ *plush.Context, gen *genny.Generator) error {
	gen.RunFn(testutilAppModifyStargate(ctx))
	return nil
}
