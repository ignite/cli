package giturl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAndRepo(t *testing.T) {
	require.Equal(t, "tendermint/starport", UserAndRepo("http://github.com/tendermint/starport/a/b"))
}
