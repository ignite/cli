package list

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

func genesisModify(opts *typed.Options, g *genny.Generator) {
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(genesisTypesModify(opts))
	g.RunFn(genesisModuleModify(opts))
	g.RunFn(genesisTestsModify(opts))
	g.RunFn(genesisTypesTestsModify(opts))
}

// Modifies the genesis.proto file to add a new field.
//
// What it depends on:
//   - Existence of a message with name "GenesisState". Adds the field there.
func genesisProtoModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}
		// Add initial import for the new type
		gogoImport := protoutil.NewImport(typed.GoGoProtoImport)
		if err = protoutil.AddImports(protoFile, true, gogoImport, opts.ProtoTypeImport()); err != nil {
			return errors.Errorf("failed while adding imports in %s: %w", path, err)
		}
		// Get next available sequence number from GenesisState.
		genesisState, err := protoutil.GetMessageByName(protoFile, typed.ProtoGenesisStateMessage)
		if err != nil {
			return errors.Errorf("failed while looking up message '%s' in %s: %w", typed.ProtoGenesisStateMessage, path, err)
		}
		seqNumber := protoutil.NextUniqueID(genesisState)
		typenameSnake, typenamePascal := opts.TypeName.Snake, opts.TypeName.PascalCase
		// Create option and List field.
		gogoOption := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
		typeList := protoutil.NewField(
			typenameSnake+"_list",
			typenamePascal,
			seqNumber,
			protoutil.Repeated(),
			protoutil.WithFieldOptions(gogoOption),
		)
		// Create count field.
		countFIeld := protoutil.NewField(typenameSnake+"_count", "uint64", seqNumber+1)
		protoutil.Append(genesisState, typeList, countFIeld)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func genesisTypesModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(f.String(), xast.WithImport("fmt"))
		if err != nil {
			return err
		}

		// add parameter to the struct into the new method.
		content, err = xast.ModifyFunction(content, "DefaultGenesis", xast.AppendFuncStruct(
			"GenesisState",
			fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
			fmt.Sprintf("[]%[1]v{}", opts.TypeName.PascalCase),
		))
		if err != nil {
			return err
		}

		templateTypesValidate := `// Check for duplicated ID in %[1]v
%[1]vIdMap := make(map[uint64]bool)
%[1]vCount := gs.Get%[2]vCount()
for _, elem := range gs.%[2]vList {
	if _, ok := %[1]vIdMap[elem.Id]; ok {
		return fmt.Errorf("duplicated id for %[1]v")
	}
	if elem.Id >= %[1]vCount {
		return fmt.Errorf("%[1]v id should be lower or equal than the last id")
	}
	%[1]vIdMap[elem.Id] = true
}`
		replacementTypesValidate := fmt.Sprintf(
			templateTypesValidate,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content, err = xast.ModifyFunction(
			content,
			"Validate",
			xast.AppendFuncCode(replacementTypesValidate),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModuleModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "keeper/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `// Set all the %[1]v
for _, elem := range genState.%[2]vList {
	if err := k.%[2]v.Set(ctx, elem.Id, elem); err != nil {
		return err
	}
}

// Set %[1]v count
if err := k.%[2]vSeq.Set(ctx, genState.%[2]vCount); err != nil {
	return err
}`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content, err := xast.ModifyFunction(
			f.String(),
			"InitGenesis",
			xast.AppendFuncCode(replacementModuleInit),
		)
		if err != nil {
			return err
		}

		templateModuleExport := `
err = k.%[1]v.Walk(ctx, nil, func(key uint64, elem types.%[2]v) (bool, error) {
		genesis.%[1]vList = append(genesis.%[1]vList, elem)
		return false, nil
})
if err != nil {
	return nil, err
}

genesis.%[1]vCount, err = k.%[1]vSeq.Peek(ctx)
if err != nil {
	return nil, err
}`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			opts.TypeName.UpperCamel,
			opts.TypeName.PascalCase,
		)
		content, err = xast.ModifyFunction(
			content,
			"ExportGenesis",
			xast.AppendFuncCode(replacementModuleExport),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTestsModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "keeper/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		replacementAssert := fmt.Sprintf(`require.EqualExportedValues(t, genesisState.%[1]vList, got.%[1]vList)
require.Equal(t, genesisState.%[1]vCount, got.%[1]vCount)`, opts.TypeName.UpperCamel)

		// add parameter to the struct into the new method.
		content, err := xast.ModifyFunction(
			f.String(),
			"TestGenesis",
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
				fmt.Sprintf("[]types.%[1]v{{ Id: 0 }, { Id: 1 }}", opts.TypeName.PascalCase),
			),
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vCount", opts.TypeName.UpperCamel),
				"2",
			),
			xast.AppendFuncCode(replacementAssert),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesTestsModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateTestDuplicated := `{
	desc:     "duplicated %[1]v",
	genState: &types.GenesisState{
		%[2]vList: []types.%[3]v{
			{
				Id: 0,
			},
			{
				Id: 0,
			},
		},
	},
	valid:    false,
}`
		replacementTestDuplicated := fmt.Sprintf(
			templateTestDuplicated,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			opts.TypeName.PascalCase,
		)

		templateTestInvalidCount := `{
	desc:     "invalid %[1]v count",
	genState: &types.GenesisState{
		%[2]vList: []types.%[3]v{
			{
				Id: 1,
			},
		},
		%[2]vCount: 0,
	},
	valid:    false,
}`
		replacementInvalidCount := fmt.Sprintf(
			templateTestInvalidCount,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
			opts.TypeName.PascalCase,
		)

		// add parameter to the struct into the new method.
		content, err := xast.ModifyFunction(
			f.String(),
			"TestGenesisState_Validate",
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
				fmt.Sprintf("[]types.%[1]v{{ Id: 0 }, { Id: 1 }}", opts.TypeName.PascalCase),
			),
			xast.AppendFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vCount", opts.TypeName.UpperCamel),
				"2",
			),
			xast.AppendFuncTestCase(replacementTestDuplicated),
			xast.AppendFuncTestCase(replacementInvalidCount),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
