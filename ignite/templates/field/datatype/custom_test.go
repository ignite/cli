package datatype_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

func TestDataCustomToProtoField(t *testing.T) {
	field := datatype.DataCustom.ToProtoField("Foo", "foo", 2)
	expected := protoutil.NewField(
		"foo",
		"Foo",
		2,
		protoutil.WithFieldOptions(
			protoutil.NewOption("gogoproto.nullable", "false", protoutil.Custom()),
		),
	)

	require.Equal(t, expected, field)
}
