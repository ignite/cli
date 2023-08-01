package protoutil

import (
	"io"
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/protofmt"
)

// ParseProtoPath opens the file denoted by path and parses it
// into a proto file.
func ParseProtoPath(path string) (pf *proto.Proto, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return proto.NewParser(f).Parse()
}

// ParseProtoFile parses the given file.
func ParseProtoFile(r io.Reader) (*proto.Proto, error) {
	return proto.NewParser(r).Parse()
}

// Print formats the proto file using proto-contrib/pkg/protofmt.
// This does have certain opinions on how formatting is done.
func Print(pf *proto.Proto) string {
	output := new(strings.Builder)
	protofmt.NewFormatter(output, "  ").Format(pf) // 2 spaces

	return output.String()
}
