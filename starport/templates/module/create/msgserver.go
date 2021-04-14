package modulecreate

import (
	"fmt"
	"github.com/tendermint/starport/starport/templates/module"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"github.com/tendermint/starport/starport/templates/typed"
)

// AddMsgServerConvention add the files and the necessary modifications to an existing module
// in order to support MsgServer convention: https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-031-msg-service.md
func AddMsgServerConvention(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(handlerPatch(opts))
	g.RunFn(codecPath(opts))

	if err := g.Box(msgServerTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName) //
	ctx.Set("modulePath", opts.ModulePath) //
	ctx.Set("appName", opts.AppName)       //
	ctx.Set("ownerName", opts.OwnerName)   //
	ctx.Set("title", strings.Title)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func handlerPatch(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add the msg server definition placeholder
		old := "func NewHandler(k keeper.Keeper) sdk.Handler {"
		new := fmt.Sprintf(`%v
%v`, old, typed.PlaceholderHandlerMsgServer)

		content := strings.Replace(f.String(), old, new, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecPath(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v

msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)`
		replacement := fmt.Sprintf(template, module.Placeholder3)
		content := strings.Replace(f.String(),  module.Placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}