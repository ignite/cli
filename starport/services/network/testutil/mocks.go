package testutil

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
)

//go:generate mockery --name CampaignClient --case underscore --output ../mocks
type CampaignClient interface {
	campaigntypes.QueryClient
}

//go:generate mockery --name ProfileClient --case underscore --output ../mocks
type ProfileClient interface {
	profiletypes.QueryClient
}

//go:generate mockery --name LaunchClient --case underscore --output ../mocks
type LaunchClient interface {
	launchtypes.QueryClient
}

//go:generate mockery --name RewardClient --case underscore --output ../mocks
type RewardClient interface {
	rewardtypes.QueryClient
}

//go:generate mockery --name AccountInfo --case underscore --output ../mocks
type AccountInfo interface {
	keyring.Info
}
