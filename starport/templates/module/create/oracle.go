package modulecreate

import (
	"fmt"
	"github.com/tendermint/starport/starport/templates/module"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xstrings"
)

// NewOracle returns the generator to scaffold the implementation of the Oracle interface inside a module
func NewOracle(replacer placeholder.Replacer, opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(moduleOracleModify(replacer, opts))

	if err := g.Box(ibcTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("title", strings.Title)
	ctx.Set("dependencies", opts.Dependencies)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func moduleOracleModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}


		// Recv packet dispatch
		templateRecv := `	
	ack, oracleResult, err := am.handleOraclePacket(ctx, modulePacket)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error())
	} else if oracleResult.Size() > 0 {
		ctx.Logger().Debug("Receive oracle packet", "result", oracleResult)
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, ack.GetBytes(), nil
	}`
		content := replacer.Replace(f.String(), module.PlaceholderOraclePacketModuleRecv, templateRecv)

		// Ack packet dispatch
		templateAck := `
	var requestID types.RequestID
	ctx, requestID = am.handleOracleAcknowledgement(ctx, ack)
	if requestID > 0 {
		ctx.Logger().Debug("Receive oracle ack", "request_id", requestID)
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, nil
	}`
		content = replacer.Replace(content, module.PlaceholderOraclePacketModuleAck, templateAck)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}