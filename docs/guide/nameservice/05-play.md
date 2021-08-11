---
order: 5
---

# Play With Your Blockchain

If you haven't already, start a blockchain node in development:

```bash
starport chain serve -r
```

We're using the optional `-r` flag to reset the blockchain's state if it has been started before. Once the `serve` command has finished building the blockchain, you will get a `nameserviced` binary installed by default in `~/go/bin`.

Open a second terminal window and use `nameserviced` to issue commands.

## Buying a New Name

Purchase a new name using the `buy-name` command. The name is `foo` and the bid is `20token`. You've hard-coded the minimal bid to `10token`, so any bid below that will result in a rejected purchase. Use the `--from` flag to specify the account, from which the transaction will be sent.

```bash
nameserviced tx nameservice buy-name foo 20token --from alice
```

`buy-name` command accepts two arguments, creates a transaction and prompts the user to sign and broadcast the transaction. Here is how an unsigned transaction looks like:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/cosmonaut.nameservice.nameservice.MsgBuyName",
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

The transaction contains only one message: `MsgBuyName`. The message `@type` matches the package name of the corresponding proto file, `proto/nameservice/tx.proto`. `creator` field is populated automatically with the address of the account broadcasting the transaction (local account `alice` has an address `cosmos1p0f...km5d`). Values of `name` and `bid` are passed as CLI arguments.

After the transaction has been broadcasted and included into a block, the blockchain returns a response. Code `0` means the transaction has been processed successfully.

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

Query the chain for a list of name and correponding values. Query commands don't need the `--from` flag, because they don't broadcast transactions and only make side-effect free requests.

```bash
nameserviced q nameservice list-whois
```

You can confirm that the name `foo` has been successfully purchased by `alice` and the current `price` has been set to `20token`.

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

## Setting a Value to the Name

Now that `alice` is an owner of the name, she can set the value to anything she wants. Use the `set-name` command to set the value to `bar`.

```bash
nameserviced tx nameservice set-name foo bar --from alice
```

```bash
nameserviced q nameservice list-whois 
```

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

## Buying an Existing Name

Use `bob`'s account to purchase a name from `alice`. The bid has to be higher than `20token` for the transaction to be processed successfully.

```bash
nameserviced tx nameservice buy-name foo 40token --from bob
```

Notice that the `creator` address has been changed to `bob`'s address. The price has also been updated to the latest bid (`40token`).

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

Use the following command to see how `alice`'s bank balance has changed after this transaction:

```bash
nameserviced q bank balances $(nameserviced keys show alice -a)
```

## Setting a Value From an Authorized Account

Try updating the value by broadcasting a transaction from `alice`'s account and notice the error being returned, because `alice` is no longer the owner the name and isn't authorized to change the value.

```bash
nameserviced tx nameservice set-name foo qoo --from alice
```

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

Congratulations 🎉. You have successfully completed the nameservice application.
You have learned how to work with module dependencies, several scaffolding methods, Cosmos SDK types, functions and so much more.
Continue your journey to learn about escrow accounts and IBC.
