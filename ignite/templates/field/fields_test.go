package field

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

func TestFieldsCustom(t *testing.T) {
	nameA, err := multiformatname.NewName("customA")
	require.NoError(t, err)
	nameB, err := multiformatname.NewName("customB")
	require.NoError(t, err)
	nameC, err := multiformatname.NewName("customC")
	require.NoError(t, err)

	fields := Fields{
		{
			Name:         nameA,
			DatatypeName: datatype.Custom,
			Datatype:     "ProductDetails",
		},
		{
			Name:         nameB,
			DatatypeName: datatype.CustomSlice,
			Datatype:     "LineItem",
		},
		{
			Name:         nameC,
			DatatypeName: datatype.String,
		},
	}

	require.Equal(t, []string{"product_details", "line_item"}, fields.Custom())
}
