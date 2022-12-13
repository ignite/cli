package cosmosutil

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// ChangeAddressPrefix returns the address with another prefix.
func ChangeAddressPrefix(address, newPrefix string) (string, error) {
	if newPrefix == "" {
		return "", errors.New("empty prefix")
	}
	_, pubKey, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	return bech32.ConvertAndEncode(newPrefix, pubKey)
}

// GetAddressPrefix returns the bech 32 prefix used by the address.
func GetAddressPrefix(address string) (string, error) {
	prefix, _, err := bech32.DecodeAndConvert(address)
	return prefix, err
}
