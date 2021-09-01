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
				{Name: "PoolRecord", Path: "testdata/liquidity/genesis.proto", Fields: map[string]string{
					"deposit_msg_states":  "DepositMsgState",
					"pool":                "Pool",
					"pool_batch":          "PoolBatch",
					"pool_metadata":       "PoolMetadata",
					"swap_msg_states":     "SwapMsgState",
					"withdraw_msg_states": "WithdrawMsgState",
				}},
				{Name: "GenesisState", Path: "testdata/liquidity/genesis.proto", Fields: map[string]string{
					"params":       "Params",
					"pool_records": "PoolRecord",
				}},
				{Name: "PoolType", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"id":                   "uint32",
					"name":                 "string",
					"min_reserve_coin_num": "uint32",
					"max_reserve_coin_num": "uint32",
					"description":          "string",
				}},
				{Name: "Params", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"max_reserve_coin_amount":    "string",
					"withdraw_fee_rate":          "bytes",
					"unit_batch_height":          "uint32",
					"pool_types":                 "PoolType",
					"min_init_deposit_amount":    "string",
					"swap_fee_rate":              "bytes",
					"max_order_amount_ratio":     "bytes",
					"init_pool_coin_mint_amount": "string",
					"pool_creation_fee":          "cosmos.base.v1beta1.Coin",
				}},
				{Name: "Pool", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"id":                      "uint64",
					"type_id":                 "uint32",
					"reserve_coin_denoms":     "string",
					"reserve_account_address": "string",
					"pool_coin_denom":         "string",
				}},
				{Name: "PoolMetadata", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"pool_id":                "uint64",
					"pool_coin_total_supply": "cosmos.base.v1beta1.Coin",
					"reserve_coins":          "cosmos.base.v1beta1.Coin",
				}},
				{Name: "PoolMetadataResponse", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"pool_coin_total_supply": "cosmos.base.v1beta1.Coin",
					"reserve_coins":          "cosmos.base.v1beta1.Coin",
				}},
				{Name: "PoolBatch", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"executed":           "bool",
					"pool_id":            "uint64",
					"index":              "uint64",
					"begin_height":       "int64",
					"deposit_msg_index":  "uint64",
					"withdraw_msg_index": "uint64",
					"swap_msg_index":     "uint64",
				}},
				{Name: "PoolBatchResponse", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"index":              "uint64",
					"begin_height":       "int64",
					"deposit_msg_index":  "uint64",
					"withdraw_msg_index": "uint64",
					"swap_msg_index":     "uint64",
					"executed":           "bool",
				}},
				{Name: "DepositMsgState", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"succeeded":     "bool",
					"to_be_deleted": "bool",
					"msg":           "MsgDepositWithinBatch",
					"msg_height":    "int64",
					"msg_index":     "uint64",
					"executed":      "bool",
				}},
				{Name: "WithdrawMsgState", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"msg":           "MsgWithdrawWithinBatch",
					"msg_height":    "int64",
					"msg_index":     "uint64",
					"executed":      "bool",
					"succeeded":     "bool",
					"to_be_deleted": "bool",
				}},
				{Name: "SwapMsgState", Path: "testdata/liquidity/liquidity.proto", Fields: map[string]string{
					"executed":                "bool",
					"succeeded":               "bool",
					"order_expiry_height":     "int64",
					"exchanged_offer_coin":    "cosmos.base.v1beta1.Coin",
					"msg":                     "MsgSwapWithinBatch",
					"msg_height":              "int64",
					"msg_index":               "uint64",
					"to_be_deleted":           "bool",
					"remaining_offer_coin":    "cosmos.base.v1beta1.Coin",
					"reserved_offer_coin_fee": "cosmos.base.v1beta1.Coin",
				}},
				{Name: "QueryLiquidityPoolRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool_id": "uint64",
				}},
				{Name: "QueryLiquidityPoolResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool": "Pool",
				}},
				{Name: "QueryLiquidityPoolBatchRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool_id": "uint64",
				}},
				{Name: "QueryLiquidityPoolBatchResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"batch": "PoolBatch",
				}},
				{Name: "QueryLiquidityPoolsRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
				}},
				{Name: "QueryLiquidityPoolsResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pools":      "Pool",
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
				}},
				{Name: "QueryParamsRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{}},
				{Name: "QueryParamsResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"params": "Params",
				}},
				{Name: "QueryPoolBatchSwapMsgsRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
					"pool_id":    "uint64",
				}},
				{Name: "QueryPoolBatchSwapMsgRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool_id":   "uint64",
					"msg_index": "uint64",
				}},
				{Name: "QueryPoolBatchSwapMsgsResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"swaps":      "SwapMsgState",
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
				}},
				{Name: "QueryPoolBatchSwapMsgResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"swap": "SwapMsgState",
				}},
				{Name: "QueryPoolBatchDepositMsgsRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool_id":    "uint64",
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
				}},
				{Name: "QueryPoolBatchDepositMsgRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"msg_index": "uint64",
					"pool_id":   "uint64",
				}},
				{Name: "QueryPoolBatchDepositMsgsResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"deposits":   "DepositMsgState",
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
				}},
				{Name: "QueryPoolBatchDepositMsgResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"deposit": "DepositMsgState",
				}},
				{Name: "QueryPoolBatchWithdrawMsgsRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool_id":    "uint64",
					"pagination": "cosmos.base.query.v1beta1.PageRequest",
				}},
				{Name: "QueryPoolBatchWithdrawMsgRequest", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"pool_id":   "uint64",
					"msg_index": "uint64",
				}},
				{Name: "QueryPoolBatchWithdrawMsgsResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"withdraws":  "WithdrawMsgState",
					"pagination": "cosmos.base.query.v1beta1.PageResponse",
				}},
				{Name: "QueryPoolBatchWithdrawMsgResponse", Path: "testdata/liquidity/query.proto", Fields: map[string]string{
					"withdraw": "WithdrawMsgState",
				}},
				{Name: "MsgCreatePool", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"pool_creator_address": "string",
					"pool_type_id":         "uint32",
					"deposit_coins":        "cosmos.base.v1beta1.Coin",
				}},
				{Name: "MsgCreatePoolRequest", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"base_req": "BaseReq",
					"msg":      "MsgCreatePool",
				}},
				{Name: "MsgCreatePoolResponse", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "MsgDepositWithinBatch", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"depositor_address": "string",
					"pool_id":           "uint64",
					"deposit_coins":     "cosmos.base.v1beta1.Coin",
				}},
				{Name: "MsgDepositWithinBatchRequest", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"pool_id":  "uint64",
					"msg":      "MsgDepositWithinBatch",
					"base_req": "BaseReq",
				}},
				{Name: "MsgDepositWithinBatchResponse", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "MsgWithdrawWithinBatch", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"pool_coin":          "cosmos.base.v1beta1.Coin",
					"withdrawer_address": "string",
					"pool_id":            "uint64",
				}},
				{Name: "MsgWithdrawWithinBatchRequest", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"base_req": "BaseReq",
					"pool_id":  "uint64",
					"msg":      "MsgWithdrawWithinBatch",
				}},
				{Name: "MsgWithdrawWithinBatchResponse", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "MsgSwapWithinBatch", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"offer_coin":             "cosmos.base.v1beta1.Coin",
					"demand_coin_denom":      "string",
					"offer_coin_fee":         "cosmos.base.v1beta1.Coin",
					"order_price":            "bytes",
					"swap_requester_address": "string",
					"pool_id":                "uint64",
					"swap_type_id":           "uint32",
				}},
				{Name: "MsgSwapWithinBatchRequest", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"base_req": "BaseReq",
					"pool_id":  "uint64",
					"msg":      "MsgSwapWithinBatch",
				}},
				{Name: "MsgSwapWithinBatchResponse", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"std_tx": "StdTx",
				}},
				{Name: "BaseReq", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"memo":           "string",
					"account_number": "uint64",
					"fees":           "cosmos.base.v1beta1.Coin",
					"gas_adjustment": "string",
					"simulate":       "bool",
					"gas":            "uint64",
					"from":           "string",
					"chain_id":       "string",
					"sequence":       "uint64",
					"timeout_height": "uint64",
					"gas_prices":     "cosmos.base.v1beta1.DecCoin",
				}},
				{Name: "Fee", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"gas":    "uint64",
					"amount": "cosmos.base.v1beta1.Coin",
				}},
				{Name: "PubKey", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"type":  "string",
					"value": "string",
				}},
				{Name: "Signature", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"pub_key":        "PubKey",
					"account_number": "uint64",
					"sequence":       "uint64",
					"signature":      "string",
				}},
				{Name: "StdTx", Path: "testdata/liquidity/tx.proto", Fields: map[string]string{
					"fee":       "Fee",
					"memo":      "string",
					"signature": "Signature",
					"msg":       "string",
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
