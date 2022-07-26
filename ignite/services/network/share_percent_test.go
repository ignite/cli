package network_test

import (
	"errors"
	"testing"

	"github.com/ignite/cli/ignite/services/network"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParseSharePercentages(t *testing.T) {
	tests := []struct {
		name     string
		shareStr string
		want     network.SharePercents
		err      error
	}{
		{
			name:     "valid share percentage",
			shareStr: "12.333%def",
			want: network.SharePercents{
				network.SampleSharePercent(t, "def", 12333, 100000),
			},
		},
		{
			name:     "valid share percentage",
			shareStr: "0.333%def",
			want: network.SharePercents{
				network.SampleSharePercent(t, "def", 333, 100000),
			},
		},
		{
			name:     "extra zeroes",
			shareStr: "12.33300%def",
			want: network.SharePercents{
				network.SampleSharePercent(t, "def", 12333, 100000),
			},
		},
		{
			name:     "100% percentage",
			shareStr: "100%def",
			want: network.SharePercents{
				network.SampleSharePercent(t, "def", 100, 100),
			},
		},
		{
			name:     "valid share percentages",
			shareStr: "12%def,10.3%abc",
			want: network.SharePercents{
				network.SampleSharePercent(t, "def", 12, 100),
				network.SampleSharePercent(t, "abc", 103, 1000),
			},
		},
		{
			name:     "share percentages greater than 100",
			shareStr: "12%def,10.3abc",
			err:      errors.New("invalid percentage format 10.3abc"),
		},
		{
			name:     "share percentages without % sign",
			shareStr: "12%def,103%abc",
			err:      errors.New("\"abc\" can not be bigger than 100"),
		},
		{
			name:     "invalid percent",
			shareStr: "12.3d3%def",
			err:      errors.New("invalid percentage format 12.3d3%def"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := network.ParseSharePercents(tt.shareStr)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestShare(t *testing.T) {
	tests := []struct {
		name    string
		percent network.SharePercent
		total   uint64
		want    sdk.Coin
		err     error
	}{
		{
			name:    "100 fraction",
			percent: network.SampleSharePercent(t, "foo", 10, 100),
			total:   10000,
			want:    sdk.NewInt64Coin("foo", 1000),
		},
		{
			name:    "1000 fraction",
			percent: network.SampleSharePercent(t, "foo", 133, 1000),
			total:   10000,
			want:    sdk.NewInt64Coin("foo", 1330),
		},
		{
			name:    "10000 fraction",
			percent: network.SampleSharePercent(t, "foo", 297, 10000),
			total:   10000,
			want:    sdk.NewInt64Coin("foo", 297),
		},
		{
			name:    "non integer share",
			percent: network.SampleSharePercent(t, "foo", 297, 10001),
			total:   10000,
			want:    sdk.NewInt64Coin("foo", 297),
			err:     errors.New("foo share from total 10000 is not integer: 296.970303"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.percent.Share(tt.total)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, result)
		})
	}
}
