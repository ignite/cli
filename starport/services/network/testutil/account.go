package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
)

const (
	TestAccountName = "test"
)

func NewTestAccount(t *testing.T, name string) cosmosaccount.Account {
	r, err := cosmosaccount.NewInMemory()
	assert.NoError(t, err)
	account, _, err := r.Create(name)
	assert.NoError(t, err)
	return account
}
