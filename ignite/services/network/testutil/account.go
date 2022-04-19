package testutil

import (
	"testing"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/stretchr/testify/assert"
)

const (
	TestAccountName = "test"
)

// NewTestAccount creates an account for test purposes using inmemory keyring backend
func NewTestAccount(t *testing.T, name string) cosmosaccount.Account {
	r, err := cosmosaccount.NewInMemory()
	assert.NoError(t, err)
	account, _, err := r.Create(name)
	assert.NoError(t, err)
	return account
}
