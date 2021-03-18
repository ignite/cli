package protoanalysis

import (
	"strings"

	"github.com/pkg/errors"
)

// Package represents a proto pkg.
type Package struct {
	// Name of the proto pkg.
	Name string

	// Path of the package in the fs.
	Path string

	// GoImportName is the go package name of proto package.
	GoImportName string

	// Messages is a list of proto messages defined in the package.
	Messages []Message

	// Services is a list of RPC services.
	Services []Service
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

	// FieldCount is the top level field count of message.
	//FieldCount int
}

// Service is an RPC service.
type Service struct {
	// Name of the services.
	Name string

	// RPC is a list of RPC funcs of the service.
	RPCFuncs []RPCFunc
}

// RPC is an RPC func.
type RPCFunc struct {
	// Name of the RPC func.
	Name string

	// RequestType is the request type of RPC func.
	RequestType string

	// ReturnsType is the response type of RPC func.
	ReturnsType string

	// HTTPRules keeps info about http rules of an RPC func.
	//
	// a single func might have multiple endpoints configures as per the spec. therefore,
	// it is an in-line slice.
	// see: https://github.com/googleapis/googleapis/blob/ca1372c/google/api/http.proto#L182-L198
	HTTPRules []HTTPRule
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
