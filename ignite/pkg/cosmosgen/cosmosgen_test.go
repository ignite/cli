package cosmosgen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/module"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
)

func TestVuexStoreModulePath(t *testing.T) {
	modulePath := VuexStoreModulePath("prefix")

	cases := []struct {
		name         string
		goModulePath string
		protoPkgName string
		want         string
	}{
		{
			name:         "github uri",
			goModulePath: "github.com/owner/app",
			protoPkgName: "owner.app.module",
			want:         "prefix/owner/app/owner.app.module/module",
		},
		{
			name:         "short uri",
			goModulePath: "domain.com/app",
			protoPkgName: "app.module",
			want:         "prefix/app/app.module/module",
		},
		{
			name:         "path",
			goModulePath: "owner/app",
			protoPkgName: "owner.app.module",
			want:         "prefix/owner/app/owner.app.module/module",
		},
		{
			name:         "name",
			goModulePath: "app",
			protoPkgName: "app.module",
			want:         "prefix/app/app.module/module",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			m := module.Module{
				GoModulePath: tt.goModulePath,
				Pkg: protoanalysis.Package{
					Name: tt.protoPkgName,
				},
			}

			require.Equal(t, tt.want, modulePath(m))
		})
	}
}
