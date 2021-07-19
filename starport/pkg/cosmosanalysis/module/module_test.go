package module

import (
	"context"
	"testing"

	"github.com/kr/pretty"
)

func TestExampleDiscover(t *testing.T) {
	pretty.Println(Discover(
		context.Background(),
		"$GOPATH/github.com/tendermint/starport/local_test/mars",
		"proto",
	))
}
