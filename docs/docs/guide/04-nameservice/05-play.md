---
sidebar_position: 5
description: Start and experiment with your nameservice blockchain and nameservice module.
---

# Play With Your Blockchain

If you haven't already, start a blockchain node in development:

```bash
ignite chain serve -r
```

The optional `-r` flag is useful in development mode since it resets the blockchain's state if it has been started before. 

After the `serve` command has finished building the blockchain, a `nameserviced` binary is installed by default in the `~/go/bin` directory.

The terminal window where the chain is started must remain open, so open a second terminal window to use `nameserviced` to run commands at the command line.

## Buy a New Name

Purchase a new name using the `buy-name` command. The name is `foo` and the bid is `20token`. 

In the keeper for the buy name transaction, you added code to the `msg_server_buy_name.go` file that hard-coded the minimum bid to `10token`. Any bid lower than that amount results in a rejected purchase. 

```bash
nameserviced tx nameservice buy-name foo 20token --from alice
```

where:

- buy-name is the command that accepts two arguments
- foo is the name 
- 20token is the bid 
- the `--from alice` flag specifies the user account that signs and broadcasts the transaction

This `buy-name` command creates a transaction and prompts the user `alice` to sign and broadcast the transaction.

Here is what an unsigned transaction looks like:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/username.nameservice.nameservice.MsgBuyName",
        "creator": "cosmos1p0fprxtpk497jvczexp96sy2w43hupeph9km5d",
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
    "fee": { "amount": [], "gas_limit": "200000", "payer": "", "granter": "" }
  },
  "signatures": []
}
```

### Buy Name Transaction Details

Look at the transaction details:

- The transaction contains only one message: `MsgBuyName`.
- The message `@type` matches the package name of the corresponding proto file, `proto/nameservice/tx.proto`. 
- The `creator` field is populated automatically with the address of the account broadcasting the transaction.
- The local account `alice` address is `cosmos1p0f...km5d`.
- Values of `name` and `bid` are passed as CLI arguments.

After the transaction is broadcast and included in a block, the blockchain returns a response where `"code": 0` means the transaction was successfully processed.

```json
{
  "height": "658",
  "txhash": "EDC1842BE4B596DDD9E2D34F2E372354F9BA5F6D2E4B3F1C2664F4FF05D433B7",
  "codespace": "",
  "code": 0,
  "data": "0A090A074275794E616D65",
  "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"BuyName\"}]}]}]",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [{ "key": "action", "value": "BuyName" }]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "47954",
  "tx": null,
  "timestamp": ""
}
```

## Query the Chain for a List of Names

Query the chain for a list of name and correponding values. Query commands don't need the `--from` flag, because they don't broadcast transactions and make only free requests.

```bash
nameserviced q nameservice list-whois
```

The response confirms that the name `foo` was successfully purchased by `alice` and the current `price` is set to `20token`.

```yaml
Whois:
- creator: cosmos1p0fprxtpk497jvczexp96sy2w43hupeph9km5d
  index: foo
  name: foo
  price: 20token
  value: ""
pagination:
  next_key: null
  total: "0"
```

## Set a Value to the Name

Now that `alice` is an owner of the name, she can set the value to anything she wants. Use the `set-name` command to set the value to `bar`:

```bash
nameserviced tx nameservice set-name foo bar --from alice
```

Query for a list of names again:

```bash
nameserviced q nameservice list-whois 
```

The response shows that `name` is now `foo`.

```yaml
Whois:
- creator: cosmos1p0fprxtpk497jvczexp96sy2w43hupeph9km5d
  index: foo
  name: foo
  price: 20token
  value: bar
pagination:
  next_key: null
  total: "0"
```

## Buy an Existing Name

Use the `bob` account to purchase an existing name from `alice`. A successful bid requires that the buy price is higher than the current value of `20token`. 

```bash
nameserviced tx nameservice buy-name foo 40token --from bob
```

In this `buy-name` command, the bid is updated to the latest bid of `40token` and the `--from bob` flag specifies that the transaction is signed by the `bob` address.

Query for a list of names again:

```bash
nameserviced q nameservice list-whois 
```

The response shows a different creator address than `alice` (it's now the address for `bob`) and the `price` is now `40token`.

```yaml
Whois:
- creator: cosmos1ku6sqpk9rgwgx98u2gs9c05aa9wrps969g0wy5
  index: foo
  name: foo
  price: 40token
  value: bar
pagination:
  next_key: null
  total: "0"
```

## Query the Bank Balance

Use the following command to see how the `alice` bank balance has changed after this transaction:

```bash
nameserviced q bank balances $(nameserviced keys show alice -a)
```

## Test an Unauthorized Transaction

Try updating the value by broadcasting a transaction from the `alice` account: 

```bash
nameserviced tx nameservice set-name foo qoo --from alice
```

An error occurs because `alice` sold the name in a previous transaction. The results show that `alice` is not the owner of the name and is therefore not authorized to change the value.

```json
{
  "height": "981",
  "txhash": "8E9951EDC5C9D76C2164BE9572B336B13CCF46653F45F54B2C1FEA702389FAE8",
  "codespace": "sdk",
  "code": 4,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: Incorrect Owner: unauthorized",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "39214",
  "tx": null,
  "timestamp": ""
}
```

```yaml
Whois:
- creator: cosmos1ku6sqpk9rgwgx98u2gs9c05aa9wrps969g0wy5
  index: foo
  name: foo
  price: 40token
  value: bar
pagination:
  next_key: null
  total: "0"
```

## Conclusion

Congratulations 🎉. You have created the nameservice module and the nameservice application.

You successfully completed these steps:

- Learned how to work with module dependencies
- Use several scaffolding methods
- Learned about Cosmos SDK types and functions
- Used the CLI to broadcast transactions , and so much more

You are now prepared to continue your journey to learn about escrow accounts and IBC.
