---
order: 5
---

# Play

## Buying a New Name

```
nameserviced tx nameservice buy-name foo 20token --from alice
```

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

```
nameserviced q nameservice list-whois
```

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

```
nameserviced tx nameservice set-name foo bar --from alice
```

```
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

```
nameserviced tx nameservice buy-name foo 40token --from bob
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

## Setting a Value from an Authorized Account

```
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