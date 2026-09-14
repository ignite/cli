package gno

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	gnodev "github.com/gnolang/gno/contribs/gnodev/pkg/dev"
	"github.com/gnolang/gno/contribs/gnodev/pkg/packages"
	"github.com/gnolang/gno/gno.land/pkg/integration"
)

const reloadTimeout = 60 * time.Second

// ServeOptions configures Serve.
type ServeOptions struct {
	// ChainID is the chain id (default: dev).
	ChainID string
	// RPCListener is the node RPC listen address (default: tcp://127.0.0.1:26657).
	RPCListener string
	// MaxGasPerBlock caps gas per block (default: 10_000_000_000).
	MaxGasPerBlock int64
}

// withDefaults fills unset ServeOptions fields.
func (o ServeOptions) withDefaults() ServeOptions {
	if o.ChainID == "" {
		o.ChainID = "dev"
	}
	if o.RPCListener == "" {
		o.RPCListener = "tcp://127.0.0.1:26657"
	}
	if o.MaxGasPerBlock == 0 {
		o.MaxGasPerBlock = 10_000_000_000
	}
	return o
}

// Serve runs an in-memory gno.land dev chain (single validator) with the gno
// packages found under dir preloaded at genesis, watching .gno files and
// reloading the chain on change. It blocks until ctx is cancelled.
//
// Serve is the integration entry point of the dev chain: it boots the real
// node and is exercised end to end rather than through unit tests.
func Serve(ctx context.Context, dir string, opts ServeOptions, out io.Writer) error {
	opts = opts.withDefaults()

	gnoRoot, err := EnsureStdlibs()
	if err != nil {
		return fmt.Errorf("preparing gno stdlibs: %w", err)
	}

	loader := packages.New(packages.Config{
		Workspace: packages.FindWorkspace(dir),
		Examples:  false,
		GnoRoot:   gnoRoot,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	nodeCfg := gnodev.DefaultNodeConfig(gnoRoot, "gno.land")
	nodeCfg.ChainID = opts.ChainID
	nodeCfg.MaxGasPerBlock = opts.MaxGasPerBlock
	nodeCfg.TMConfig.RPC.ListenAddress = opts.RPCListener
	// errLog captures node errors (e.g. a preloaded package that fails to
	// compile) so they can be surfaced after startup instead of being buried
	// in node logs.
	errLog := &errorLog{}
	nodeCfg.Logger = slog.New(newErrorLogHandler(errLog, slog.NewTextHandler(out, nil)))
	nodeCfg.Reload = loader.Reload
	nodeCfg.NoReplay = true

	// preload the package in dir (if any) so it is deployed at genesis.
	paths := preloadPaths(dir)

	node, err := gnodev.NewDevNode(ctx, nodeCfg, paths...)
	if err != nil {
		return fmt.Errorf("starting dev node: %w", err)
	}
	defer node.Close()

	printServeBanner(out, opts, paths)

	// report packages that failed to load at genesis (the dev chain skips
	// failing genesis txs, so a broken package would otherwise be silent).
	if failures := errLog.errors(); len(failures) > 0 {
		fmt.Fprintf(out, "\n❌ %d package transaction(s) failed at genesis; the chain is running WITHOUT them:\n", len(failures))
		for _, e := range failures {
			fmt.Fprintf(out, "   - %s\n", e)
		}
		fmt.Fprintln(out, "   Check imports (only stdlibs and your workspace packages are available) and syntax, then fix and save to reload.")
	}

	return watchLoop(ctx, node.Reload, dir, out)
}

// printServeBanner prints the dev chain summary.
func printServeBanner(out io.Writer, opts ServeOptions, paths []string) {
	fmt.Fprintf(out, `⭐ gno.land dev chain started.

   Chain ID:   %s
   RPC:        %s
   Deploy key: %s (%s, pre-funded)
   Mnemonic:   %s
`, opts.ChainID, opts.RPCListener, integration.DefaultAccount_Address,
		integration.DefaultAccount_Name, integration.DefaultAccount_Seed)

	if len(paths) == 0 {
		fmt.Fprintln(out, "\n   No gno package found in the current directory.")
		fmt.Fprintln(out, "   Scaffold one with `ignite scaffold realm <name>` then restart,")
		fmt.Fprintln(out, "   or deploy at runtime with `ignite chain deploy`.")
	} else {
		fmt.Fprintln(out, "\n   Loaded packages:")
		for _, p := range paths {
			fmt.Fprintf(out, "     - %s\n", p)
		}
	}
	fmt.Fprintln(out, "\n   Watching for .gno file changes (reload chain on save)...")
}

// preloadPaths returns the module path of the package in dir (if any).
func preloadPaths(dir string) []string {
	gnomodPath := filepath.Join(dir, "gnomod.toml")
	if _, err := os.Stat(gnomodPath); err != nil {
		return nil
	}
	mod, err := parseGnoMod(gnomodPath)
	if err != nil || mod == "" {
		return nil
	}
	return []string{mod}
}

// watchLoop watches .gno and gnomod.toml files under dir and calls reload on
// change. It returns when ctx is done.
func watchLoop(ctx context.Context, reload func(context.Context) error, dir string, out io.Writer) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := watchDirRecursive(fw, dir); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-fw.Events:
			if !ok {
				return nil
			}
			reloadOnChange(ctx, reload, event, out)
		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(out, "watch error: %v\n", err)
		}
	}
}

// reloadOnChange reloads the chain when a gno source event fires.
func reloadOnChange(ctx context.Context, reload func(context.Context) error, event fsnotify.Event, out io.Writer) {
	if !isGnoFile(event.Name) {
		return
	}
	if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) == 0 {
		return
	}
	fmt.Fprintf(out, "🔁 %s changed, reloading chain...\n", event.Name)
	reloadCtx, cancel := context.WithTimeout(ctx, reloadTimeout)
	defer cancel()
	if err := reload(reloadCtx); err != nil {
		fmt.Fprintf(out, "❌ reload failed: %v\n", err)
	} else {
		fmt.Fprintln(out, "✅ chain reloaded")
	}
}

// watchDirRecursive adds dir and its subdirectories (skipping dot dirs) to
// the watcher.
func watchDirRecursive(fw *fsnotify.Watcher, dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil //nolint:nilerr // skip unreadable entries
		}
		if !d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") && path != dir {
			return filepath.SkipDir
		}
		return fw.Add(path)
	})
}

// isGnoFile reports whether path is a gno source or module file.
func isGnoFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasSuffix(base, ".gno") || base == "gnomod.toml"
}
