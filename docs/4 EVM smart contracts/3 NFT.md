# Create NFTs on the EVM with the Cosmos SDK (ERC721)

Non-Fungible Token or in short NFT are tokens that are non-divisible and live on the blockchain. Those tokens can represent a variety of things, from pieces of artwork, to houses, to ownership of cars or digital items in games. The usage is from individuals to companies or institutions that are leverage the power of NFTs. 

NFT are defined by the contract proposal ERC721. The contract proposal and definition can be read on https://eips.ethereum.org/EIPS/eip-721 

The NFT contract will implement the following functions

```
totalSupply() - Total amount of emitted tokens.
balanceOf( _owner ) - Amount of tokens in specific _owner’s wallet.
ownerOf( _tokenId ) - Returns wallet address of the specific tokens owner.
transfer( _to, _tokenId ) - Transfers token with _tokenId from senders wallet to specific wallet.
takeOwnership( _tokenId ) - Claims the ownership of a given token ID
approve( _to, _tokenId ) - Approves another address to claim for the ownership of the given token ID
```

As we have done in the [EVM ERC20 & ERC777 Tutorial](../../04_usecases/02_erc20/02_erc20.md), we will be building upon the implementation from OpenZeppelin. They provide us with two additional functions that make working with the contract easier

```
tokensOf( _owner ) - Returns list of token ID’s of specific _owner
approvedFor( _tokenId ) - Returns the approved address to take ownership of a given token ID
```

## ERC20 on the Cosmos SDK with EVM

The `EVM` is a module that allows to replicate the Ethereum Virtual Machine on any Cosmos SDK application. Ethermint is a reference project from Tendermint and partners (ChainSafe) that implements the EVM module. We can load solidity contracts onto the Cosmos SDK blockchain with the EVM module. Let us examine an example usage of Ethermint with the ERC721 Token.

Working with Ethermint can come in different forms. You can either work with the official [Ethermint chain](https://ethermint.zone/), fork from the [Ethermint Codebase](https://github.com/ChainSafe/ethermint), or use the EVM Module in your own Cosmos SDK application. The last option is what we will be looking into. We have created our Cosmos SDK template with Starport and we can just add the Ethereum Virtual Machine to our application and allow smart contracts to be created in Solity and uploaded to the blockchain.

[Read how to integrate the `evm` module into your Cosmos SDK application.](05_extras/05_01_cosmos_entermint/05_01_cosmos_ethermint.md)

On GitHub, you can also find the `ethapp` application which you can fork and get started right away with an example `evm` implementation following the steps as outlined in the tutorial but with RPC integration, Makefile, and initialisation script of the blockchain.

```bash
git clone https://github.com/Tosch110/ethapp && cd ethapp
make install
./init.sh
ethappd start
```

This application is very close to the starport scaffold and has an empty module in the `ethapp/x` folder where you can start integrating your blockchain logic as in any starport scaffolded application.

## Truffle

For working with our contract we will be using the Truffle Suite. We will be mostly working with Solidity, JavaScript and your command line interface. To get started, make sure to have NPM and Truffle installed.

```bash
npm install -g truffle
```

Let's create a new folder to work with our contract

```bash
mkdir nft && cd nft
truffle init
npm init
```

This will bootstrap our smart contract working environment. Let us configure the network we are working on.
`truffle-config.js`

```javascript
networks: {
    development: {
      host: "127.0.0.1",
      port: 8545,
      network_id: "*" // Match any network id
    }
  }
```

For the compiler, we will use the following settings in this tutorial in the same file

```javascript
  // Configure your compilers
  compilers: {
    solc: {
      version: "0.4.24",    // Fetch exact version from solc-bin (default: truffle's version)
      // docker: true,        // Use "0.5.1" you've installed locally with docker (default: false)
      // settings: {          // See the solidity docs for advice about optimization and evmVersion
      //  optimizer: {
      //    enabled: false,
      //    runs: 200
      //  },
      //  evmVersion: "byzantium"
      // }
    },
```

This will be the connection our rest-server uses. In order to setup the rest server on your Cosmos SDK blockchain, the following command will get it started on the right port:

`ethappcli rest-server --laddr "tcp://localhost:8545" --trace --chain-id ethapp-1 --unlock-key me`

## Smart Contract Code

From here our setup is complete, we can now jump into the smart contract code. Since we are working with OpenZeppelin reference implementation, we need to install their components.

```bash
npm i zeppelin-solidity --save-dev
```

In our `contracts` directory we add the solidity code for the NFT ERC721 Token that we would like to address, create the `contracts/MyToken.sol` file and add

```solidity
pragma solidity ^0.4.24;

import "../node_modules/zeppelin-solidity/contracts/token/ERC721/ERC721Token.sol";

contract MyNFT is ERC721Token {
    constructor (string _name, string _symbol) public
        ERC721Token(_name, _symbol)
    {
    }

    function mintUniqueTokenTo(
        address _to,
        uint256 _tokenId,
        string  _tokenURI
    ) public
    {
        super._mint(_to, _tokenId);
        super._setTokenURI(_tokenId, _tokenURI);
    }
}
```

in the `migrations` directory we add the file for the deployment, `2_deploy_contract.js`.

```javascript
const MyNFT = artifacts.require("./MyNFT.sol");

module.exports = async function(deployer) {
  await deployer.deploy(MyNFT, "MyNFT", "MyNFT")
  const erc721 = await MyNFT.deployed()
};
```

From this, we launch and deploy our contract. When the blockchain is running and our rest-server is setup, we can now migrate using truffle. 
In the root of our nft directory, we run to migrate truffle

```bash
truffle migrate --network development
```

When everything has been setup correctly, the output should be similar to

```bash
Starting migrations...
======================
> Network name:    'development'
> Network id:      1
> Block gas limit: 4294967295 (0xffffffff)


1_initial_migration.js
======================

   Replacing 'Migrations'
   ----------------------
   > transaction hash:    0xcd69d3f372ef8eba1b907e6bd147fe304b8596da09a5401a7ddaca8d61954983
   > Blocks: 0            Seconds: 0
   > contract address:    0x26299431295c347b462bEa52e5798D35B412baB6
   > block number:        144
   > block timestamp:     1601544063
   > account:             0x365247AA0fDC939007c84D614d12059046BdD929
   > balance:             2.8278025
   > gas used:            143403 (0x2302b)
   > gas price:           20 gwei
   > value sent:          0 ETH
   > total cost:          0.00286806 ETH


   > Saving migration to chain.
   > Saving artifacts
   -------------------------------------
   > Total cost:          0.00286806 ETH


2_deploy_contract.js
====================

   Deploying 'MyNFT'
   -----------------
   > transaction hash:    0xec5c88568f63cdd1a2d49f04eead516a87fd495b29b161b05f8f9019a511274c
   > Blocks: 0            Seconds: 4
   > contract address:    0xe498F3DF4343Dc6ab666C63B531E0E4727ba4783
   > block number:        146
   > block timestamp:     1601544073
   > account:             0x365247AA0fDC939007c84D614d12059046BdD929
   > balance:             2.5589235
   > gas used:            1683560 (0x19b068)
   > gas price:           20 gwei
   > value sent:          0 ETH
   > total cost:          0.0336712 ETH


   > Saving migration to chain.
   > Saving artifacts
   -------------------------------------
   > Total cost:           0.0336712 ETH


Summary
=======
> Total deployments:   2
> Final cost:          0.03653926 ETH
```

