// Package protoanalysis provides a toolset for analyzing proto files and packages.
package protoanalysis

import (
	"os"
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

	// GoImportName is the go package name of proto package.
	GoImportName string
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

		isPkgExists = func(pkg Package) bool {
			for _, epkg := range pkgs {
				if pkg == epkg {
					return true
				}
			}
			return false
		}
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

			if !isPkgExists(pkg) {
				pkgs = append(pkgs, pkg)
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

	var pkg Package

	proto.Walk(
		def,
		proto.WithPackage(func(p *proto.Package) { pkg.Name = p.Name }),
		proto.WithOption(func(o *proto.Option) {
			if o.Name != optionGoPkg {
				return
			}
			pkg.GoImportName = o.Constant.Source
		}))

	return pkg, nil
}

// SearchProto recursively finds all proto files under path.
func SearchProto(path string) ([]string, error) {
	return zglob.Glob(GlobPattern(path))
}

// GlobPattern returns a recursive glob search pattern to find all proto files under path.
func GlobPattern(path string) string { return path + "/**/*.proto" }
