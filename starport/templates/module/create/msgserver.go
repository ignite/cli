package modulecreate

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"github.com/tendermint/starport/starport/templates/module"
	"github.com/tendermint/starport/starport/templates/typed"
)

const msgServiceImport = "github.com/cosmos/cosmos-sdk/types/msgservice"

// AddMsgServerConventionToLegacyModule add the files and the necessary modifications to an existing module that doesn't support MsgServer convention
// https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-031-msg-service.md
func AddMsgServerConventionToLegacyModule(opts *MsgServerOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(handlerPatch(opts.ModuleName))
	g.RunFn(codecPath(opts.ModuleName))

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

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func handlerPatch(moduleName string) genny.RunFn {
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
		content := strings.Replace(f.String(), old, new, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecPath(moduleName string) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", moduleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add msgservice import
		old := "import("
		new := fmt.Sprintf(`%v
%v`, old, msgServiceImport)
		content := strings.Replace(f.String(), old, new, 1)

		// Add RegisterMsgServiceDesc method call
		template := `%[1]v

msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)`
		replacement := fmt.Sprintf(template, module.Placeholder3)
		content = strings.Replace(content, module.Placeholder3, replacement, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
