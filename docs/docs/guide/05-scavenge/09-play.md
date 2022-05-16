---
sidebar_position: 9
---

# Play with your blockchain

Now that you have made the required updates to you new chain, it's time to test it.

## Start the blockchain

To start your blockchain in development, run the following command:

```bash
ignite chain serve
```

The `serve` command builds the chain binary file, initializes a data directory, and starts a node in development. 

For all subsequent commands, you use a terminal window that is different from the window you started the chain in.

## Create a scavenge

Follow the instructions and submit a new scavenge. 

The first parameter is the `solution`, but you should also know what the actual question is that your solution solves (the `description`).

You can make the challenge question something family-friendly like `What's brown and sticky?`. Of course, the only solution to this question is `A stick`.

Next, you must specify the `reward`. Give away `100token` as a reward for solving the scavenge.

Now you have all the pieces needed to create your message. Piece them all together in the command and add the flag `--from` so the CLI knows who is sending it:

```bash
scavenged tx scavenge submit-scavenge "A stick" "What's brown and sticky?" 100token --from alice
```

The results:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/username.scavenge.scavenge.MsgSubmitScavenge",
        "creator": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh",
        "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
        "description": "What's brown and sticky?",
        "reward": "100token"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": { "amount": [], "gas_limit": "200000", "payer": "", "granter": "" }
  },
  "signatures": []
}
```

```json
{
  "height": "229",
  "txhash": "CE401E1F95FC583355BF6ABB823A4655185E2983CACE7C430E22CC7B573152DD",
  "codespace": "",
  "code": 0,
  "data": "0A100A0E43726561746553636176656E6765",
  "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"SubmitScavenge\"},{\"key\":\"sender\",\"value\":\"cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg\"},{\"key\":\"sender\",\"value\":\"cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh\"},{\"key\":\"amount\",\"value\":\"100token\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            { "key": "action", "value": "SubmitScavenge" },
            {
              "key": "sender",
              "value": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh"
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
              "value": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh"
            },
            { "key": "amount", "value": "100token" }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "65320",
  "tx": null,
  "timestamp": ""
}
```

Replace the `txhash` with your `txhash` from the previous output:

```bash
scavenged q tx CE401E1F95FC583355BF6ABB823A4655185E2983CACE7C430E22CC7B573152DD --output json
```

The results:

```json
{
  "height": "229",
  "txhash": "CE401E1F95FC583355BF6ABB823A4655185E2983CACE7C430E22CC7B573152DD",
  "codespace": "",
  "code": 0,
  "data": "0A100A0E43726561746553636176656E6765",
  "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"SubmitScavenge\"},{\"key\":\"sender\",\"value\":\"cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos13aupkh5020l9u6qquf7lvtcxhtr5jjama2kwyg\"},{\"key\":\"sender\",\"value\":\"cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh\"},{\"key\":\"amount\",\"value\":\"100token\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            { "key": "action", "value": "SubmitScavenge" },
            {
              "key": "sender",
              "value": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh"
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
              "value": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh"
            },
            { "key": "amount", "value": "100token" }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "65320",
  "tx": {
    "@type": "/cosmos.tx.v1beta1.Tx",
    "body": {
      "messages": [
        {
          "@type": "/username.scavenge.scavenge.MsgSubmitScavenge",
          "creator": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh",
          "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
          "description": "What's brown and sticky?",
          "reward": "100token"
        }
      ],
      "memo": "",
      "timeout_height": "0",
      "extension_options": [],
      "non_critical_extension_options": []
    },
    "auth_info": {
      "signer_infos": [
        {
          "public_key": {
            "@type": "/cosmos.crypto.secp256k1.PubKey",
            "key": "ApRuim5kLByq9AqJJ9dEF5rFCkAbhIehEcPzSouM92p6"
          },
          "mode_info": { "single": { "mode": "SIGN_MODE_DIRECT" } },
          "sequence": "1"
        }
      ],
      "fee": { "amount": [], "gas_limit": "200000", "payer": "", "granter": "" }
    },
    "signatures": [
      "8W5MkgV8oWpB6UWRGVKuimfPyb1OutG8KPXTIneM6WIvy4YHToG3GUXFpUrh+CxPXmlDh5gIfeR4+nFfUuQXng=="
    ]
  },
  "timestamp": "2021-07-09T10:24:52Z"
}
```

## Query for a list of scavenges

```bash
scavenged q scavenge list-scavenge --output json
```

The results:

```json
{
  "Scavenge": [
    {
      "creator": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh",
      "index": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
      "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
      "solution": "",
      "description": "What's brown and sticky?",
      "reward": "100token",
      "scavenger": ""
    }
  ],
  "pagination": { "next_key": null, "total": "0" }
}
```

## Commit a solution

Follow the instructions and submit the answer as a commit on behalf of `bob`:

```bash
scavenged tx scavenge commit-solution "A stick" --from bob
```

The results:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/username.scavenge.scavenge.MsgCommitSolution",
        "creator": "cosmos1gkheudhhjsvq0s8fxt7p6pwe0k3k30kepcnz9p",
        "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
        "solutionScavengerHash": "461d54ec0bbb1d696a79af80d7f63e4c6df262d76309423da37189453eaec127"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": { "amount": [], "gas_limit": "200000", "payer": "", "granter": "" }
  },
  "signatures": []
}
```

## Query for a list of commits

```bash
scavenged q scavenge list-commit --output json
```

The results:

```json
{
  "Commit": [
    {
      "creator": "cosmos1gkheudhhjsvq0s8fxt7p6pwe0k3k30kepcnz9p",
      "index": "461d54ec0bbb1d696a79af80d7f63e4c6df262d76309423da37189453eaec127",
      "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
      "solutionScavengerHash": "461d54ec0bbb1d696a79af80d7f63e4c6df262d76309423da37189453eaec127"
    }
  ],
  "pagination": { "next_key": null, "total": "0" }
}
```

You don't need to put the `solutionHash` because it can be generated by hashing the actual solution.

Since all you need is the solution again, you can send and confirm the final message:

## Reveal a solution

```bash
scavenged tx scavenge reveal-solution "A stick" --from bob
```

The results:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/username.scavenge.scavenge.MsgRevealSolution",
        "creator": "cosmos1gkheudhhjsvq0s8fxt7p6pwe0k3k30kepcnz9p",
        "solution": "A stick"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": { "amount": [], "gas_limit": "200000", "payer": "", "granter": "" }
  },
  "signatures": []
}
```

## Query for a list of solved scavenges

```bash
scavenged q scavenge list-scavenge --output json
```

The results:

```json
{
  "Scavenge": [
    {
      "creator": "cosmos1wzgkalxjhaqtznrzzp0xy5jgkxx82xaa660jxh",
      "index": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
      "solutionHash": "2f9457a6e8fb202f9e10389a143a383106268c460743dd59d723c0f82d9ba906",
      "solution": "A stick",
      "description": "What's brown and sticky?",
      "reward": "100token",
      "scavenger": "cosmos1gkheudhhjsvq0s8fxt7p6pwe0k3k30kepcnz9p"
    }
  ],
  "pagination": { "next_key": null, "total": "0" }
}
```

## Commit a solution again, expect to get an error

```bash
scavenged tx scavenge commit-solution "A stick" --from bob
```

The results:

```json
{
  "height": "665",
  "txhash": "EFA43A3C08BD1D77E597D57E60CD7B4D2E8E442F49BA88C85CC9EEC86E992B75",
  "codespace": "sdk",
  "code": 18,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: Commit with that hash already exists: invalid request",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "41086",
  "tx": null,
  "timestamp": ""
}
```

Congratulations, you have built and tested a scavenge hunt blockchain!
