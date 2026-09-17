package gno

import (
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestExtractSignaturesRealm(t *testing.T) {
	dir, _, err := Scaffold(KindRealm, "idlrealm", ScaffoldOptions{Path: filepath.Join(t.TempDir(), "idlrealm")})
	assert.NilError(t, err)

	sig, err := ExtractSignatures(dir)
	assert.NilError(t, err)
	assert.Equal(t, "gno.land/r/idlrealm", sig.PkgPath)
	assert.Equal(t, "idlrealm", sig.Name)

	byName := map[string]FuncSig{}
	for _, fn := range sig.Funcs {
		byName[fn.Name] = fn
	}
	assert.Assert(t, byName["Increment"].Realm, "Increment takes a realm param")
	assert.Assert(t, !byName["Get"].Realm, "Get is read-only")
	assert.Equal(t, "string", argType(byName["Render"], "path"))
}

func TestExtractSignaturesErrors(t *testing.T) {
	_, err := ExtractSignatures(t.TempDir())
	assert.ErrorContains(t, err, "gnomod.toml")
}

func TestGenerateRealmClient(t *testing.T) {
	sig := &PackageSignatures{
		Name:    "counter",
		PkgPath: "gno.land/r/demo/counter",
		Funcs: []FuncSig{
			{Name: "Increment", Realm: true},
			{Name: "Get"},
			{Name: "Set", Args: []FuncArg{{Name: "value", Type: "int"}}, Realm: true},
		},
	}
	src := GenerateRealmClient(sig)

	for _, want := range []string{
		`import {`,
		`  GnoWallet,`,
		`  parseGnoReturns,`,
		`} from "@gnolang/gno-js-client";`,
		`const pkgPath = "gno.land/r/demo/counter"`,
		"export interface CounterRealm",
		"export const counterRealm = (instance: GnoWallet)",
		`instance.callMethod(pkgPath, "Increment", [], "commit", send)`,
		`instance.callMethod(pkgPath, "Set", [String(value)], "commit", send)`,
		`provider.evaluateExpression(pkgPath, "Get" + "(" + [].join(", ") + ")")`,
		"parseGnoReturns(res)",
		"value: number",
		"GnoWallet.addRealm(counterRealm)",
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
func argType(fn FuncSig, name string) string {
	for _, a := range fn.Args {
		if a.Name == name {
			return a.Type
		}
	}
	return ""
}
