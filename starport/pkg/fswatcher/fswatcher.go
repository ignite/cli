// Package fswatcher provides functionalities to watch changes on the
// filesystem.
package fswatcher

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"time"

	wt "github.com/radovskyb/watcher"
)

type watcher struct {
	wt           *wt.Watcher
	workdir      string
	ignoreHidden bool
	ignoreExts   []string
	onChange     func()
	interval     time.Duration
	ctx          context.Context
	done         *sync.WaitGroup
}

// Option used to configure watcher.
type Option func(*watcher)

// Workdir to set as a root to paths needs to be watched.
func Workdir(path string) Option {
	return func(w *watcher) {
		w.workdir = path
	}
}

// OnChange sets a hook that executed on every change on filesystem.
func OnChange(hook func()) Option {
	return func(w *watcher) {
		w.onChange = hook
	}
}

// PollingInterval overwrites default polling interval to check filesystem changes.
func PollingInterval(d time.Duration) Option {
	return func(w *watcher) {
		w.interval = d
	}
}

// IgnoreHidden ignores hidden(dot) files.
func IgnoreHidden() Option {
	return func(w *watcher) {
		w.ignoreHidden = true
	}
}

// IgnoreExt ignores files with matching file extensions.
func IgnoreExt(exts ...string) Option {
	return func(w *watcher) {
		w.ignoreExts = exts
	}
}

// Watch starts watching changes on the paths. options are used to configure the
// behaviour of watch operation.
func Watch(ctx context.Context, paths []string, options ...Option) error {
	wt := wt.New()
	wt.SetMaxEvents(1)

	w := &watcher{
		wt:       wt,
		onChange: func() {},
		interval: time.Millisecond * 300,
		done:     &sync.WaitGroup{},
		ctx:      ctx,
	}
	for _, o := range options {
		o(w)
	}

	// ignore hidden paths.
	w.wt.IgnoreHiddenFiles(w.ignoreHidden)

	// add paths to watch
	w.addPaths(paths...)

	// start watching.
	w.done.Add(1)
	go w.listen()
	if err := w.wt.Start(w.interval); err != nil {
		return err
	}
	w.done.Wait()
	return nil
}

func (w *watcher) listen() {
	defer w.done.Done()
	for {
		select {
		case e := <-w.wt.Event:
			if !w.isFileIgnored(e.Path) {
				w.onChange()
			}
		case <-w.wt.Closed:
			return
		case <-w.ctx.Done():
			w.wt.Close()
		}
	}
}

func (w *watcher) addPaths(paths ...string) {
	for _, path := range paths {
		w.wt.AddRecursive(filepath.Join(w.workdir, path))
	}
}

func (w *watcher) isFileIgnored(path string) bool {
	for _, ext := range w.ignoreExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
