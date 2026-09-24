package accountview_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ignite/cli/v30/ignite/pkg/cliui/view/accountview"
)

// ansiRegex matches ANSI escape sequences. Views render styled strings since
// the lipgloss v2 upgrade; colors are stripped at the output writer.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestAccountString(t *testing.T) {
	tests := []struct {
		name    string
		account accountview.Account
		want    string
	}{
		{
			name:    "new account (mnemonic available) to string is not idented",
			account: accountview.NewAccount("alice", "cosmos193he38n21khnmb2", accountview.WithMnemonic("person estate daughter box chimney clay bronze ring story truck make excess ring frame desk start food leader sleep predict item rifle stem boy")),
			want:    "✔ Added account alice with address cosmos193he38n21khnmb2 and mnemonic:\nperson estate daughter box chimney clay bronze ring story truck make excess ring frame desk start food leader sleep predict item rifle stem boy\n",
		},
		{
			name:    "existent account to string is not idented",
			account: accountview.NewAccount("alice", "cosmos193he38n21khnmb2"),
			want:    "👤 alice's account address: cosmos193he38n21khnmb2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ansiRegex.ReplaceAllString(tt.account.String(), "")

			assert.NotEmpty(t, result)
			assert.Equal(t, tt.want, result)
		})
	}
}
