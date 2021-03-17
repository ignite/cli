# Project Scaffold Reference

The `starport app` command scaffolds a project. By default, the Cosmos SDK version is Stargate. <!-- what is a project? compared to a "blockchain" or "app" -->

## Address prefix

You can change the way addresses look in your blockchain.

On the Cosmos SDK Hub, addresses have a `cosmos` prefix, like `cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

To specify a custom address prefix on the command line, use the `--address-prefix` flag. For example, to change the blockchain prefix to moonlight:

```
starport app github.com/foo/bar --address-prefix moonlight
```

To change the address prefix for subsequent blockchain builds:

1. Change the `AccountAddressPrefix` variable in the `/app/prefix.go` file. Be sure to preserve other variables in the file.
2. To recognize the new prefix, change the `VUE_APP_ADDRESS_PREFIX` variable in `/vue/.env`.
