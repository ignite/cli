package xgenny

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packd"

	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/randstr"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
)

type Runner struct {
	*genny.Runner
	ctx     context.Context
	tracer  *placeholder.Tracer
	tmpPath string
	path    string
	sm      SourceModification
}

// NewRunner is a xgenny Runner with a logger.
func NewRunner(ctx context.Context, appPath string) *Runner {
	var (
		runner  = genny.WetRunner(ctx)
		tmpPath = filepath.Join(os.TempDir(), randstr.Runes(5))
	)
	runner.FileFn = func(f genny.File) (genny.File, error) {
		return wetFileFn(f, tmpPath, appPath)
	}
	return &Runner{
		ctx:     ctx,
		Runner:  runner,
		tracer:  placeholder.New(),
		path:    appPath,
		tmpPath: tmpPath,
	}
}

func (r *Runner) Tracer() *placeholder.Tracer {
	return r.tracer
}

func (r *Runner) ApplyModifications() (SourceModification, error) {
	return r.sm, xos.CopyFolder(r.tmpPath, r.path)
}

// Run all generators into a temp folder for we can apply the modifications later.
func (r *Runner) Run(gens ...*genny.Generator) error {
	// execute the modification with a wet runner
	for _, gen := range gens {
		if err := r.With(gen); err != nil {
			return err
		}
		if err := r.Runner.Run(); err != nil {
			return err
		}
	}
	if err := r.tracer.Err(); err != nil {
		return err
	}

	// fetch the source modification
	sm := NewSourceModification()
	for _, file := range r.Results().Files {
		fileName := file.Name()
		_, err := os.Stat(fileName)

		//nolint:gocritic
		if os.IsNotExist(err) {
			// if the file doesn't exist in the source, it means it has been created by the runner
			sm.AppendCreatedFiles(fileName)
		} else if err != nil {
			return err
		} else {
			// the file has been modified by the runner
			sm.AppendModifiedFiles(fileName)
		}
	}
	r.sm = sm
	return nil
}

func wetFileFn(f genny.File, tmpPath, appPath string) (genny.File, error) {
	if d, ok := f.(genny.Dir); ok {
		if err := os.MkdirAll(d.Name(), d.Perm); err != nil {
			return f, err
		}
		return d, nil
	}

	var err error
	if !filepath.IsAbs(appPath) {
		appPath, err = filepath.Abs(appPath)
		if err != nil {
			return f, err
		}
	}

	name := f.Name()
	if !filepath.IsAbs(name) {
		name = filepath.Join(appPath, name)
	}
	relPath, err := filepath.Rel(appPath, name)
	if err != nil {
		return f, err
	}

	dstPath := filepath.Join(tmpPath, relPath)
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return f, err
	}
	ff, err := os.Create(dstPath)
	if err != nil {
		return f, err
	}
	defer ff.Close()
	if _, err := io.Copy(ff, f); err != nil {
		return f, err
	}
	return f, nil
}

// Box will mount each file in the Box and wrap it, already existing files are ignored.
func Box(g *genny.Generator, box packd.Walker) error {
	return box.Walk(func(path string, bf packd.File) error {
		f := genny.NewFile(path, bf)
		f, err := g.Transform(f)
		if err != nil {
			return err
		}
		filePath := strings.TrimSuffix(f.Name(), ".plush")
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			// path doesn't exist. move on.
			g.File(f)
			return nil
		}
		return err
	})
}
