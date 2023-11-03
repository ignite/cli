package protoanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackage_ModuleName(t *testing.T) {
	tests := []struct {
		name string
		p    Package
		want string
	}{
		{
			name: "test single name",
			p:    Package{Name: "staking"},
			want: "staking",
		},
		{
			name: "test two names",
			p:    Package{Name: "cosmos.staking"},
			want: "staking",
		},
		{
			name: "test three name",
			p:    Package{Name: "cosmos.ignite.staking"},
			want: "staking",
		},
		{
			name: "test with the version 1",
			p:    Package{Name: "cosmos.staking.v1"},
			want: "staking",
		},
		{
			name: "test with the version 2",
			p:    Package{Name: "cosmos.staking.v2"},
			want: "staking",
		},
		{
			name: "test with the version 10",
			p:    Package{Name: "cosmos.staking.v10"},
			want: "staking",
		},
		{
			name: "test with the version 1 beta 1",
			p:    Package{Name: "cosmos.staking.v1beta1"},
			want: "staking",
		},
		{
			name: "test with the version 1 beta 2",
			p:    Package{Name: "cosmos.staking.v1beta2"},
			want: "staking",
		},
		{
			name: "test with the version 2 beta 1",
			p:    Package{Name: "cosmos.staking.v2beta1"},
			want: "staking",
		},
		{
			name: "test with the version 2 beta 2",
			p:    Package{Name: "cosmos.staking.v2beta2"},
			want: "staking",
		},
		{
			name: "test with the version 3 alpha 5",
			p:    Package{Name: "cosmos.staking.v3alpha5"},
			want: "staking",
		},
		{
			name: "test with the wrong version",
			p:    Package{Name: "cosmos.staking.v3bank5"},
			want: "v3bank5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.ModuleName()
			require.Equal(t, tt.want, got)
		})
	}
}
