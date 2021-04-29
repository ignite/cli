package localspn

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetupSPN(t *testing.T) {
	t.SkipNow()
	_, err := SetupSPN(context.TODO(), WithBranch("develop"))
	require.NoError(t, err)
}
