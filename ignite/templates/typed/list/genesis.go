package list

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
	"github.com/ignite/cli/v29/ignite/templates/typed"
)

func genesisModify(replacer placeholder.Replacer, opts *typed.Options, g *genny.Generator) {
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(genesisTypesModify(opts))
	g.RunFn(genesisModuleModify(replacer, opts))
	g.RunFn(genesisTestsModify(replacer, opts))
	g.RunFn(genesisTypesTestsModify(replacer, opts))
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
		typenameLower, typenameUpper := opts.TypeName.LowerCamel, opts.TypeName.UpperCamel
		// Create option and List field.
		gogoOption := protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom())
		typeList := protoutil.NewField(
			typenameLower+"List",
			typenameUpper,
			seqNumber,
			protoutil.Repeated(),
			protoutil.WithFieldOptions(gogoOption),
		)
		// Create count field.
		countFIeld := protoutil.NewField(typenameLower+"Count", "uint64", seqNumber+1)
		protoutil.Append(genesisState, typeList, countFIeld)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func genesisTypesModify(opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(f.String(), xast.WithLastImport("fmt"))
		if err != nil {
			return err
		}

		// add parameter to the struct into the new method.
		content, err = xast.ModifyFunction(content, "DefaultGenesis", xast.AppendInsideFuncStruct(
			"GenesisState",
			fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
			fmt.Sprintf("[]%[1]v{}", opts.TypeName.UpperCamel),
			-1,
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

func genesisModuleModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateModuleInit := `// Set all the %[2]v
for _, elem := range genState.%[3]vList {
	if err := k.%[3]v.Set(ctx, elem.Id, elem); err != nil {
		return err
	}
}

// Set %[2]v count
if err := k.%[3]vSeq.Set(ctx, genState.%[3]vCount); err != nil {
	return err
}
%[1]v`
		replacementModuleInit := fmt.Sprintf(
			templateModuleInit,
			typed.PlaceholderGenesisModuleInit,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementModuleInit)

		templateModuleExport := `
err = k.%[2]v.Walk(ctx, nil, func(key uint64, elem types.%[2]v) (bool, error) {
		genesis.%[2]vList = append(genesis.%[2]vList, elem)
		return false, nil
})
if err != nil {
	return nil, err
}

genesis.%[2]vCount, err = k.%[2]vSeq.Peek(ctx)
if err != nil {
	return nil, err
}

%[1]v`
		replacementModuleExport := fmt.Sprintf(
			templateModuleExport,
			typed.PlaceholderGenesisModuleExport,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, typed.PlaceholderGenesisModuleExport, replacementModuleExport)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTestsModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "keeper/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateState := `%[2]vList: []types.%[2]v{
		{
			Id: 0,
		},
		{
			Id: 1,
		},
	},
	%[2]vCount: 2,
	%[1]v`
		replacementValid := fmt.Sprintf(
			templateState,
			module.PlaceholderGenesisTestState,
			opts.TypeName.UpperCamel,
		)
		content := replacer.Replace(f.String(), module.PlaceholderGenesisTestState, replacementValid)

		templateAssert := `require.ElementsMatch(t, genesisState.%[2]vList, got.%[2]vList)
require.Equal(t, genesisState.%[2]vCount, got.%[2]vCount)
%[1]v`
		replacementTests := fmt.Sprintf(
			templateAssert,
			module.PlaceholderGenesisTestAssert,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, module.PlaceholderGenesisTestAssert, replacementTests)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesTestsModify(replacer placeholder.Replacer, opts *typed.Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis_test.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// add parameter to the struct into the new method.
		content, err := xast.ModifyFunction(
			f.String(),
			"TestGenesisState_Validate",
			xast.AppendInsideFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vList", opts.TypeName.UpperCamel),
				fmt.Sprintf("[]types.%[1]v{{ Id: 0 }, { Id: 1 }}", opts.TypeName.UpperCamel),
				-1,
			),
			xast.AppendInsideFuncStruct(
				"GenesisState",
				fmt.Sprintf("%[1]vCount", opts.TypeName.UpperCamel),
				"2",
				-1,
			),
		)
		if err != nil {
			return err
		}

		templateTests := `{
	desc:     "duplicated %[2]v",
	genState: &types.GenesisState{
		%[3]vList: []types.%[3]v{
			{
				Id: 0,
			},
			{
				Id: 0,
			},
		},
	},
	valid:    false,
},
{
	desc:     "invalid %[2]v count",
	genState: &types.GenesisState{
		%[3]vList: []types.%[3]v{
			{
				Id: 1,
			},
		},
		%[3]vCount: 0,
	},
	valid:    false,
},
%[1]v`
		replacementTests := fmt.Sprintf(
			templateTests,
			module.PlaceholderTypesGenesisTestcase,
			opts.TypeName.LowerCamel,
			opts.TypeName.UpperCamel,
		)
		content = replacer.Replace(content, module.PlaceholderTypesGenesisTestcase, replacementTests)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
