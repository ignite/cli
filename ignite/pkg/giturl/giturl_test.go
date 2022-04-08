package giturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parsed, err := Parse("http://github.com/tendermint/starport/a/b")
	require.NoError(t, err)
	require.Equal(t, "github.com", parsed.Host)
	require.Equal(t, "tendermint", parsed.User)
	require.Equal(t, "starport", parsed.Repo)
	require.Equal(t, "tendermint/starport", parsed.UserAndRepo())
}
