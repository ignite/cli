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
	"github.com/spf13/cobra"
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

	// Execute will be invoked by ignite when a plugin Command is executed.
	// It is global for all commands declared in Manifest, if you have declared
	// multiple commands, use cmd.Path to distinguish them.
	Execute(cmd ExecutedCommand) error

	// ExecuteHookPre is invoked by ignite when a command specified by the Hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPre(hook ExecutedHook) error
	// ExecuteHookPost is invoked by ignite when a command specified by the hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPost(hook ExecutedHook) error
	// ExecuteHookCleanUp is invoked by ignite when a command specified by the
	// hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
	// execution status of the command and hooks.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookCleanUp(hook ExecutedHook) error
}

// Manifest represents the plugin behavior.
type Manifest struct {
	Name string
	// Commands contains the commands that will be added to the list of ignite
	// commands. Each commands are independent, for nested commands use the
	// inner Commands field.
	Commands []Command
	// Hooks contains the hooks that will be attached to the existing ignite
	// commands.
	Hooks []Hook
	// SharedHost enables sharing a single plugin server across all running instances
	// of a plugin. Useful if a plugin adds or extends long running commands
	//
	// Example: if a plugin defines a hook on `ignite chain serve`, a plugin server is instanciated
	// when the command is run. Now if you want to interact with that instance from commands
	// defined in that plugin, you need to enable `SharedHost`, or else the commands will just
	// instantiate separate plugin servers.
	//
	// When enabled, all plugins of the same `Path` loaded from the same configuration will
	// attach it's rpc client to a an existing rpc server.
	//
	// If a plugin instance has no other running plugin servers, it will create one and it will be the host.
	SharedHost bool `yaml:"shared_host"`
}

// ImportCobraCommand allows to hydrate m with a standard root cobra commands.
func (m *Manifest) ImportCobraCommand(c *cobra.Command, placeCommandUnder string) {
	m.Commands = append(m.Commands, convertCobraCommand(c, placeCommandUnder))
}

func convertCobraCommand(c *cobra.Command, placeCommandUnder string) Command {
	cmd := Command{
		Use:               c.Use,
		Aliases:           c.Aliases,
		Short:             c.Short,
		Long:              c.Long,
		Hidden:            c.Hidden,
		PlaceCommandUnder: placeCommandUnder,
		Flags:             convertPFlags(c),
	}
	for _, c := range c.Commands() {
		cmd.Commands = append(cmd.Commands, convertCobraCommand(c, ""))
	}
	return cmd
}

// Command represents a plugin command.
type Command struct {
	// Same as cobra.Command.Use
	Use string
	// Same as cobra.Command.Aliases
	Aliases []string
	// Same as cobra.Command.Short
	Short string
	// Same as cobra.Command.Long
	Long string
	// Same as cobra.Command.Hidden
	Hidden bool
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

// PlaceCommandUnderFull returns a normalized p.PlaceCommandUnder, by adding
// the `ignite ` prefix if not present.
func (c Command) PlaceCommandUnderFull() string {
	return commandFull(c.PlaceCommandUnder)
}

func commandFull(cmdPath string) string {
	const rootCmdName = "ignite"
	if !strings.HasPrefix(cmdPath, rootCmdName) {
		cmdPath = rootCmdName + " " + cmdPath
	}
	return strings.TrimSpace(cmdPath)
}

// ToCobraCommand turns Command into a cobra.Command so it can be added to a
// parent command.
func (c Command) ToCobraCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     c.Use,
		Aliases: c.Aliases,
		Short:   c.Short,
		Long:    c.Long,
		Hidden:  c.Hidden,
	}
	for _, f := range c.Flags {
		err := f.feedFlagSet(cmd)
		if err != nil {
			return nil, err
		}
	}
	return cmd, nil
}

// Hook represents a user defined action within a plugin.
type Hook struct {
	// Name identifies the hook for the client to invoke the correct hook
	// must be unique
	Name string
	// PlaceHookOn indicates the command to register the hooks for
	PlaceHookOn string
}

// PlaceHookOnFull returns a normalized p.PlaceCommandUnder, by adding the
// `ignite ` prefix if not present.
func (h Hook) PlaceHookOnFull() string {
	return commandFull(h.PlaceHookOn)
}

// ExecutedCommand represents a plugin command under execution.
type ExecutedCommand struct {
	// Use is copied from Command.Use
	Use string
	// Path contains the command path, e.g. `ignite scaffold foo`
	Path string
	// Args are the command arguments
	Args []string
	// Full list of args taken from os.Args
	OSArgs []string
	// With contains the plugin config parameters
	With map[string]string

	flags  *pflag.FlagSet
	pflags *pflag.FlagSet
}

// ExecutedHook represents a plugin hook under execution.
type ExecutedHook struct {
	// ExecutedCommand gives access to the command attached by the hook.
	ExecutedCommand ExecutedCommand
	// Hook is a copy of the original Hook defined in the Manifest.
	Hook
}

// Flags gives access to the commands' flags, like cobra.Command.Flags.
func (c *ExecutedCommand) Flags() *pflag.FlagSet {
	if c.flags == nil {
		c.flags = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	}
	return c.flags
}

// PersistentFlags gives access to the commands' persistent flags, like
// cobra.Command.PersistentFlags.
func (c *ExecutedCommand) PersistentFlags() *pflag.FlagSet {
	if c.pflags == nil {
		c.pflags = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	}
	return c.pflags
}

// SetFlags set the flags.
// As a plugin developer, you probably don't need to use it.
func (c *ExecutedCommand) SetFlags(cmd *cobra.Command) {
	c.flags = cmd.Flags()
	c.pflags = cmd.PersistentFlags()
}

// Flag is a serializable representation of pflag.Flag.
type Flag struct {
	Name      string // name as it appears on command line
	Shorthand string // one-letter abbreviated flag
	Usage     string // help message
	DefValue  string // default value (as text); for usage message
	Type      FlagType
	Value     string
	// Persistent indicates wether or not the flag is propagated on children
	// commands
	Persistent bool
}

// FlagType represents the pflag.Flag.Value.Type().
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

// feedFlagSet fills flagger with f.
func (f Flag) feedFlagSet(fgr flagger) error {
	fs := fgr.Flags()
	if f.Persistent {
		fs = fgr.PersistentFlags()
	}
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

// gobCommandContext is the same as ExecutedCommand but without GobDecode
// attached, which avoids infinite loops.
type gobCommandContext ExecutedCommand

// GobEncode implements gob.Encoder.
// It actually encodes a gobCommandContext struct built from c.
func (c ExecutedCommand) GobEncode() ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(gobCommandContextFlags{
		CommandContext: gobCommandContext(c),
		Flags:          convertPFlags(&c),
	})
	return b.Bytes(), err
}

// flagger matches both cobra.Command and Command.
type flagger interface {
	Flags() *pflag.FlagSet
	PersistentFlags() *pflag.FlagSet
}

func convertPFlags(fgr flagger) []Flag {
	var ff []Flag
	if fgr.Flags() != nil {
		fgr.Flags().VisitAll(func(pf *pflag.Flag) {
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
	if fgr.PersistentFlags() != nil {
		fgr.PersistentFlags().VisitAll(func(pf *pflag.Flag) {
			ff = append(ff, Flag{
				Name:       pf.Name,
				Shorthand:  pf.Shorthand,
				Usage:      pf.Usage,
				DefValue:   pf.DefValue,
				Value:      pf.Value.String(),
				Type:       FlagType(pf.Value.Type()),
				Persistent: true,
			})
		})
	}
	return ff
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
		err := f.feedFlagSet(c)
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

// InterfaceRPC is an implementation that talks over RPC.
type InterfaceRPC struct{ client *rpc.Client }

// Manifest implements Interface.Manifest.
func (g *InterfaceRPC) Manifest() (Manifest, error) {
	var resp Manifest
	return resp, g.client.Call("Plugin.Manifest", new(interface{}), &resp)
}

// Execute implements Interface.Commands.
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

// InterfaceRPCServer is the RPC server that InterfaceRPC talks to, conforming to
// the requirements of net/rpc.
type InterfaceRPCServer struct {
	// This is the real implementation
	Impl Interface
}

func (s *InterfaceRPCServer) Manifest(_ interface{}, resp *Manifest) error {
	var err error
	*resp, err = s.Impl.Manifest()
	return err
}

func (s *InterfaceRPCServer) Execute(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.Execute(args["executedCommand"].(ExecutedCommand))
}

func (s *InterfaceRPCServer) ExecuteHookPre(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.ExecuteHookPre(args["executedHook"].(ExecutedHook))
}

func (s *InterfaceRPCServer) ExecuteHookPost(args map[string]interface{}, _ *interface{}) error {
	return s.Impl.ExecuteHookPost(args["executedHook"].(ExecutedHook))
}

func (s *InterfaceRPCServer) ExecuteHookCleanUp(args map[string]interface{}, _ *interface{}) error {
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

func (InterfacePlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &InterfaceRPC{client: c}, nil
}
