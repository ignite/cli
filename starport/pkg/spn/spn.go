package spn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/go-bip39"
	genesistypes "github.com/tendermint/spn/x/genesis/types"
	"github.com/tendermint/starport/starport/pkg/jsondoc"
	"github.com/tendermint/starport/starport/pkg/xurl"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var spn = "spn"
var homedir = os.ExpandEnv("$HOME/spnd")

// Account represents an account on SPN.
type Account struct {
	Name     string
	Address  string
	Mnemonic string
}

// Client is client to interact with SPN.
type Client struct {
	kr        keyring.Keyring
	factory   tx.Factory
	clientCtx client.Context
	out       *bytes.Buffer
}

type options struct {
	keyringBackend string
}

// Option configures Client options.
type Option func(*options)

// Keyring uses given keyring type as storage.
func Keyring(keyring string) Option {
	return func(c *options) {
		c.keyringBackend = keyring
	}
}

// New creates a new SPN Client with nodeAddress of a full SPN node.
// by default, OS is used as keyring backend.
func New(nodeAddress string, option ...Option) (*Client, error) {
	opts := &options{
		keyringBackend: keyring.BackendOS,
	}
	for _, o := range option {
		o(opts)
	}
	kr, err := keyring.New(types.KeyringServiceName(), opts.keyringBackend, homedir, os.Stdin)
	if err != nil {
		return nil, err
	}

	client, err := rpchttp.New(xurl.TCP(nodeAddress), "/websocket")
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	clientCtx := NewClientCtx(kr, client, out)
	factory := NewFactory(clientCtx)
	return &Client{
		kr:        kr,
		factory:   factory,
		clientCtx: clientCtx,
		out:       out,
	}, nil
}

// AccountGet retrieves an account by name from the keyring.
func (c *Client) AccountGet(accountName string) (Account, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return Account{}, err
	}
	return toAccount(info), nil
}

// AccountList returns a list of accounts.
func (c *Client) AccountList() ([]Account, error) {
	var accounts []Account
	infos, err := c.kr.List()
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		accounts = append(accounts, toAccount(info))
	}
	return accounts, nil
}

// AccountCreate creates an account by name and mnemonic (optional) in the keyring.
func (c *Client) AccountCreate(accountName, mnemonic string) (Account, error) {
	if mnemonic == "" {
		entropySeed, err := bip39.NewEntropy(256)
		if err != nil {
			return Account{}, err
		}
		mnemonic, err = bip39.NewMnemonic(entropySeed)
		if err != nil {
			return Account{}, err
		}
	}
	algos, _ := c.kr.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algos)
	if err != nil {
		return Account{}, err
	}
	hdPath := hd.CreateHDPath(types.GetConfig().GetCoinType(), 0, 0).String()
	info, err := c.kr.NewAccount(accountName, mnemonic, "", hdPath, algo)
	if err != nil {
		return Account{}, err
	}
	account := toAccount(info)
	account.Mnemonic = mnemonic
	return account, nil
}

func toAccount(info keyring.Info) Account {
	ko, _ := keyring.Bech32KeyOutput(info)
	return Account{
		Name:    ko.Name,
		Address: ko.Address,
	}
}

// AccountExport exports an account in the keyring by name and an encryption password into privateKey.
// password later can be used to decrypt the privateKey.
func (c *Client) AccountExport(accountName, password string) (privateKey string, err error) {
	return c.kr.ExportPrivKeyArmor(accountName, password)
}

// AccountImport imports an account to the keyring by account name, privateKey and decryption password.
func (c *Client) AccountImport(accountName, privateKey, password string) error {
	return c.kr.ImportPrivKey(accountName, privateKey, password)
}

// ChainCreate creates a new chain.
func (c *Client) ChainCreate(ctx context.Context, accountName, chainID string, genesis []byte, sourceURL, sourceHash string) error {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}
	return c.broadcast(clientCtx, genesistypes.NewMsgChainCreate(
		chainID,
		clientCtx.GetFromAddress(),
		sourceURL,
		sourceHash,
		genesis,
	))
}

func (c *Client) buildClientCtx(accountName string) (client.Context, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return client.Context{}, err
	}
	return c.clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress()), nil
}

func (c *Client) broadcast(clientCtx client.Context, msg types.Msg) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	c.out.Reset()
	if err := tx.BroadcastTx(clientCtx, c.factory, msg); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("make sure that your SPN account has enough balance")
		}
		return err
	}
	out := struct {
		Code int    `json:"code"`
		Log  string `json:"raw_log"`
	}{}
	if err := json.NewDecoder(c.out).Decode(&out); err != nil {
		return err
	}
	if out.Code > 0 {
		return fmt.Errorf("SPN error with '%d' code: %s", out.Code, out.Log)
	}
	return nil
}

// Chain represents a chain in Genesis module of SPN.
type Chain struct {
	URL     string
	Hash    string
	Genesis jsondoc.Doc
	Peers   []string
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
	currentGenesisReq := &genesistypes.QueryCurrentGenesisRequest{
		ChainID: chainID,
	}
	currentGenesisRes, err := q.CurrentGenesis(ctx, currentGenesisReq)
	if err != nil {
		return Chain{}, err
	}

	return Chain{
		URL:     res.Chain.SourceURL,
		Hash:    res.Chain.SourceHash,
		Genesis: currentGenesisRes.Genesis,
		Peers:   res.Chain.Peers,
	}, nil
}

// ProposalStatus keeps a proposal's status state.
type ProposalStatus string

const (
	ProposalPending  = "pending"
	ProposalApproved = "approved"
	ProposalRejected = "rejected"
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
	Gentx         jsondoc.Doc
	PublicAddress string
}

// ProposalList lists proposals on a chain by status.
func (c *Client) ProposalList(ctx context.Context, acccountName, chainID string, status ProposalStatus) ([]Proposal, error) {
	var proposals []Proposal
	var spnProposals []*genesistypes.Proposal

	queryClient := genesistypes.NewQueryClient(c.clientCtx)

	switch status {
	case ProposalPending:
		res, err := queryClient.PendingProposals(ctx, &genesistypes.QueryPendingProposalsRequest{
			ChainID: chainID,
		})
		if err != nil {
			return nil, err
		}
		spnProposals = res.Proposals
	case ProposalApproved:
		res, err := queryClient.ApprovedProposals(ctx, &genesistypes.QueryApprovedProposalsRequest{
			ChainID: chainID,
		})
		if err != nil {
			return nil, err
		}
		spnProposals = res.Proposals
	case ProposalRejected:
		res, err := queryClient.RejectedProposals(ctx, &genesistypes.QueryRejectedProposalsRequest{
			ChainID: chainID,
		})
		if err != nil {
			return nil, err
		}
		spnProposals = res.Proposals
	}

	for _, gp := range spnProposals {
		proposal, err := c.toProposal(*gp)
		if err != nil {
			return nil, err
		}

		proposals = append(proposals, proposal)
	}

	return proposals, nil
}

var toStatus = map[genesistypes.ProposalState_Status]ProposalStatus{
	genesistypes.ProposalState_PENDING:  ProposalPending,
	genesistypes.ProposalState_APPROVED: ProposalApproved,
	genesistypes.ProposalState_REJECTED: ProposalRejected,
}

func (c *Client) toProposal(proposal genesistypes.Proposal) (Proposal, error) {
	p := Proposal{
		ID:     int(proposal.ProposalInformation.ProposalID),
		Status: toStatus[proposal.ProposalState.GetStatus()],
	}
	switch payload := proposal.Payload.(type) {
	case *genesistypes.Proposal_AddAccountPayload:
		p.Account = &ProposalAddAccount{
			Address: payload.AddAccountPayload.Address.String(),
			Coins:   payload.AddAccountPayload.Coins,
		}

	case *genesistypes.Proposal_AddValidatorPayload:
		p.Validator = &ProposalAddValidator{
			PublicAddress: payload.AddValidatorPayload.Peer,
		}

		// Marshal gentx
		gentx, err := c.clientCtx.JSONMarshaler.MarshalJSON(payload.AddValidatorPayload.GenTx)
		if err != nil {
			return Proposal{}, err
		}
		p.Validator.Gentx = gentx
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

// ProposeAddAccount proposes to add a validator to chain.
func (c *Client) ProposeAddAccount(ctx context.Context, accountName, chainID string, account ProposalAddAccount) error {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}

	addr, err := types.AccAddressFromBech32(account.Address)
	if err != nil {
		return err
	}

	// Create the proposal payload
	payload := genesistypes.NewProposalAddAccountPayload(
		addr,
		account.Coins,
	)

	msg := genesistypes.NewMsgProposalAddAccount(
		chainID,
		clientCtx.GetFromAddress(),
		payload,
	)

	return c.broadcast(clientCtx, msg)
}

// ProposeAddValidator proposes to add a validator to chain.
func (c *Client) ProposeAddValidator(ctx context.Context, accountName, chainID string, validator ProposalAddValidator) error {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}

	// Read the gentx
	var gentx txtypes.Tx
	err = clientCtx.JSONMarshaler.UnmarshalJSON(validator.Gentx, &gentx)
	if err != nil {
		return err
	}

	// Create the proposal payload
	payload := genesistypes.NewProposalAddValidatorPayload(
		gentx,
		validator.PublicAddress,
	)

	msg := genesistypes.NewMsgProposalAddValidator(
		chainID,
		clientCtx.GetFromAddress(),
		payload,
	)

	return c.broadcast(clientCtx, msg)
}

// ProposalApprove approves a proposal by id.
func (c *Client) ProposalApprove(ctx context.Context, accountName, chainID string, id int) error {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}

	// Create approve message
	msg := genesistypes.NewMsgApprove(chainID, int32(id), clientCtx.GetFromAddress())

	return c.broadcast(clientCtx, msg)
}

// ProposalReject rejects a proposal by id.
func (c *Client) ProposalReject(ctx context.Context, accountName, chainID string, id int) error {
	clientCtx, err := c.buildClientCtx(accountName)
	if err != nil {
		return err
	}

	// Create reject message
	msg := genesistypes.NewMsgReject(chainID, int32(id), clientCtx.GetFromAddress())

	return c.broadcast(clientCtx, msg)
}
