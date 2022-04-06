package network

import (
	"context"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/starport/starport/pkg/cosmoserror"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/services/network/mocks"
	"github.com/tendermint/starport/starport/services/network/testutil"
)

const (
	TestDenom        = "stake"
	TestAmountString = "95000000"
	TestAmountInt    = int64(95000000)

	TestAccountRequestID          = uint64(1)
	TestGenesisValidatorRequestID = uint64(2)
)

func stubNetworkForJoin() Network {
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("GenesisValidator", mock.Anything, &launchtypes.QueryGetGenesisValidatorRequest{
		Address:  testutil.Address,
		LaunchID: testutil.TestLaunchID,
	}).Return(nil, cosmoserror.ErrNotFound)
	launchQueryMock.On("VestingAccount", mock.Anything, &launchtypes.QueryGetVestingAccountRequest{
		Address:  testutil.Address,
		LaunchID: testutil.TestLaunchID,
	}).Return(nil, cosmoserror.ErrNotFound)
	launchQueryMock.On("GenesisAccount", mock.Anything, &launchtypes.QueryGetGenesisAccountRequest{
		Address:  testutil.Address,
		LaunchID: testutil.TestLaunchID,
	}).Return(nil, cosmoserror.ErrNotFound)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On(
		"BroadcastTx",
		testutil.AccountName,
		mock.AnythingOfType("*types.MsgRequestAddValidator"),
	).Return(testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{
		RequestID:    TestGenesisValidatorRequestID,
		AutoApproved: false,
	}), nil)
	networkClientMock.On(
		"BroadcastTx",
		testutil.AccountName,
		mock.AnythingOfType("*types.MsgRequestAddAccount"),
	).Return(testutil.NewResponse(&launchtypes.MsgRequestAddAccountResponse{
		RequestID:    TestAccountRequestID,
		AutoApproved: false,
	}), nil)
	return Network{
		cosmos:      networkClientMock,
		account:     testutil.GetTestAccount(),
		launchQuery: launchQueryMock,
	}
}

func TestJoin(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.Nil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
}

func TestJoinWithCustomGentx(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxJson, err := gentx.JSON()
	assert.Nil(t, err)
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, gentxPath)
	require.Nil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 0)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 1)
	network.cosmos.(*mocks.CosmosClient).AssertCalled(t,
		"BroadcastTx",
		mock.Anything,
		&launchtypes.MsgRequestAddValidator{
			Creator:        testutil.Address,
			LaunchID:       testutil.TestLaunchID,
			ValAddress:     testutil.Address,
			GenTx:          gentxJson,
			ConsPubKey:     []byte{},
			SelfDelegation: sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)),
			Peer: launchtypes.Peer{
				Id: testutil.TestNodeID,
				Connection: &launchtypes.Peer_TcpAddress{
					TcpAddress: testutil.TestPublicAddress,
				},
			},
		},
	)
}

func TestJoinValidatorAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("GenesisValidator", mock.Anything, &launchtypes.QueryGetGenesisValidatorRequest{
		Address:  testutil.Address,
		LaunchID: testutil.TestLaunchID,
	}).Return(nil, nil)
	network.launchQuery = launchQueryMock
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.NotNil(t, joinErr)
	require.Errorf(t, joinErr, "validator %s already exist", testutil.Address)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
}

func TestJoinValidatorExistenceCheckFailed(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("GenesisValidator", mock.Anything, &launchtypes.QueryGetGenesisValidatorRequest{
		Address:  testutil.Address,
		LaunchID: testutil.TestLaunchID,
	}).Return(nil, errors.New("failed to get validator"))
	network.launchQuery = launchQueryMock
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
}

func TestJoinAddValidatorTxFailed(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On(
		"BroadcastTx",
		testutil.AccountName,
		mock.AnythingOfType("*types.MsgRequestAddValidator"),
	).Return(testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{}), errors.New("failed to add validator"))
	network.cosmos = networkClientMock
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
}

func TestJoinWithAccountRequest(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	coin := sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt))

	joinErr := network.Join(
		context.Background(),
		chain,
		testutil.TestLaunchID,
		testutil.TestPublicAddress,
		"",
		WithAccountRequest(sdk.NewCoins(coin)),
	)
	require.Nil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
}

func TestJoinWithAccountRequestAndAccountExistsInGenesis(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesis.AddAccount(testutil.Address)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	coin := sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt))

	joinErr := network.Join(
		context.Background(),
		chain,
		testutil.TestLaunchID,
		testutil.TestPublicAddress,
		"",
		WithAccountRequest(sdk.NewCoins(coin)),
	)
	require.NotNil(t, joinErr)
	require.Errorf(t, joinErr, "account %s already exist", testutil.Address)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
}

func TestJoinWithAccountRequestFailedToCreateAccount(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On(
		"BroadcastTx",
		testutil.AccountName,
		mock.AnythingOfType("*types.MsgRequestAddAccount"),
	).Return(testutil.NewResponse(&launchtypes.MsgRequestAddAccountResponse{}), errors.New("failed to create account"))
	network.cosmos = networkClientMock
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath), testutil.WithDefaultGentxPath(gentxPath))

	coin := sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt))

	joinErr := network.Join(
		context.Background(),
		chain,
		testutil.TestLaunchID,
		testutil.TestPublicAddress,
		"",
		WithAccountRequest(sdk.NewCoins(coin)),
	)
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
}

func TestJoinFailedToReadNodeID(t *testing.T) {
	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithNodeIDFail())

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
}

func TestJoinFailedToReadDefaultGentxPath(t *testing.T) {
	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithDefaultGentxPathFail())

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 0)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
}

func TestJoinFailedToReadGenesis(t *testing.T) {
	tmp := t.TempDir()
	gentx := testutil.NewGentx(testutil.Address, TestDenom, TestAmountString, "")
	gentxPath, err := gentx.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithGenesisPathFail(), testutil.WithDefaultGentxPath(gentxPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, "")
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
}

func TestJoinFailedToReadCustomGentx(t *testing.T) {
	gentxPath := "invalid/path"

	tmp := t.TempDir()
	genesis := testutil.NewGenesis(testutil.TestChainChainID)
	genesisPath, err := genesis.SaveTo(tmp)
	assert.Nil(t, err)

	network := stubNetworkForJoin()
	chain := testutil.NewChainMock(testutil.WithGenesisPath(genesisPath))

	joinErr := network.Join(context.Background(), chain, testutil.TestLaunchID, testutil.TestPublicAddress, gentxPath)
	require.NotNil(t, joinErr)
	chain.AssertNumberOfCalls(t, "NodeID", 1)
	chain.AssertNumberOfCalls(t, "GenesisPath", 1)
	chain.AssertNumberOfCalls(t, "DefaultGentxPath", 0)
	network.launchQuery.(*mocks.LaunchClient).AssertNumberOfCalls(t, "GenesisValidator", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)

}
