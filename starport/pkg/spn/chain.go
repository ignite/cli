package spn

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	genesistypes "github.com/tendermint/spn/x/genesis/types"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
	"time"
)

// chainListOptions holds chain listing options.
type chainListOptions struct {
	paginationKey   []byte
	paginationLimit uint64
}

// ChainListOption configures chain listing options.
type ChainListOption func(*chainListOptions)

// PaginateChainListing sets pagination for chain listing.
func PaginateChainListing(key []byte, limit uint64) ChainListOption {
	return func(o *chainListOptions) {
		o.paginationKey = key
		o.paginationLimit = limit
	}
}

// ChainList lists chain summaries
func (c *Client) ChainList(ctx context.Context, accountName string, options ...ChainListOption) (chains []Chain, nextPageKey []byte, err error) {
	o := &chainListOptions{}
	for _, apply := range options {
		apply(o)
	}

	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return nil, nil, err
	}

	q := genesistypes.NewQueryClient(clientCtx)
	chainList, err := q.ListChains(ctx, &genesistypes.QueryListChainsRequest{
		Pagination: &query.PageRequest{
			Key:   o.paginationKey,
			Limit: o.paginationLimit,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	for _, c := range chainList.Chains {
		chains = append(chains, toChain(c))
	}

	return chains, chainList.Pagination.NextKey, nil
}

// ChainCreate creates a new chain.
func (c *Client) ChainCreate(ctx context.Context, accountName, chainID string, sourceURL, sourceHash string) error {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}
	return c.broadcast(ctx, clientCtx, genesistypes.NewMsgChainCreate(
		chainID,
		clientCtx.GetFromAddress(),
		sourceURL,
		sourceHash,
	))
}

// GenesisAccount represents a genesis account inside a chain with its allocated coins.
type GenesisAccount struct {
	Address types.AccAddress
	Coins   types.Coins
}

// Chain represents a chain in Genesis module of SPN.
type Chain struct {
	ChainID   string
	Creator   string
	URL       string
	Hash      string
	CreatedAt time.Time
}

// ShowChain shows chain info.
func (c *Client) ShowChain(ctx context.Context, accountName, chainID string) (Chain, error) {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return Chain{}, err
	}

	// Query the chain from spnd
	q := genesistypes.NewQueryClient(clientCtx)
	res, err := q.ShowChain(ctx, &genesistypes.QueryShowChainRequest{
		ChainID: chainID,
	})
	if err != nil {
		return Chain{}, err
	}

	return toChain(res.Chain), nil
}

// toChain converts proto chain to Chain type.
func toChain(chain *genesistypes.Chain) Chain {
	return Chain{
		ChainID:   chain.ChainID,
		Creator:   chain.Creator,
		URL:       chain.SourceURL,
		Hash:      chain.SourceHash,
		CreatedAt: time.Unix(chain.CreatedAt, 0),
	}
}

// LaunchInformation keeps the chain's launch information.
type LaunchInformation struct {
	GenesisAccounts []GenesisAccount
	GenTxs          []jsondoc.Doc
	Peers           []string
}

// LaunchInformation retrieves chain's launch information.
func (c *Client) LaunchInformation(ctx context.Context, accountName, chainID string) (LaunchInformation, error) {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return LaunchInformation{}, err
	}

	// Query the chain from spnd
	q := genesistypes.NewQueryClient(clientCtx)
	res, err := q.LaunchInformation(ctx, &genesistypes.QueryLaunchInformationRequest{
		ChainID: chainID,
	})
	if err != nil {
		return LaunchInformation{}, err
	}

	// Get the genesis accounts
	var genesisAccounts []GenesisAccount
	for _, addAccountProposalPayload := range res.LaunchInformation.Accounts {
		genesisAccount := GenesisAccount{
			Address: addAccountProposalPayload.Address,
			Coins:   addAccountProposalPayload.Coins,
		}

		genesisAccounts = append(genesisAccounts, genesisAccount)
	}

	return LaunchInformation{
		GenesisAccounts: genesisAccounts,
		GenTxs:          jsondoc.ToDocs(res.LaunchInformation.GenTxs),
		Peers:           res.LaunchInformation.Peers,
	}, nil
}

// SimulatedLaunchInformation retrieves chain's simulated launch information.
func (c *Client) SimulatedLaunchInformation(ctx context.Context, accountName, chainID string, proposalIDs []int) (LaunchInformation, error) {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return LaunchInformation{}, err
	}

	// Convert proposal ids to int32
	var proposalIDs32 []int32
	for _, proposalID := range proposalIDs {
		proposalIDs32 = append(proposalIDs32, int32(proposalID))
	}

	// Query the chain from spnd
	q := genesistypes.NewQueryClient(clientCtx)
	res, err := q.SimulatedLaunchInformation(ctx, &genesistypes.QuerySimulatedLaunchInformationRequest{
		ChainID:     chainID,
		ProposalIDs: proposalIDs32,
	})
	if err != nil {
		return LaunchInformation{}, err
	}

	// Get the genesis accounts
	var genesisAccounts []GenesisAccount
	for _, addAccountProposalPayload := range res.LaunchInformation.Accounts {
		genesisAccount := GenesisAccount{
			Address: addAccountProposalPayload.Address,
			Coins:   addAccountProposalPayload.Coins,
		}

		genesisAccounts = append(genesisAccounts, genesisAccount)
	}

	return LaunchInformation{
		GenesisAccounts: genesisAccounts,
		GenTxs:          jsondoc.ToDocs(res.LaunchInformation.GenTxs),
		Peers:           res.LaunchInformation.Peers,
	}, nil
}
