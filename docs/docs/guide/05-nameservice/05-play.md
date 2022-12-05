---
sidebar_position: 5
description: Start and experiment with your nameservice blockchain and nameservice module.
---

# Play With Your Blockchain

If you haven't already, start a blockchain node in development:

```bash
ignite chain serve -r
```

The optional `-r` flag is useful in development mode since it resets the
blockchain's state if it has been started before.

After the `serve` command has finished building the blockchain, a `nameserviced`
binary is installed by default in the `~/go/bin` directory.

The terminal window where the chain is started must remain open, so open a
second terminal window to use `nameserviced` to run commands at the command
line.

## Buy a New Name

Purchase a new name using the `buy-name` command. The name is `foo` and the bid
is `20token`.

In the keeper for the buy name transaction, you added code to the
`msg_server_buy_name.go` file that hard-coded the minimum bid to `10token`. Any
bid lower than that amount results in a rejected purchase.

```bash
nameserviced tx nameservice buy-name foo 20token --from alice
```

where:

- `buy-name` is the command that accepts two arguments
- `foo` is the name
- `20token` is the bid
- the `--from alice` flag specifies the user account that signs and broadcasts
  the transaction

This `buy-name` command creates a transaction and prompts the user `alice` to
sign and broadcast the transaction.

Here is what an unsigned transaction looks like:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/nameservice.nameservice.MsgBuyName",
        "creator": "cosmos1vh7akm79jdh87ytkgk2qld5t4m3nwegt5ge248",
        "name": "foo",
        "bid": "20token"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    },
    "tip": null
  },
  "signatures": []
}
```

### Buy Name Transaction Details

Look at the transaction details:

- The transaction contains only one message: `MsgBuyName`.
- The message `@type` matches the package name of the corresponding proto file,
  `proto/nameservice/nameservice/tx.proto`.
- The `creator` field is populated automatically with the address of the account
  broadcasting the transaction.
- The local account `alice` address is `cosmos1vh7...e248`.
- Values of `name` and `bid` are passed as CLI arguments.

After the transaction is broadcast and included in a block, the blockchain
returns a response where `"code": 0` means the transaction was successfully
processed.

```json
{
  "height": "160",
  "txhash": "38401948E1535F068CED99D640B10F0F08AB57F6B7781ABFEC128B073ADC89C7",
  "codespace": "",
  "code": 0,
  "data": "122D0A2B2F6E616D65736572766963652E6E616D65736572766963652E4D73674275794E616D65526573706F6E7365",
  "raw_log": "",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "78793",
  "tx": null,
  "timestamp": "",
  "events": []
}
```

Contents of `raw_log`, `logs`, and `events` are omitted.

## Query the Chain for a List of Names

Query the chain for a list of name and corresponding values. Query commands
don't need the `--from` flag, because they don't broadcast transactions and make
only free requests.

```bash
nameserviced q nameservice list-whois
```

The response confirms that the name `foo` was successfully purchased by `alice`
and the current `price` is set to`20token`.

```yaml
pagination:
  next_key: null
  total: "0"
whois:
- index: foo
  name: foo
  owner: cosmos1vh7akm79jdh87ytkgk2qld5t4m3nwegt5ge248
  price: 20token
  value: ""
```

## Set a Value to the Name

Now that `alice` is an owner of the name, she can set the value to anything she
wants. Use the `set-name` command to set the value to `bar`:

```bash
nameserviced tx nameservice set-name foo bar --from alice
```

Query for a list of names again:

```bash
nameserviced q nameservice list-whois 
```

The response shows that `name` is now `foo`.

```yaml
pagination:
  next_key: null
  total: "0"
whois:
- index: foo
  name: foo
  owner: cosmos1vh7akm79jdh87ytkgk2qld5t4m3nwegt5ge248
  price: 20token
  value: bar
```

## Buy an Existing Name

Use the `bob` account to purchase an existing name from `alice`. A successful
bid requires that the buy price is higher than the current value of `20token`.

```bash
nameserviced tx nameservice buy-name foo 40token --from bob
```

In this `buy-name` command, the bid is updated to the latest bid of `40token`
and the `--from bob` flag specifies that the transaction is signed by the `bob`
address.

Query for a list of names again:

```bash
nameserviced q nameservice list-whois 
```

The response shows a different creator address than `alice` (it's now the
address for `bob`) and the `price` is now`40token`.

```yaml
pagination:
  next_key: null
  total: "0"
whois:
- index: foo
  name: foo
  owner: cosmos1hgtyllz4cdrg2umnngg8f9t3tvjaccwc8e9736
  price: 40token
  value: bar
```

## Query the Bank Balance

Use the following command to see how the `alice` bank balance has changed after
this transaction:

```bash
nameserviced q bank balances $(nameserviced keys show alice -a)
```

## Test an Unauthorized Transaction

Try updating the value by broadcasting a transaction from the `alice` account:

```bash
nameserviced tx nameservice set-name foo qoo --from alice
```

An error occurs because `alice` sold the name in a previous transaction. The
results show that `alice` is not the owner of the name and is therefore not
authorized to change the value.

```json
{
  "height": "1167",
  "txhash": "9D5D28A9C61DCFF4FCC833EFB35A9A5C45CC561116448671269A0EA058712663",
  "codespace": "sdk",
  "code": 4,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: Incorrect Owner: unauthorized",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "44669",
  "tx": null,
  "timestamp": "",
  "events": []
}
```

```bash
nameserviced q nameservice list-whois 
```

```yaml
pagination:
  next_key: null
  total: "0"
whois:
- index: foo
  name: foo
  owner: cosmos1hgtyllz4cdrg2umnngg8f9t3tvjaccwc8e9736
  price: 40token
  value: bar
```

## Conclusion

Congratulations 🎉. You have created the nameservice module and the nameservice
application.

You successfully completed these steps:

- Learned how to work with module dependencies
- Use several scaffolding methods
- Learned about Cosmos SDK types and functions
- Used the CLI to broadcast transactions, and so much more

You are now prepared to continue your journey to learn about escrow accounts and
IBC.
