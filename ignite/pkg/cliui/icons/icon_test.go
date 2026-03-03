package icons

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIconsAreInitialized(t *testing.T) {
	require.NotEmpty(t, Earth)
	require.NotEmpty(t, CD)
	require.NotEmpty(t, User)
	require.NotEmpty(t, Tada)
	require.NotEmpty(t, Survey)
	require.NotEmpty(t, Announcement)

	require.Contains(t, OK, "✔")
	require.Contains(t, NotOK, "✘")
	require.Contains(t, Bullet, "⋆")
	require.Contains(t, Info, "𝓲")
}
