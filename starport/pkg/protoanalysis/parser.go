package protoanalysis

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/emicklei/proto"
	"github.com/mattn/go-zglob"
)

const optionGoPkg = "go_package"

type parser struct {
	m        sync.Mutex // protects following.
	packages []*pkg
}

func newParser() *parser {
	return &parser{}
}

type pkg struct {
	name  string
	dir   string
	files []file
}

func (p *pkg) options() (o []*proto.Option) {
	for _, f := range p.files {
		o = append(o, f.options...)
	}

	return
}

func (p *pkg) messages() (m []*proto.Message) {
	for _, f := range p.files {
		m = append(m, f.messages...)
	}

	return
}

func (p *pkg) services() (s []*proto.Service) {
	for _, f := range p.files {
		s = append(s, f.services...)
	}

	return
}

type file struct {
	path     string
	pkg      *proto.Package
	options  []*proto.Option
	messages []*proto.Message
	services []*proto.Service
}

func (p *parser) parse(ctx context.Context, pattern string) error {
	paths, err := zglob.Glob(pattern)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if err := p.parseFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) parseFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	def, err := proto.NewParser(f).Parse()
	if err != nil {
		return err
	}

	p.m.Lock()
	defer p.m.Unlock()

	var pkgName string

	proto.Walk(
		def,
		proto.WithPackage(func(p *proto.Package) { pkgName = p.Name }),
	)

	var pp *pkg
	for _, v := range p.packages {
		if pkgName == v.name {
			pp = v
			break
		}
	}
	if pp == nil {
		pp = &pkg{
			name: pkgName,
			dir:  filepath.Dir(path),
		}
		p.packages = append(p.packages, pp)
	}

	pf := file{
		path: path,
	}

	proto.Walk(
		def,
		proto.WithPackage(func(p *proto.Package) { pf.pkg = p }),
		proto.WithOption(func(o *proto.Option) { pf.options = append(pf.options, o) }),
		proto.WithMessage(func(m *proto.Message) { pf.messages = append(pf.messages, m) }),
		proto.WithService(func(s *proto.Service) { pf.services = append(pf.services, s) }),
	)

	pp.files = append(pp.files, pf)

	return nil
}
