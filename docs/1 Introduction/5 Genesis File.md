# Genesis File

The first block in a blockchain, block 0, is traditionally called the "genesis" or "genesis block".

The genesis block is for all blockchains the starting point of history. As per definition of blockchain, each block contains a hash of all transactions that it embeds as well as a pointer to the previous block. Obviously, as the starting point - the genesis block does not have a pointer to any previous block.

The genesis block is usually the only block in the blockchain that cannot be found on same P2P network that you are about to start, it has to be shared in different ways - we will have a look at ways to share the genesis file in another tutorial.

Because it is the starting point for a blockchain, especially in Proof-of-Stake blockchains, it contains a list of initial addresses and balances. Furthermore, most of the times the genesis block defines which network you are using.

With Starport you will create a genesis file from your `config.yml`, it typically looks similar to this:

```json
{
  "genesis_time": "2020-09-03T20:39:19.245733Z",
  "chain_id": "blog",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "app_hash": "",
  "app_state": {
    "supply": {
      "supply": []
    },
    "blog": {},
    "genutil": {
      "gentxs": [
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "cosmos-sdk/MsgCreateValidator",
                "value": {
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
                  "delegator_address": "cosmos1al6ytsuyhq3e0v0mcrq88sqwguqe02n2yjxuwj",
                  "validator_address": "cosmosvaloper1al6ytsuyhq3e0v0mcrq88sqwguqe02n2pxjfzp",
                  "pubkey": "cosmosvalconspub1zcjduepqg7sqkpgeqd0hnd025fztrlmzk0f3lrc8ea90cpfhz4cjq5m0h2rqkaxgd2",
                  "value": {
                    "denom": "stake",
                    "amount": "100000000"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySecp256k1",
                  "value": "A4hHCwk6n8Dq0TAg+BA2vz8wbFsB0gg8WEtkbVnsi3P9"
                },
                "signature": "tRvI8ypTM2z45jXJgFVQ2aJ0Q4Dz3KZW92MtYDl6OlgMXA4EA99eQPl6gzVskcOM3OB+OsKzmkX4AyyVrm4qOQ=="
              }
            ],
            "memo": "3bea7947b66f99b1a6585c67456191859838b709@192.168.2.191:26656"
          }
        }
      ]
    },
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
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "cosmos1al6ytsuyhq3e0v0mcrq88sqwguqe02n2yjxuwj",
            "coins": [
              {
                "denom": "stake",
                "amount": "100000000"
              },
              {
                "denom": "token",
                "amount": "1000"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "cosmos1s25rfpwsm23ene6krurgghre6u8dnyguku52nk",
            "coins": [
              {
                "denom": "token",
                "amount": "500"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        }
      ]
    },
    "bank": {
      "send_enabled": true
    },
    "staking": {
      "params": {
        "unbonding_time": "1814400000000000",
        "max_validators": 100,
        "max_entries": 7,
        "historical_entries": 0,
        "bond_denom": "stake"
      },
      "last_total_power": "0",
      "last_validator_powers": null,
      "validators": null,
      "delegations": null,
      "unbonding_delegations": null,
      "redelegations": null,
      "exported": false
    },
    "params": null
  }
}
```

Looking closely at the genesis file, you can observe that it contains the initial state parameters for your blockchain application, furthermore it contains definitions and parameters of modules that you are using.

Apart from module definition and configuration, the genesis file holds the addresses for the initial stakeholders and validators of a blockchain. These reside in the `gentx` parameter, which is part of the `genutil`. When starting a blockchain, these validators should be part of the network in order to get the network running. Or at least 66% of the validators should be available in order to start the BFT consensus.

In order to setup your genesis file correctly, it is important to understand the `config.yml`, which is discussed in more depth [here](4%20Configuration.md).

## Summary

- The genesis block is the first block of a blockchain.
- The genesis block contains initial stakeholders and starting validators.
