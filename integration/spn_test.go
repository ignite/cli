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
	spnValidator   = "validator"
)

func TestChainCreateAndJoin(t *testing.T) {
	spnClient, err := spn.New(
		"http://0.0.0.0:26657",
		"http://0.0.0.0:1317",
		"http://0.0.0.0:4500",
		spn.Keyring(keyring.BackendMemory),
	)

	// initialize network builder and create accounts
	nb, err := networkbuilder.New(spnClient)
	require.NoError(t, err)
	_, err = nb.AccountCreate(spnCoordinator, "")
	require.NoError(t, err)
	_, err = nb.AccountCreate(spnValidator, "")
	require.NoError(t, err)
	err = nb.AccountUse(spnCoordinator)
	require.NoError(t, err)

	chainID := "mars"
	chainSource := "https://github.com/cosmos/gaia"
	chainHome, err := os.MkdirTemp("", "spn-chain-home")
	require.NoError(t, err)

	// initialize the chain for spn
	sourceOption := networkbuilder.SourceRemote(chainSource)
	initOptions := []networkbuilder.InitOption{
		networkbuilder.InitializationHomePath(chainHome),
		networkbuilder.InitializationKeyringBackend(chaincmd.KeyringBackendTest),
	}
	blockchain, err := nb.Init(context.TODO(), chainID, sourceOption, initOptions...)
	require.NoError(t, err)
	defer blockchain.Cleanup()

	// can create the chain
	err = blockchain.Create(context.TODO())
	require.NoError(t, err)

	// fetch the chain launch information
	chainInfo, err := nb.ShowChain(context.TODO(), chainID)
	require.NoError(t, err)
	require.Equal(t, chainID, chainInfo.ChainID)
	launchInfo, err := nb.LaunchInformation(context.TODO(), chainID)
	require.NoError(t, err)
	require.Empty(t, launchInfo.Peers)
	require.Empty(t, launchInfo.GenTxs)
	require.Empty(t, launchInfo.GenesisAccounts)

	// can join a chain
	err = nb.AccountUse(spnValidator)
	require.NoError(t, err)

	account, err := blockchain.CreateAccount(context.TODO(), chain.Account{Name: "alice"})
	require.NoError(t, err)
	account.Coins = "1000token,1000000000stake"

	ip, err := ipify.GetIp()
	require.NoError(t, err)
	peer := fmt.Sprintf("%s:26656", ip)

	proposal := networkbuilder.Proposal{
		Validator: chain.Validator{
			Name:                    "alice",
			Moniker:                 "alice",
			StakingAmount:           "100000000stake",
			CommissionRate:          "0.10",
			CommissionMaxRate:       "0.20",
			CommissionMaxChangeRate: "0.01",
			MinSelfDelegation:       "1",
			GasPrices:               "0.025stake",
		},
		Meta: networkbuilder.ProposalMeta{
			Website:  "https://cosmos.network",
			Identity: "alice",
			Details:  "foo",
		},
	}
	gentx, err := blockchain.IssueGentx(context.TODO(), account, proposal)
	require.NoError(t, err)

	err = blockchain.Join(
		context.TODO(),
		&account,
		account.Address,
		peer,
		gentx,
		sdk.NewCoin("stake", sdk.NewInt(100000000)),
	)
	require.NoError(t, err)

	// check pending proposals
	proposals, err := nb.ProposalList(
		context.TODO(),
		chainID,
		spn.ProposalListStatus(spn.ProposalStatusPending),
		spn.ProposalListType(spn.ProposalTypeAll),
	)
	require.Len(t, proposals, 2)

	// can verify the proposals
	err = nb.AccountUse(spnCoordinator)
	require.NoError(t, err)
	err = nb.VerifyProposals(context.TODO(), chainID, []int{0, 1}, ioutil.Discard)
	require.NoError(t, err)

	// can approve the proposals
	_, broadcast, err := nb.SubmitReviewals(
		context.TODO(),
		chainID,
		spn.ApproveProposal(0),
		spn.ApproveProposal(1),
	)
	require.NoError(t, err)
	err = broadcast()
	require.NoError(t, err)

	// check approved proposals
	proposals, err = nb.ProposalList(
		context.TODO(),
		chainID,
		spn.ProposalListStatus(spn.ProposalStatusApproved),
		spn.ProposalListType(spn.ProposalTypeAll),
	)
	require.Len(t, proposals, 2)

	// check launch information
	launchInfo, err = nb.LaunchInformation(context.TODO(), chainID)
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
}
