package testdata

import (
	"encoding/hex"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
)

func NewResponse(data codec.ProtoMarshaler) cosmosclient.Response {
	marshaler := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	anyEncoded, _ := prototypes.MarshalAny(data)
	txData := &sdktypes.TxMsgData{Data: []*sdktypes.MsgData{
		{
			Data: anyEncoded.Value,
			//TODO: Find a better way
			MsgType: strings.TrimSuffix(anyEncoded.TypeUrl, "Response"),
		},
	}}
	encodedTxData, _ := marshaler.Marshal(txData)
	resp := cosmosclient.Response{
		TxResponse: &sdk.TxResponse{
			Data: hex.EncodeToString(encodedTxData),
		},
	}
	resp.SetCodec(marshaler)
	return resp
}
