package cosmosutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyPeerFormat(t *testing.T) {
	tests := []struct {
		peer string
		want bool
	}{
		{
			peer: "nodeid@peer",
			want: true,
		},
		{
			peer: "@peer",
			want: false,
		},
		{
			peer: "nodeid@",
			want: false,
		},
		{
			peer: "nodeid",
			want: false,
		},
		{
			peer: "@",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run("verifying "+tt.peer, func(t *testing.T) {
			got := VerifyPeerFormat(tt.peer)
			require.Equal(t, tt.want, got)
		})
	}
}
