# CW20 - Fungible Tokens with CosmWasm

The `wasm` adaptation of the ERC20 contract is called the CW20, it encapsulates the same logic as the ERC20|ERC777 contract and allows for token creation on the blockchain with the `wasm` module. 

## Upload the contract to a starport chain

The `wasm` module uses Rust language instead of Solidity, therefore it is necessary to have Rust installed on the machine you are planning to compile and upload the smart contract. In the following steps, we will download the CW20 code, test it, deploy it for production, host it on a blockchain and use it.
Requirements are:

- [Starport](https://github.com/tendermint/starport)
- [Install Rust](https://www.rust-lang.org/)
- Install Cargo generate
`cargo install cargo-generate --features vendored-openssl`
- [Docker](https://www.docker.com/)

If you have followed along the commands in the last chapters, you should have a blockchain with `wasm` enabled in your `myapp` folder. Let's revisit the commands necessary:

```bash
starport app github.com/username/myapp
cd myapp
```

```bash
starport module import wasm
```

You will be greeted by a success message like `🎉 Imported module 'wasm'.`

Our blockchain application has been setup. Now we can run the blockchain daemon.

```bash
starport serve
```

Now we will upload a contract to our chain and interact with it. We will be using `cargo` to clone the CW20 from the Wasm examples github repository.

```bash
# Check Rust version, should be minimum 1.45.2+
rustc --version
# Check if cargo has been installed correctly
cargo --version
# Check if wasm has been added to the rustup list
rustup target list --installed
# If wasm32 is not listed above, run 
rustup target add wasm32-unknown-unknown
# This will get the contract code from cosm wasm examples directory into the name folder
git clone https://github.com/CosmWasm/cosmwasm-plus
# Enter the just created folder
cd cosmwasm-plus

# This will produce a wasm build in of the example contracts in cosmwasm-plus
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/workspace-optimizer:0.10.4
```

Now we have our contract wasm file in our target folder. In order to upload it to the blockchain, we will copy the contact to our current folder.

```bash
# Move to the root of the myapp directory again
cd ..
# Copies the contract to our current location
cp ./cosmwasm-plus/artifacts/cw20_base.wasm .
```

We now have our `cw20_base.wasm` file in our current folder with the bytecode accordingly. We can upload the contract with the following command to our blockchain `myapp`

```bash
myappcli tx wasm store cw20_base.wasm --from user1 --gas 1500000
```

You will get asked to sign the transaction, insert `y` and press Enter to confirm.

:tada: congratulations, the CW20 has been uploaded to the blockchain.
We can now check if it is listed with

```bash
myappcli query wasm list-code
```

The output will look something similar to

```
> myappcli query wasm list-code
[
  {
    "id": 1,
    "creator": "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf",
    "data_hash": "9238275A43DFE0DE390F2431F825D33AD010D07B1ECD0B9505C838F37B08B07D",
    "source": "",
    "builder": ""
  }
]
```

## Create your token with the CW20

On the output of the `list-code` command you can check the creator, id and hash of your contract. The id of the contract is what we are going to use in order to `instantiate` our contract. 

Having uploaded the code it is now available to instantiate and create our first token with the CW20 smart contract. 
For demonstration purposes, let's name our CW20 Token XRP and print an arbitrary amount for testing purposes.

```bash
myappcli tx wasm instantiate 1 '{ "name": "xrp", "symbol": "XRP", "decimals": 8, "initial_balances": [ { "address": "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf", "amount": 10000000000 } ]}' --from user1 --label xrp --gas 1000000
```

When everything is successfull, you will be greeted with a message to confirm the transaction y|N - input y and enter.
The output then looks similar to this

```
confirm transaction before signing and broadcasting [y/N]: y
{
  "height": "0",
  "txhash": "D58E53C4E97A9A687E7E818B68FFB3E810B1661FCAFF72594D7CD3BAB3DB4760",
  "raw_log": "[]"
}
```

You can query for transaction details with the following command, using the txhash displayed above:

```bash
myappcli query tx D58E53C4E97A9A687E7E818B68FFB3E810B1661FCAFF72594D7CD3BAB3DB4760
```

Checking transactions this way can be very convinient for debugging or researching more information.
Congratulations, you have created your first CW20 Token.

As execerise you should try out to send tokens from one account to another, from here on there are also CW20 standards that allow for staking, escrow or more features.

Don't forget to checkout the documentation at https://github.com/CosmWasm/cosmwasm-plus/tree/master/packages/cw20
