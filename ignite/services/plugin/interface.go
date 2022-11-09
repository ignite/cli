package plugin

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/pflag"
)

func init() {
	gob.Register(Manifest{})
	gob.Register(ExecutedCommand{})
	gob.Register(ExecutedHook{})
}

// An ignite plugin must implements the Plugin interface.
//
//go:generate mockery --srcpkg . --name Interface --structname PluginInterface --filename interface.go --with-expecter
type Interface interface {
	// Manifest declares the plugin's Command(s) and Hook(s).
	Manifest() (Manifest, error)

	// Execute will be invoked by ignite when a plugin commands is executed.
	// cmd is the executed command (one of the those returned by Commands method)
	// args is the command line arguments passed behing the command.
	Execute(cmd ExecutedCommand) error

	// ExecuteHookPre is invoked by Ignite when a command specified by the hook
	// path is invoked is global for all hooks registered to a plugin context on
	// the hook being invoked is given by the `hook` parameter.
	ExecuteHookPre(hook ExecutedHook) error
	// ExecuteHookPost is invoked by Ignite when a command specified by the hook
	// path is invoked is global for all hooks registered to a plugin context on
	// the hook being invoked is given by the `hook` parameter.
	ExecuteHookPost(hook ExecutedHook) error
	// ExecuteHookCleanUp is invoked right before the command is done executing
	// will be called regardless of execution status of the command and hooks.
	ExecuteHookCleanUp(hook ExecutedHook) error
}

type Manifest struct {
	Name     string
	Commands []Command
	Hooks    []Hook
}

// Command represents a plugin command.
type Command struct {
	// Same as cobra.Command.Use
	Use string
	// Same as cobra.Command.Short
	Short string
	// Same as cobra.Command.Long
	Long string
	// Flags holds the list of command flags
	Flags []Flag
	// PlaceCommandUnder indicates where the command should be placed.
	// For instance `ignite scaffold` will place the command at the
	// `scaffold` command.
	// An empty value is interpreted as `ignite` (==root).
	PlaceCommandUnder string
	// List of sub commands
	Commands []Command
}

// Hook represents a user defined action within a plugin
type Hook struct {
	// Name identifies the hook for the client to invoke the correct hook
	// must be unique
	Name string
	// PlaceHookOn indicates the command to register the hooks for
	PlaceHookOn string
}

type ExecutedCommand struct {
	Use string
	// Path contains the command path, e.g. `ignite scaffold foo`
	Path  string
	Args  []string
	flags *pflag.FlagSet
	With  map[string]string
}

type ExecutedHook struct {
	ExecutedCommand
	Hook
}

func (c *ExecutedCommand) Flags() *pflag.FlagSet {
	if c.flags == nil {
		c.flags = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	}
	return c.flags
}

func (c *ExecutedCommand) SetFlags(fs *pflag.FlagSet) {
	c.flags = fs
}

type Flag struct {
	Name      string // name as it appears on command line
	Shorthand string // one-letter abbreviated flag
	Usage     string // help message
	DefValue  string // default value (as text); for usage message
	Type      FlagType
	Value     string
	// TODO add Persistent field ?
}

type FlagType string

const (
	// NOTE(tb): we declare only the main used cobra flag types for simplicity
	// If a plugin receives an unhandled type, it will output an error.
	FlagTypeString      FlagType = "string"
	FlagTypeInt         FlagType = "int"
	FlagTypeUint        FlagType = "uint"
	FlagTypeInt64       FlagType = "int64"
	FlagTypeUint64      FlagType = "uint64"
	FlagTypeBool        FlagType = "bool"
	FlagTypeStringSlice FlagType = "stringSlice"
)

func (f Flag) FeedFlagSet(fs *pflag.FlagSet) error {
	switch f.Type {
	case FlagTypeBool:
		defVal, _ := strconv.ParseBool(f.DefValue)
		fs.BoolP(f.Name, f.Shorthand, defVal, f.Usage)
		fs.Set(f.Name, f.Value)
	case FlagTypeInt:
		defVal, _ := strconv.Atoi(f.DefValue)
		fs.IntP(f.Name, f.Shorthand, defVal, f.Usage)
		fs.Set(f.Name, f.Value)
	case FlagTypeUint:
		defVal, _ := strconv.ParseUint(f.DefValue, 10, 64)
		fs.UintP(f.Name, f.Shorthand, uint(defVal), f.Usage)
		fs.Set(f.Name, f.Value)
	case FlagTypeInt64:
		defVal, _ := strconv.ParseInt(f.DefValue, 10, 64)
		fs.Int64P(f.Name, f.Shorthand, defVal, f.Usage)
		fs.Set(f.Name, f.Value)
	case FlagTypeUint64:
		defVal, _ := strconv.ParseUint(f.DefValue, 10, 64)
		fs.Uint64P(f.Name, f.Shorthand, defVal, f.Usage)
		fs.Set(f.Name, f.Value)
	case FlagTypeString:
		fs.StringP(f.Name, f.Shorthand, f.DefValue, f.Usage)
		fs.Set(f.Name, f.Value)
	case FlagTypeStringSlice:
		s := strings.Trim(f.DefValue, "[]")
		defValue := strings.Fields(s)
		fs.StringSliceP(f.Name, f.Shorthand, defValue, f.Usage)
		fs.Set(f.Name, strings.Trim(f.Value, "[]"))
	default:
		return fmt.Errorf("flagset unmarshal: unhandled flag type %q in flag %#v", f.Type, f)
	}
	return nil
}

// gobCommandFlags is used to gob encode/decode Command.
// Command can't be encoded because :
// - flags is unexported (because we want to expose it via the Flags() method,
// like a regular cobra.Command)
// - flags type is *pflag.FlagSet which is also full of unexported fields.
type gobCommandContextFlags struct {
	CommandContext gobCommandContext
	Flags          []Flag
}

type gobCommandContext ExecutedCommand

// GobEncode implements gob.Encoder.
// It actually encodes a gobCommandContext struct built from c.
func (c ExecutedCommand) GobEncode() ([]byte, error) {
	var ff []Flag
	if c.flags != nil {
		c.flags.VisitAll(func(pf *pflag.Flag) {
			ff = append(ff, Flag{
				Name:      pf.Name,
				Shorthand: pf.Shorthand,
				Usage:     pf.Usage,
				DefValue:  pf.DefValue,
				Value:     pf.Value.String(),
				Type:      FlagType(pf.Value.Type()),
			})
		})
	}
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(gobCommandContextFlags{
		CommandContext: gobCommandContext(c),
		Flags:          ff,
	})
	return b.Bytes(), err
}

// GobDecode implements gob.Decoder.
// It actually decodes a gobCommandContext struct and fills c with it.
func (c *ExecutedCommand) GobDecode(bz []byte) error {
	var gb gobCommandContextFlags
	err := gob.NewDecoder(bytes.NewReader(bz)).Decode(&gb)
	if err != nil {
		return err
	}
	*c = ExecutedCommand(gb.CommandContext)
	for _, f := range gb.Flags {
		err := f.FeedFlagSet(c.Flags())
		if err != nil {
			return err
		}
	}
	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func HandshakeConfig() plugin.HandshakeConfig {
	return handshakeConfig
}

// Here is an implementation that talks over RPC
type InterfaceRPC struct{ client *rpc.Client }

// Manifest implements Interface.Manifest
func (g *InterfaceRPC) Manifest() (Manifest, error) {
	var resp Manifest
	return resp, g.client.Call("Plugin.Manifest", new(interface{}), &resp)
}

// Execute implements Interface.Commands
func (g *InterfaceRPC) Execute(c ExecutedCommand) error {
	var resp interface{}
	return g.client.Call("Plugin.Execute", map[string]interface{}{
		"executedCommand": c,
	}, &resp)
}

func (g *InterfaceRPC) ExecuteHookPre(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookPre", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

func (g *InterfaceRPC) ExecuteHookPost(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookPost", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

func (g *InterfaceRPC) ExecuteHookCleanUp(hook ExecutedHook) error {
	var resp interface{}
	return g.client.Call("Plugin.ExecuteHookCleanUp", map[string]interface{}{
		"executedHook": hook,
	}, &resp)
}

// Here is the RPC server that InterfaceRPC talks to, conforming to
// the requirements of net/rpc
type InterfaceRPCServer struct {
	// This is the real implementation
	Impl Interface
}

func (s *InterfaceRPCServer) Manifest(args interface{}, resp *Manifest) error {
	var err error
	*resp, err = s.Impl.Manifest()
	return err
}

func (s *InterfaceRPCServer) Execute(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.Execute(args["executedCommand"].(ExecutedCommand))
}

func (s *InterfaceRPCServer) ExecuteHookPre(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.ExecuteHookPre(args["executedHook"].(ExecutedHook))
}

func (s *InterfaceRPCServer) ExecuteHookPost(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.ExecuteHookPost(args["executedHook"].(ExecutedHook))
}

func (s *InterfaceRPCServer) ExecuteHookCleanUp(args map[string]interface{}, resp *interface{}) error {
	return s.Impl.ExecuteHookCleanUp(args["executedHook"].(ExecutedHook))
}

// This is the implementation of plugin.Interface so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a InterfaceRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return InterfaceRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type InterfacePlugin struct {
	// Impl Injection
	Impl Interface
}

func (p *InterfacePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &InterfaceRPCServer{Impl: p.Impl}, nil
}

func (InterfacePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &InterfaceRPC{client: c}, nil
}
