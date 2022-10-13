package postgres_test

import (
	"bytes"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
)

func TestSchemasWalk(t *testing.T) {
	// Arrange: Scripts data by version
	data := map[uint]string{
		1: "/* TEST-V1 */",
		2: "/* TEST-V2 */",
	}

	// Arrange: Script argument matchers
	matchByDataV1 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(data[1]))
	})
	matchByDataV2 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(data[2]))
	})

	// Prepare the walk function mock
	m := mock.Mock{}
	m.Test(t)
	m.On("fn", uint64(1), matchByDataV1).Return(nil)
	m.On("fn", uint64(2), matchByDataV2).Return(nil)

	fn := func(version uint64, script []byte) error {
		return m.MethodCalled("fn", version, script).Error(0)
	}

	// Arrange: A new schema that contains SQL scripts for three versions
	fs := fstest.MapFS{
		"schemas/1.sql": &fstest.MapFile{Data: []byte(data[1])},
		"schemas/2.sql": &fstest.MapFile{Data: []byte(data[2])},
	}
	s := postgres.NewSchemas(fs, "")

	// Act
	err := s.WalkFrom(1, fn)

	// Assert
	require.NoError(t, err)
	m.AssertExpectations(t)
}

func TestSchemasWalkOrder(t *testing.T) {
	// Arrange: Scripts data by version
	data := map[uint]string{
		1:  "/* TEST-V1 */",
		2:  "/* TEST-V2 */",
		10: "/* TEST-V10 */",
	}

	// Arrange: Script argument matchers
	matchByDataV1 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(data[1]))
	})
	matchByDataV2 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(data[2]))
	})
	matchByDataV10 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(data[10]))
	})

	// Arrange: Walk function mock
	m := mock.Mock{}
	m.Test(t)
	m.On("fn", uint64(1), matchByDataV1).Return(nil)
	m.On("fn", uint64(2), matchByDataV2).Return(nil)
	m.On("fn", uint64(10), matchByDataV10).Return(nil)

	var versions []uint64

	fn := func(ver uint64, script []byte) error {
		versions = append(versions, ver)

		return m.MethodCalled("fn", ver, script).Error(0)
	}

	// Arrange: A new schema that contains SQL scripts for three versions
	fs := fstest.MapFS{
		"schemas/1.sql":  &fstest.MapFile{Data: []byte(data[1])},
		"schemas/2.sql":  &fstest.MapFile{Data: []byte(data[2])},
		"schemas/10.sql": &fstest.MapFile{Data: []byte(data[10])},
	}
	s := postgres.NewSchemas(fs, "")

	// Act
	err := s.WalkFrom(1, fn)

	// Assert
	require.NoError(t, err)
	require.IsIncreasing(t, versions)
	m.AssertExpectations(t)
}

func TestScriptBuilder(t *testing.T) {
	// Arrange
	s1 := "SCRIPT-1;"
	s2 := "SCRIPT-2;"
	c1 := "COMMAND-1"
	c2 := "COMMAND-2"

	b := postgres.ScriptBuilder{}
	b.BeginTX()
	b.AppendScript([]byte(s1))
	b.AppendScript([]byte(s2))
	b.AppendCommand(c1)
	b.AppendCommand(c2)
	b.CommitTX()

	want := fmt.Sprintf("BEGIN;%s%s%s;%s;COMMIT;", s1, s2, c1, c2)

	// Act
	script := b.Bytes()

	// Assert
	require.EqualValues(t, []byte(want), script)
}
