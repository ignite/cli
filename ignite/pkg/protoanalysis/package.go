package protoanalysis

import (
	"strings"

	"github.com/pkg/errors"
)

type Packages []Package

func (p Packages) Files() Files {
	var files []File
	for _, pkg := range p {
		files = append(files, pkg.Files...)
	}
	return files
}

// Package represents a proto pkg.
type Package struct {
	// Name of the proto pkg.
	Name string

	// Path of the package in the fs.
	Path string

	// Files is a list of .proto files in the package.
	Files Files

	// GoImportName is the go package name of proto package.
	GoImportName string

	// Messages is a list of proto messages defined in the package.
	Messages []Message

	// Services is a list of RPC services.
	Services []Service
}

type Files []File

type File struct {
	// Path of the file.
	Path string

	// Dependencies is a list of imported .proto files in this package.
	Dependencies []string
}

func (f Files) Paths() []string {
	var paths []string
	for _, ff := range f {
		paths = append(paths, ff.Path)
	}
	return paths
}

// MessageByName finds a message by its name inside Package.
func (p Package) MessageByName(name string) (Message, error) {
	for _, message := range p.Messages {
		if message.Name == name {
			return message, nil
		}
	}
	return Message{}, errors.New("no message found")
}

// GoImportPath retrieves the Go import path.
func (p Package) GoImportPath() string {
	return strings.Split(p.GoImportName, ";")[0]
}

// Message represents a proto message.
type Message struct {
	// Name of the message.
	Name string

	// Path of the file where message is defined at.
	Path string

	// HighestFieldNumber is the highest field number among fields of the message
	// This allows to determine new field number when writing to proto message
	HighestFieldNumber int

	// Fields contains message's field names and types
	Fields map[string]string
}

// Service is an RPC service.
type Service struct {
	// Name of the services.
	Name string

	// RPC is a list of RPC funcs of the service.
	RPCFuncs []RPCFunc
}

// RPCFunc is an RPC func.
type RPCFunc struct {
	// Name of the RPC func.
	Name string

	// RequestType is the request type of RPC func.
	RequestType string

	// ReturnsType is the response type of RPC func.
	ReturnsType string

	// HTTPRules keeps info about http rules of an RPC func.
	// spec:
	//   https://github.com/googleapis/googleapis/blob/master/google/api/http.proto.
	HTTPRules []HTTPRule

	// Paginated indicates that the RPC function is using pagination.
	Paginated bool
}

// HTTPRule keeps info about a configured http rule of an RPC func.
type HTTPRule struct {
	// Params is a list of parameters defined in the http endpoint itself.
	Params []string

	// HasQuery indicates if there is a request query.
	HasQuery bool

	// HasBody indicates if there is a request payload.
	HasBody bool
}
