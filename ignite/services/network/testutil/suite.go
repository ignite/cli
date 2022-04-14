package testutil

import (
	"testing"

	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/services/network/mocks"
)

type Suite struct {
	Account           cosmosaccount.Account
	ChainMock         *mocks.Chain
	CosmosClientMock  *mocks.CosmosClient
	LaunchQueryMock   *mocks.LaunchClient
	CampaignQueryMock *mocks.CampaignClient
	ProfileQueryMock  *mocks.ProfileClient
	RewardClient      *mocks.RewardClient
}

func (s *Suite) AssertAllMocks(t *testing.T) {
	s.ChainMock.AssertExpectations(t)
	s.LaunchQueryMock.AssertExpectations(t)
	s.CosmosClientMock.AssertExpectations(t)
	s.CampaignQueryMock.AssertExpectations(t)
	s.RewardClient.AssertExpectations(t)
	s.CosmosClientMock.AssertExpectations(t)
}

func NewSuite(account cosmosaccount.Account) Suite {
	return Suite{
		Account:           account,
		ChainMock:         new(mocks.Chain),
		CosmosClientMock:  new(mocks.CosmosClient),
		LaunchQueryMock:   new(mocks.LaunchClient),
		CampaignQueryMock: new(mocks.CampaignClient),
		ProfileQueryMock:  new(mocks.ProfileClient),
		RewardClient:      new(mocks.RewardClient),
	}
}
