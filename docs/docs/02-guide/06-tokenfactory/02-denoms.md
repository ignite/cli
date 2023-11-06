# Understanding Denoms in Cosmos SDK and Ignite

## What is a Denom?

**Denom** stands for `denomination` and represents the name of a token within the Cosmos SDK and Ignite. In the Cosmos ecosystem, denoms play a crucial role in identifying and managing tokens.

In Ignite, the configuration of your blockchain, including the specification of denoms, is set in the `config.yml` file within your blockchain directory. This file allows the definition of various denoms before initializing your blockchain.

Common examples of denoms include formats like `token` or `stake`.

## Usage of Denoms

In the Cosmos SDK, assets are represented as a `Coins` type, which combines an amount with a denom. The amount is flexible, allowing for a wide range of values. Accounts in the Cosmos SDK, including both basic and module accounts, maintain balances comprised of these `Coins`.

The `x/bank` module is pivotal in the Cosmos SDK as it tracks all account balances and the total supply of tokens in the application.

### Key Points on Denoms and Balances:

- **Fixed Denomination Unit:** The Cosmos SDK treats the amount of a balance as a single, fixed unit of denomination, regardless of the denom itself.
- **Client and App Flexibility:** While clients and apps built on Cosmos SDK chains can define arbitrary denomination units, all transactions and operations in the Cosmos SDK ultimately use these fixed units.
- **Example:** On the Cosmos Hub (Gaia), the common assumption is 1 ATOM = 10^6 uatom, and operations are based on these units of 10^6.

## Denoms and IBC (Inter-Blockchain Communication)

One of the primary uses of IBC is the transfer of tokens between blockchains. This process involves creating a token `voucher` on the target blockchain upon receiving tokens from a source chain.

### Characteristics of IBC Voucher Tokens:

- **Naming Convention:** IBC voucher tokens are denoted with a naming syntax that starts with `ibc/`. This convention helps in identifying and managing IBC tokens on a blockchain.
- **Native vs. Voucher Tokens:** With IBC, a native token on one blockchain can be referenced as a `voucher` token on another. These tokens are differentiated by their `denom` names.

For a comprehensive understanding of IBC denoms and their application, refer to [Understand IBC Denoms with Gaia](https://tutorials.cosmos.network/tutorials/6-ibc-dev/), which provides detailed insights into the format and utilization of voucher tokens in the IBC context.
