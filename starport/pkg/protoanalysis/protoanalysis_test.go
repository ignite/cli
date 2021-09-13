package protoanalysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

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
				{Name: "PoolRecord", Path: "testdata/liquidity/genesis.proto"},
				{Name: "GenesisState", Path: "testdata/liquidity/genesis.proto"},
				{Name: "PoolType", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "Params", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "Pool", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "PoolMetadata", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "PoolMetadataResponse", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "PoolBatch", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "PoolBatchResponse", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "DepositMsgState", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "WithdrawMsgState", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "SwapMsgState", Path: "testdata/liquidity/liquidity.proto"},
				{Name: "QueryLiquidityPoolRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryLiquidityPoolResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryLiquidityPoolBatchRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryLiquidityPoolBatchResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryLiquidityPoolsRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryLiquidityPoolsResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryParamsRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryParamsResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchSwapMsgsRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchSwapMsgRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchSwapMsgsResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchSwapMsgResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchDepositMsgsRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchDepositMsgRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchDepositMsgsResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchDepositMsgResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchWithdrawMsgsRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchWithdrawMsgRequest", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchWithdrawMsgsResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "QueryPoolBatchWithdrawMsgResponse", Path: "testdata/liquidity/query.proto"},
				{Name: "MsgCreatePool", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgCreatePoolRequest", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgCreatePoolResponse", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgDepositWithinBatch", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgDepositWithinBatchRequest", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgDepositWithinBatchResponse", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgWithdrawWithinBatch", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgWithdrawWithinBatchRequest", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgWithdrawWithinBatchResponse", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgSwapWithinBatch", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgSwapWithinBatchRequest", Path: "testdata/liquidity/tx.proto"},
				{Name: "MsgSwapWithinBatchResponse", Path: "testdata/liquidity/tx.proto"},
				{Name: "BaseReq", Path: "testdata/liquidity/tx.proto"},
				{Name: "Fee", Path: "testdata/liquidity/tx.proto"},
				{Name: "PubKey", Path: "testdata/liquidity/tx.proto"},
				{Name: "Signature", Path: "testdata/liquidity/tx.proto"},
				{Name: "StdTx", Path: "testdata/liquidity/tx.proto"},
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
