package module

import (
	"context"
	"testing"

	"github.com/kr/pretty"
)

func TestExampleDiscover(t *testing.T) {
	pretty.Println(Discover(context.Background(), "/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/mars"))
}
