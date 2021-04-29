package localspn

import (
	"context"
	"testing"
)

func TestSetupSPN(t *testing.T) {
	SetupSPN(context.TODO(), WithBranch("develop"))
}
