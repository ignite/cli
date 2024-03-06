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
	Path     string
	ctx      context.Context
	tracer   *placeholder.Tracer
	TempPath string
}

// NewRunner is a xgenny Runner with a logger.
func NewRunner(ctx context.Context, appPath string) *Runner {
	var (
		runner  = genny.WetRunner(ctx)
		tmpPath = filepath.Join(os.TempDir(), randstr.Runes(5))
	)
	r := &Runner{
		ctx:      ctx,
		Runner:   runner,
		Path:     appPath,
		TempPath: tmpPath,
		tracer:   placeholder.New(),
	}
	runner.FileFn = func(f genny.File) (genny.File, error) {
		return wetFileFn(r, f)
	}
	return r
}

func (r *Runner) Tracer() *placeholder.Tracer {
	return r.tracer
}

// ApplyModifications copy all modifications from the temporary folder to the target path.
func (r *Runner) ApplyModifications() (SourceModification, error) {
	sm := NewSourceModification()
	if _, err := os.Stat(r.TempPath); os.IsNotExist(err) {
		return sm, nil
	}

	err := xos.CopyFolder(r.TempPath, r.Path)
	if err != nil {
		return sm, nil
	}

	// fetch the source modification
	for _, file := range r.Results().Files {
		fileName := file.Name()
		_, err := os.Stat(fileName)
		switch {
		case os.IsNotExist(err):
			sm.AppendCreatedFiles(fileName) // if the file doesn't exist in the source, it means it has been created by the runner
		case err != nil:
			return sm, err
		default:
			sm.AppendModifiedFiles(fileName) // the file has been modified by the runner
		}
	}
	return sm, os.RemoveAll(r.TempPath)
}

// RunAndApply run the generators and apply the modifications to the target path.
func (r *Runner) RunAndApply(gens ...*genny.Generator) (SourceModification, error) {
	if err := r.Run(gens...); err != nil {
		return SourceModification{}, err
	}
	return r.ApplyModifications()
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
	return nil
}

func wetFileFn(runner *Runner, f genny.File) (genny.File, error) {
	if d, ok := f.(genny.Dir); ok {
		if err := os.MkdirAll(d.Name(), d.Perm); err != nil {
			return f, err
		}
		return d, nil
	}

	var err error
	if !filepath.IsAbs(runner.Path) {
		runner.Path, err = filepath.Abs(runner.Path)
		if err != nil {
			return f, err
		}
	}

	name := f.Name()
	if !filepath.IsAbs(name) {
		name = filepath.Join(runner.Path, name)
	}
	relPath, err := filepath.Rel(runner.Path, name)
	if err != nil {
		return f, err
	}

	dstPath := filepath.Join(runner.TempPath, relPath)
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
			// Path doesn't exist. move on.
			g.File(f)
			return nil
		}
		return err
	})
}
