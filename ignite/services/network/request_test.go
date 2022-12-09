package network

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/services/network/networktypes"
	"github.com/ignite/cli/ignite/services/network/testutil"
)

func TestRequestParamChange(t *testing.T) {
	t.Run("successfully send request", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			module         = "module"
			param          = "param"
			value          = []byte("value")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewParamChange(
						testutil.LaunchID,
						module,
						param,
						value,
					),
				),
			).
			Return(testutil.NewResponse(&launchtypes.MsgSendRequestResponse{
				RequestID:    0,
				AutoApproved: false,
			}), nil).
			Once()

		sendRequestError := network.SendParamChangeRequest(context.Background(), testutil.LaunchID, module, param, value)
		require.NoError(t, sendRequestError)
		suite.AssertAllMocks(t)
	})
}
