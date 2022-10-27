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
				{Name: "PoolRecord", Path: "testdata/liquidity/genesis.proto", HighestFieldNumber: 6, Fields: map[string]string{
					"deposit_msg_states":  "DepositMsgState",
					"pool":                "Pool",
					"pool_batch":          "PoolBatch",
					"pool_metadata":       "PoolMetadata",
					"swap_msg_states":     "SwapMsgState",
					"withdraw_msg_states": "WithdrawMsgState",
				}},
				{Name: "GenesisState", Path: "testdata/liquidity/genesis.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"params":       "Params",
					"pool_records": "PoolRecord",
				}},
				{Name: "PoolType", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 5, Fields: map[string]string{
					"description":          "string",
					"id":                   "uint32",
					"max_reserve_coin_num": "uint32",
					"min_reserve_coin_num": "uint32",
					"name":                 "string",
				}},
				{Name: "Params", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 9, Fields: map[string]string{
					"init_pool_coin_mint_amount": "string",
					"max_order_amount_ratio":     "bytes",
					"max_reserve_coin_amount":    "string",
					"min_init_deposit_amount":    "string",
					"pool_creation_fee":          "cosmos.base.v1beta1.Coin",
					"pool_types":                 "PoolType",
					"swap_fee_rate":              "bytes",
					"unit_batch_height":          "uint32",
					"withdraw_fee_rate":          "bytes",
				}},
				{Name: "Pool", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 5, Fields: map[string]string{
					"id":                      "uint64",
					"pool_coin_denom":         "string",
					"reserve_account_address": "string",
					"reserve_coin_denoms":     "string",
					"type_id":                 "uint32",
				}},
				{Name: "PoolMetadata", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 3, Fields: map[string]string{
					"pool_coin_total_supply": "cosmos.base.v1beta1.Coin",
					"pool_id":                "uint64",
					"reserve_coins":          "cosmos.base.v1beta1.Coin",
				}},
				{Name: "PoolMetadataResponse", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pool_coin_total_supply": "cosmos.base.v1beta1.Coin",
					"reserve_coins":          "cosmos.base.v1beta1.Coin",
				}},
				{Name: "PoolBatch", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 7, Fields: map[string]string{
					"begin_height":       "int64",
					"deposit_msg_index":  "uint64",
					"executed":           "bool",
					"index":              "uint64",
					"pool_id":            "uint64",
					"swap_msg_index":     "uint64",
					"withdraw_msg_index": "uint64",
				}},
				{Name: "PoolBatchResponse", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 6, Fields: map[string]string{
					"begin_height":       "int64",
					"deposit_msg_index":  "uint64",
					"executed":           "bool",
					"index":              "uint64",
					"swap_msg_index":     "uint64",
					"withdraw_msg_index": "uint64",
				}},
				{Name: "DepositMsgState", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 6, Fields: map[string]string{
					"executed":      "bool",
					"msg":           "MsgDepositWithinBatch",
					"msg_height":    "int64",
					"msg_index":     "uint64",
					"succeeded":     "bool",
					"to_be_deleted": "bool",
				}},
				{Name: "WithdrawMsgState", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 6, Fields: map[string]string{
					"executed":      "bool",
					"msg":           "MsgWithdrawWithinBatch",
					"msg_height":    "int64",
					"msg_index":     "uint64",
					"succeeded":     "bool",
					"to_be_deleted": "bool",
				}},
				{Name: "SwapMsgState", Path: "testdata/liquidity/liquidity.proto", HighestFieldNumber: 10, Fields: map[string]string{
					"exchanged_offer_coin":    "cosmos.base.v1beta1.Coin",
					"executed":                "bool",
					"msg":                     "MsgSwapWithinBatch",
					"msg_height":              "int64",
					"msg_index":               "uint64",
					"order_expiry_height":     "int64",
					"remaining_offer_coin":    "cosmos.base.v1beta1.Coin",
					"reserved_offer_coin_fee": "cosmos.base.v1beta1.Coin",
					"succeeded":               "bool",
					"to_be_deleted":           "bool",
				}},
				{Name: "QueryLiquidityPoolRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"pool_id": "uint64",
				}},
				{Name: "QueryLiquidityPoolResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"pool": "Pool",
				}},
				{Name: "QueryLiquidityPoolBatchRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"pool_id": "uint64",
				}},
				{Name: "QueryLiquidityPoolBatchResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"batch": "PoolBatch",
				}},
				{Name: "QueryLiquidityPoolsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
				}},
				{Name: "QueryLiquidityPoolsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
					"pools":      "Pool",
				}},
				{Name: "QueryParamsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 0, Fields: map[string]string{}},
				{Name: "QueryParamsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"params": "Params",
				}},
				{Name: "QueryPoolBatchSwapMsgsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
					"pool_id":    "uint64",
				}},
				{Name: "QueryPoolBatchSwapMsgRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"msg_index": "uint64",
					"pool_id":   "uint64",
				}},
				{Name: "QueryPoolBatchSwapMsgsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
					"swaps":      "SwapMsgState",
				}},
				{Name: "QueryPoolBatchSwapMsgResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"swap": "SwapMsgState",
				}},
				{Name: "QueryPoolBatchDepositMsgsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
					"pool_id":    "uint64",
				}},
				{Name: "QueryPoolBatchDepositMsgRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"msg_index": "uint64",
					"pool_id":   "uint64",
				}},
				{Name: "QueryPoolBatchDepositMsgsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"deposits":   "DepositMsgState",
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
				}},
				{Name: "QueryPoolBatchDepositMsgResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"deposit": "DepositMsgState",
				}},
				{Name: "QueryPoolBatchWithdrawMsgsRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
					"pool_id":    "uint64",
				}},
				{Name: "QueryPoolBatchWithdrawMsgRequest", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"msg_index": "uint64",
					"pool_id":   "uint64",
				}},
				{Name: "QueryPoolBatchWithdrawMsgsResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
					"withdraws":  "WithdrawMsgState",
				}},
				{Name: "QueryPoolBatchWithdrawMsgResponse", Path: "testdata/liquidity/query.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"withdraw": "WithdrawMsgState",
				}},
				{Name: "MsgCreatePool", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 4, Fields: map[string]string{
					"deposit_coins":        "cosmos.base.v1beta1.Coin",
					"pool_creator_address": "string",
					"pool_type_id":         "uint32",
				}},
				{Name: "MsgCreatePoolRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"base_req": "BaseReq",
					"msg":      "MsgCreatePool",
				}},
				{Name: "MsgCreatePoolResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "MsgDepositWithinBatch", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3, Fields: map[string]string{
					"deposit_coins":     "cosmos.base.v1beta1.Coin",
					"depositor_address": "string",
					"pool_id":           "uint64",
				}},
				{Name: "MsgDepositWithinBatchRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3, Fields: map[string]string{
					"base_req": "BaseReq",
					"msg":      "MsgDepositWithinBatch",
					"pool_id":  "uint64",
				}},
				{Name: "MsgDepositWithinBatchResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "MsgWithdrawWithinBatch", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3, Fields: map[string]string{
					"pool_coin":          "cosmos.base.v1beta1.Coin",
					"pool_id":            "uint64",
					"withdrawer_address": "string",
				}},
				{Name: "MsgWithdrawWithinBatchRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3, Fields: map[string]string{
					"base_req": "BaseReq",
					"msg":      "MsgWithdrawWithinBatch",
					"pool_id":  "uint64",
				}},
				{Name: "MsgWithdrawWithinBatchResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "MsgSwapWithinBatch", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 7, Fields: map[string]string{
					"demand_coin_denom":      "string",
					"offer_coin":             "cosmos.base.v1beta1.Coin",
					"offer_coin_fee":         "cosmos.base.v1beta1.Coin",
					"order_price":            "bytes",
					"pool_id":                "uint64",
					"swap_requester_address": "string",
					"swap_type_id":           "uint32",
				}},
				{Name: "MsgSwapWithinBatchRequest", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 3, Fields: map[string]string{
					"base_req": "BaseReq",
					"msg":      "MsgSwapWithinBatch",
					"pool_id":  "uint64",
				}},
				{Name: "MsgSwapWithinBatchResponse", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 1, Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "BaseReq", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 11, Fields: map[string]string{
					"account_number": "uint64",
					"chain_id":       "string",
					"fees":           "cosmos.base.v1beta1.Coin",
					"from":           "string",
					"gas":            "uint64",
					"gas_adjustment": "string",
					"gas_prices":     "cosmos.base.v1beta1.DecCoin",
					"memo":           "string",
					"sequence":       "uint64",
					"simulate":       "bool",
					"timeout_height": "uint64",
				}},
				{Name: "Fee", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"amount": "cosmos.base.v1beta1.Coin",
					"gas":    "uint64",
				}},
				{Name: "PubKey", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 2, Fields: map[string]string{
					"type":  "string",
					"value": "string",
				}},
				{Name: "Signature", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 4, Fields: map[string]string{
					"account_number": "uint64",
					"pub_key":        "PubKey",
					"sequence":       "uint64",
					"signature":      "string",
				}},
				{Name: "StdTx", Path: "testdata/liquidity/tx.proto", HighestFieldNumber: 4, Fields: map[string]string{
					"fee":       "Fee",
					"memo":      "string",
					"msg":       "string",
					"signature": "Signature",
				}},
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
