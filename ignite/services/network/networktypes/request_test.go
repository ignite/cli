package networktypes_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/ignite/cli/ignite/services/network/networktypes"
)

func TestVerifyAddValidatorRequest(t *testing.T) {
	gentx := []byte(`{
  "body": {
    "messages": [
      {
        "delegator_address": "cosmos1dd246yq6z5vzjz9gh8cff46pll75yyl8ygndsj",
        "pubkey": {
          "@type": "/cosmos.crypto.ed25519.PubKey",
          "key": "aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs="
        },
        "validator_address": "cosmosvaloper1dd246yq6z5vzjz9gh8cff46pll75yyl8pu8cup",
        "value": {
          "amount": "95000000",
          "denom": "stake"
        }
      }
    ]
  }
}`)
	pk, err := base64.StdEncoding.DecodeString("aeQLCJOjXUyB7evOodI4mbrshIt3vhHGlycJDbUkaMs=")
	require.NoError(t, err)

	tests := []struct {
		name string
		req  *launchtypes.RequestContent_GenesisValidator
		want error
	}{
		{
			name: "valid genesis validator request",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
					GenTx:          gentx,
					ConsPubKey:     ed25519.PubKey(pk),
					SelfDelegation: sdk.NewCoin("stake", sdk.NewInt(95000000)),
					Peer:           launchtypes.NewPeerConn("nodeid", "127.163.0.1:2446"),
				},
			},
		},
		{
			name: "invalid peer host",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
					GenTx:          gentx,
					ConsPubKey:     ed25519.PubKey(pk),
					SelfDelegation: sdk.NewCoin("stake", sdk.NewInt(95000000)),
					Peer:           launchtypes.NewPeerConn("nodeid", "122.114.800.11"),
				},
			},
			want: fmt.Errorf("the peer address id:\"nodeid\" tcpAddress:\"122.114.800.11\"  doesn't match the peer format <host>:<port>"),
		},
		{
			name: "invalid gentx",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
					GenTx:          []byte(`{}`),
					ConsPubKey:     ed25519.PubKey(pk),
					SelfDelegation: sdk.NewCoin("stake", sdk.NewInt(95000000)),
					Peer:           launchtypes.NewPeerConn("nodeid", "127.163.0.1:2446"),
				},
			},
			want: fmt.Errorf("cannot parse gentx the gentx cannot be parsed"),
		},
		{
			name: "invalid self delegation denom",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
					GenTx:          gentx,
					ConsPubKey:     ed25519.PubKey(pk),
					SelfDelegation: sdk.NewCoin("foo", sdk.NewInt(95000000)),
					Peer:           launchtypes.NewPeerConn("nodeid", "127.163.0.1:2446"),
				},
			},
			want: fmt.Errorf("the self delegation 95000000foo doesn't match the one inside the gentx 95000000stake"),
		},
		{
			name: "invalid self delegation value",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
					GenTx:          gentx,
					ConsPubKey:     ed25519.PubKey(pk),
					SelfDelegation: sdk.NewCoin("stake", sdk.NewInt(3)),
					Peer:           launchtypes.NewPeerConn("nodeid", "127.163.0.1:2446"),
				},
			},
			want: fmt.Errorf("the self delegation 3stake doesn't match the one inside the gentx 95000000stake"),
		},
		{
			name: "invalid consensus pub key",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g",
					GenTx:          gentx,
					ConsPubKey:     ed25519.PubKey("cosmos1gkheudhhjsvq0s8fxt7p6pwe0k3k30kepcnz9p="),
					SelfDelegation: sdk.NewCoin("stake", sdk.NewInt(95000000)),
					Peer:           launchtypes.NewPeerConn("nodeid", "127.163.0.1:2446"),
				},
			},
			want: fmt.Errorf("the consensus pub key PubKeyEd25519{636F736D6F7331676B6865756468686A737671307338667874377036707765306B336B33306B6570636E7A39703D} doesn't match the one inside the gentx PubKeyEd25519{69E40B0893A35D4C81EDEBCEA1D23899BAEC848B77BE11C69727090DB52468CB}"),
		},
		{
			name: "invalid validator address",
			req: &launchtypes.RequestContent_GenesisValidator{
				GenesisValidator: &launchtypes.GenesisValidator{
					Address:        "spn1gkheudhhjsvq0s8fxt7p6pwe0k3k30keaytytm",
					GenTx:          gentx,
					ConsPubKey:     ed25519.PubKey(pk),
					SelfDelegation: sdk.NewCoin("stake", sdk.NewInt(95000000)),
					Peer:           launchtypes.NewPeerConn("nodeid", "127.163.0.1:2446"),
				},
			},
			want: fmt.Errorf("the validator address spn1gkheudhhjsvq0s8fxt7p6pwe0k3k30keaytytm doesn't match the one inside the gentx spn1dd246yq6z5vzjz9gh8cff46pll75yyl8c5tt7g"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := networktypes.VerifyAddValidatorRequest(tt.req)
			if tt.want != nil {
				require.Error(t, err)
				require.Equal(t, tt.want.Error(), err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
