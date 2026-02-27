package xast

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendAutoCLIRPCCommand(t *testing.T) {
	t.Run("append to inline rpc options", func(t *testing.T) {
		content := `package module

import autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Params"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "UpdateParams"},
			},
		},
	}
}
`

		got, err := AppendAutoCLIRPCCommand(content, "Query", `&autocliv1.RpcCommandOptions{RpcMethod: "ListBook"}`)
		require.NoError(t, err)
		require.Contains(t, got, `RpcMethod: "Params"`)
		require.Contains(t, got, `RpcMethod: "ListBook"`)
	})

	t.Run("append when rpc options use variables", func(t *testing.T) {
		content := `package module

import autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	queryRPCOptions := []*autocliv1.RpcCommandOptions{
		{RpcMethod: "Params"},
	}

	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			RpcCommandOptions: queryRPCOptions,
		},
	}
}
`

		got, err := AppendAutoCLIRPCCommand(content, "Query", `&autocliv1.RpcCommandOptions{RpcMethod: "ListBook"}`)
		require.NoError(t, err)

		appendStmt := `queryRPCOptions = append(queryRPCOptions, &autocliv1.RpcCommandOptions{RpcMethod: "ListBook"})`
		require.Contains(t, got, appendStmt)
		require.True(
			t,
			strings.Index(got, appendStmt) < strings.Index(got, "return &autocliv1.ModuleOptions"),
			"append statement must be inserted before return",
		)
	})
}
