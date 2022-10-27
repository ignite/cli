package localfs

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	wt "github.com/radovskyb/watcher"
)

type watcher struct {
	wt            *wt.Watcher
	workdir       string
	ignoreHidden  bool
	ignoreFolders bool
	ignoreExts    []string
	onChange      func()
	interval      time.Duration
	ctx           context.Context
	done          *sync.WaitGroup
}

// WatcherOption used to configure watcher.
type WatcherOption func(*watcher)

// WatcherWorkdir to set as a root to paths needs to be watched.
func WatcherWorkdir(path string) WatcherOption {
	return func(w *watcher) {
		w.workdir = path
	}
}

// WatcherOnChange sets a hook that executed on every change on filesystem.
func WatcherOnChange(hook func()) WatcherOption {
	return func(w *watcher) {
		w.onChange = hook
	}
}

// WatcherPollingInterval overwrites default polling interval to check filesystem changes.
func WatcherPollingInterval(d time.Duration) WatcherOption {
	return func(w *watcher) {
		w.interval = d
	}
}

// WatcherIgnoreHidden ignores hidden(dot) files.
func WatcherIgnoreHidden() WatcherOption {
	return func(w *watcher) {
		w.ignoreHidden = true
	}
}

func WatcherIgnoreFolders() WatcherOption {
	return func(w *watcher) {
		w.ignoreFolders = true
	}
}

// WatcherIgnoreExt ignores files with matching file extensions.
func WatcherIgnoreExt(exts ...string) WatcherOption {
	return func(w *watcher) {
		w.ignoreExts = exts
	}
}

// Watch starts watching changes on the paths. options are used to configure the
// behaviour of watch operation.
func Watch(ctx context.Context, paths []string, options ...WatcherOption) error {
	w := &watcher{
		wt:       wt.New(),
		onChange: func() {},
		interval: time.Millisecond * 300,
		done:     &sync.WaitGroup{},
		ctx:      ctx,
	}
	w.wt.SetMaxEvents(1)

	for _, o := range options {
		o(w)
	}

	w.wt.AddFilterHook(func(info os.FileInfo, fullPath string) error {
		if info.IsDir() && w.ignoreFolders {
			return wt.ErrSkip
		}
		if w.isFileIgnored(fullPath) {
			return wt.ErrSkip
		}

		return nil
	})

	// ignore hidden paths.
	w.wt.IgnoreHiddenFiles(w.ignoreHidden)

	// add paths to watch
	if err := w.addPaths(paths...); err != nil {
		return err
	}

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
		case <-w.wt.Event:
			w.onChange()
		case <-w.wt.Closed:
			return
		case <-w.ctx.Done():
			w.wt.Close()
		}
	}
}

func (w *watcher) addPaths(paths ...string) error {
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(w.workdir, path)
		}

		// Ignoring paths that don't exist
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			continue
		}

		if err := w.wt.AddRecursive(path); err != nil {
			return err
		}
	}

	return nil
}

func (w *watcher) isFileIgnored(path string) bool {
	for _, ext := range w.ignoreExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
