package starportconf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseUserConfig(t *testing.T) {
	confyml := `
accounts:
  - name: me
    coins: ["1000token", "100000000stake"]
  - name: you
    coins: ["5000token"]
`

	conf, err := ParseUserConfig(strings.NewReader(confyml))
	require.NoError(t, err)
	require.Equal(t, UserConfig{
		Accounts: []Account{
			{
				Name:  "me",
				Coins: []string{"1000token", "100000000stake"},
			},
			{
				Name:  "you",
				Coins: []string{"5000token"},
			},
		},
	}, conf)
}
