---
order: 23
---

# Run REST routes

Now that you tested your CLI queries and transactions, time to test same things in the REST server. Leave the `nameserviced` that you had running earlier and start by gathering your addresses:

```bash
$ nameservicecli keys show jack --address
$ nameservicecli keys show alice --address
```

Now its time to start the `rest-server` in another terminal window:

```bash
$ nameservicecli rest-server --chain-id namechain --trust-node
```

Then you can construct and run the following queries:

> NOTE: Be sure to substitute your password and buyer/owner addresses for the ones listed below!

```bash
# Get the sequence and account numbers for jack to construct the below requests
curl -s http://localhost:1317/auth/accounts/$(nameservicecli keys show jack -a)
# > {"type":"auth/Account","value":{"address":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6","coins":[{"denom":"jackCoin","amount":"1000"},{"denom":"nametoken","amount":"1010"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"A9YxyEbSWzLr+IdK/PuMUYmYToKYQ3P/pM8SI1Bxx3wu"},"account_number":"0","sequence":"1"}}

# Get the sequence and account numbers for alice to construct the below requests
curl -s http://localhost:1317/auth/accounts/$(nameservicecli keys show alice -a)
# > {"type":"auth/Account","value":{"address":"cosmos1h7ztnf2zkf4558hdxv5kpemdrg3tf94hnpvgsl","coins":[{"denom":"aliceCoin","amount":"1000"},{"denom":"nametoken","amount":"980"}],"public_key":{"type":"tendermint/PubKeySecp256k1","value":"Avc7qwecLHz5qb1EKDuSTLJfVOjBQezk0KSPDNybLONJ"},"account_number":"1","sequence":"2"}}

# Buy another name for jack, first create the raw transaction
# NOTE: Be sure to specialize this request for your specific environment, also the "buyer" and "from" should be the same address
curl -X POST -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show jack -a)'","chain_id":"namechain"},"name":"jack1.id","amount":"5nametoken","buyer":"'$(nameservicecli keys show jack -a)'"}' > unsignedTx.json

# Then sign this transaction
# NOTE: In a real environment the raw transaction should be signed on the client side. Also the sequence needs to be adjusted, depending on what the query of alice's account has shown.
nameservicecli tx sign unsignedTx.json --from jack --offline --chain-id namechain --sequence 1 --account-number 0 > signedTx.json

# And finally broadcast the signed transaction
nameservicecli tx broadcast signedTx.json
# > { "height": "266", "txhash": "C041AF0CE32FBAE5A4DD6545E4B1F2CB786879F75E2D62C79D690DAE163470BC", "logs": [  {   "msg_index": "0",   "success": true,   "log": ""  } ],"gas_wanted":"200000", "gas_used": "41510", "tags": [  {   "key": "action",   "value": "buy_name"  } ]}

# Set the data for that name that jack just bought
# NOTE: Be sure to specialize this request for your specific environment, also the "owner" and "from" should be the same address
curl -X PUT -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show jack -a)'","chain_id":"namechain"},"name":"jack1.id","value":"8.8.4.4","owner":"'$(nameservicecli keys show jack -a)'"}' > unsignedTx.json
# > {"check_tx":{"gasWanted":"200000","gasUsed":"1242"},"deliver_tx":{"log":"Msg 0: ","gasWanted":"200000","gasUsed":"1352","tags":[{"key":"YWN0aW9u","value":"c2V0X25hbWU="}]},"hash":"B4DF0105D57380D60524664A2E818428321A0DCA1B6B2F091FB3BEC54D68FAD7","height":"26"}

# Again we need to sign and broadcast
nameservicecli tx sign unsignedTx.json --from jack --offline --chain-id namechain --sequence 2 --account-number 0 > signedTx.json
nameservicecli tx broadcast signedTx.json

# Query the value for the name jack just set
$ curl -s http://localhost:1317/nameservice/whois/jack1.id/resolve
# 8.8.4.4

# Query whois for the name jack just bought
$ curl -s http://localhost:1317/nameservice/whois/jack1.id
# > {"value":"8.8.8.8","owner":"cosmos127qa40nmq56hu27ae263zvfk3ey0tkapwk0gq6","price":[{"denom":"STAKE","amount":"10"}]}

# Alice buys name from jack
$ curl -X POST -s http://localhost:1317/nameservice/whois --data-binary '{"base_req":{"from":"'$(nameservicecli keys show alice -a)'","chain_id":"namechain"},"name":"jack1.id","amount":"10nametoken","buyer":"'$(nameservicecli keys show alice -a)'"}' > unsignedTx.json

# Again we need to sign and broadcast
# NOTE: The account number has changed to 1 and the sequence is now 2, according to the query of alice's account
nameservicecli tx sign unsignedTx.json --from alice --offline --chain-id namechain --sequence 2 --account-number 1 > signedTx.json
nameservicecli tx broadcast signedTx.json
# > { "height": "1515", "txhash": "C9DCC423E10E7E5E40A549057A4AA060DA6D6A885A394F6ED5C0E40AEE984A77", "logs": [  {   "msg_index": "0",   "success": true,   "log": ""  } ],"gas_wanted": "200000", "gas_used": "42375", "tags": [  {   "key": "action",   "value": "buy_name"  } ]}

# Now, Alice no longer needs the name she bought from jack and hence deletes it
# NOTE: Only the owner can delete the name. Since she is one, she can delete the name she bought from jack
$ curl -XDELETE -s http://localhost:1317/nameservice/names --data-binary '{"base_req":{"from":"'$(nameservicecli keys show alice -a)'","chain_id":"namechain"},"name":"jack1.id","owner":"'$(nameservicecli keys show alice -a)'"}' > unsignedTx.json

# And a final time sign and broadcast
# NOTE: The account number is still 1, but the sequence is changed to 3, according to the query of alice's account
nameservicecli tx sign unsignedTx.json --from alice --offline --chain-id namechain --sequence 3 --account-number 1 > signedTx.json
nameservicecli tx broadcast signedTx.json

# Query whois for the name Alice just deleted
$ curl -s http://localhost:1317/nameservice/names/jack1.id/whois
# > {"value":"","owner":"","price":[{"denom":"STAKE","amount":"1"}]}
```

### Request Schemas:

#### `POST /nameservice/names` BuyName Request Body:
```json
{
  "base_req": {
    "name": "string",
    "chain_id": "string",
    "gas": "string,not_req",
    "gas_adjustment": "string,not_req",
  },
  "name": "string",
  "amount": "string",
  "buyer": "string"
}
```

#### `PUT /nameservice/names` SetName Request Body:
```json
{
  "base_req": {
    "name": "string",
    "chain_id": "string",
    "gas": "string,not_req",
    "gas_adjustment": "strin,not_reqg"
  },
  "name": "string",
  "value": "string",
  "owner": "string"
}
```

#### `DELETE /nameservice/names` DeleteName Request Body:
```json
{
  "base_req": {
    "name": "string",
    "chain_id": "string",
    "gas": "string,not_req",
    "gas_adjustment": "strin,not_reqg"
  },
  "name": "string",
  "owner": "string"
}
```

### [Back to start of tutorial](./README.md)
