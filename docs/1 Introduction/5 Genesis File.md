# Genesis File

The first block in a blockchain, block 0, is traditionally called the "genesis" or "genesis block".

The genesis block is for all blockchains the starting point of history. As per definition of blockchain, each block contains a hash of all transactions that it embeds as well as a pointer to the previous block. Obviously, as the starting point - the genesis block does not have a pointer to any previous block.

The genesis block is usually the only block in the blockchain that cannot be found on same P2P network that you are about to start, it has to be shared in different ways - we will have a look at ways to share the genesis file in another tutorial.

Because it is the starting point for a blockchain, especially in Proof-of-Stake blockchains, it contains a list of initial addresses and balances. Furthermore, most of the times the genesis block defines which network you are using.

With Starport you will create a genesis file from your `config.yml`, it typically looks similar to this:

```json
{
  "genesis_time": "2021-01-15T12:43:55.453718Z",
  "chain_id": "blog",
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
      "max_bytes": "1048576"
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
          "address": "cosmos1ml7mett44vnvmxzku24ncujjl5gpztp4p3kag6",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "cosmos17ahmvzxgvkl0ccacnpj0nwstalww68z0jrz89q",
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
          "address": "cosmos1ml7mett44vnvmxzku24ncujjl5gpztp4p3kag6",
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
        },
        {
          "address": "cosmos17ahmvzxgvkl0ccacnpj0nwstalww68z0jrz89q",
          "coins": [
            {
              "denom": "token",
              "amount": "500"
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
    "coinz": {},
    "crisis": {
      "constant_fee": {
        "amount": "1000",
        "denom": "stake"
      }
    },
    "distribution": {
      "delegator_starting_infos": [],
      "delegator_withdraw_infos": [],
      "fee_pool": {
        "community_pool": []
      },
      "outstanding_rewards": [],
      "params": {
        "base_proposer_reward": "0.010000000000000000",
        "bonus_proposer_reward": "0.040000000000000000",
        "community_tax": "0.020000000000000000",
        "withdraw_addr_enabled": true
      },
      "previous_proposer": "",
      "validator_accumulated_commissions": [],
      "validator_current_rewards": [],
      "validator_historical_rewards": [],
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
                "delegator_address": "cosmos1ml7mett44vnvmxzku24ncujjl5gpztp4p3kag6",
                "validator_address": "cosmosvaloper1ml7mett44vnvmxzku24ncujjl5gpztp4y9zgyf",
                "pubkey": {
                  "@type": "/cosmos.crypto.ed25519.PubKey",
                  "key": "oHynTiWsRxbdwWAkJLFBk/T8TWiFATFrKWfxemP2AWw="
                },
                "value": {
                  "denom": "stake",
                  "amount": "100000000"
                }
              }
            ],
            "memo": "9e5f3e69aab836337554f3f8699aec8ee7978360@192.168.178.23:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "Aidt5eKiDmfXiXy7el+zI+i0VyD34e3pOoOzI3ljWXKt"
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
              "amount": [],
              "gas_limit": "200000",
              "payer": "",
              "granter": ""
            }
          },
          "signatures": [
            "/bELv1NNgp7941Ux8Zz6T0x5yKqKMjbX2Q+TzPB7IAMdr680fkZGMGOm2F0bX9hK2lapCp+3RnGgmxhMG/bHUw=="
          ]
        }
      ]
    },
    "gov": {
      "deposit_params": {
        "max_deposit_period": "172800s",
        "min_deposit": [
          {
            "amount": "10000000",
            "denom": "stake"
          }
        ]
      },
      "deposits": [],
      "proposals": [],
      "starting_proposal_id": "1",
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000"
      },
      "votes": [],
      "voting_params": {
        "voting_period": "172800s"
      }
    },
    "ibc": {
      "channel_genesis": {
        "ack_sequences": [],
        "acknowledgements": [],
        "channels": [],
        "commitments": [],
        "receipts": [],
        "recv_sequences": [],
        "send_sequences": []
      },
      "client_genesis": {
        "clients": [],
        "clients_consensus": [],
        "create_localhost": false
      },
      "connection_genesis": {
        "client_connection_paths": [],
        "connections": []
      }
    },
    "mint": {
      "minter": {
        "annual_provisions": "0.000000000000000000",
        "inflation": "0.130000000000000000"
      },
      "params": {
        "blocks_per_year": "6311520",
        "goal_bonded": "0.670000000000000000",
        "inflation_max": "0.200000000000000000",
        "inflation_min": "0.070000000000000000",
        "inflation_rate_change": "0.130000000000000000",
        "mint_denom": "stake"
      }
    },
    "params": null,
    "slashing": {
      "missed_blocks": [],
      "params": {
        "downtime_jail_duration": "600s",
        "min_signed_per_window": "0.500000000000000000",
        "signed_blocks_window": "100",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": []
    },
    "staking": {
      "delegations": [],
      "exported": false,
      "last_total_power": "0",
      "last_validator_powers": [],
      "params": {
        "bond_denom": "stake",
        "historical_entries": 100,
        "max_entries": 7,
        "max_validators": 100,
        "unbonding_time": "1814400s"
      },
      "redelegations": [],
      "unbonding_delegations": [],
      "validators": []
    },
    "transfer": {
      "denom_traces": [],
      "params": {
        "receive_enabled": true,
        "send_enabled": true
      },
      "port_id": "transfer"
    },
    "upgrade": {},
    "vesting": {}
  }
}
```

Looking closely at the genesis file, you can observe that it contains the initial state parameters for your blockchain application, furthermore it contains definitions and parameters of modules that you are using.

Apart from module definition and configuration, the genesis file holds the addresses for the initial stakeholders and validators of a blockchain. These reside in the `gentx` parameter, which is part of the `genutil`. When starting a blockchain, these validators should be part of the network in order to get the network running. Or at least 66% of the validators should be available in order to start the BFT consensus.

In order to setup your genesis file correctly, it is important to understand the `config.yml`, which is discussed in more depth [here](4%20Configuration.md).

## Summary

- The genesis block is the first block of a blockchain.
- The genesis block contains initial stakeholders and starting validators.
