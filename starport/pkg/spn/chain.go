package spn

import (
	"context"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	genesistypes "github.com/tendermint/spn/x/genesis/types"
)

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