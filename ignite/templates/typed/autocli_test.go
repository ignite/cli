package typed

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const autoCLITestContent = `package foo

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
				},
			},
		},
	}
}
`

func TestAppendAutoCLIQueryOptionsSkipsDuplicates(t *testing.T) {
	option := `{
		RpcMethod: "ListBook",
		Use: "list-book",
		Short: "List all books",
	}`

	content, err := AppendAutoCLIQueryOptions(autoCLITestContent, option, option)
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "ListBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "Params"`))
}

func TestAppendAutoCLIQueryOptionsIsIdempotent(t *testing.T) {
	options := []string{
		`{
			RpcMethod: "ListBook",
			Use: "list-book",
			Short: "List all books",
		}`,
		`{
			RpcMethod: "GetBook",
			Use: "get-book [id]",
			Short: "Gets a book",
		}`,
	}

	content, err := AppendAutoCLIQueryOptions(autoCLITestContent, options...)
	require.NoError(t, err)

	content, err = AppendAutoCLIQueryOptions(content, options...)
	require.NoError(t, err)

	require.Equal(t, 1, strings.Count(content, `RpcMethod: "ListBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "GetBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "Params"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "UpdateParams"`))
	require.Greater(t, strings.Index(content, `RpcMethod: "ListBook"`), strings.Index(content, `RpcMethod: "Params"`))
	require.Greater(t, strings.Index(content, `RpcMethod: "GetBook"`), strings.Index(content, `RpcMethod: "ListBook"`))
}

func TestAppendAutoCLITxOptionsSkipsDuplicates(t *testing.T) {
	option := `{
		RpcMethod: "CreateBook",
		Use: "create-book [title]",
		Short: "Create a new book",
	}`

	content, err := AppendAutoCLITxOptions(autoCLITestContent, option, option)
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "CreateBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "UpdateParams"`))
}

func TestAppendAutoCLITxOptionsSkipsExistingMethods(t *testing.T) {
	options := []string{
		`{
			RpcMethod: "UpdateParams",
			Use: "update-params",
			Short: "should be skipped",
		}`,
		`{
			RpcMethod: "CreateBook",
			Use: "create-book [title]",
			Short: "Create a new book",
		}`,
	}

	content, err := AppendAutoCLITxOptions(autoCLITestContent, options...)
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "UpdateParams"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "CreateBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "Params"`))
}

func TestAppendAutoCLIQueryOptionsErrors(t *testing.T) {
	t.Run("missing AutoCLIOptions", func(t *testing.T) {
		_, err := AppendAutoCLIQueryOptions(`package foo`, `{RpcMethod:"ListBook"}`)
		require.Error(t, err)
		require.Contains(t, err.Error(), `function "AutoCLIOptions" not found`)
	})

	t.Run("missing Query field", func(t *testing.T) {
		content := strings.Replace(autoCLITestContent, "Query:", "MissingQuery:", 1)
		_, err := AppendAutoCLIQueryOptions(content, `{RpcMethod:"ListBook"}`)
		require.Error(t, err)
		require.Contains(t, err.Error(), `field "Query" not found in ModuleOptions`)
	})

	t.Run("missing RpcCommandOptions field", func(t *testing.T) {
		content := strings.Replace(autoCLITestContent, "RpcCommandOptions:", "MissingRpcOptions:", 1)
		_, err := AppendAutoCLIQueryOptions(content, `{RpcMethod:"ListBook"}`)
		require.Error(t, err)
		require.Contains(t, err.Error(), `field "RpcCommandOptions" not found in "Query" service descriptor`)
	})

	t.Run("invalid option expression", func(t *testing.T) {
		_, err := AppendAutoCLIQueryOptions(autoCLITestContent, `invalid(`)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse autocli option expression")
	})
}

func TestAppendAutoCLIQueryOptionsFormatting(t *testing.T) {
	options := []string{
		`{
			RpcMethod: "ListBook",
			Use: "list-book",
			Short: "List all books",
		}`,
		`{
			RpcMethod: "GetBook",
			Use: "get-book [id]",
			Short: "Gets a book",
			Alias: []string{"show-book"},
			PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
		}`,
	}

	content, err := AppendAutoCLIQueryOptions(autoCLITestContent, options...)
	require.NoError(t, err)
	require.NotContains(t, content, "&autocliv1.RpcCommandOptions")
	require.NotContains(t, content, "}, {")
	require.Contains(t, content, "RpcCommandOptions: []*autocliv1.RpcCommandOptions{")

	normalized := strings.NewReplacer(" ", "", "\t", "", "\n", "").Replace(content)
	require.Contains(t, normalized, "RpcMethod:\"ListBook\",")
	require.Contains(t, normalized, "RpcMethod:\"GetBook\",")
}
