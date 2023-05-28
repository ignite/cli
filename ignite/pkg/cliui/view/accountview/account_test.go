package accountview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_String(t *testing.T) {
	t.Run("new account (mnemonic available) to string is not idented", func(t *testing.T) {
		cases := []struct {
			account  Account
			expected string
		}{
			{
				NewAccount("test", "test_address", WithMnemonic("test_mnemonic")),
				"✔ Added account \x1b[1mtest\x1b[0m with address test_address and mnemonic:\ntest_mnemonic\n",
			},
			{
				NewAccount("alice", "cosmos193he38n21khnmb2", WithMnemonic("person estate daughter box chimney clay bronze ring story truck make excess ring frame desk start food leader sleep predict item rifle stem boy")),
				"✔ Added account \x1b[1malice\x1b[0m with address cosmos193he38n21khnmb2 and mnemonic:\nperson estate daughter box chimney clay bronze ring story truck make excess ring frame desk start food leader sleep predict item rifle stem boy\n",
			},
		}

		for _, tc := range cases {
			output := tc.account.String()

			assert.NotEmpty(t, output)
			assert.Equal(t, tc.expected, output)
		}
	})
	t.Run("existent account to string is not idented", func(t *testing.T) {
		cases := []struct {
			account  Account
			expected string
		}{
			{
				NewAccount("test", "test_address"),
				"👤 \x1b[1mtest\x1b[0m's account address: test_address\n",
			},
			{
				NewAccount("alice", "cosmos193he38n21khnmb2"),
				"👤 \x1b[1malice\x1b[0m's account address: cosmos193he38n21khnmb2\n",
			},
		}

		for _, tc := range cases {
			output := tc.account.String()

			assert.NotEmpty(t, output)
			assert.Equal(t, tc.expected, output)
		}
	})
}
