package xgenny_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/xgenny"
)

var (
	modifiedExample = []string{"mfoo", "mbar", "mfoobar"}
	createdExample  = []string{"cfoo", "cbar", "cfoobar"}
)

func sourceModificationExample() xgenny.SourceModification {
	sourceModification := xgenny.NewSourceModification()
	sourceModification.AppendModifiedFiles(modifiedExample...)
	sourceModification.AppendCreatedFiles(createdExample...)
	return sourceModification
}

func TestNewSourceModification(t *testing.T) {
	sm := xgenny.NewSourceModification()
	require.Empty(t, sm.ModifiedFiles())
	require.Empty(t, sm.CreatedFiles())
}

func TestModifiedFiles(t *testing.T) {
	sm := sourceModificationExample()
	require.Len(t, sm.ModifiedFiles(), len(modifiedExample))
	require.Subset(t, sm.ModifiedFiles(), modifiedExample)
}

func TestCreatedFiles(t *testing.T) {
	sm := sourceModificationExample()
	require.Len(t, sm.CreatedFiles(), len(createdExample))
	require.Subset(t, sm.CreatedFiles(), createdExample)
}

func TestAppendModifiedFiles(t *testing.T) {
	sm := sourceModificationExample()
	sm.AppendModifiedFiles("foo1")
	require.Len(t, sm.ModifiedFiles(), len(modifiedExample)+1)
	require.Contains(t, sm.ModifiedFiles(), "foo1")

	// Do not append a existing element
	sm.AppendModifiedFiles("foo1")
	require.Len(t, sm.ModifiedFiles(), len(modifiedExample)+1)
	sm.AppendCreatedFiles("foo2")
	sm.AppendModifiedFiles("foo2")
	require.Len(t, sm.ModifiedFiles(), len(modifiedExample)+1)
}

func TestAppendCreatedFiles(t *testing.T) {
	sm := sourceModificationExample()
	sm.AppendCreatedFiles("foo1")
	require.Len(t, sm.CreatedFiles(), len(createdExample)+1)
	require.Contains(t, sm.CreatedFiles(), "foo1")

	// Do not append a existing element
	sm.AppendCreatedFiles("foo1")
	require.Len(t, sm.CreatedFiles(), len(createdExample)+1)
	sm.AppendModifiedFiles("foo2")
	sm.AppendCreatedFiles("foo2")
	require.Len(t, sm.ModifiedFiles(), len(modifiedExample)+1)
}

func TestMerge(t *testing.T) {
	sm1 := xgenny.NewSourceModification()
	sm2 := xgenny.NewSourceModification()

	sm1.AppendModifiedFiles("foo1", "foo2", "foo3")
	sm2.AppendModifiedFiles("foo3", "foo4", "foo5")
	sm1.AppendCreatedFiles("bar1", "bar2", "bar3")
	sm2.AppendCreatedFiles("foo1", "bar2", "bar3")

	sm1.Merge(sm2)
	require.Len(t, sm1.ModifiedFiles(), 5)
	require.Len(t, sm1.CreatedFiles(), 3)
	require.Subset(t, sm1.ModifiedFiles(), []string{"foo1", "foo2", "foo3", "foo4", "foo5"})
	require.Subset(t, sm1.CreatedFiles(), []string{"bar1", "bar2", "bar3"})
}
