package protoanalysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNestedMessages(t *testing.T) {
	packages, err := Parse(context.Background(), nil, "testutil/nested_messages")
	require.NoError(t, err)

	pkg := packages[0]
	require.Equal(t, "A", pkg.Messages[0].Name)
	require.Equal(t, "A_B", pkg.Messages[1].Name)
	require.Equal(t, "A_B_C", pkg.Messages[2].Name)
}

func TestLiquidity(t *testing.T) {
	packages, err := Parse(context.Background(), nil, "testutil/liquidity")
	require.NoError(t, err)

	expected := Packages{
		{
			Name: "tendermint.liquidity",
			Path: "testutil/liquidity",
			Files: Files{
				{
					Path:         "testutil/liquidity/genesis.proto",
					Dependencies: []string{"liquidity.proto", "gogoproto/gogo.proto"},
				},
				{
					Path:         "testutil/liquidity/liquidity.proto",
					Dependencies: []string{"tx.proto", "gogoproto/gogo.proto", "cosmos_proto/coin.proto", "protoc-gen-openapiv2/options/annotations.proto"},
				},
				{
					Path:         "testutil/liquidity/msg.proto",
					Dependencies: []string{"google/api/annotations.proto", "protoc-gen-openapiv2/options/annotations.proto", "tx.proto"},
				},
				{
					Path:         "testutil/liquidity/query.proto",
					Dependencies: []string{"gogoproto/gogo.proto", "liquidity.proto", "google/api/annotations.proto", "cosmos_proto/pagination.proto", "protoc-gen-openapiv2/options/annotations.proto"},
				},
				{
					Path:         "testutil/liquidity/tx.proto",
					Dependencies: []string{"gogoproto/gogo.proto", "cosmos_proto/coin.proto", "protoc-gen-openapiv2/options/annotations.proto"},
				},
			},
			GoImportName: "github.com/tendermint/liquidity/x/liquidity/types",
			Messages: []Message{
				{Name: "PoolRecord", Path: "testutil/liquidity/genesis.proto", HighestFieldNumber: 6},
				{Name: "GenesisState", Path: "testutil/liquidity/genesis.proto", HighestFieldNumber: 2},
				{Name: "PoolType", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 5},
				{Name: "Params", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 9},
				{Name: "Pool", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 5},
				{Name: "PoolMetadata", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 3},
				{Name: "PoolMetadataResponse", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 2},
				{Name: "PoolBatch", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 7},
				{Name: "PoolBatchResponse", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 6},
				{Name: "DepositMsgState", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 6},
				{Name: "WithdrawMsgState", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 6},
				{Name: "SwapMsgState", Path: "testutil/liquidity/liquidity.proto", HighestFieldNumber: 10},
				{Name: "QueryLiquidityPoolRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolBatchRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolBatchResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolsRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolsResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryParamsRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 0},
				{Name: "QueryParamsResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryPoolBatchSwapMsgsRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchSwapMsgRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchSwapMsgsResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchSwapMsgResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryPoolBatchDepositMsgsRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchDepositMsgRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchDepositMsgsResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchDepositMsgResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryPoolBatchWithdrawMsgsRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchWithdrawMsgRequest", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchWithdrawMsgsResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchWithdrawMsgResponse", Path: "testutil/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "MsgCreatePool", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 4},
				{Name: "MsgCreatePoolRequest", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 2},
				{Name: "MsgCreatePoolResponse", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "MsgDepositWithinBatch", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgDepositWithinBatchRequest", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgDepositWithinBatchResponse", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "MsgWithdrawWithinBatch", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgWithdrawWithinBatchRequest", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgWithdrawWithinBatchResponse", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "MsgSwapWithinBatch", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 7},
				{Name: "MsgSwapWithinBatchRequest", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgSwapWithinBatchResponse", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "BaseReq", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 11},
				{Name: "Fee", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 2},
				{Name: "PubKey", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 2},
				{Name: "Signature", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 4},
				{Name: "StdTx", Path: "testutil/liquidity/tx.proto", HighestFieldNumber: 4},
			},
			Services: []Service{
				{
					Name: "MsgApi",
					RPCFuncs: []RPCFunc{
						{
							Name:        "CreatePoolApi",
							RequestType: "MsgCreatePoolRequest",
							ReturnsType: "MsgCreatePoolResponse",
							HTTPRules: []HTTPRule{
								{
									Params:  []string{"test"},
									HasBody: true,
								},
							},
						},
						{
							Name:        "DepositWithinBatchApi",
							RequestType: "MsgDepositWithinBatchRequest",
							ReturnsType: "MsgDepositWithinBatchResponse",
							HTTPRules: []HTTPRule{
								{
									Params:  []string{"pool_id"},
									HasBody: true,
								},
							},
						},
						{
							Name:        "WithdrawWithinBatchApi",
							RequestType: "MsgWithdrawWithinBatchRequest",
							ReturnsType: "MsgWithdrawWithinBatchResponse",
							HTTPRules: []HTTPRule{
								{
									Params:  []string{"pool_id"},
									HasBody: true,
								},
							},
						},
						{
							Name:        "SwapApi",
							RequestType: "MsgSwapWithinBatchRequest",
							ReturnsType: "MsgSwapWithinBatchResponse",
							HTTPRules: []HTTPRule{
								{
									Params:   []string{"pool_id"},
									HasQuery: true,
									HasBody:  true,
								},
							},
						},
					},
				},
				{
					Name: "Query",
					RPCFuncs: []RPCFunc{
						{
							Name:        "LiquidityPools",
							RequestType: "QueryLiquidityPoolsRequest",
							ReturnsType: "QueryLiquidityPoolsResponse",
							HTTPRules: []HTTPRule{
								{
									HasQuery: true,
								},
							},
						},
						{
							Name:        "LiquidityPool",
							RequestType: "QueryLiquidityPoolRequest",
							ReturnsType: "QueryLiquidityPoolResponse",
							HTTPRules: []HTTPRule{
								{
									Params: []string{"pool_id"},
								},
							},
						},
						{
							Name:        "LiquidityPoolBatch",
							RequestType: "QueryLiquidityPoolBatchRequest",
							ReturnsType: "QueryLiquidityPoolBatchResponse",
							HTTPRules: []HTTPRule{
								{
									Params: []string{"pool_id"},
								},
							},
						},
						{
							Name:        "PoolBatchSwapMsgs",
							RequestType: "QueryPoolBatchSwapMsgsRequest",
							ReturnsType: "QueryPoolBatchSwapMsgsResponse",
							HTTPRules: []HTTPRule{
								{
									Params:   []string{"pool_id"},
									HasQuery: true,
								},
							},
						},
						{
							Name:        "PoolBatchSwapMsg",
							RequestType: "QueryPoolBatchSwapMsgRequest",
							ReturnsType: "QueryPoolBatchSwapMsgResponse",
							HTTPRules: []HTTPRule{
								{
									Params: []string{"pool_id", "msg_index"},
								},
							},
						},
						{
							Name:        "PoolBatchDepositMsgs",
							RequestType: "QueryPoolBatchDepositMsgsRequest",
							ReturnsType: "QueryPoolBatchDepositMsgsResponse",
							HTTPRules: []HTTPRule{
								{
									Params:   []string{"pool_id"},
									HasQuery: true,
								},
							},
						},
						{
							Name:        "PoolBatchDepositMsg",
							RequestType: "QueryPoolBatchDepositMsgRequest",
							ReturnsType: "QueryPoolBatchDepositMsgResponse",
							HTTPRules: []HTTPRule{
								{
									Params: []string{"pool_id", "msg_index"},
								},
							},
						},
						{
							Name:        "PoolBatchWithdrawMsgs",
							RequestType: "QueryPoolBatchWithdrawMsgsRequest",
							ReturnsType: "QueryPoolBatchWithdrawMsgsResponse",
							HTTPRules: []HTTPRule{
								{
									Params:   []string{"pool_id"},
									HasQuery: true,
								},
							},
						},
						{
							Name:        "PoolBatchWithdrawMsg",
							RequestType: "QueryPoolBatchWithdrawMsgRequest",
							ReturnsType: "QueryPoolBatchWithdrawMsgResponse",
							HTTPRules: []HTTPRule{
								{
									Params: []string{"pool_id", "msg_index"},
								},
							},
						},
						{
							Name:        "Params",
							RequestType: "QueryParamsRequest",
							ReturnsType: "QueryParamsResponse",
							HTTPRules: []HTTPRule{
								{},
							},
						},
					},
				},
				{
					Name: "Msg",
					RPCFuncs: []RPCFunc{
						{
							Name:        "CreatePool",
							RequestType: "MsgCreatePool",
							ReturnsType: "MsgCreatePoolResponse",
						},
						{
							Name:        "DepositWithinBatch",
							RequestType: "MsgDepositWithinBatch",
							ReturnsType: "MsgDepositWithinBatchResponse",
						},
						{
							Name:        "WithdrawWithinBatch",
							RequestType: "MsgWithdrawWithinBatch",
							ReturnsType: "MsgWithdrawWithinBatchResponse",
						},
						{
							Name:        "Swap",
							RequestType: "MsgSwapWithinBatch",
							ReturnsType: "MsgSwapWithinBatchResponse",
						},
					},
				},
			},
		},
	}

	require.Equal(t, expected, packages)
}
