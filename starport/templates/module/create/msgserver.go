package modulecreate

import (
	"fmt"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/plushhelpers"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"github.com/tendermint/starport/starport/templates/module"
	"github.com/tendermint/starport/starport/templates/typed"
)

const msgServiceImport = `"github.com/cosmos/cosmos-sdk/types/msgservice"`

// AddMsgServerConventionToLegacyModule add the files and the necessary modifications to an existing module that doesn't support MsgServer convention
// https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-031-msg-service.md
func AddMsgServerConventionToLegacyModule(replacer placeholder.Replacer, opts *MsgServerOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(handlerPatch(replacer, opts.ModuleName))
	g.RunFn(codecPath(replacer, opts.ModuleName))

	if err := g.Box(msgServerTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func handlerPatch(replacer placeholder.Replacer, moduleName string) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", moduleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add the msg server definition placeholder
		old := "func NewHandler(k keeper.Keeper) sdk.Handler {"
		new := fmt.Sprintf(`%v
%v`, old, typed.PlaceholderHandlerMsgServer)
		content := replacer.ReplaceOnce(f.String(), old, new)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecPath(replacer placeholder.Replacer, moduleName string) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", moduleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add msgservice import
		old := "import ("
		new := fmt.Sprintf(`%v
%v`, old, msgServiceImport)
		content := replacer.Replace(f.String(), old, new)

		// Add RegisterMsgServiceDesc method call
		template := `%[1]v

msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)`
		replacement := fmt.Sprintf(template, module.Placeholder3)
		content = replacer.Replace(content, module.Placeholder3, replacement)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
