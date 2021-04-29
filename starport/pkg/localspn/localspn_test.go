package localspn

import (
	"context"
	"testing"
)

func TestStartSPN(t *testing.T) {
	StartSPN(context.TODO(), WithBranch("develop"))
}
