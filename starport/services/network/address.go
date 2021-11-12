package network

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// SetSPNPrefix change the address prefix to the spn prefix
func SetSPNPrefix(address string) (string, error) {
	_, pubKey, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	return bech32.ConvertAndEncode(SPNAddressPrefix, pubKey)
}
