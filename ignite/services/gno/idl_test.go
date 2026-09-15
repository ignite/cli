package gno

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateIDLRealm(t *testing.T) {
	dir, _, err := Scaffold(KindRealm, "idlrealm", ScaffoldOptions{Path: filepath.Join(t.TempDir(), "idlrealm")})
	assert.NilError(t, err)

	idl, err := GenerateIDL(dir)
	assert.NilError(t, err)
	assert.Equal(t, "gno.land/r/idlrealm", idl.PkgPath)
	assert.Equal(t, "realm", idl.Kind)
	assert.Equal(t, "idlrealm", idl.Name)

	byName := map[string]IDLInstruction{}
	for _, instr := range idl.Instructions {
		byName[instr.Name] = instr
	}
	assert.Assert(t, byName["Increment"].Realm, "Increment takes a realm param")
	assert.Assert(t, !byName["Get"].Realm, "Get is read-only")
	assert.Equal(t, "string", argType(byName["Render"], "path"))
}

func TestGenerateIDLErrors(t *testing.T) {
	_, err := GenerateIDL(t.TempDir())
	assert.ErrorContains(t, err, "gnomod.toml")
}

func TestIDLJSONRoundTrip(t *testing.T) {
	idl := &IDL{
		Version: idlVersion,
		Name:    "x",
		PkgPath: "gno.land/r/x",
		Kind:    "realm",
		Instructions: []IDLInstruction{
			{Name: "Set", Args: []IDLArg{{Name: "value", Type: "int"}}, Realm: true},
		},
	}
	data, err := idl.WriteIDLJSON()
	assert.NilError(t, err)

	var back IDL
	assert.NilError(t, json.Unmarshal(data, &back))
	assert.Equal(t, idl.PkgPath, back.PkgPath)
	assert.Equal(t, 1, len(back.Instructions))
	assert.Equal(t, "value", back.Instructions[0].Args[0].Name)
}

func TestGenerateTSClient(t *testing.T) {
	idl := &IDL{
		Version: "1.0",
		Name:    "counter",
		PkgPath: "gno.land/r/demo/counter",
		Kind:    "realm",
		Instructions: []IDLInstruction{
			{Name: "Increment", Realm: true},
			{Name: "Get"},
			{Name: "Set", Args: []IDLArg{{Name: "value", Type: "int"}}, Realm: true},
		},
	}
	src := GenerateTSClient(idl)

	for _, want := range []string{
		"export class CounterClient",
		`this.signer.signAndBroadcast("gno.land/r/demo/counter", "Increment", [])`,
		`signAndBroadcast("gno.land/r/demo/counter", "Set", [JSON.stringify(String(value))])`,
		`"gno.land/r/demo/counter.Get" + "(" + [].join(", ") + ")"`,
		"value: number",
		"vm/qeval_json",
	} {
		assert.Assert(t, strings.Contains(src, want), "missing: "+want)
	}
}

func TestTSTypeMapping(t *testing.T) {
	tt := []struct {
		in   string
		want string
	}{
		{"string", "string"},
		{"int64", "number"},
		{"bool", "boolean"},
		{"[]string", "string[]"},
		{"foo.Bar", "any"},
	}
	for _, tc := range tt {
		assert.Equal(t, tc.want, tsType(tc.in), tc.in)
	}
}

func TestTsPascal(t *testing.T) {
	assert.Equal(t, "Counter", tsPascal("counter"))
	assert.Equal(t, "A", tsPascal("a"))
}

// argType returns the declared type of an argument by name.
func argType(instr IDLInstruction, name string) string {
	for _, a := range instr.Args {
		if a.Name == name {
			return a.Type
		}
	}
	return ""
}
