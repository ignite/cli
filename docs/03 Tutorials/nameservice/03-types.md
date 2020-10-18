---
order: 3
---

# Types

First thing we're going to do is create a type in the `/x/nameservice` folder with the `starport` tool using the below command:


```bash
starport type name value price
```

Currently, we're only passing two values when scaffolding the `whois` type, because there are additional fields that we will be replacing. In all uses of `whois`, we'll be removing the auto-generated `ID` field, and replacing the `Creator` field with `Owner` to reflect ownership of the name.

## `types.go`

Now we can continue with creating a module. The `./x/nameservice/types` folder will hold custom types for the module, and you should see a few files that have already been scaffolded, including `MsgCreateWhois.go`, `MsgDeleteWhois.go`, `MsgSetWhois.go`, and `TypeWhois.go`.

## Whois

Each name will have three pieces of data associated with it.

- Value - The value that a name resolves to. This is just an arbitrary string, but in the future you can modify this to require it fitting a specific format, such as an IP address, DNS Zone file, or blockchain address.
- Owner - The address of the current owner of the name
- Price - The price you will need to pay in order to buy the name

To start your SDK module, define your `nameservice.Whois` struct in the `./x/nameservice/types/TypeWhois.go` file:

<<< @/nameservice/nameservice/x/nameservice/types/TypeWhois.go

As mentioned in the [Design doc](./app-design.md), if a name does not already have an owner, we want to initialize it with some MinPrice.
