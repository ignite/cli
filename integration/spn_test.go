package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rdegges/go-ipify"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/spn"
	"github.com/tendermint/starport/starport/services/chain"
	"github.com/tendermint/starport/starport/services/networkbuilder"
)

var (
	spnCoordinator = "coordinator"
	spnValidator1  = "validator1"
	spnValidator2  = "validator2"
)

func initializeNetworkBuilder() (*networkbuilder.Builder, error) {
	spnClient, err := spn.New(
		"http://0.0.0.0:26657",
		"http://0.0.0.0:1317",
		"http://0.0.0.0:4500",
		spn.Keyring(keyring.BackendMemory),
	)

	// initialize network builder and create accounts
	nb, err := networkbuilder.New(spnClient)
	if err != nil {
		return nil, err
	}

	// create some accounts
	if _, err = nb.AccountCreate(spnCoordinator, ""); err != nil {
		return nil, err
	}
	if _, err = nb.AccountCreate(spnValidator1, ""); err != nil {
		return nil, err
	}
	if _, err = nb.AccountCreate(spnValidator2, ""); err != nil {
		return nil, err
	}
	if err = nb.AccountUse(spnCoordinator); err != nil {
		return nil, err
	}

	return nb, nil
}

func initializeGaia(
	ctx context.Context,
	t *testing.T,
	nb *networkbuilder.Builder,
	chainID string,
) (*networkbuilder.Blockchain, error) {
	chainSource := "https://github.com/cosmos/gaia"
	chainHome, err := os.MkdirTemp("", "spn-chain-home")
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { os.RemoveAll(chainHome) })

	// initialize the chain for spn
	sourceOption := networkbuilder.SourceRemote(chainSource)
	initOptions := []networkbuilder.InitOption{
		networkbuilder.InitializationHomePath(chainHome),
		networkbuilder.InitializationKeyringBackend(chaincmd.KeyringBackendTest),
	}
	blockchain, err := nb.Init(ctx, chainID, sourceOption, initOptions...)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { blockchain.Cleanup() })

	return blockchain, nil
}

func TestApproveProposals(t *testing.T) {
	ctx := context.Background()

	nb, err := initializeNetworkBuilder()
	require.NoError(t, err)

	chainID := "mars"
	blockchain, err := initializeGaia(ctx, t, nb, chainID)
	require.NoError(t, err)

	// can create the chain
	err = blockchain.Create(ctx)
	require.NoError(t, err)

	// fetch the chain launch information
	chainInfo, err := nb.ShowChain(ctx, chainID)
	require.NoError(t, err)
	require.Equal(t, chainID, chainInfo.ChainID)
	launchInfo, err := nb.LaunchInformation(ctx, chainID)
	require.NoError(t, err)
	require.Empty(t, launchInfo.Peers)
	require.Empty(t, launchInfo.GenTxs)
	require.Empty(t, launchInfo.GenesisAccounts)

	// can join a chain
	err = nb.AccountUse(spnValidator1)
	require.NoError(t, err)

	account, err := blockchain.CreateAccount(ctx, chain.Account{Name: "alice"})
	require.NoError(t, err)
	account.Coins = "1000token,1000000000stake"

	ip, err := ipify.GetIp()
	require.NoError(t, err)
	peer := fmt.Sprintf("%s:26656", ip)

	proposal := proposalMock("alice")
	gentx, err := blockchain.IssueGentx(ctx, account, proposal)
	require.NoError(t, err)

	err = blockchain.Join(
		ctx,
		&account,
		account.Address,
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(100000000)),
	)
	require.NoError(t, err)

	// check pending proposals
	proposals, err := nb.ProposalList(
		ctx,
		chainID,
		spn.ProposalListStatus(spn.ProposalStatusPending),
		spn.ProposalListType(spn.ProposalTypeAll),
	)
	require.Len(t, proposals, 2)

	// can verify the proposals
	err = nb.AccountUse(spnCoordinator)
	require.NoError(t, err)
	err = nb.VerifyProposals(ctx, chainID, []int{0, 1}, ioutil.Discard)
	require.NoError(t, err)

	// can approve the proposals
	_, broadcast, err := nb.SubmitReviewals(
		ctx,
		chainID,
		spn.ApproveProposal(0),
		spn.ApproveProposal(1),
	)
	require.NoError(t, err)
	err = broadcast()
	require.NoError(t, err)

	// check approved proposals
	proposals, err = nb.ProposalList(
		ctx,
		chainID,
		spn.ProposalListStatus(spn.ProposalStatusApproved),
		spn.ProposalListType(spn.ProposalTypeAll),
	)
	require.Len(t, proposals, 2)

	// check launch information
	launchInfo, err = nb.LaunchInformation(ctx, chainID)
	require.NoError(t, err)
	require.Len(t, launchInfo.Peers, 1)
	require.Len(t, launchInfo.GenTxs, 1)
	require.Len(t, launchInfo.GenesisAccounts, 1)
	require.Contains(t, launchInfo.Peers[0], peer)
	require.Equal(t, gentx, launchInfo.GenTxs[0])
	require.Equal(t, spn.GenesisAccount{
		account.Address,
		sdk.NewCoins(
			sdk.NewCoin("token", sdk.NewInt(1000)),
			sdk.NewCoin("stake", sdk.NewInt(1000000000)),
		),
	}, launchInfo.GenesisAccounts[0])

	// can reject proposals
	err = nb.AccountUse(spnValidator2)
	require.NoError(t, err)

	account, err = blockchain.CreateAccount(ctx, chain.Account{Name: "bob"})
	require.NoError(t, err)
	account.Coins = "1000token,1000000000stake"

	proposal = proposalMock("bob")
	gentx, err = blockchain.IssueGentx(ctx, account, proposal)
	require.NoError(t, err)

	err = blockchain.Join(
		ctx,
		&account,
		account.Address,
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(100000000)),
	)
	require.NoError(t, err)

	err = nb.AccountUse(spnCoordinator)
	require.NoError(t, err)
	_, broadcast, err = nb.SubmitReviewals(
		ctx,
		chainID,
		spn.ApproveProposal(2),
		spn.ApproveProposal(3),
	)
	require.NoError(t, err)
	err = broadcast()
	require.NoError(t, err)

	// check rejected proposals
	proposals, err = nb.ProposalList(
		ctx,
		chainID,
		spn.ProposalListStatus(spn.ProposalStatusRejected),
		spn.ProposalListType(spn.ProposalTypeAll),
	)
	require.Len(t, proposals, 2)
}

func proposalMock(name string) networkbuilder.Proposal {
	return networkbuilder.Proposal{
		Validator: chain.Validator{
			Name:                    name,
			Moniker:                 name,
			StakingAmount:           "100000000stake",
			CommissionRate:          "0.10",
			CommissionMaxRate:       "0.20",
			CommissionMaxChangeRate: "0.01",
			MinSelfDelegation:       "1",
			GasPrices:               "0.025stake",
		},
		Meta: networkbuilder.ProposalMeta{
			Website:  "https://cosmos.network",
			Identity: name,
			Details:  "foo",
		},
	}
}
