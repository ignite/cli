package protoutil

import (
	"io"
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/protofmt"
)

// TODO: Maybe allow a cache?

// ParseProtoFile opens the file denoted by path and parses it
// into a proto file.
func ParseProtoPath(path string) (pf *proto.Proto, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return proto.NewParser(f).Parse()
}

// PraseProtoFile parses the given file.
func ParseProtoFile(r io.Reader) (*proto.Proto, error) {
	return proto.NewParser(r).Parse()
}

// String() formats the proto file using proto-contrib/pkg/protofmt.
// This does have certain opinions on how formatting is done.
func Printer(pf *proto.Proto) string {
	output := new(strings.Builder)
	protofmt.NewFormatter(output, "  ").Format(pf) // 2 spaces

	return output.String()
}

// parseStringProto takes a string, parses it into a proto.File, and returns a ProtoFile.
// Nodes can be created easily (newnode) by wrapping them correctly. (e.g field in a message)
func parseStringProto(s string) (*proto.Proto, error) {
	p, err := proto.NewParser(strings.NewReader(s)).Parse()
	if err != nil {
		return nil, err
	}

	return p, nil
}
