package network

import (
	"encoding/json"
	"errors"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	GentxInfo struct {
		DelegatorAddress string
		ValidatorAddress string
		SelfDelegation   sdk.Coin
	}
	StargateGentx struct {
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
	ChainGenesis struct {
		AppState struct {
			Auth struct {
				Accounts []struct {
					Address       string `json:"address"`
					AccountNumber uint64 `json:"account_number"`
					Sequence      uint64 `json:"sequence"`
				} `json:"accounts"`
			} `json:"auth"`
		} `json:"app_state"`
	}
)

func (g ChainGenesis) HasAccount(address string) bool {
	for _, account := range g.AppState.Auth.Accounts {
		if account.Address == address {
			return true
		}
	}
	return false
}

func ParseGenesis(genesisPath string) (genesis ChainGenesis, err error) {
	genesisFile, err := os.ReadFile(genesisPath)
	if err != nil {
		return genesis, errors.New("cannot open genesis file: " + err.Error())
	}

	if err := json.Unmarshal(genesisFile, &genesis); err != nil {
		return genesis, err
	}
	return
}

func ParseGentx(gentxPath string) (info GentxInfo, gentx []byte, err error) {
	gentx, err = os.ReadFile(gentxPath)
	if err != nil {
		return info, gentx, errors.New("cannot open gentx file: " + err.Error())
	}

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
	info.ValidatorAddress = stargateGentx.Body.Messages[0].ValidatorAddress
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
