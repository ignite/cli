package xgenny

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/randstr"
	"github.com/ignite/cli/v29/ignite/pkg/xos"
)

type Runner struct {
	*genny.Runner
	ctx     context.Context
	tracer  *placeholder.Tracer
	results []genny.File
	tmpPath string
}

// NewRunner is a xgenny Runner with a logger.
func NewRunner(ctx context.Context, root string) *Runner {
	var (
		runner  = genny.WetRunner(ctx)
		tmpPath = filepath.Join(os.TempDir(), randstr.Runes(5))
	)
	runner.Root = root
	r := &Runner{
		ctx:     ctx,
		Runner:  runner,
		tmpPath: tmpPath,
		tracer:  placeholder.New(),
		results: make([]genny.File, 0),
	}
	runner.FileFn = wetFileFn(r)
	return r
}

func (r *Runner) Tracer() *placeholder.Tracer {
	return r.tracer
}

type (
	OverwriteCallback func(_, _, duplicated []string) error

	// ApplyOption holds the ApplyModifications options.
	applyOptions struct {
		preRun  OverwriteCallback
		postRun OverwriteCallback
	}

	// ApplyOption configures the ApplyModifications options.
	ApplyOption func(r *applyOptions)
)

// ApplyPreRun sets pre-runner for the ApplyModifications function.
func ApplyPreRun(preRun OverwriteCallback) ApplyOption {
	return func(o *applyOptions) {
		o.preRun = preRun
	}
}

// ApplyPostRun sets pos-runner for the ApplyModifications function.
func ApplyPostRun(postRun OverwriteCallback) ApplyOption {
	return func(o *applyOptions) {
		o.postRun = postRun
	}
}

// ApplyModifications copy all modifications from the temporary folder to the target path.
func (r *Runner) ApplyModifications(options ...ApplyOption) (SourceModification, error) {
	opts := applyOptions{}
	for _, apply := range options {
		apply(&opts)
	}

	// fetch the source modification
	sm := NewSourceModification()
	for _, file := range r.results {
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
	r.results = make([]genny.File, 0)

	if _, err := os.Stat(r.tmpPath); os.IsNotExist(err) {
		return sm, nil
	}

	duplicatedFiles, err := xos.ValidateFolderCopy(r.tmpPath, r.Root, sm.ModifiedFiles()...)
	if err != nil {
		return sm, err
	}

	if opts.preRun != nil {
		if err := opts.preRun(sm.CreatedFiles(), sm.ModifiedFiles(), duplicatedFiles); err != nil {
			return sm, err
		}
	}

	// Create the target path and copy the content from the temporary folder.
	if err := os.MkdirAll(r.Root, os.ModePerm); err != nil {
		return sm, err
	}

	if err := xos.CopyFolder(r.tmpPath, r.Root); err != nil {
		return sm, err
	}

	if err := os.RemoveAll(r.tmpPath); err != nil {
		return sm, err
	}

	if opts.postRun != nil {
		if err := opts.postRun(sm.CreatedFiles(), sm.ModifiedFiles(), duplicatedFiles); err != nil {
			return sm, err
		}
	}
	return sm, nil
}

// RunAndApply run the generators and apply the modifications to the target path.
func (r *Runner) RunAndApply(gens *genny.Generator, options ...ApplyOption) (SourceModification, error) {
	if err := r.Run(gens); err != nil {
		return SourceModification{}, err
	}
	return r.ApplyModifications(options...)
}

// Run all generators into a temp folder for we can apply the modifications later.
func (r *Runner) Run(gens ...*genny.Generator) error {
	// execute the modification with a wet runner
	for _, gen := range gens {
		if err := r.Runner.With(gen); err != nil {
			return err
		}
		if err := r.Runner.Run(); err != nil {
			return err
		}
	}
	r.results = append(r.results, r.Results().Files...)
	return r.tracer.Err()
}

func wetFileFn(runner *Runner) func(genny.File) (genny.File, error) {
	return func(f genny.File) (genny.File, error) {
		if d, ok := f.(genny.Dir); ok {
			if err := os.MkdirAll(d.Name(), d.Perm); err != nil {
				return f, err
			}
			return d, nil
		}

		var err error
		if !filepath.IsAbs(runner.Root) {
			runner.Root, err = filepath.Abs(runner.Root)
			if err != nil {
				return f, err
			}
		}

		name := f.Name()
		if !filepath.IsAbs(name) {
			name = filepath.Join(runner.Root, name)
		}
		relPath, err := filepath.Rel(runner.Root, name)
		if err != nil {
			return f, err
		}

		dstPath := filepath.Join(runner.tmpPath, relPath)
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
}
