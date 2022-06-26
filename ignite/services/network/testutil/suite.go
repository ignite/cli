package testutil

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/ignite/cli/ignite/services/network/mocks"
)

// Suite is a mocks container, used to write less code for tests setup
type Suite struct {
	ChainMock         *mocks.Chain
	CosmosClientMock  *mocks.CosmosClient
	LaunchQueryMock   *mocks.LaunchClient
	CampaignQueryMock *mocks.CampaignClient
	ProfileQueryMock  *mocks.ProfileClient
	RewardClient      *mocks.RewardClient
	StakingClient     *mocks.StakingClient
}

// AssertAllMocks asserts all suite mocks expectations
func (s *Suite) AssertAllMocks(t *testing.T) {
	s.ChainMock.AssertExpectations(t)
	s.ProfileQueryMock.AssertExpectations(t)
	s.LaunchQueryMock.AssertExpectations(t)
	s.CosmosClientMock.AssertExpectations(t)
	s.CampaignQueryMock.AssertExpectations(t)
	s.RewardClient.AssertExpectations(t)
	s.StakingClient.AssertExpectations(t)
}

// NewSuite creates new suite with mocks
func NewSuite() Suite {
	cosmos := new(mocks.CosmosClient)
	cosmos.On("Context").Return(client.Context{})
	return Suite{
		ChainMock:         new(mocks.Chain),
		CosmosClientMock:  cosmos,
		LaunchQueryMock:   new(mocks.LaunchClient),
		CampaignQueryMock: new(mocks.CampaignClient),
		ProfileQueryMock:  new(mocks.ProfileClient),
		RewardClient:      new(mocks.RewardClient),
		StakingClient:     new(mocks.StakingClient),
	}
}
