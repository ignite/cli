package cosmosutil_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosutil"
)

func TestChangePrefix(t *testing.T) {
	tests := []struct {
		name    string
		address string
		prefix  string
		want    string
		wantErr bool
	}{
		{
			name:    "cosmos address to spn address",
			address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
			prefix:  "spn",
			want:    "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
		},
		{
			name:    "cosmos address to spn address 2",
			address: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
			prefix:  "spn",
			want:    "spn1mmlqwyqk7neqegffp99q86eckpm4pjahdcne08",
		},
		{
			name:    "cosmos validator address",
			address: "cosmosvaloper1mmlqwyqk7neqegffp99q86eckpm4pjah5sl2dw",
			prefix:  "spn",
			want:    "spn1mmlqwyqk7neqegffp99q86eckpm4pjahdcne08",
		},
		{
			name:    "mars address to earth address",
			address: "mars1c6ac48k2ur8tl3tf0cpntlw5068kvp8xf4xq37",
			prefix:  "earth",
			want:    "earth1c6ac48k2ur8tl3tf0cpntlw5068kvp8x0xyl2v",
		},
		{
			name:    "invalid bech32 address",
			address: "mars1c6ac48k2ur9tl3tf0cpntlw5068kvp8xf4xq37",
			prefix:  "spn",
			wantErr: true,
		},
		{
			name:    "empty target prefix",
			address: "mars1c6ac48k2ur8tl3tf0cpntlw5068kvp8xf4xq37",
			prefix:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cosmosutil.ChangeAddressPrefix(tt.address, tt.prefix)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetPrefix(t *testing.T) {
	prefix, err := cosmosutil.GetAddressPrefix("cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj")
	require.Equal(t, "cosmos", prefix)
	require.NoError(t, err)

	prefix, err = cosmosutil.GetAddressPrefix("mars1c6ac48k2ur8tl3tf0cpntlw5068kvp8xf4xq37")
	require.Equal(t, "mars", prefix)
	require.NoError(t, err)

	// invalid bech32 address
	_, err = cosmosutil.GetAddressPrefix("mars1c6ac48k2ur9tl3tf0cpntlw5068kvp8xf4xq37")
	require.Error(t, err)
}
