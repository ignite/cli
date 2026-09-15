package testing

func setSysParamStrings(module, submodule, name string, value []string)
func setSysParamBool(module, submodule, name string, value bool)
func setSysParamUint64(module, submodule, name string, value uint64)
func setSysParamInt64(module, submodule, name string, value int64)

// SetSysParamStrings seeds a []string-typed param in the test
// store, bypassing the realm-caller gate that production
// sys/params helpers enforce. Test-only — this file is loaded by
// gno's test runner only (not by the chain runtime).
//
// Use to seed state that realm code reads via prms.GetSysParam*,
// e.g. validator-set fixtures for r/sys/validators tests.
func SetSysParamStrings(module, submodule, name string, value []string) {
	setSysParamStrings(module, submodule, name, value)
}

// SetSysParamBool seeds a bool-typed param in the test store. See
// SetSysParamStrings.
func SetSysParamBool(module, submodule, name string, value bool) {
	setSysParamBool(module, submodule, name, value)
}

// SetSysParamUint64 seeds a uint64-typed param in the test store.
// See SetSysParamStrings.
func SetSysParamUint64(module, submodule, name string, value uint64) {
	setSysParamUint64(module, submodule, name, value)
}

// SetSysParamInt64 seeds an int64-typed param in the test store.
// See SetSysParamStrings.
func SetSysParamInt64(module, submodule, name string, value int64) {
	setSysParamInt64(module, submodule, name, value)
}
