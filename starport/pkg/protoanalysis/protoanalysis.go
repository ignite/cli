// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emicklei/proto"
	"github.com/mattn/go-zglob"
	"golang.org/x/sync/errgroup"
)

const (
	optionGoPkg = "go_package"
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

// Service is an RPC service.
type Service struct {
	// Name of the services.
	Name string

	// RPCFuncs is a list of RPC funcs of the service.
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

// Message represents a proto message.
type Message struct {
	// Name of the message.
	Name string

	// Path of the file where message is defined at.
	Path string
}

func (p Package) GoImportPath() string {
	return strings.Split(p.GoImportName, ";")[0]
}

// DiscoverPackages recursively discovers proto files, parses them, and returns info about
// each found package.
func DiscoverPackages(path string) ([]Package, error) {
	files, err := SearchProto(path)
	if err != nil {
		return nil, err
	}

	var (
		// m protects pkgs.
		m    sync.Mutex
		pkgs []Package
	)

	g := &errgroup.Group{}

	for _, path := range files {
		path := path

		g.Go(func() error {
			pkg, err := Parse(path)
			if err != nil {
				return err
			}

			m.Lock()
			defer m.Unlock()

			var (
				exists bool
				index  int
			)
			for i, epkg := range pkgs {
				if epkg.Name == pkg.Name {
					exists = true
					index = i
				}
			}
			if !exists {
				pkgs = append(pkgs, pkg)
			} else {
				pkgs[index].Messages = append(pkgs[index].Messages, pkg.Messages...)
				pkgs[index].Services = append(pkgs[index].Services, pkg.Services...)
			}

			return nil
		})
	}

	return pkgs, g.Wait()
}

// Parse parses a proto file residing at path.
func Parse(path string) (Package, error) {
	f, err := os.Open(path)
	if err != nil {
		return Package{}, err
	}
	defer f.Close()

	def, err := proto.NewParser(f).Parse()
	if err != nil {
		return Package{}, err
	}

	pkg := Package{
		Path: filepath.Dir(path),
	}

	proto.Walk(
		def,
		proto.WithPackage(func(p *proto.Package) { pkg.Name = p.Name }),
		proto.WithOption(func(o *proto.Option) {
			if o.Name != optionGoPkg {
				return
			}
			pkg.GoImportName = o.Constant.Source
		}),
		proto.WithMessage(func(m *proto.Message) {
			pkg.Messages = append(pkg.Messages, Message{
				Name: m.Name,
				Path: path,
			})
		}),
		proto.WithService(func(s *proto.Service) {
			sv := Service{
				Name: s.Name,
			}

			for _, el := range s.Elements {
				rpc, ok := el.(*proto.RPC)
				if !ok {
					continue
				}
				sv.RPCFuncs = append(sv.RPCFuncs, RPCFunc{
					Name:        rpc.Name,
					RequestType: rpc.RequestType,
					ReturnsType: rpc.ReturnsType,
				})
			}

			pkg.Services = append(pkg.Services, sv)
		}),
	)

	return pkg, nil
}

// SearchProto recursively finds all proto files under path.
func SearchProto(path string) ([]string, error) {
	return zglob.Glob(GlobPattern(path))
}

// GlobPattern returns a recursive glob search pattern to find all proto files under path.
func GlobPattern(path string) string { return path + "/**/*.proto" }
