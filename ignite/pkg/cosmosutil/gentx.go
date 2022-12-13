package cosmosutil

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const GentxFilename = "gentx.json"

type (
	// GentxInfo represents the basic info about gentx file.
	GentxInfo struct {
		DelegatorAddress string
		PubKey           ed25519.PubKey
		SelfDelegation   sdk.Coin
		Memo             string
	}

	// Gentx represents the gentx file.
	Gentx struct {
		Body struct {
			Messages []struct {
				DelegatorAddress string `json:"delegator_address"`
				ValidatorAddress string `json:"validator_address"`
				PubKey           struct {
					Type string `json:"@type"`
					Key  string `json:"key"`
				} `json:"pubkey"`
				Value struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"value"`
			} `json:"messages"`
			Memo string `json:"memo"`
		} `json:"body"`
	}
)

// GentxFromPath returns GentxInfo from the json file.
func GentxFromPath(path string) (info GentxInfo, gentx []byte, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return info, gentx, errors.New("chain home folder is not initialized yet: " + path)
	}

	gentx, err = os.ReadFile(path)
	if err != nil {
		return info, gentx, err
	}

	info, err = ParseGentx(gentx)
	return info, gentx, err
}

// ParseGentx returns GentxInfo and the gentx file in bytes.
func ParseGentx(gentxBz []byte) (info GentxInfo, err error) {
	// Try parsing gentx
	var gentx Gentx
	if err := json.Unmarshal(gentxBz, &gentx); err != nil {
		return info, fmt.Errorf("unmarshal gentx: %w", err)
	}
	if gentx.Body.Messages == nil {
		return info, errors.New("the gentx cannot be parsed")
	}

	if len(gentx.Body.Messages) != 1 {
		return info, errors.New("add validator gentx must contain 1 message")
	}

	info.Memo = gentx.Body.Memo
	info.DelegatorAddress = gentx.Body.Messages[0].DelegatorAddress

	pb := gentx.Body.Messages[0].PubKey.Key
	info.PubKey, err = base64.StdEncoding.DecodeString(pb)
	if err != nil {
		return info, fmt.Errorf("invalid validator public key %w", err)
	}

	amount, ok := sdkmath.NewIntFromString(gentx.Body.Messages[0].Value.Amount)
	if !ok {
		return info, errors.New("the self-delegation inside the gentx is invalid")
	}

	info.SelfDelegation = sdk.NewCoin(
		gentx.Body.Messages[0].Value.Denom,
		amount,
	)

	return info, nil
}
