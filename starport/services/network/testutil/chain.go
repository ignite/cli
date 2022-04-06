package testutil

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/tendermint/starport/starport/services/network/mocks"
)

const (
	TestChainSourceHash = "testhash"
	TestChainSourceURL  = "http://example.com/test"
	TestChainName       = "test"
	TestChainChainID    = "test-1"
	TestPublicAddress   = "1.2.3.4"
	TestNodeID          = "9b1f4adbfb0c0b513040d914bfb717303c0eaa71"
)

const (
	TestLaunchID   = uint64(1)
	TestCampaignID = uint64(1)
	TestMainnetID  = uint64(1)
)

type ChainMockOption func(*chainMockOptions)

type chainMockOptions struct {
	chainIDMustFail          bool
	genesisPath              string
	genesisPathMustFail      bool
	defaultGentxPath         string
	defaultGentxPathMustFail bool
	nodeIDMustFail           bool
	cacheBinaryMustFail      bool
}

func WithChainIDFail() ChainMockOption {
	return func(options *chainMockOptions) {
		options.chainIDMustFail = true
	}
}

func WithGenesisPath(path string) ChainMockOption {
	return func(options *chainMockOptions) {
		options.genesisPath = path
	}
}

func WithGenesisPathFail() ChainMockOption {
	return func(options *chainMockOptions) {
		options.genesisPathMustFail = true
	}
}

func WithDefaultGentxPath(path string) ChainMockOption {
	return func(options *chainMockOptions) {
		options.defaultGentxPath = path
	}
}

func WithDefaultGentxPathFail() ChainMockOption {
	return func(options *chainMockOptions) {
		options.defaultGentxPathMustFail = true
	}
}

func WithNodeIDFail() ChainMockOption {
	return func(options *chainMockOptions) {
		options.nodeIDMustFail = true
	}
}

func WithCacheBinaryFail() ChainMockOption {
	return func(options *chainMockOptions) {
		options.cacheBinaryMustFail = true
	}
}

func NewChainMock(options ...ChainMockOption) *mocks.Chain {
	o := chainMockOptions{}
	for _, apply := range options {
		apply(&o)
	}
	chainMock := new(mocks.Chain)
	chainMock.On("SourceHash").Return(TestChainSourceHash)
	chainMock.On("SourceURL").Return(TestChainSourceURL)
	chainMock.On("Name").Return(TestChainName)

	if o.cacheBinaryMustFail {
		chainMock.On("CacheBinary", TestLaunchID).Return(errors.New("failed to cache binary"))
	} else {
		chainMock.On("CacheBinary", TestLaunchID).Return(nil)
	}

	if o.chainIDMustFail {
		chainMock.On("ChainID").Return("", errors.New("failed to get chainID"))
	} else {
		chainMock.On("ChainID").Return(TestChainChainID, nil)
	}

	if o.genesisPathMustFail {
		chainMock.On("GenesisPath").Return("", errors.New("failed to get genesis path"))
	} else if o.genesisPath != "" {
		chainMock.On("GenesisPath").Return(o.genesisPath, nil)
	}

	if o.defaultGentxPathMustFail {
		chainMock.On("DefaultGentxPath").Return("", errors.New("failed to get default gentx path"))
	} else if o.defaultGentxPath != "" {
		chainMock.On("DefaultGentxPath").Return(o.defaultGentxPath, nil)
	}

	if o.nodeIDMustFail {
		chainMock.On("NodeID", mock.Anything).Return("", errors.New("failed to get node id"))
	} else {
		chainMock.On("NodeID", mock.Anything).Return(TestNodeID, nil)
	}

	return chainMock
}
