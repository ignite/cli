package protoanalysis

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLiquidity(t *testing.T) {
	packages, err := Parse(context.Background(), PatternRecursive("testdata/liquidity"))
	require.NoError(t, err)

	var expected []Package

	err = json.Unmarshal([]byte(`
[
	{
		"Name": "tendermint.liquidity",
		"Path": "testdata/liquidity",
		"GoImportName": "github.com/tendermint/liquidity/x/liquidity/types",
		"Messages": [
			{
				"Name": "PoolType",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "Params",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "Pool",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "PoolMetadata",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "PoolMetadataResponse",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "PoolBatch",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "PoolBatchResponse",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "DepositMsgState",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "WithdrawMsgState",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "SwapMsgState",
				"Path": "testdata/liquidity/liquidity.proto"
			},
			{
				"Name": "MsgCreatePool",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgCreatePoolRequest",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgCreatePoolResponse",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgDepositWithinBatch",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgDepositWithinBatchRequest",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgDepositWithinBatchResponse",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgWithdrawWithinBatch",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgWithdrawWithinBatchRequest",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgWithdrawWithinBatchResponse",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgSwapWithinBatch",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgSwapWithinBatchRequest",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "MsgSwapWithinBatchResponse",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "BaseReq",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "Fee",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "PubKey",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "Signature",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "StdTx",
				"Path": "testdata/liquidity/tx.proto"
			},
			{
				"Name": "QueryLiquidityPoolRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryLiquidityPoolResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryLiquidityPoolBatchRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryLiquidityPoolBatchResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryLiquidityPoolsRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryLiquidityPoolsResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryParamsRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryParamsResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchSwapMsgsRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchSwapMsgRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchSwapMsgsResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchSwapMsgResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchDepositMsgsRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchDepositMsgRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchDepositMsgsResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchDepositMsgResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchWithdrawMsgsRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchWithdrawMsgRequest",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchWithdrawMsgsResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "QueryPoolBatchWithdrawMsgResponse",
				"Path": "testdata/liquidity/query.proto"
			},
			{
				"Name": "PoolRecord",
				"Path": "testdata/liquidity/genesis.proto"
			},
			{
				"Name": "GenesisState",
				"Path": "testdata/liquidity/genesis.proto"
			}
		],
		"Services": [
			{
				"Name": "MsgApi",
				"RPCFuncs": [
					{
						"Name": "CreatePoolApi",
						"RequestType": "MsgCreatePoolRequest",
						"ReturnsType": "MsgCreatePoolResponse",
						"HTTPRules": [
							{
								"Params": [
									"test"
								],
								"HasQuery": false,
								"HasBody": true
							}
						]
					},
					{
						"Name": "DepositWithinBatchApi",
						"RequestType": "MsgDepositWithinBatchRequest",
						"ReturnsType": "MsgDepositWithinBatchResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": false,
								"HasBody": true
							}
						]
					},
					{
						"Name": "WithdrawWithinBatchApi",
						"RequestType": "MsgWithdrawWithinBatchRequest",
						"ReturnsType": "MsgWithdrawWithinBatchResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": false,
								"HasBody": true
							}
						]
					},
					{
						"Name": "SwapApi",
						"RequestType": "MsgSwapWithinBatchRequest",
						"ReturnsType": "MsgSwapWithinBatchResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": true,
								"HasBody": true
							}
						]
					}
				]
			},
			{
				"Name": "Msg",
				"RPCFuncs": [
					{
						"Name": "CreatePool",
						"RequestType": "MsgCreatePool",
						"ReturnsType": "MsgCreatePoolResponse",
						"HTTPRules": null
					},
					{
						"Name": "DepositWithinBatch",
						"RequestType": "MsgDepositWithinBatch",
						"ReturnsType": "MsgDepositWithinBatchResponse",
						"HTTPRules": null
					},
					{
						"Name": "WithdrawWithinBatch",
						"RequestType": "MsgWithdrawWithinBatch",
						"ReturnsType": "MsgWithdrawWithinBatchResponse",
						"HTTPRules": null
					},
					{
						"Name": "Swap",
						"RequestType": "MsgSwapWithinBatch",
						"ReturnsType": "MsgSwapWithinBatchResponse",
						"HTTPRules": null
					}
				]
			},
			{
				"Name": "Query",
				"RPCFuncs": [
					{
						"Name": "LiquidityPools",
						"RequestType": "QueryLiquidityPoolsRequest",
						"ReturnsType": "QueryLiquidityPoolsResponse",
						"HTTPRules": [
							{
								"Params": null,
								"HasQuery": true,
								"HasBody": false
							}
						]
					},
					{
						"Name": "LiquidityPool",
						"RequestType": "QueryLiquidityPoolRequest",
						"ReturnsType": "QueryLiquidityPoolResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": false,
								"HasBody": false
							}
						]
					},
					{
						"Name": "LiquidityPoolBatch",
						"RequestType": "QueryLiquidityPoolBatchRequest",
						"ReturnsType": "QueryLiquidityPoolBatchResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": false,
								"HasBody": false
							}
						]
					},
					{
						"Name": "PoolBatchSwapMsgs",
						"RequestType": "QueryPoolBatchSwapMsgsRequest",
						"ReturnsType": "QueryPoolBatchSwapMsgsResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": true,
								"HasBody": false
							}
						]
					},
					{
						"Name": "PoolBatchSwapMsg",
						"RequestType": "QueryPoolBatchSwapMsgRequest",
						"ReturnsType": "QueryPoolBatchSwapMsgResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id",
									"msg_index"
								],
								"HasQuery": false,
								"HasBody": false
							}
						]
					},
					{
						"Name": "PoolBatchDepositMsgs",
						"RequestType": "QueryPoolBatchDepositMsgsRequest",
						"ReturnsType": "QueryPoolBatchDepositMsgsResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": true,
								"HasBody": false
							}
						]
					},
					{
						"Name": "PoolBatchDepositMsg",
						"RequestType": "QueryPoolBatchDepositMsgRequest",
						"ReturnsType": "QueryPoolBatchDepositMsgResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id",
									"msg_index"
								],
								"HasQuery": false,
								"HasBody": false
							}
						]
					},
					{
						"Name": "PoolBatchWithdrawMsgs",
						"RequestType": "QueryPoolBatchWithdrawMsgsRequest",
						"ReturnsType": "QueryPoolBatchWithdrawMsgsResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id"
								],
								"HasQuery": true,
								"HasBody": false
							}
						]
					},
					{
						"Name": "PoolBatchWithdrawMsg",
						"RequestType": "QueryPoolBatchWithdrawMsgRequest",
						"ReturnsType": "QueryPoolBatchWithdrawMsgResponse",
						"HTTPRules": [
							{
								"Params": [
									"pool_id",
									"msg_index"
								],
								"HasQuery": false,
								"HasBody": false
							}
						]
					},
					{
						"Name": "Params",
						"RequestType": "QueryParamsRequest",
						"ReturnsType": "QueryParamsResponse",
						"HTTPRules": [
							{
								"Params": null,
								"HasQuery": false,
								"HasBody": false
							}
						]
					}
				]
			}
		]
	}
]
`), &expected)
	require.NoError(t, err)

	require.Equal(t, expected, packages)
}
