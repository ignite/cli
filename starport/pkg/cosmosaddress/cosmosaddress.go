// Package cosmosaddress implements helper methods to interact with Cosmos-SDK address
package cosmosaddress

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// ChangePrefix returns the address with another prefix
func ChangePrefix(address, newPrefix string) (string, error) {
	if newPrefix == "" {
		return "", errors.New("empty prefix")
	}
	_, pubKey, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", err
	}
	return bech32.ConvertAndEncode(newPrefix, pubKey)
}

// GetPrefix returns the bech 32 prefix used by the address
func GetPrefix(address string) (string, error) {
	prefix, _, err := bech32.DecodeAndConvert(address)
	return prefix, err
}
