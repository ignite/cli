package modulecreate

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

const msgServiceImport = `"github.com/cosmos/cosmos-sdk/types/msgservice"`

// AddMsgServerConventionToLegacyModule add the files and the necessary modifications to an existing module that doesn't support MsgServer convention
// https://github.com/cosmos/cosmos-sdk/blob/main/docs/architecture/adr-031-msg-service.md
func AddMsgServerConventionToLegacyModule(replacer placeholder.Replacer, opts *MsgServerOptions) (*genny.Generator, error) {
	subFs, err := fs.Sub(fsMsgServer, "files/msgserver")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	g.RunFn(codecPath(replacer, opts.ModuleName))

	if err := g.OnlyFS(subFs, nil, nil); err != nil {
		return g, err
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("protoVer", opts.ProtoVer)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName, opts.ProtoVer))

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))

	return g, nil
}

func codecPath(replacer placeholder.Replacer, moduleName string) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", moduleName, "types/codec.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add msgservice import
		oldImport := "import ("
		newImport := fmt.Sprintf(`%v
%v`, oldImport, msgServiceImport)
		content := replacer.Replace(f.String(), oldImport, newImport)

		// Add RegisterMsgServiceDesc method call
		content, err = xast.ModifyFunction(
			content,
			"RegisterInterfaces",
			xast.AppendFuncCode("msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)"),
		)
		if err != nil {
			return err
		}
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
