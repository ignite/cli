package cosmosutil

import (
	"encoding/json"
	"errors"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var GentxFilename = "gentx.json"

type (
	// GentxInfo represents the basic info about gentx file
	GentxInfo struct {
		DelegatorAddress string
		PubKey           []byte
		SelfDelegation   sdk.Coin
	}
	// StargateGentx represents the stargate gentx file
	StargateGentx struct {
		Body struct {
			Messages []struct {
				DelegatorAddress string `json:"delegator_address"`
				ValidatorAddress string `json:"validator_address"`
				PubKey           struct {
					Key string `json:"key"`
				} `json:"pubkey"`
				Value struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"value"`
			} `json:"messages"`
		} `json:"body"`
	}
)

// GentxFromPath returns GentxInfo from the json file
func GentxFromPath(path string) (info GentxInfo, gentx []byte, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return info, gentx, errors.New("chain home folder is not initialized yet: " + path)
	}

	gentx, err = os.ReadFile(path)
	if err != nil {
		return info, gentx, err
	}
	return ParseGentx(gentx)
}

// ParseGentx returns GentxInfo and the gentx file in bytes
func ParseGentx(gentx []byte) (info GentxInfo, file []byte, err error) {
	// Try parsing Stargate gentx
	var stargateGentx StargateGentx
	if err := json.Unmarshal(gentx, &stargateGentx); err != nil {
		return info, gentx, err
	}
	if stargateGentx.Body.Messages == nil {
		return info, gentx, errors.New("the gentx cannot be parsed")
	}

	// This is a stargate gentx
	if len(stargateGentx.Body.Messages) != 1 {
		return info, gentx, errors.New("add validator gentx must contain 1 message")
	}

	info.DelegatorAddress = stargateGentx.Body.Messages[0].DelegatorAddress
	info.PubKey = []byte(stargateGentx.Body.Messages[0].PubKey.Key)

	amount, ok := sdk.NewIntFromString(stargateGentx.Body.Messages[0].Value.Amount)
	if !ok {
		return info, gentx, errors.New("the self-delegation inside the gentx is invalid")
	}

	info.SelfDelegation = sdk.NewCoin(
		stargateGentx.Body.Messages[0].Value.Denom,
		amount,
	)

	return info, gentx, nil
}
