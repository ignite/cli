package params

// SetSysParam*(module, submodule, name, value) is used for directly setting params in
// ExecContext.Params, or SysParam.
// Access to "sys/params" must be restricted.

func SetSysParamString(module, submodule, name string, val string) {
	setSysParamString(module, submodule, name, val)
}

func SetSysParamBool(module, submodule, name string, val bool) {
	setSysParamBool(module, submodule, name, val)
}

func SetSysParamInt64(module, submodule, name string, val int64) {
	setSysParamInt64(module, submodule, name, val)
}

func SetSysParamUint64(module, submodule, name string, val uint64) {
	setSysParamUint64(module, submodule, name, val)
}

func SetSysParamBytes(module, submodule, name string, val []byte) {
	setSysParamBytes(module, submodule, name, val)
}

func SetSysParamStrings(module, submodule, name string, val []string) {
	if val == nil {
		val = []string{}
	}
	setSysParamStrings(module, submodule, name, val)
}

func UpdateSysParamStrings(module, submodule, name string, val []string, add bool) {
	if val == nil {
		val = []string{}
	}
	updateSysParamStrings(module, submodule, name, val, add)
}

// GetSysParam*(module, submodule, name) reads back a previously-Set value.
// Returns (value, true) if the key is set, (zero, false) otherwise. The
// found-vs-empty distinction matters for callers like v3.init() that need
// to tell "never seeded" apart from "seeded as empty list".
//
// Access to "sys/params" is restricted to gno.land/r/sys/params, mirroring
// the write side. External realms read params via RPC instead.

func GetSysParamString(module, submodule, name string) (string, bool) {
	return getSysParamString(module, submodule, name)
}

func GetSysParamBool(module, submodule, name string) (bool, bool) {
	return getSysParamBool(module, submodule, name)
}

func GetSysParamInt64(module, submodule, name string) (int64, bool) {
	return getSysParamInt64(module, submodule, name)
}

func GetSysParamUint64(module, submodule, name string) (uint64, bool) {
	return getSysParamUint64(module, submodule, name)
}

func GetSysParamBytes(module, submodule, name string) ([]byte, bool) {
	return getSysParamBytes(module, submodule, name)
}

func GetSysParamStrings(module, submodule, name string) ([]string, bool) {
	return getSysParamStrings(module, submodule, name)
}

func setSysParamString(module, submodule, name string, val string)
func setSysParamBool(module, submodule, name string, val bool)
func setSysParamInt64(module, submodule, name string, val int64)
func setSysParamUint64(module, submodule, name string, val uint64)
func setSysParamBytes(module, submodule, name string, val []byte)
func setSysParamStrings(module, submodule, name string, val []string)
func updateSysParamStrings(module, submodule, name string, val []string, add bool)

func getSysParamString(module, submodule, name string) (string, bool)
func getSysParamBool(module, submodule, name string) (bool, bool)
func getSysParamInt64(module, submodule, name string) (int64, bool)
func getSysParamUint64(module, submodule, name string) (uint64, bool)
func getSysParamBytes(module, submodule, name string) ([]byte, bool)
func getSysParamStrings(module, submodule, name string) ([]string, bool)
