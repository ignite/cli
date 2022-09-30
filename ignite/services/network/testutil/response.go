package testutil

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/runtime/protoiface"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
)

// NewResponse creates cosmosclient.Response object from proto struct
// for using as a return result for a cosmosclient mock
func NewResponse(data protoiface.MessageV1) cosmosclient.Response {
	marshaler := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	anyEncoded, _ := codectypes.NewAnyWithValue(data)

	txData := &sdk.TxMsgData{MsgResponses: []*codectypes.Any{anyEncoded}}

	encodedTxData, _ := marshaler.Marshal(txData)
	resp := cosmosclient.Response{
		Codec: marshaler,
		TxResponse: &sdk.TxResponse{
			Data: hex.EncodeToString(encodedTxData),
		},
	}
	return resp
}
