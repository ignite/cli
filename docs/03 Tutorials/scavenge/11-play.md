---
order: 11
---

# Play

Your application is running! That's great but who cares unless you can play with it. The first command you will want to try is creating a new scavenge. Since our user `user1` has way more `token` token than the user `user2`, let's create the scavenge from their account.

You can begin by running `scavengecli tx scavenge --help` to see all the commands we created for your new module. You should see the following options:
```bash
$ scavengecli tx scavenge
scavenge transactions subcommands

Usage:
  scavengecli tx scavenge [flags]
  scavengecli tx scavenge [command]

Available Commands:
  commit-solution Commits a solution for scavenge
  create-scavenge Creates a new scavenge
  reveal-solution Reveals a solution for scavenge

Flags:
  -h, --help   help for scavenge

Global Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/user/.scavengecli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "scavengecli tx scavenge [command] --help" for more information about a command.
```

We want to use the `create-scavenge` command so let's check the help screen for it as well like `scavengecli scavenge create-scavenge --help`. It should look like:

```bash
$ scavengecli tx scavenge create-scavenge --help
Creates a new scavenge

Usage:
  scavengecli tx scavenge create-scavenge [description] [solution] [reward] [flags]

Flags:
  -a, --account-number uint      The account number of the signing account (offline mode only)
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async|block) (default "sync")
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate required gas automatically (default 200000) (default "200000")
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices to determine the transaction fee (e.g. 10uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible and the node operates offline)
  -h, --help                     help for create-scavenge
      --indent                   Add indent to JSON response
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
      --ledger                   Use a connected Ledger device
      --memo string              Memo to send along with transaction
      --node string              <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --trust-node               Trust connected full node (don't verify proofs for responses) (default true)
  -y, --yes                      Skip tx broadcasting prompt confirmation

Global Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/b/.scavengecli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

Let's follow the instructions and create a new scavenge. The first parameter we need is the `reward`. Let's give away `69token` as a reward for solving our scavenge (nice).

Next we should list our `solution`, but probably we should also know what the actual quesiton is that our solution solves (our `description`). How about our challenge question be something family friendly like: `What's brown and sticky?`. Of course the only solution to this question is: `A stick`.

Now we have all the pieces needed to create our Message. Let's piece them all together, adding the flag `--from` so the CLI knows who is sending it:

```bash
scavengecli tx scavenge create-scavenge "What's brown and sticky?" "A stick" 69token --from user1
```

After confirming the message looks correct and signing it, you should see something like the following:

```json
{
  "height": "0",
  "txhash": "1319B81F05735A36BAEFAE8D8A308674E6DE054058CE8A69983069BB3CAECA5A",
  "raw_log": "[]"
}
```

This tells you that the message was accepted into the app. Whether the message failed afterwards can not be told from this screen. However, the section under `txhash` is like a receipt for this interaction. To see if it was successfully processed after being successfully included you can run the following command:

```bash
scavengecli q tx <txhash>
```

But replace the `<txhash>` with your own. You should see something similar to this afterwards:
```json
{
  "height": "55",
  "txhash": "1319B81F05735A36BAEFAE8D8A308674E6DE054058CE8A69983069BB3CAECA5A",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreateScavenge\"},{\"key\":\"sender\",\"value\":\"cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809\"},{\"key\":\"module\",\"value\":\"scavenge\"},{\"key\":\"action\",\"value\":\"CreateScavenge\"},{\"key\":\"sender\",\"value\":\"cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809\"},{\"key\":\"description\",\"value\":\"What's brown and sticky?\"},{\"key\":\"solutionHash\",\"value\":\"2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906\"},{\"key\":\"reward\",\"value\":\"69token\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg\"},{\"key\":\"sender\",\"value\":\"cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809\"},{\"key\":\"amount\",\"value\":\"69token\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "CreateScavenge"
            },
            {
              "key": "sender",
              "value": "cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809"
            },
            {
              "key": "module",
              "value": "scavenge"
            },
            {
              "key": "action",
              "value": "CreateScavenge"
            },
            {
              "key": "sender",
              "value": "cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809"
            },
            {
              "key": "description",
              "value": "What's brown and sticky?"
            },
            {
              "key": "solutionHash",
              "value": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906"
            },
            {
              "key": "reward",
              "value": "69token"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg"
            },
            {
              "key": "sender",
              "value": "cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809"
            },
            {
              "key": "amount",
              "value": "69token"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "200000",
  "gas_used": "56807",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "scavenge/CreateScavenge",
          "value": {
            "creator": "cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809",
            "description": "What's brown and sticky?",
            "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
            "reward": [
              {
                "denom": "token",
                "amount": "69"
              }
            ]
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
            "value": "AqEF1aTpf+8b+3rAvY6FGH725ezlFFU920cAuGumVAYI"
          },
          "signature": "dc5jz6VielT+sxx8VMr9hDD/UwNPDBd+79DJcaeMMZBdGWQer7G41OmZJ0VCULZ3jUIDy6P8xCYiWcFjAd3Yvg=="
        }
      ],
      "memo": ""
    }
  },
  "timestamp": "2020-09-22T01:22:03Z"
}
```

Here you can see all the events we defined within our `Handler` that describes exactly what happened when this message was processed. Since our message was formatted correctly and since the user `user1` had enough `token` to pay the bounty, our `Scavenge` was accepted. You can also see what the solution looks like now that it has been hashed:

```json
{
    "key": "solutionHash",
    "value": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906"
}
```

Since we know the solution to this question and since we have another user at hand that can submit it, let's begin the process of committing and revealing that solution.

First we should check the CLI command for `commit-solution` by running `scavengecli tx scavenge commit-solution --help` in order to see:
```bash
$ scavengecli tx scavenge commit-solution --help
Commits a solution for scavenge

Usage:
  scavengecli tx scavenge commit-solution [solution] [flags]

Flags:
  -a, --account-number uint      The account number of the signing account (offline mode only)
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async|block) (default "sync")
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate required gas automatically (default 200000) (default "200000")
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices to determine the transaction fee (e.g. 10uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible and the node operates offline)
  -h, --help                     help for commit-solution
      --indent                   Add indent to JSON response
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
      --ledger                   Use a connected Ledger device
      --memo string              Memo to send along with transaction
      --node string              <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --trust-node               Trust connected full node (don't verify proofs for responses) (default true)
  -y, --yes                      Skip tx broadcasting prompt confirmation

Global Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/billy/.scavengecli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
Let's follow the instructions and submit the answer as a commit on behalf of `user2`:
```bash
scavengecli tx scavenge commit-solution "A stick" --from user2 -y
```

We don't need to put the `solutionHash` because it can be generated by hashing our actual solution. This time we're passing the `-y` to auto-confirm the transaction. Afterwards, we should see our `txhash` again. To confirm the `txhash` let's look at it again with `scavengecli q tx <txhash>`. This time you should see something like:
```json
{
  "height": "105",
  "txhash": "66738350C719094A3854BFB6F4FC33536923DB0F9A3A7FD1E1A2C53B81CD1BF1",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"CreateCommit\"},{\"key\":\"module\",\"value\":\"scavenge\"},{\"key\":\"action\",\"value\":\"CommitSolution\"},{\"key\":\"sender\",\"value\":\"cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq\"},{\"key\":\"solutionHash\",\"value\":\"2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906\"},{\"key\":\"solutionScavengerHash\",\"value\":\"6eeeb9af010478a6972efd670de7f5b4335c00edc27931d30f880b0f873de8c5\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "CreateCommit"
            },
            {
              "key": "module",
              "value": "scavenge"
            },
            {
              "key": "action",
              "value": "CommitSolution"
            },
            {
              "key": "sender",
              "value": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq"
            },
            {
              "key": "solutionHash",
              "value": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906"
            },
            {
              "key": "solutionScavengerHash",
              "value": "6eeeb9af010478a6972efd670de7f5b4335c00edc27931d30f880b0f873de8c5"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "200000",
  "gas_used": "46856",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "scavenge/CommitSolution",
          "value": {
            "scavenger": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq",
            "solutionhash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
            "solutionScavengerHash": "6eeeb9af010478a6972efd670de7f5b4335c00edc27931d30f880b0f873de8c5"
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
            "value": "AreVRxbXB/ZS0D43BXLnITFCSf0ZvbRwP+7mCR5U7srD"
          },
          "signature": "/SSD708BkIE3TeGAd6DXfLMfgM3CITCn9Zav8Jv09lJr3R9++f2tZc4VrapxtHCjPT3a7gJhiO5TTUA2lIBQSg=="
        }
      ],
      "memo": ""
    }
  },
  "timestamp": "2020-09-22T01:26:15Z"
}
```
You'll notice that the `solutionHash` matches the one before. We've also created a new hash for the `solutionScavengerHash` which is the combination of the solution and our account address. We can make sure the commit has been made by querying it directly as well:

```bash
scavengecli q scavenge get-commit "6eeeb9af010478a6972efd670de7f5b4335c00edc27931d30f880b0f873de8c5"
```

Hopefully you should see something like:
```json
{
  "scavenger": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq",
  "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
  "solutionScavengerHash": "6eeeb9af010478a6972efd670de7f5b4335c00edc27931d30f880b0f873de8c5"
}
```
This confirms that your commit was successfully submitted and is awaiting the follow-up reveal. To make that command let's first check the `--help` command using `scavengecli tx scavenge reveal-solution --help`. This should show the following screen:

```bash
$ scavengecli tx scavenge reveal-solution --help
Reveals a solution for scavenge

Usage:
  scavengecli tx scavenge reveal-solution [solution] [flags]

Flags:
  -a, --account-number uint      The account number of the signing account (offline mode only)
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async|block) (default "sync")
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate required gas automatically (default 200000) (default "200000")
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices to determine the transaction fee (e.g. 10uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible and the node operates offline)
  -h, --help                     help for reveal-solution
      --indent                   Add indent to JSON response
      --keyring-backend string   Select keyring's backend (os|file|test) (default "os")
      --ledger                   Use a connected Ledger device
      --memo string              Memo to send along with transaction
      --node string              <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --trust-node               Trust connected full node (don't verify proofs for responses) (default true)
  -y, --yes                      Skip tx broadcasting prompt confirmation

Global Flags:
      --chain-id string   Chain ID of tendermint node
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/Users/b/.scavengecli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

Since all we need is the solution again let's send and confirm our final message:
```bash
scavengecli tx scavenge reveal-solution "A stick" --from user2
```
We can gather the `txhash` and query it again using `scavengecli q tx <txhash>` to reveal:
```json
{
  "height": "163",
  "txhash": "8EBE3492115EBF3B5BCBA9B19EA67BCF499E02BBEC6116613E7703FE080F368A",
  "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"RevealSolution\"},{\"key\":\"sender\",\"value\":\"cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg\"},{\"key\":\"module\",\"value\":\"scavenge\"},{\"key\":\"action\",\"value\":\"SolveScavenge\"},{\"key\":\"sender\",\"value\":\"cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq\"},{\"key\":\"solutionHash\",\"value\":\"2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906\"},{\"key\":\"description\",\"value\":\"What's brown and sticky?\"},{\"key\":\"solution\",\"value\":\"A stick\"},{\"key\":\"scavenger\",\"value\":\"cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq\"},{\"key\":\"reward\",\"value\":\"69token\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq\"},{\"key\":\"sender\",\"value\":\"cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg\"},{\"key\":\"amount\",\"value\":\"69token\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "RevealSolution"
            },
            {
              "key": "sender",
              "value": "cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg"
            },
            {
              "key": "module",
              "value": "scavenge"
            },
            {
              "key": "action",
              "value": "SolveScavenge"
            },
            {
              "key": "sender",
              "value": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq"
            },
            {
              "key": "solutionHash",
              "value": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906"
            },
            {
              "key": "description",
              "value": "What's brown and sticky?"
            },
            {
              "key": "solution",
              "value": "A stick"
            },
            {
              "key": "scavenger",
              "value": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq"
            },
            {
              "key": "reward",
              "value": "69token"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq"
            },
            {
              "key": "sender",
              "value": "cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg"
            },
            {
              "key": "amount",
              "value": "69token"
            }
          ]
        }
      ]
    }
  ],
  "gas_wanted": "200000",
  "gas_used": "56048",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "scavenge/RevealSolution",
          "value": {
            "scavenger": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq",
            "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
            "solution": "A stick"
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
            "value": "AreVRxbXB/ZS0D43BXLnITFCSf0ZvbRwP+7mCR5U7srD"
          },
          "signature": "b+gP9WULpUm325aqQTfr9VgZT8uYLkJ4LKEz3DMMjPNrJ5ncyHHRv2E/dtkS5z3BBN38TR371Sk/WZkybR48Rg=="
        }
      ],
      "memo": ""
    }
  },
  "timestamp": "2020-09-22T01:31:08Z"
}
```
You'll notice that the final event that was submitted was a transfer. This shows the movement of the reward into the account of the user `user1`. To confirm `user2` now has `69token` more you can query their account balance as follows:
```bash
scavengecli q account $(scavengecli keys show user2 -a)
```
This should show a healthy account balance of `569token` since `user2` began with `500token`:
```json
{
  "type": "cosmos-sdk/Account",
  "value": {
    "address": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq",
    "coins": [
      {
        "denom": "token",
        "amount": "569"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "AreVRxbXB/ZS0D43BXLnITFCSf0ZvbRwP+7mCR5U7srD"
    },
    "account_number": "3",
    "sequence": "2"
  }
}
```
If you'd like to take a look at the completed scavenge you can first query all scavenges with:
```bash
scavengecli q scavenge list 
```
To see the specific one just use 
```bash
scavengecli q scavenge get-scavenge 2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906
```
Which should show you that the scavenge has in fact been completed:
```json
{
  "creator": "cosmos13t3zdk8shwf68cqtdrqr380ucqfr7q9l7vd809",
  "description": "What's brown and sticky?",
  "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
  "reward": [
    {
      "denom": "token",
      "amount": "69"
    }
  ],
  "solution": "A stick",
  "scavenger": "cosmos1wr67fr5ad0vqe83mvzs3nkd2x4phacryk8zgzq"
}
```

<img src="./img/carmen.jpg" style="margin:auto;display:block;">

---

**Thanks for joining me** in building a deterministic state machine and using it as a game. I hope you can see that even such a simple app can be extremely powerful as it contains digital scarcity. 

If you'd like to keep going, consider trying to expand on the capabilities of this application by doing one of the following:
* Allow the `Creator` of a `Scavenge` to edit or delete a scavenge.
* Create a query that lists all commits.

If you're interested in learning more about the Cosmos SDK check out the rest of our [docs](https://docs.cosmos.network) or join our [forum](https://forum.cosmos.network).

Topics to look out for in future tutorials are:
* [Communication between applications (IBC)](https://cosmos.network/ibc/)
* [Digital Collectibles (NFTs)](https://github.com/cosmos/modules)
* [Using the Ethereum Virtual Machine (EVM) as a module within an application](https://github.com/chainsafe/ethermint)

If you have any questions or comments feel free to open an issue on this tutorial's [github](https://github.com/cosmos/sdk-tutorials).


If you'd like to stay in touch with me follow my github at [@okwme](https://github.com/okwme) or twitter at [@billyrennekamp](https://twitter.com/billyrennekamp).