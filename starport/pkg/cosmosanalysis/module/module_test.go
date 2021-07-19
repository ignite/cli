package module

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kr/pretty"
)

func TestExampleDiscover(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	t.Run("test folder without proto files", func(t *testing.T) {
		pretty.Println(Discover(
			context.Background(),
			filepath.Join(goPath, "src/github.com/tendermint/starport/local_test/mars"),
			"proto",
		))
	})
	t.Run("test folder with proto files", func(t *testing.T) {
		pretty.Println(Discover(
			context.Background(),
			filepath.Join(goPath, "src/github.com/tendermint/starport"),
			"starport/pkg/protoc/data/include/google/protobuf",
		))
	})
}
