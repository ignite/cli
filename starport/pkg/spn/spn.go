package spn

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	chattypes "github.com/tendermint/spn/x/chat/types"
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
func New(nodeAddress string, option ...Option) (Client, error) {
	opts := &options{
		keyringBackend: keyring.BackendOS,
	}
	for _, o := range option {
		o(opts)
	}
	kr, err := keyring.New(types.KeyringServiceName(), opts.keyringBackend, homedir, os.Stdin)
	if err != nil {
		return Client{}, err
	}

	client, err := rpchttp.New(xurl.TCP(nodeAddress), "/websocket")
	if err != nil {
		return Client{}, err
	}
	clientCtx := NewClientCtx(kr, client)
	factory := NewFactory(clientCtx)
	return Client{
		kr:        kr,
		factory:   factory,
		clientCtx: clientCtx,
	}, nil
}

// AccountGet retrieves an account by name from the keyring.
func (c Client) AccountGet(accountName string) (Account, error) {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return Account{}, err
	}
	return toAccount(info), nil
}

// AccountList returns a list of accounts.
func (c Client) AccountList() ([]Account, error) {
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

// AccountCreate creates an account by name in the keyring.
func (c Client) AccountCreate(accountName string) (Account, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return Account{}, err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return Account{}, err
	}
	algos, _ := c.kr.SupportedAlgorithms()
	if err != nil {
		return Account{}, err
	}
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
func (c Client) AccountExport(accountName, password string) (privateKey string, err error) {
	return c.kr.ExportPrivKeyArmor(accountName, password)
}

// AccountImport imports an account to the keyring by account name, privateKey and decryption password.
func (c Client) AccountImport(accountName, privateKey, password string) error {
	return c.kr.ImportPrivKey(accountName, privateKey, password)
}

// ChainCreate creates a new chain.
// TODO right now this uses chat module, use genesis.
func (c Client) ChainCreate(ctx context.Context, accountName, chainID, genesis, sourceURL, sourceHash string) error {
	info, err := c.kr.Key(accountName)
	if err != nil {
		return err
	}
	clientCtx := c.clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress())
	msg, err := chattypes.NewMsgCreateChannel(
		clientCtx.GetFromAddress(),
		chainID,
		sourceURL,
		[]byte(genesis),
	)
	if err != nil {
		return err
	}
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.BroadcastTx(clientCtx, c.factory, msg)
}

// Chain represents a chain in Genesis module of SPN.
type Chain struct {
	URL     string
	Hash    string
	Genesis interface{}
}

// TODO ChainGet shows chain info.
func (c Client) ChainGet(ctx context.Context, accountName, chainID string) (Chain, error) {
	return Chain{
		URL:     "https://github.com/tendermint/spn",
		Hash:    "df49c9256dfcbd0096fd0a8acdd4907ba3332cd5",
		Genesis: mockGenesis,
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
	Gentx         interface{}
	PublicAddress string
}

// ProposalList lists proposals on a chain by status.
func (c Client) ProposalList(ctx context.Context, acocuntName, chainID string, status ProposalStatus) ([]Proposal, error) {
	return []Proposal{
		{
			Status: ProposalPending,
			Account: &ProposalAddAccount{
				"comos123",
				[]types.Coin{
					types.NewInt64Coin("token", 10),
					types.NewInt64Coin("stake", 20),
				},
			},
		},
		{
			Status: ProposalPending,
			Validator: &ProposalAddValidator{
				"agentx",
				"aurl",
			},
		},
	}, nil
}

// ProposalGet retrieves a proposal on a chain by id.
func (c Client) ProposalGet(ctx context.Context, accountName, chainID string, id int) (Proposal, error) {
	return Proposal{
		Status: ProposalPending,
		Validator: &ProposalAddValidator{
			"agentx",
			"aurl",
		},
	}, nil
}

// ProposeAddAccount proposes to add an account to chain.
func (c Client) ProposeAddAccount(ctx context.Context, accountName, chainID string, account ProposalAddAccount) error {
	return nil
}

// ProposeAddValidator proposes to add a validator to chain.
func (c Client) ProposeAddValidator(ctx context.Context, accountName, chainID string, validator ProposalAddValidator) error {
	return nil
}

// ProposalApprove approves a proposal by id.
func (c Client) ProposalApprove(ctx context.Context, accountName, chainID string, id int) error {
	return nil
}

// ProposalReject rejects a proposal by id.
func (c Client) ProposalReject(ctx context.Context, accountName, chainID string, id int) error {
	return nil
}

var mockGenesis = []byte(`{
  "genesis_time": "2020-11-11T14:32:55.850301112Z",
  "chain_id": "spn",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_num": 50
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {}
  },
  "app_hash": "",
  "app_state": {
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      },
      "accounts": [
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "cosmos15yhep24n3y2c5c4edp3uhrtuvxf0c0nuggddyq",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        }
      ]
    },
    "bank": {
      "params": {
        "send_enabled": [],
        "default_send_enabled": true
      },
      "balances": [
        {
          "address": "cosmos15yhep24n3y2c5c4edp3uhrtuvxf0c0nuggddyq",
          "coins": [
            {
              "denom": "stake",
              "amount": "100000000"
            },
            {
              "denom": "token",
              "amount": "1000"
            }
          ]
        }
      ],
      "supply": [],
      "denom_metadata": []
    },
    "capability": {
      "index": "1",
      "owners": []
    },
    "chat": {},
    "crisis": {
      "constant_fee": {
        "denom": "stake",
        "amount": "1000"
      }
    },
    "distribution": {
      "params": {
        "community_tax": "0.020000000000000000",
        "base_proposer_reward": "0.010000000000000000",
        "bonus_proposer_reward": "0.040000000000000000",
        "withdraw_addr_enabled": true
      },
      "fee_pool": {
        "community_pool": []
      },
      "delegator_withdraw_infos": [],
      "previous_proposer": "",
      "outstanding_rewards": [],
      "validator_accumulated_commissions": [],
      "validator_historical_rewards": [],
      "validator_current_rewards": [],
      "delegator_starting_infos": [],
      "validator_slash_events": []
    },
    "evidence": {
      "evidence": []
    },
    "genutil": {
      "gen_txs": [
        {
          "body": {
            "messages": [
              {
                "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                "description": {
                  "moniker": "mynode",
                  "identity": "",
                  "website": "",
                  "security_contact": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.100000000000000000",
                  "max_rate": "0.200000000000000000",
                  "max_change_rate": "0.010000000000000000"
                },
                "min_self_delegation": "1",
                "delegator_address": "cosmos15yhep24n3y2c5c4edp3uhrtuvxf0c0nuggddyq",
                "validator_address": "cosmosvaloper15yhep24n3y2c5c4edp3uhrtuvxf0c0nuduecgn",
                "pubkey": "cosmosvalconspub1zcjduepqamy2dk057cultaf5hehaskgy7kzj8pdfreh7l3t6w9zktw8u37hqd99nu6",
                "value": {
                  "denom": "stake",
                  "amount": "95000000"
                }
              }
            ],
            "memo": "1117892966a4b8399c355cf022c9f1a8c221a85e@192.168.1.20:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "AiPbhlk6HDOBSrW5oyO2flEG+t0EoDFzoILmaZySk88g"
                },
                "mode_info": {
                  "single": {
                    "mode": "SIGN_MODE_DIRECT"
                  }
                },
                "sequence": "0"
              }
            ],
            "fee": {
              "amount": [
                {
                  "denom": "stake",
                  "amount": "5000"
                }
              ],
              "gas_limit": "200000",
              "payer": "",
              "granter": ""
            }
          },
          "signatures": [
            "y71CjK/eKpNNuvPQ0wPmNHvqTNcN4+vam6jfBo+Dt/Mceq7wJD29g9h9QE8/4GW51c+gATm4f4z805yZlLZq1g=="
          ]
        }
      ]
    },
    "gov": {
      "starting_proposal_id": "1",
      "deposits": [],
      "votes": [],
      "proposals": [],
      "deposit_params": {
        "min_deposit": [
          {
            "denom": "stake",
            "amount": "10000000"
          }
        ],
        "max_deposit_period": "172800s"
      },
      "voting_params": {
        "voting_period": "172800s"
      },
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000"
      }
    },
    "ibc": {
      "client_genesis": {
        "clients": [],
        "clients_consensus": [],
        "create_localhost": true
      },
      "connection_genesis": {
        "connections": [],
        "client_connection_paths": []
      },
      "channel_genesis": {
        "channels": [],
        "acknowledgements": [],
        "commitments": [],
        "send_sequences": [],
        "recv_sequences": [],
        "ack_sequences": []
      }
    },
    "identity": {},
    "mint": {
      "minter": {
        "inflation": "0.130000000000000000",
        "annual_provisions": "0.000000000000000000"
      },
      "params": {
        "mint_denom": "stake",
        "inflation_rate_change": "0.130000000000000000",
        "inflation_max": "0.200000000000000000",
        "inflation_min": "0.070000000000000000",
        "goal_bonded": "0.670000000000000000",
        "blocks_per_year": "6311520"
      }
    },
    "params": null,
    "slashing": {
      "params": {
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600s",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": [],
      "missed_blocks": []
    },
    "staking": {
      "params": {
        "unbonding_time": "1814400s",
        "max_validators": 100,
        "max_entries": 7,
        "historical_entries": 100,
        "bond_denom": "stake"
      },
      "last_total_power": "0",
      "last_validator_powers": [],
      "validators": [],
      "delegations": [],
      "unbonding_delegations": [],
      "redelegations": [],
      "exported": false
    },
    "transfer": {
      "port_id": "transfer",
      "denom_traces": [],
      "params": {
        "send_enabled": true,
        "receive_enabled": true
      }
    },
    "upgrade": {}
  }
}`)
