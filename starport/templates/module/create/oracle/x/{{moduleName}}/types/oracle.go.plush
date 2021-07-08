package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgOracleData{}

	// OracleResultStoreKeyPrefix is a prefix for storing result
	OracleResultStoreKeyPrefix = []byte{0xff}

	// LastOracleIDKey is the key for the last request id
	LastOracleIDKey = []byte{0x01}
)

type (
	// OracleScriptID is the type-safe unique identifier type for oracle scripts.
	OracleScriptID int64

	// RequestID is the type-safe unique identifier type for data requests.
	RequestID int64
)

func NewMsgOracleData(
	creator string,
	oracleScriptID OracleScriptID,
	sourceChannel string,
	calldata *CallData,
	askCount uint64,
	minCount uint64,
	feeLimit sdk.Coins,
	requestKey string,
	prepareGas uint64,
	executeGas uint64,
) *MsgOracleData {
	return &MsgOracleData{
		Creator:        creator,
		OracleScriptID: int64(oracleScriptID),
		SourceChannel:  sourceChannel,
		Calldata:       calldata,
		AskCount:       askCount,
		MinCount:       minCount,
		FeeLimit:       feeLimit,
		RequestKey:     requestKey,
		PrepareGas:     prepareGas,
		ExecuteGas:     executeGas,
	}
}

func (m *MsgOracleData) Route() string {
	return RouterKey
}

func (m *MsgOracleData) Type() string {
	return "OracleData"
}

func (m *MsgOracleData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgOracleData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgOracleData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// ResultStoreKey is a function to generate key for each result in store
func ResultStoreKey(requestID RequestID) []byte {
	return append(OracleResultStoreKeyPrefix, int64ToBytes(int64(requestID))...)
}

func int64ToBytes(num int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(num))
	return result
}
