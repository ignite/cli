---
sidebar_position: 5
description: Other commands for coordination.
---

# Other commands for coordination

Ignite CLI offers various other commands to coordinate chain launches that can be used by coordinators, validators, or other participants.

The requests follow the same logic as the request for validator participation; they must be approved by the chain coordinator to be effective in the genesis.

---

## Request a genesis account

Any participant can request a genesis account with an associated balance for the chain.
The participant must provide an address with a comma-separated list of token balances.

Any prefix can be used for the Bech32 address, it is automatically converted into `spn` on the Ignite Chain.

```
ignite n request add-account 3 spn1pe5h2gelhu8aukmrnj0clmec56aspxzuxcy99y 1000stake
```

**Output**

```
Source code fetched
Blockchain set up
⋆ Request 10 to add account to the network has been submitted!
```
---

## Request to remove a genesis account

Any participant can request to remove a genesis account from the chain genesis.
It might be the case if, for example, a user suggests an account balance that is so high it could harm the network.
The participant must provide the address of the account.

Any prefix can be used for the Bech32 address, it is automatically converted into `spn` on the Ignite Chain.

```
ignite n request remove-account 3 spn1pe5h2gelhu8aukmrnj0clmec56aspxzuxcy99y
```

**Output**

```
Request 11 to remove account from the network has been submitted!
```
---

## Request to remove a genesis validator

Any participant can request to remove a genesis validator (gentx) from the chain genesis.
It might be the case if, for example, a chain failed to launch because of some validators, and they must be removed from genesis.
The participant must provide the address of the validator account (same format as genesis account).

Any prefix can be used for the Bech32 address, it is automatically converted into `spn` on the Ignite Chain.

The request removes only the gentx from the genesis but not the associated account balance.

```
ignite n request remove-validator 429 spn1pe5h2gelhu8aukmrnj0clmec56aspxzuxcy99y
```

**Output**

```
Request 12 to remove validator from the network has been submitted!
```
---