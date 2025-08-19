---
sidebar_position: 2
description: Configuration File Example.
title: Configuration File Example
---

## Configuration File Example

```yaml title="config.yml"
include: (string list) # Include incorporate a separate config.yml file directly in your current config file.
validation: (string) # Specifies the type of validation the blockchain uses (e.g., sovereign).
version: (uint) # Defines the configuration version number.
build: # Contains build configuration options.
  main: (string) # Path to the main build file.
  binary: (string) # Path to the binary file.
  ldflags: (string list) # List of custom linker flags for building the binary.
  proto: # Contains proto build configuration options.
    path: (string) # Relative path where the application&#39;s proto files are located.
accounts: (list) # Lists the options for setting up Cosmos Accounts.
  name: (string) # Local name associated with the Account&#39;s key pair.
  coins: (string list) # List of token balances for the account.
  mnemonic: (string) # Mnemonic phrase for the account.
  address: (string) # Address of the account.
  cointype: (string) # Coin type number for HD derivation (default is 118).
  account_number: (string) # Account number for HD derivation (must be ≤ 2147483647).
  address_index: (string) # Address index number for HD derivation (must be ≤ 2147483647).
faucet: # Configuration for the faucet.
  name: (string) # Name of the faucet account.
  coins: (string list) # Types and amounts of coins the faucet distributes.
  coins_max: (string list) # Maximum amounts of coins that can be transferred to a single user.
  rate_limit_window: (string) # Timeframe after which the limit will be refreshed.
  host: (string) # Host address of the faucet server.
  port: (uint) # Port number for the faucet server.
  tx_fee: (string) # Tx fee the faucet needs to pay for each transaction.
client: # Configures client code generation.
  typescript: # Relative path where the application&#39;s Typescript files are located.
    path: (string) # Relative path where the application&#39;s Typescript files are located.
  composables: # Configures Vue 3 composables code generation.
    path: (string) # Relative path where the application&#39;s composable files are located.
  openapi: # Configures OpenAPI spec generation for the API.
    path: (string) # Relative path where the application&#39;s OpenAPI files are located.
genesis: (key/value) # Custom genesis block modifications. Follow the nesting of the genesis file here to access all the parameters.
default_denom: (string) # Default staking denom (default is stake).
validators: (list) # Contains information related to the list of validators and settings.
  name: (string) # Name of the validator.
  bonded: (string) # Amount staked by the validator.
  app: (key/value) # Overwrites the appd&#39;s config/app.toml configurations.
  config: (key/value) # Overwrites the appd&#39;s config/config.toml configurations.
  client: (key/value) # Overwrites the appd&#39;s config/client.toml configurations.
  home: (string) # Overwrites the default home directory used for the application.
  gentx: # Overwrites the appd&#39;s config/gentx.toml configurations.
    amount: (string) # Amount for the current Gentx.
    moniker: (string) # Optional moniker for the validator.
    keyring-backend: (string) # Backend for the keyring.
    chain-id: (string) # Network chain ID.
    commission-max-change-rate: (string) # Maximum commission change rate percentage per day.
    commission-max-rate: (string) # Maximum commission rate percentage (e.g., 0.01 = 1%).
    commission-rate: (string) # Initial commission rate percentage (e.g., 0.01 = 1%).
    details: (string) # Optional details about the validator.
    security-contact: (string) # Optional security contact email for the validator.
    website: (string) # Optional website for the validator.
    account-number: (int) # Account number of the signing account (offline mode only).
    broadcast-mode: (string) # Transaction broadcasting mode (sync|async|block) (default is &#39;sync&#39;).
    dry-run: (bool) # Simulates the transaction without actually performing it, ignoring the --gas flag.
    fee-account: (string) # Account that pays the transaction fees instead of the signer.
    fee: (string) # Fee to pay with the transaction (e.g.: 10uatom).
    from: (string) # Name or address of the private key used to sign the transaction.
    gas: (string) # Gas limit per transaction; set to &#39;auto&#39; to calculate sufficient gas automatically (default is 200000).
    gas-adjustment: (string) # Factor to multiply against the estimated gas (default is 1).
    gas-prices: (string) # Gas prices in decimal format to determine the transaction fee (e.g., 0.1uatom).
    generate-only: (bool) # Creates an unsigned transaction and writes it to STDOUT.
    identity: (string) # Identity signature (e.g., UPort or Keybase).
    ip: (string) # Node&#39;s public IP address (default is &#39;192.168.1.64&#39;).
    keyring-dir: (string) # Directory for the client keyring; defaults to the &#39;home&#39; directory if omitted.
    ledger: (bool) # Uses a connected Ledger device if true.
    min-self-delegation: (string) # Minimum self-delegation required for the validator.
    node: (string) # &lt;host&gt;:&lt;port&gt; for the Tendermint RPC interface (default &#39;tcp://localhost:26657&#39;)
    node-id: (string) # Node&#39;s NodeID
    note: (string) # Adds a description to the transaction (formerly --memo).
    offline: (bool) # Operates in offline mode, disallowing any online functionality.
    output: (string) # Output format (text|json) (default &#39;json&#39;).
    output-document: (string) # Writes the genesis transaction JSON document to the specified file instead of the default location.
    pubkey: (string) # Protobuf JSON encoded public key of the validator.
    sequence: (uint) # Sequence number of the signing account (offline mode only).
    sign-mode: (string) # Chooses sign mode (direct|amino-json), an advanced feature.
    timeout-height: (uint) # Sets a block timeout height to prevent the transaction from being committed past a certain height.
```