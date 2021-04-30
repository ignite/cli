package starport_network_test

import (
	"context"
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

var verificationError networkbuilder.VerificationError

func TestVerifyProposals(t *testing.T) {
	ctx := context.Background()

	// create a chain
	nb, err := initializeNetworkBuilder()
	require.NoError(t, err)
	chainID := "mercury"
	blockchain, err := initializeGaia(ctx, t, nb, chainID, false)
	require.NoError(t, err)

	err = blockchain.Create(ctx)
	require.NoError(t, err)

	// create join proposals (AddAccount + AddValidator)
	err = nb.AccountUse(spnValidator1)
	require.NoError(t, err)

	account, err := blockchain.CreateAccount(ctx, chain.Account{Name: "alice"})
	require.NoError(t, err)
	account.Coins = "1000token,1000000000stake"

	proposal := proposalMock("alice")
	gentx, err := blockchain.IssueGentx(ctx, account, proposal)
	require.NoError(t, err)

	peer := peerMock()
	err = blockchain.Join(
		ctx,
		&account,
		account.Address,
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(100000000)),
	)
	require.NoError(t, err)

	// can verify join proposals
	err = nb.VerifyProposals(ctx, chainID, []int{0, 1}, ioutil.Discard)
	require.NoError(t, err)

	// can't verify a gentx without a genesis account
	err = nb.VerifyProposals(ctx, chainID, []int{1}, ioutil.Discard)
	require.Error(t, err)
	require.ErrorAs(t, err, &verificationError)

	// can't verify a proposal with an invalid account address
	invalidAddressAccount := account
	invalidAddressAccount.Address = "invalid address"
	err = blockchain.Join(
		ctx,
		&invalidAddressAccount,
		invalidAddressAccount.Address,
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(100000000)),
	)
	require.NoError(t, err)
	err = nb.VerifyProposals(ctx, chainID, []int{2, 3}, ioutil.Discard)
	require.Error(t, err)
	require.ErrorAs(t, err, &verificationError)

	// can't verify a proposal with an invalid self delegation
	err = blockchain.Join(
		ctx,
		&account,
		account.Address,
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(10000000000)), // bigger than account coins
	)
	require.NoError(t, err)
	err = nb.VerifyProposals(ctx, chainID, []int{4, 5}, ioutil.Discard)
	require.Error(t, err)
	require.ErrorAs(t, err, &verificationError)

	// can't verify a proposal with an invalid validator address
	err = blockchain.Join(
		ctx,
		&account,
		"cosmos1pmxhse92uugjm3dsltr6p0cfmweprt70q8qykq",
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(100000000)),
	)
	require.NoError(t, err)
	err = nb.VerifyProposals(ctx, chainID, []int{6, 7}, ioutil.Discard)
	require.Error(t, err)
	require.ErrorAs(t, err, &verificationError)
}
