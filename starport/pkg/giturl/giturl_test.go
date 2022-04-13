package giturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		url   string
		want  []string
		error bool
	}{
		{
			name: "url",
			url:  "http://github.com/tendermint/starport",
			want: []string{"github.com", "tendermint", "starport", "tendermint/starport"},
		},
		{
			name:  "invalid url",
			url:   "http://github.com/tendermint",
			error: true,
		},
		{
			name: "url without scheme",
			url:  "github.com/tendermint/starport",
			want: []string{"github.com", "tendermint", "starport", "tendermint/starport"},
		},
		{
			name:  "invalid url without scheme",
			url:   "github.com/tendermint",
			error: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			g, err := Parse(tt.url)

			if tt.error {
				require.ErrorIs(t, err, ErrInvalidURL)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, []string{
					g.Host,
					g.User,
					g.Repo,
					g.UserAndRepo(),
				})
			}
		})
	}
}
