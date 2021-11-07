package network

import (
	"encoding/json"
	"errors"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	GentxInfo struct {
		DelegatorAddress, ValidatorAddress string
		SelfDelegation                     sdk.Coin
	}
	stargateGentx struct {
		Body struct {
			Messages []struct {
				DelegatorAddress string `json:"delegator_address"`
				ValidatorAddress string `json:"validator_address"`
				Value            struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"value"`
			} `json:"messages"`
		} `json:"body"`
	}
)

func ParseGentx(gentxPath string) (info GentxInfo, err error) {
	gentx, err := os.ReadFile(gentxPath)
	if err != nil {
		return info, errors.New("cannot open gentx file: " + err.Error())
	}

	// Try parsing Stargate gentx
	var stargateGentx stargateGentx
	if err := json.Unmarshal(gentx, &stargateGentx); err != nil {
		return info, err
	}
	if stargateGentx.Body.Messages == nil {
		return info, errors.New("the gentx cannot be parsed")
	}

	// This is a stargate gentx
	if len(stargateGentx.Body.Messages) != 1 {
		return info, errors.New("add validator gentx must contain 1 message")
	}
	info.DelegatorAddress = stargateGentx.Body.Messages[0].DelegatorAddress
	info.ValidatorAddress = stargateGentx.Body.Messages[0].ValidatorAddress
	amount, ok := sdk.NewIntFromString(stargateGentx.Body.Messages[0].Value.Amount)
	if !ok {
		return info, errors.New("the self-delegation inside the gentx is invalid")
	}
	info.SelfDelegation = sdk.NewCoin(
		stargateGentx.Body.Messages[0].Value.Denom,
		amount,
	)

	return info, nil
}
