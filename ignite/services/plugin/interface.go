package plugin

import (
	"context"

	v1 "github.com/ignite/cli/ignite/services/plugin/grpc/v1"
)

// Flag type aliases.
const (
	FlagTypeString      = v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED
	FlagTypeInt         = v1.Flag_TYPE_FLAG_INT
	FlagTypeUint        = v1.Flag_TYPE_FLAG_UINT
	FlagTypeInt64       = v1.Flag_TYPE_FLAG_INT64
	FlagTypeUint64      = v1.Flag_TYPE_FLAG_UINT64
	FlagTypeBool        = v1.Flag_TYPE_FLAG_BOOL
	FlagTypeStringSlice = v1.Flag_TYPE_FLAG_STRING_SLICE
)

// Type aliases for the current plugin version.
type (
	Command         = v1.Command
	Dependency      = v1.Dependency
	ExecutedCommand = v1.ExecutedCommand
	ExecutedHook    = v1.ExecutedHook
	Flag            = v1.Flag
	FlagType        = v1.Flag_Type
	Hook            = v1.Hook
	Manifest        = v1.Manifest
)

// An ignite plugin must implements the Plugin interface.
//
//go:generate mockery --srcpkg . --name Interface --structname PluginInterface --filename interface.go --with-expecter
type Interface interface {
	// Manifest declares the plugin's Command(s) and Hook(s).
	Manifest(context.Context) (*Manifest, error)

	// Execute will be invoked by ignite when a plugin Command is executed.
	// It is global for all commands declared in Manifest, if you have declared
	// multiple commands, use cmd.Path to distinguish them.
	// The analyzer argument can be used by plugins to get chain app analysis info.
	Execute(context.Context, *ExecutedCommand, Analyzer) error

	// ExecuteHookPre is invoked by ignite when a command specified by the Hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	// The analyzer argument can be used by plugins to get chain app analysis info.
	ExecuteHookPre(context.Context, *ExecutedHook, Analyzer) error

	// ExecuteHookPost is invoked by ignite when a command specified by the hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	// The analyzer argument can be used by plugins to get chain app analysis info.
	ExecuteHookPost(context.Context, *ExecutedHook, Analyzer) error

	// ExecuteHookCleanUp is invoked by ignite when a command specified by the
	// hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
	// execution status of the command and hooks.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	// The analyzer argument can be used by plugins to get chain app analysis info.
	ExecuteHookCleanUp(context.Context, *ExecutedHook, Analyzer) error
}

// Analyzer defines the interface for plugins to get chain app code analysis info.
//
//go:generate mockery --srcpkg . --name Analyzer --structname PluginAnalyzer --filename interface.go --with-expecter
type Analyzer interface {
	// Dependencies returns the app dependencies.
	Dependencies(context.Context) ([]*Dependency, error)
}
