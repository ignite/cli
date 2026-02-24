package modulecreate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddIBCModuleRoute(t *testing.T) {
	content := readFixture(t, "../../app/files/app/ibc.go.plush")

	modified, err := addIBCModuleRoute(content, "blog")
	require.NoError(t, err)

	routeCall := "ibcRouter=ibcRouter.AddRoute(blogmoduletypes.ModuleName,blogmodule.NewIBCModule(app.appCodec,app.BlogKeeper))"
	require.Equal(t, 1, strings.Count(normalizedExpr(modified), routeCall))

	modified, err = addIBCModuleRoute(modified, "blog")
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(normalizedExpr(modified), routeCall))
}
