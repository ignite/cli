package testutil

import "github.com/tendermint/starport/starport/pkg/cosmosaccount"

const (
	TestAccountName = "test"
)

func NewTestAccount(name string) (cosmosaccount.Account, error) {
	r, err := cosmosaccount.NewInMemory()
	if err != nil {
		return cosmosaccount.Account{}, err
	}
	account, _, err := r.Create(name)
	return account, err
}
