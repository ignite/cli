package postgres_test

import (
	"bytes"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/ignite-hq/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSchemasWalk(t *testing.T) {
	// Arrange
	dataV1 := "/* TEST-V1 */"
	dataV2 := "/* TEST-V2 */"

	// Define script argument matchers
	matchByDataV1 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(dataV1))
	})
	matchByDataV2 := mock.MatchedBy(func(script []byte) bool {
		return bytes.Contains(script, []byte(dataV2))
	})

	// Prepare the walk function mock
	m := mock.Mock{}
	m.Test(t)
	m.On("fn", uint64(1), matchByDataV1).Return(nil)
	m.On("fn", uint64(2), matchByDataV2).Return(nil)

	fn := func(version uint64, script []byte) error {
		return m.MethodCalled("fn", version, script).Error(0)
	}

	// Create a new schema that contains SQL scripts for two versions
	fs := fstest.MapFS{
		"schemas/1.sql": &fstest.MapFile{Data: []byte(dataV1)},
		"schemas/2.sql": &fstest.MapFile{Data: []byte(dataV2)},
	}
	s := postgres.NewSchemas(fs, "")

	// Act
	err := s.WalkFrom(1, fn)

	// Assert
	require.NoError(t, err)
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
