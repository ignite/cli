package giturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parsed, err := Parse("http://github.com/tendermint/starport/a/b")
	require.NoError(t, err)
	require.Equal(t, "github.com", parsed.Host)
	require.Equal(t, "ignite-hq", parsed.User)
	require.Equal(t, "cli", parsed.Repo)
	require.Equal(t, "ignite-hq/cli", parsed.UserAndRepo())
}
