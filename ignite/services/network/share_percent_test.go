package network

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSharePercentages(t *testing.T) {
	tests := []struct {
		name     string
		shareStr string
		want     SharePercents
		err      error
	}{
		{
			name:     "valid share percentage",
			shareStr: "12.333%def",
			want: SharePercents{
				{
					denom:       "def",
					nominator:   12333,
					denominator: 100000,
				},
			},
		},
		{
			name:     "valid share percentage",
			shareStr: "0.333%def",
			want: SharePercents{
				{
					denom:       "def",
					nominator:   333,
					denominator: 100000,
				},
			},
		},
		{
			name:     "extra zeroes",
			shareStr: "12.33300%def",
			want: SharePercents{
				{
					denom:       "def",
					nominator:   12333,
					denominator: 100000,
				},
			},
		},
		{
			name:     "100% percentage",
			shareStr: "100%def",
			want: SharePercents{
				{
					denom:       "def",
					nominator:   100,
					denominator: 100,
				},
			},
		},
		{
			name:     "valid share percentages",
			shareStr: "12%def,10.3%abc",
			want: SharePercents{
				{
					denom:       "def",
					nominator:   12,
					denominator: 100,
				},
				{
					denom:       "abc",
					nominator:   103,
					denominator: 1000,
				},
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
			result, err := ParseSharePercents(tt.shareStr)
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
		percent SharePercent
		total   uint64
		want    sdk.Coin
		err     error
	}{
		{
			name: "100 fraction",
			percent: SharePercent{
				denom:       "foo",
				nominator:   10,
				denominator: 100,
			},
			total: 10000,
			want:  sdk.NewInt64Coin("foo", 1000),
		},
		{
			name: "1000 fraction",
			percent: SharePercent{
				denom:       "foo",
				nominator:   133,
				denominator: 1000,
			},
			total: 10000,
			want:  sdk.NewInt64Coin("foo", 1330),
		},
		{
			name: "10000 fraction",
			percent: SharePercent{
				denom:       "foo",
				nominator:   297,
				denominator: 10000,
			},
			total: 10000,
			want:  sdk.NewInt64Coin("foo", 297),
		},
		{
			name: "non integer share",
			percent: SharePercent{
				denom:       "foo",
				nominator:   297,
				denominator: 10001,
			},
			total: 10000,
			want:  sdk.NewInt64Coin("foo", 297),
			err:   errors.New("foo share from total 10000 is not integer: 296.970303"),
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

func SampleSharePercent(t *testing.T, denom string, nominator, denominator uint64) SharePercent {
	sp, err := NewSharePercent(denom, nominator, denominator)
	assert.NoError(t, err)
	return sp
}
