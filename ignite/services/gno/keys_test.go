package gno

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

// withTestGnoHome points the gno home at a temp dir for the test duration.
func withTestGnoHome(t *testing.T) string {
	t.Helper()
	home := filepath.Join(t.TempDir(), "gnohome")
	t.Setenv("GNOHOME", home)
	err := os.MkdirAll(home, 0o755)
	assert.NilError(t, err)
	return home
}

func TestKeyLifecycle(t *testing.T) {
	withTestGnoHome(t)

	info, mnemonic, err := CreateKey("alice", "", "", 0, 0)
	assert.NilError(t, err)
	assert.Equal(t, "alice", info.Name)
	assert.Assert(t, len(mnemonic) > 0, "mnemonic should be generated")
	assert.Assert(t, len(info.Address) > 0, "address should be set")

	// show by name
	shown, err := ShowKey("alice")
	assert.NilError(t, err)
	assert.Equal(t, info.Address, shown.Address)

	// list contains the key
	list, err := ListKeys()
	assert.NilError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "alice", list[0].Name)

	// export + import round trip
	armorFile := filepath.Join(t.TempDir(), "alice.asc")
	_, err = ExportKey("alice", "", armorFile)
	assert.NilError(t, err)
	err = DeleteKey("alice", "")
	assert.NilError(t, err)
	imported, err := ImportKey("alice2", armorFile, "")
	assert.NilError(t, err)
	assert.Equal(t, info.Address, imported.Address)

	// recover from mnemonic
	recovered, err := RecoverKey("alice-recovered", mnemonic, "", 0, 0)
	assert.NilError(t, err)
	assert.Equal(t, info.Address, recovered.Address)
}
