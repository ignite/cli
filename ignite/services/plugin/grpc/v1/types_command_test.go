package v1_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/ignite/cli/v29/ignite/services/plugin/grpc/v1"
)

func TestCommandToCobraCommand(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		pcmd    = v1.Command{
			Use:     "new",
			Aliases: []string{"n"},
			Short:   "short",
			Long:    "long",
			Hidden:  true,
			Flags: []*v1.Flag{
				{
					Name:         "bool",
					Shorthand:    "b",
					DefaultValue: "true",
					Value:        "true",
					Usage:        "a bool",
					Type:         v1.Flag_TYPE_FLAG_BOOL,
				},
				{
					Name:         "string",
					DefaultValue: "hello",
					Value:        "hello",
					Usage:        "a string",
					Type:         v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED,
					Persistent:   true,
				},
			},
			Commands: []*v1.Command{
				{
					Use:     "sub",
					Aliases: []string{"s"},
					Short:   "sub short",
					Long:    "sub long",
				},
			},
		}
	)

	cmd, err := pcmd.ToCobraCommand()

	require.NoError(err)
	require.NotNil(cmd)
	assert.Empty(cmd.Commands()) // subcommands aren't converted
	assert.Equal(pcmd.Use, cmd.Use)
	assert.Equal(pcmd.Short, cmd.Short)
	assert.Equal(pcmd.Long, cmd.Long)
	assert.Equal(pcmd.Aliases, cmd.Aliases)
	assert.Equal(pcmd.Hidden, cmd.Hidden)
	for _, f := range pcmd.Flags {
		if f.Persistent {
			assert.NotNil(cmd.PersistentFlags().Lookup(f.Name), "missing pflag %s", f.Name)
		} else {
			assert.NotNil(cmd.Flags().Lookup(f.Name), "missing flag %s", f.Name)
		}
	}
}

func TestCommandPath(t *testing.T) {
	cases := []struct {
		name, wantPath string
		cmd            *v1.Command
	}{
		{
			name: "relative path",
			cmd: &v1.Command{
				PlaceCommandUnder: "chain",
			},
			wantPath: "ignite chain",
		},
		{
			name: "full path",
			cmd: &v1.Command{
				PlaceCommandUnder: "ignite chain",
			},
			wantPath: "ignite chain",
		},
		{
			name: "path with spaces",
			cmd: &v1.Command{
				PlaceCommandUnder: " ignite scaffold  ",
			},
			wantPath: "ignite scaffold",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.cmd.Path()
			require.Equal(t, tc.wantPath, path)
		})
	}
}

func TestExecutedCommandImportFlags(t *testing.T) {
	// Arrange
	execCmd := &v1.ExecutedCommand{}
	wantFlags := []*v1.Flag{
		{
			Name:         "foo",
			Shorthand:    "f",
			Usage:        "foo usage",
			DefaultValue: "bar",
			Value:        "baz",
			Type:         v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED,
		}, {
			Name:         "test",
			Shorthand:    "t",
			Usage:        "test usage",
			DefaultValue: "1",
			Value:        "42",
			Type:         v1.Flag_TYPE_FLAG_INT,
			Persistent:   true,
		},
	}

	cmd := cobra.Command{}
	cmd.Flags().StringP("foo", "f", "bar", "foo usage")
	cmd.PersistentFlags().IntP("test", "t", 1, "test usage")
	err := cmd.ParseFlags([]string{"--foo", "baz", "--test", "42"})
	require.NoError(t, err)

	// Act
	execCmd.ImportFlags(&cmd)

	// Assert
	require.Equal(t, wantFlags, execCmd.Flags)
}

func TestExecutedCommandNewFlags(t *testing.T) {
	// Arrange
	execCmd := &v1.ExecutedCommand{
		Flags: []*v1.Flag{
			{
				Name:         "bool",
				Shorthand:    "b",
				Usage:        "bool usage",
				DefaultValue: "false",
				Value:        "true",
				Type:         v1.Flag_TYPE_FLAG_BOOL,
			},
			{
				Name:         "int",
				Shorthand:    "i",
				Usage:        "int usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_INT,
			},
			{
				Name:         "uint",
				Shorthand:    "u",
				Usage:        "uint usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_UINT,
			},
			{
				Name:         "int64",
				Shorthand:    "j",
				Usage:        "int64 usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_INT64,
			},
			{
				Name:         "uint64",
				Shorthand:    "k",
				Usage:        "uint64 usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_UINT64,
			},
			{
				Name:         "string",
				Shorthand:    "s",
				Usage:        "string usage",
				DefaultValue: "",
				Value:        "hello",
				Type:         v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED,
			},
			{
				Name:         "string-slice",
				Shorthand:    "l",
				Usage:        "string slice usage",
				DefaultValue: "[]",
				Value:        "[]",
				Type:         v1.Flag_TYPE_FLAG_STRING_SLICE,
			},
			{
				Name:       "persistent",
				Persistent: true,
			},
		},
	}

	wantFlags := make(map[string]pflag.Flag)
	for _, f := range execCmd.Flags {
		wantFlags[f.Name] = pflag.Flag{
			Name:      f.Name,
			Shorthand: f.Shorthand,
			Usage:     f.Usage,
			DefValue:  f.DefaultValue,
		}
	}

	var (
		flagCount int

		// Persistent flag should not be included
		wantFlagCount = len(execCmd.Flags) - 1
	)

	// Act
	flags, err := execCmd.NewFlags()

	// Assert
	require.NoError(t, err)

	flags.VisitAll(func(f *pflag.Flag) {
		flag, ok := wantFlags[f.Name]

		require.True(t, ok, fmt.Sprintf("missing flag: %s", f.Name))
		require.Equal(t, flag.Name, f.Name)
		require.Equal(t, flag.Shorthand, f.Shorthand)
		require.Equal(t, flag.Usage, f.Usage)
		require.Equal(t, flag.DefValue, f.DefValue)

		flagCount++
	})

	require.Equal(t, wantFlagCount, flagCount)
}

func TestExecutedCommandNewPersistentFlags(t *testing.T) {
	// Arrange
	execCmd := &v1.ExecutedCommand{
		Flags: []*v1.Flag{
			{
				Name:         "bool",
				Shorthand:    "b",
				Usage:        "bool usage",
				DefaultValue: "false",
				Value:        "true",
				Type:         v1.Flag_TYPE_FLAG_BOOL,
				Persistent:   true,
			},
			{
				Name:         "int",
				Shorthand:    "i",
				Usage:        "int usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_INT,
				Persistent:   true,
			},
			{
				Name:         "uint",
				Shorthand:    "u",
				Usage:        "uint usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_UINT,
				Persistent:   true,
			},
			{
				Name:         "int64",
				Shorthand:    "j",
				Usage:        "int64 usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_INT64,
				Persistent:   true,
			},
			{
				Name:         "uint64",
				Shorthand:    "k",
				Usage:        "uint64 usage",
				DefaultValue: "0",
				Value:        "42",
				Type:         v1.Flag_TYPE_FLAG_UINT64,
				Persistent:   true,
			},
			{
				Name:         "string",
				Shorthand:    "s",
				Usage:        "string usage",
				DefaultValue: "",
				Value:        "hello",
				Type:         v1.Flag_TYPE_FLAG_STRING_UNSPECIFIED,
				Persistent:   true,
			},
			{
				Name:         "string-slice",
				Shorthand:    "l",
				Usage:        "string slice usage",
				DefaultValue: "[]",
				Value:        "[]",
				Type:         v1.Flag_TYPE_FLAG_STRING_SLICE,
				Persistent:   true,
			},
			{
				Name: "non-persistent",
			},
		},
	}

	wantFlags := make(map[string]pflag.Flag)
	for _, f := range execCmd.Flags {
		wantFlags[f.Name] = pflag.Flag{
			Name:      f.Name,
			Shorthand: f.Shorthand,
			Usage:     f.Usage,
			DefValue:  f.DefaultValue,
		}
	}

	var (
		flagCount int

		// Non persistent flag should not be included
		wantFlagCount = len(execCmd.Flags) - 1
	)

	// Act
	flags, err := execCmd.NewPersistentFlags()

	// Assert
	require.NoError(t, err)

	flags.VisitAll(func(f *pflag.Flag) {
		flag, ok := wantFlags[f.Name]

		require.True(t, ok, fmt.Sprintf("missing flag: %s", f.Name))
		require.Equal(t, flag.Name, f.Name)
		require.Equal(t, flag.Shorthand, f.Shorthand)
		require.Equal(t, flag.Usage, f.Usage)
		require.Equal(t, flag.DefValue, f.DefValue)

		flagCount++
	})

	require.Equal(t, wantFlagCount, flagCount)
}
