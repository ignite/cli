package testutil

import (
	"encoding/hex"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"google.golang.org/protobuf/runtime/protoiface"
)

func NewResponse(data protoiface.MessageV1) cosmosclient.Response {
	marshaler := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	anyEncoded, _ := prototypes.MarshalAny(data)
	txData := &sdk.TxMsgData{Data: []*sdk.MsgData{
		{
			Data: anyEncoded.Value,
			//TODO: Find a better way
			MsgType: strings.TrimSuffix(anyEncoded.TypeUrl, "Response"),
		},
	}}
	encodedTxData, _ := marshaler.Marshal(txData)
	resp := cosmosclient.Response{
		Codec: marshaler,
		TxResponse: &sdk.TxResponse{
			Data: hex.EncodeToString(encodedTxData),
		},
	}
	return resp
}
