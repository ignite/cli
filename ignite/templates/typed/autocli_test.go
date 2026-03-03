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
	option := `&autocliv1.RpcCommandOptions{
		RpcMethod: "ListBook",
		Use: "list-book",
		Short: "List all books",
	}`

	content, err := AppendAutoCLIQueryOptions(autoCLITestContent, option, option)
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "ListBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "Params"`))
}

func TestAppendAutoCLITxOptionsSkipsDuplicates(t *testing.T) {
	option := `&autocliv1.RpcCommandOptions{
		RpcMethod: "CreateBook",
		Use: "create-book [title]",
		Short: "Create a new book",
	}`

	content, err := AppendAutoCLITxOptions(autoCLITestContent, option, option)
	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "CreateBook"`))
	require.Equal(t, 1, strings.Count(content, `RpcMethod: "UpdateParams"`))
}
