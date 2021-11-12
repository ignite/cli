package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetSPNPrefix(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
		wantErr bool
	}{
		{
			name:    "cosmos address 1",
			address: "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
			want:    "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
		}, {
			name:    "cosmos address 2",
			address: "cosmos1mmlqwyqk7neqegffp99q86eckpm4pjah3ytlpa",
			want:    "spn1mmlqwyqk7neqegffp99q86eckpm4pjahdcne08",
		}, {
			name:    "cosmos validator address",
			address: "cosmosvaloper1mmlqwyqk7neqegffp99q86eckpm4pjah5sl2dw",
			want:    "spn1mmlqwyqk7neqegffp99q86eckpm4pjahdcne08",
		}, {
			name:    "mars address",
			address: "mars1az6r084s0l95p7cj5es72f3cnhwdmz09y4zhna",
			want:    "spn1az6r084s0l95p7cj5es72f3cnhwdmz0995rggu",
		}, {
			name:    "invalid address",
			address: "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetSPNPrefix(tt.address)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
