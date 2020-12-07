package spn

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	genesistypes "github.com/tendermint/spn/x/genesis/types"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
)

// ChainSummary represents the summary of a chain in Genesis module of SPN.
type ChainSummary struct {
	ChainID            string
	Source             string
	TotalValidators    int
	ApprovedValidators int
	TotalProposals     int
	ApprovedProposals  int
}

// ChainList lists chain summaries
func (c *Client) ChainList(ctx context.Context, accountName string) ([]ChainSummary, error) {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return []ChainSummary{}, err
	}

	// List the chains
	q := genesistypes.NewQueryClient(clientCtx)
	req := &genesistypes.QueryListChainsRequest{}
	chainList, err := q.ListChains(ctx, req)
	if err != nil {
		return []ChainSummary{}, err
	}

	// Get the summary of each chain
	chainSummaries := make([]ChainSummary, len(chainList.Chains))
	g, ctx := errgroup.WithCancel(ctx)

	for i, chain := range chainList.Chains {
		i, chain := i, chain // https://golang.org/doc/faq#closures_and_goroutines
		chainSummariesGroup.Go(func() error {
			chainSummaryGroup := new(errgroup.Group)

			var chainSummary ChainSummary
			chainSummary.ChainID = chain.ChainID
			chainSummary.Source = chain.SourceURL

			// Get the number of validators
			chainSummaryGroup.Go(func() error {
				reqValidators := &genesistypes.QueryListProposalsRequest{
					ChainID: chain.ChainID,
					Status:  genesistypes.ProposalStatus_ANY_STATUS,
					Type:    genesistypes.ProposalType_ADD_VALIDATOR,
				}
				resValidators, err := q.ListProposals(ctx, reqValidators)
				if err != nil {
					return err
				}
				chainSummary.TotalValidators = len(resValidators.Proposals)
				return nil
			})

			// Get the number of approved validators
			chainSummaryGroup.Go(func() error {
				reqApprovedValidators := &genesistypes.QueryListProposalsRequest{
					ChainID: chain.ChainID,
					Status:  genesistypes.ProposalStatus_APPROVED,
					Type:    genesistypes.ProposalType_ADD_VALIDATOR,
				}
				resApprovedValidators, err := q.ListProposals(ctx, reqApprovedValidators)
				if err != nil {
					return err
				}
				chainSummary.ApprovedValidators = len(resApprovedValidators.Proposals)
				return nil
			})

			// Get the number of proposals
			chainSummaryGroup.Go(func() error {
				reqProposals := &genesistypes.QueryListProposalsRequest{
					ChainID: chain.ChainID,
					Status:  genesistypes.ProposalStatus_ANY_STATUS,
					Type:    genesistypes.ProposalType_ANY_TYPE,
				}
				resProposals, err := q.ListProposals(ctx, reqProposals)
				if err != nil {
					return err
				}
				chainSummary.TotalProposals = len(resProposals.Proposals)
				return nil
			})

			// Get the number of approved proposals
			chainSummaryGroup.Go(func() error {
				reqApprovedProposals := &genesistypes.QueryListProposalsRequest{
					ChainID: chain.ChainID,
					Status:  genesistypes.ProposalStatus_APPROVED,
					Type:    genesistypes.ProposalType_ANY_TYPE,
				}
				resApprovedProposals, err := q.ListProposals(ctx, reqApprovedProposals)
				if err != nil {
					return err
				}
				chainSummary.ApprovedProposals = len(resApprovedProposals.Proposals)
				return nil
			})

			if err := chainSummaryGroup.Wait(); err != nil {
				return err
			}

			chainSummaries[i] = chainSummary

			return nil
		})
	}

	if err := chainSummariesGroup.Wait(); err != nil {
		return []ChainSummary{}, err
	}

	return chainSummaries, nil
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

	return Chain{
		ChainID:   res.Chain.ChainID,
		Creator:   res.Chain.Creator,
		URL:       res.Chain.SourceURL,
		Hash:      res.Chain.SourceHash,
		CreatedAt: time.Unix(res.Chain.CreatedAt, 0),
	}, nil
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
	for _, addAccountProposalPayload := range res.Accounts {
		genesisAccount := GenesisAccount{
			Address: addAccountProposalPayload.Address,
			Coins:   addAccountProposalPayload.Coins,
		}

		genesisAccounts = append(genesisAccounts, genesisAccount)
	}

	return LaunchInformation{
		GenesisAccounts: genesisAccounts,
		GenTxs:          jsondoc.ToDocs(res.GenTxs),
		Peers:           res.Peers,
	}, nil
}
