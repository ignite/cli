package testdata

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
)

const (
	Mnemonic    = "pottery potato snake address catalog original alarm memory float try blood walnut aunt recall fortune dance leg oppose pact cup leg chapter above priority"
	Address     = "spn146z74afxwhe839zenzcsfyj60cwvv3x384nwz0"
	AccountName = "test"
)

func GetTestAccount() cosmosaccount.Account {
	kb := keyring.NewInMemory()
	path := hd.CreateHDPath(sdktypes.GetConfig().GetCoinType(), 0, 0).String()
	acc, _ := kb.NewAccount(AccountName, Mnemonic, "", path, hd.Secp256k1)
	return cosmosaccount.Account{
		Name: "test",
		Info: acc,
	}
}
