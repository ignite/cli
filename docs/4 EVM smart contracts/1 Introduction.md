# Smart Contracts

Smart Contracts allow for transactions to execute code and have a deterministic output. 

In blockchain, smart contracts do not always need to be "smart", the most important characteristic that these contracts have are to be __deterministic in outcome__. Every machine that replicates a smart contract and a specific input, need to be able to confirm the same output. In digital form, smart contracts are the brainchild of Nick Szabo (Szabo 1994) and go back to 1994. While the first blockchain emerged with Bitcoin in 2009. But as blockchain offers the first deterministic public ledger, in combination those two ideas form a powerful partnership. On blockchains smart contracts are a growing segment and probably will continue to do so for a long time.

Generally speaking, Smart Contracts are pieces of code that get uploaded to the blockchain. If the blockchain supports a "smart contract virtual machine" then transactions can refer these contracts and interact with it. The most known example of an interpreter, that processes the input of smart contract code, is the "ethereum virtual machine" (EVM). 

<img src="evm_structure.jpg" width="200" alt="Ethereum Virtual Machine" />
###### Source: https://ethereum.stackexchange.com/questions/268/ethereum-block-architecture/6413#6413

## Smart Contracts and the Cosmos SDK

With the Cosmos SDK you can use smart contracts in various ways, either you add the whole EVM to your blockchain with adding the module `wasm`, which would allow users to upload any smart contract to your blockchain - or you can use specified modules that allow only specified smart contracts to be executed by the user. There are a variety of modules that allow for NFTs, ERC20 like Tokens, DeFi contracts and more. 
The biggest difference in using a whole smart contract virtual machine or only a subset of smart contracts is the security and programming language parameters. The virtual machines tend to require more statically typed languages like Solidity, Rust or Ocaml, while implementing your smart contract in a module would allow for Go or even JavaScript to be used to create smart contract designs.

We will be looking into creating and uploading a Smart Contract with the `wasm` or `evm` module. In later chapters, we will also look closer into creating our own "smart modules" that use particular ideas of smart contracts and embedd them into a Cosmos SDK blockchain.

## Summary

- Smart contracts need to be deterministic in output.
- The most Smart contracts run on the Ethereum virtual machine.
- Starport and Cosmos support smart contracts but also smart modules.
