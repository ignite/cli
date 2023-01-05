package debugger

import (
	"context"
	"fmt"
	"net"

	"github.com/go-delve/delve/pkg/logflags"
	"github.com/go-delve/delve/pkg/terminal"
	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/debugger"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/go-delve/delve/service/rpccommon"
	"golang.org/x/sync/errgroup"
)

const (
	// DefaultAddress defines the default debug server address.
	DefaultAddress = "127.0.0.1:30500"

	// DefaultWorkingDir defines the default directory to use as
	// working dir when running the app binary that will be debugged.
	DefaultWorkingDir = "."
)

// Option configures debugging.
type Option func(*debuggerOptions)

type debuggerOptions struct {
	disconnectChan                 chan struct{}
	address, workingDir            string
	listener                       net.Listener
	binaryArgs                     []string
	clientRunHook, serverStartHook func()
}

// Address sets the address for the debug server.
func Address(address string) Option {
	return func(o *debuggerOptions) {
		o.address = address
	}
}

// DisconnectChannel sets the channel used by the server to signal when the client disconnects.
func DisconnectChannel(c chan struct{}) Option {
	return func(o *debuggerOptions) {
		o.disconnectChan = c
	}
}

// Listener sets a custom listener to serve requests.
func Listener(l net.Listener) Option {
	return func(o *debuggerOptions) {
		o.listener = l
	}
}

// WorkingDir sets the working directory of the new process.
func WorkingDir(path string) Option {
	return func(o *debuggerOptions) {
		o.workingDir = path
	}
}

// BinaryArgs sets command line argument for the new process.
func BinaryArgs(args ...string) Option {
	return func(o *debuggerOptions) {
		o.binaryArgs = args
	}
}

// ClientRunHook sets a function to be executed right before debug client is run.
func ClientRunHook(fn func()) Option {
	return func(o *debuggerOptions) {
		o.clientRunHook = fn
	}
}

// ServerStartHook sets a function to be executed right before debug server starts.
func ServerStartHook(fn func()) Option {
	return func(o *debuggerOptions) {
		o.serverStartHook = fn
	}
}

// Start starts a debug server.
func Start(ctx context.Context, binaryPath string, options ...Option) (err error) {
	o := applyDebuggerOptions(options...)

	listener := o.listener
	if listener == nil {
		var c net.ListenConfig

		listener, err = c.Listen(ctx, "tcp", o.address)
		if err != nil {
			return err
		}

		defer listener.Close()
	}

	if err = disableDelveLogging(); err != nil {
		return err
	}

	server := rpccommon.NewServer(&service.Config{
		Listener:           listener,
		AcceptMulti:        false,
		APIVersion:         2,
		CheckLocalConnUser: true,
		DisconnectChan:     o.disconnectChan,
		ProcessArgs:        append([]string{binaryPath}, o.binaryArgs...),
		Debugger: debugger.Config{
			WorkingDir: o.workingDir,
			Backend:    "default",
		},
	})

	if o.serverStartHook != nil {
		o.serverStartHook()
	}

	if err = server.Run(); err != nil {
		return fmt.Errorf("failed to run debug server: %w", err)
	}

	defer server.Stop()

	// Wait until the context is done or the connected client disconnects
	select {
	case <-ctx.Done():
	case <-o.disconnectChan:
	}

	return nil
}

// Run runs a debug client.
func Run(ctx context.Context, binaryPath string, options ...Option) error {
	listener, conn := service.ListenerPipe()
	defer listener.Close()

	o := applyDebuggerOptions(options...)

	options = append(options, Listener(listener))
	g, ctx := errgroup.WithContext(ctx)

	// Start the debugger server
	g.Go(func() error {
		return Start(ctx, binaryPath, options...)
	})

	// Start the debug client
	g.Go(func() error {
		client := rpc2.NewClientFromConn(conn)
		term := terminal.New(client, nil)

		if o.clientRunHook != nil {
			o.clientRunHook()
		}

		_, err := term.Run()
		return err
	})

	return g.Wait()
}

func applyDebuggerOptions(options ...Option) debuggerOptions {
	o := debuggerOptions{
		address:        DefaultAddress,
		workingDir:     DefaultWorkingDir,
		disconnectChan: make(chan struct{}),
	}
	for _, apply := range options {
		apply(&o)
	}
	return o
}

func disableDelveLogging() error {
	if err := logflags.Setup(false, "", ""); err != nil {
		return err
	}
	return nil
}
