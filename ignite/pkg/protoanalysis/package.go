package protoanalysis

import (
	"regexp"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type (
	// Packages represents slice of Package.
	Packages []Package

	PkgName string

	// Package represents a proto pkg.
	Package struct {
		// Name of the proto pkg.
		Name string `json:"name,omitempty"`

		// Path of the package in the fs.
		Path string `json:"path,omitempty"`

		// Files is a list of .proto files in the package.
		Files Files `json:"files,omitempty"`

		// GoImportName is the go package name of proto package.
		GoImportName string `json:"go_import_name,omitempty"`

		// Messages is a list of proto messages defined in the package.
		Messages []Message `json:"messages,omitempty"`

		// Services is a list of RPC services.
		Services []Service `json:"services,omitempty"`
	}
)

var regexBetaVersion = regexp.MustCompile("^v[0-9]+(beta|alpha)[0-9]+")

// ModuleName retrieves the single module name of the package.
func (p Package) ModuleName() (name string) {
	names := strings.Split(p.Name, ".")
	for i := len(names) - 1; i >= 0; i-- {
		name = names[i]
		if !semver.IsValid(name) && !regexBetaVersion.MatchString(name) {
			break
		}
	}
	return
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

// Files retrieves the files from the package list.
func (p Packages) Files() Files {
	var files []File
	for _, pkg := range p {
		files = append(files, pkg.Files...)
	}
	return files
}

type (
	Files []File

	File struct {
		// Path of the file.
		Path string `json:"path,omitempty"`

		// Dependencies is a list of imported proto packages.
		Dependencies []string `json:"dependencies,omitempty"`
	}
)

// Paths retrieves the list of paths from the files.
func (f Files) Paths() []string {
	var paths []string
	for _, ff := range f {
		paths = append(paths, ff.Path)
	}
	return paths
}

type (
	// Message represents a proto message.
	Message struct {
		// Name of the message.
		Name string `json:"name,omitempty"`

		// Path of the proto file where the message is defined.
		Path string `json:"path,omitempty"`

		// HighestFieldNumber is the highest field number among fields of the message.
		// This allows to determine new field number when writing to proto message.
		HighestFieldNumber int `json:"highest_field_number,omitempty"`

		// Fields contains message's field names and types.
		Fields map[string]string `json:"fields,omitempty"`
	}

	// Service is an RPC service.
	Service struct {
		// Name of the services.
		Name string `json:"name,omitempty"`

		// RPCFuncs is a list of RPC funcs of the service.
		RPCFuncs []RPCFunc `json:"functions,omitempty"`
	}

	// RPCFunc is an RPC func.
	RPCFunc struct {
		// Name of the RPC func.
		Name string `json:"name,omitempty"`

		// RequestType is the request type of RPC func.
		RequestType string `json:"request_type,omitempty"`

		// ReturnsType is the response type of RPC func.
		ReturnsType string `json:"return_type,omitempty"`

		// HTTPRules keeps info about http rules of an RPC func.
		// spec:
		//   https://github.com/googleapis/googleapis/blob/master/google/api/http.proto.
		HTTPRules []HTTPRule `json:"http_rules,omitempty"`

		// Paginated indicates that the RPC function is using pagination.
		Paginated bool `json:"paginated,omitempty"`
	}

	// HTTPRule keeps info about a configured http rule of an RPC func.
	HTTPRule struct {
		// Params is a list of parameters defined in the HTTP endpoint itself.
		Params []string `json:"params,omitempty"`

		// HasQuery indicates if there is a request query.
		HasQuery bool `json:"has_query,omitempty"`

		// HasBody indicates if there is a request payload.
		HasBody bool `json:"has_body,omitempty"`
	}
)
