package spn

import (
	"context"
	"errors"

	"github.com/cosmos/cosmos-sdk/types"
	genesistypes "github.com/tendermint/spn/x/genesis/types"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
)

// ProposalType represents the type of the proposal
type ProposalType string

const (
	ProposalTypeAll          ProposalType = ""
	ProposalTypeAddAccount   ProposalType = "add-account"
	ProposalTypeAddValidator ProposalType = "add-validator"
)

// ProposalStatus represents the status of the proposal
type ProposalStatus string

const (
	ProposalStatusAll      ProposalStatus = ""
	ProposalStatusPending  ProposalStatus = "pending"
	ProposalStatusApproved ProposalStatus = "approved"
	ProposalStatusRejected ProposalStatus = "rejected"
)

// Proposal represents a proposal.
type Proposal struct {
	ID        int                   `yaml:",omitempty"`
	Status    ProposalStatus        `yaml:",omitempty"`
	Account   *ProposalAddAccount   `yaml:",omitempty"`
	Validator *ProposalAddValidator `yaml:",omitempty"`
}

// ProposalAddAccount used to propose adding an account.
type ProposalAddAccount struct {
	Address string
	Coins   types.Coins
}

// ProposalAddValidator used to propose adding a validator.
type ProposalAddValidator struct {
	Gentx            jsondoc.Doc
	ValidatorAddress string
	SelfDelegation   types.Coin
	P2PAddress       string
}

var statusFromSPN = map[genesistypes.ProposalStatus]ProposalStatus{
	genesistypes.ProposalStatus_PENDING:  ProposalStatusPending,
	genesistypes.ProposalStatus_APPROVED: ProposalStatusApproved,
	genesistypes.ProposalStatus_REJECTED: ProposalStatusRejected,
}

var statusToSPN = map[ProposalStatus]genesistypes.ProposalStatus{
	ProposalStatusAll:      genesistypes.ProposalStatus_ANY_STATUS,
	ProposalStatusPending:  genesistypes.ProposalStatus_PENDING,
	ProposalStatusApproved: genesistypes.ProposalStatus_APPROVED,
	ProposalStatusRejected: genesistypes.ProposalStatus_REJECTED,
}

var proposalTypeToSPN = map[ProposalType]genesistypes.ProposalType{
	ProposalTypeAll:          genesistypes.ProposalType_ANY_TYPE,
	ProposalTypeAddAccount:   genesistypes.ProposalType_ADD_ACCOUNT,
	ProposalTypeAddValidator: genesistypes.ProposalType_ADD_VALIDATOR,
}

// proposalListOptions holds proposal listing options.
type proposalListOptions struct {
	typ    ProposalType
	status ProposalStatus
}

// ProposalListOption configures proposal listing options.
type ProposalListOption func(*proposalListOptions)

// ProposalListStatus sets proposal status filter for proposal listing.
func ProposalListStatus(status ProposalStatus) ProposalListOption {
	return func(o *proposalListOptions) {
		o.status = status
	}
}

// ProposalListType sets proposal type filter for proposal listing.
func ProposalListType(typ ProposalType) ProposalListOption {
	return func(o *proposalListOptions) {
		o.typ = typ
	}
}

// ProposalList lists proposals on a chain by status.
func (c *Client) ProposalList(ctx context.Context, acccountName, chainID string, options ...ProposalListOption) ([]Proposal, error) {
	o := &proposalListOptions{}
	for _, apply := range options {
		apply(o)
	}

	// Get spn proposal status
	spnStatus, ok := statusToSPN[o.status]
	if !ok {
		return nil, errors.New("unrecognized status")
	}

	// Get spn proposal type
	spnType, ok := proposalTypeToSPN[o.typ]
	if !ok {
		return nil, errors.New("unrecognized type")
	}

	var proposals []Proposal
	var spnProposals []*genesistypes.Proposal

	queryClient := genesistypes.NewQueryClient(c.clientCtx)

	// Send query
	res, err := queryClient.ListProposals(ctx, &genesistypes.QueryListProposalsRequest{
		ChainID: chainID,
		Status:  spnStatus,
		Type:    spnType,
	})
	if err != nil {
		return nil, err
	}
	spnProposals = res.Proposals

	// Format proposals
	for _, gp := range spnProposals {
		proposal, err := c.toProposal(*gp)
		if err != nil {
			return nil, err
		}

		proposals = append(proposals, proposal)
	}

	return proposals, nil
}

func (c *Client) toProposal(proposal genesistypes.Proposal) (Proposal, error) {
	p := Proposal{
		ID:     int(proposal.ProposalInformation.ProposalID),
		Status: statusFromSPN[proposal.ProposalState.GetStatus()],
	}
	switch payload := proposal.Payload.(type) {
	case *genesistypes.Proposal_AddAccountPayload:
		p.Account = &ProposalAddAccount{
			Address: payload.AddAccountPayload.Address.String(),
			Coins:   payload.AddAccountPayload.Coins,
		}

	case *genesistypes.Proposal_AddValidatorPayload:
		p.Validator = &ProposalAddValidator{
			P2PAddress: payload.AddValidatorPayload.Peer,
			Gentx:      payload.AddValidatorPayload.GenTx,
		}
	}

	return p, nil
}

func (c *Client) ProposalGet(ctx context.Context, accountName, chainID string, id int) (Proposal, error) {
	queryClient := genesistypes.NewQueryClient(c.clientCtx)

	// Query the proposal
	param := &genesistypes.QueryShowProposalRequest{
		ChainID:    chainID,
		ProposalID: int32(id),
	}
	res, err := queryClient.ShowProposal(ctx, param)
	if err != nil {
		return Proposal{}, err
	}

	return c.toProposal(*res.Proposal)
}

// ProposalOption configures Proposal to set a spesific type of proposal.
type ProposalOption func(*Proposal)

// AddAccountProposal creates an add account proposal option.
func AddAccountProposal(address string, coins types.Coins) ProposalOption {
	return func(p *Proposal) {
		p.Account = &ProposalAddAccount{address, coins}
	}
}

// AddValidatorProposal creates an add validator proposal option.
func AddValidatorProposal(gentx jsondoc.Doc, validatorAddress string, selfDelegation types.Coin, p2pAddress string) ProposalOption {
	return func(p *Proposal) {
		p.Validator = &ProposalAddValidator{gentx, validatorAddress, selfDelegation, p2pAddress}
	}
}

// Propose proposes given proposals in batch for chainID by using SPN accountName.
func (c *Client) Propose(ctx context.Context, accountName, chainID string, proposals ...ProposalOption) error {
	if len(proposals) == 0 {
		return errors.New("at least one proposal required")
	}

	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}

	var msgs []types.Msg

	for _, p := range proposals {
		var proposal Proposal
		p(&proposal)

		switch {
		case proposal.Account != nil:
			addr, err := types.AccAddressFromBech32(proposal.Account.Address)
			if err != nil {
				return err
			}

			// Create the proposal payload
			payload := genesistypes.NewProposalAddAccountPayload(
				addr,
				proposal.Account.Coins,
			)

			msgs = append(msgs, genesistypes.NewMsgProposalAddAccount(
				chainID,
				clientCtx.GetFromAddress(),
				payload,
			))

		case proposal.Validator != nil:
			// Get the validator address
			addr, err := types.AccAddressFromBech32(proposal.Validator.ValidatorAddress)
			if err != nil {
				return err
			}
			validatorAddress := types.ValAddress(addr)

			// Create the proposal payload
			payload := genesistypes.NewProposalAddValidatorPayload(
				proposal.Validator.Gentx,
				validatorAddress,
				proposal.Validator.SelfDelegation,
				proposal.Validator.P2PAddress,
			)

			msgs = append(msgs, genesistypes.NewMsgProposalAddValidator(
				chainID,
				clientCtx.GetFromAddress(),
				payload,
			))
		}
	}

	return c.broadcast(ctx, clientCtx, msgs...)
}
