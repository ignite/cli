package spn

import (
	"context"
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

// Represent a genesis account inside a chain with its allocated coins
type GenesisAccount struct {
	Address types.AccAddress
	Coins   types.Coins
}

// Chain represents a chain in Genesis module of SPN.
type Chain struct {
	URL             string
	Hash            string
	Peers           []string
	GenesisAccounts []GenesisAccount
	GenTxs          [][]byte
	CreatedAt       time.Time
}

// ChainGet shows chain info.
func (c *Client) ChainGet(ctx context.Context, accountName, chainID string) (Chain, error) {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return Chain{}, err
	}

	// Query the chain from spnd
	q := genesistypes.NewQueryClient(clientCtx)
	params := &genesistypes.QueryShowChainRequest{
		ChainID: chainID,
	}
	res, err := q.ShowChain(ctx, params)
	if err != nil {
		return Chain{}, err
	}

	// Get the updated genesis
	launchInformationReq := &genesistypes.QueryLaunchInformationRequest{
		ChainID: chainID,
	}
	launchInformationRes, err := q.LaunchInformation(ctx, launchInformationReq)
	if err != nil {
		return Chain{}, err
	}

	// Get the genesis accounts
	var genesisAccounts []GenesisAccount
	for _, addAccountProposalPayload := range launchInformationRes.Accounts {
		genesisAccount := GenesisAccount{
			Address: addAccountProposalPayload.Address,
			Coins:   addAccountProposalPayload.Coins,
		}

		genesisAccounts = append(genesisAccounts, genesisAccount)
	}

	return Chain{
		URL:             res.Chain.SourceURL,
		Hash:            res.Chain.SourceHash,
		Peers:           launchInformationRes.Peers,
		GenesisAccounts: genesisAccounts,
		GenTxs:          launchInformationRes.GenTxs,
		CreatedAt:       time.Unix(res.Chain.CreatedAt, 0),
	}, nil
}