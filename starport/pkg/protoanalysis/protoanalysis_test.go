package protoanalysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNestedMessages(t *testing.T) {
	packages, err := Parse(context.Background(), nil, "testdata/nested_messages")
	require.NoError(t, err)

	pkg := packages[0]
	require.Equal(t, "A", pkg.Messages[0].Name)
	require.Equal(t, "A_B", pkg.Messages[1].Name)
	require.Equal(t, "A_B_C", pkg.Messages[2].Name)
}

func TestLiquidity(t *testing.T) {
	packages, err := Parse(context.Background(), nil, "testdata/liquidity")
	require.NoError(t, err)

	expected := Packages{
		{
			Name: "tendermint.liquidity",
			Path: "testdata/liquidity",
			Files: Files{
				{
					Path:         "testdata/liquidity/genesis.proto",
					Dependencies: []string{"liquidity.proto", "gogoproto/gogo.proto"},
				},
				{
					Path:         "testdata/liquidity/liquidity.proto",
					Dependencies: []string{"tx.proto", "gogoproto/gogo.proto", "cosmos_proto/coin.proto", "protoc-gen-openapiv2/options/annotations.proto"},
				},
				{
					Path:         "testdata/liquidity/msg.proto",
					Dependencies: []string{"google/api/annotations.proto", "protoc-gen-openapiv2/options/annotations.proto", "tx.proto"},
				},
				{
					Path:         "testdata/liquidity/query.proto",
					Dependencies: []string{"gogoproto/gogo.proto", "liquidity.proto", "google/api/annotations.proto", "cosmos_proto/pagination.proto", "protoc-gen-openapiv2/options/annotations.proto"},
				},
				{
					Path:         "testdata/liquidity/tx.proto",
					Dependencies: []string{"gogoproto/gogo.proto", "cosmos_proto/coin.proto", "protoc-gen-openapiv2/options/annotations.proto"},
				},
			},
			GoImportName: "github.com/tendermint/liquidity/x/liquidity/types",
			Messages: []Message{
				{Name: "PoolRecord", Path: "testdata/liquidity/genesis.proto", HighestFieldNumber: 6},
				{Name: "GenesisState", Path: "testdata/liquidity/genesis.proto", HighestFieldNumber: 2},
				{Name: "PoolType", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 5},
				{Name: "Params", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 9},
				{Name: "Pool", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 5},
				{Name: "PoolMetadata", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 3},
				{Name: "PoolMetadataResponse", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 2},
				{Name: "PoolBatch", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 7},
				{Name: "PoolBatchResponse", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 6},
				{Name: "DepositMsgState", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 6},
				{Name: "WithdrawMsgState", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 6},
				{Name: "SwapMsgState", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 10},
				{Name: "QueryLiquidityPoolRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolBatchRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolBatchResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryLiquidityPoolsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryParamsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 0},
				{Name: "QueryParamsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryPoolBatchSwapMsgsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchSwapMsgRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchSwapMsgsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchSwapMsgResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryPoolBatchDepositMsgsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchDepositMsgRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchDepositMsgsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchDepositMsgResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "QueryPoolBatchWithdrawMsgsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchWithdrawMsgRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchWithdrawMsgsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2},
				{Name: "QueryPoolBatchWithdrawMsgResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1},
				{Name: "MsgCreatePool", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 4},
				{Name: "MsgCreatePoolRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 2},
				{Name: "MsgCreatePoolResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "MsgDepositWithinBatch", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgDepositWithinBatchRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgDepositWithinBatchResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "MsgWithdrawWithinBatch", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgWithdrawWithinBatchRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgWithdrawWithinBatchResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "MsgSwapWithinBatch", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 7},
				{Name: "MsgSwapWithinBatchRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3},
				{Name: "MsgSwapWithinBatchResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1},
				{Name: "BaseReq", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 11},
				{Name: "Fee", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 2},
				{Name: "PubKey", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 2},
				{Name: "Signature", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 4},
				{Name: "StdTx", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 4},
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
