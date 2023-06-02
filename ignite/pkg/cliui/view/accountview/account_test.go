package accountview_test

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cliui/view/accountview"
	"github.com/stretchr/testify/assert"
)

func TestAccountString(t *testing.T) {
	tests := []struct {
		name    string
		account accountview.Account
		want    string
	}{
		{
			name:    "new account (mnemonic available) to string is not idented",
			account: accountview.NewAccount("alice", "cosmos193he38n21khnmb2", accountview.WithMnemonic("person estate daughter box chimney clay bronze ring story truck make excess ring frame desk start food leader sleep predict item rifle stem boy")),
			want:    "✔ Added account \x1b[1malice\x1b[0m with address cosmos193he38n21khnmb2 and mnemonic:\nperson estate daughter box chimney clay bronze ring story truck make excess ring frame desk start food leader sleep predict item rifle stem boy\n",
		},
		{
			name:    "existent account to string is not idented",
			account: accountview.NewAccount("alice", "cosmos193he38n21khnmb2"),
			want:    "👤 \x1b[1malice\x1b[0m's account address: cosmos193he38n21khnmb2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.String()

			assert.NotEmpty(t, result)
			assert.Equal(t, tt.want, result)
		})
	}
}
