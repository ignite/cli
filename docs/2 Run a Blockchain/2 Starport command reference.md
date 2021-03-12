# Starport command reference

The Starport tool CLI provides command line help. <!-- should we change "application" and "app" to "blockchain" in the command descriptions? I'm confused at the interchangeble terms -->

For example, to list all Starport commands, type `starport`:

```txt
Usage:
  starport [command]

Available Commands:
  app         Generates an empty application
  build       Builds an app and installs binaries
  chain       Relay connects blockchains via IBC protocol
  faucet      Send coins to an account
  help        Help about any command
  module      Manage cosmos modules for your app
  network     Create and start blockchains collaboratively
  serve       Launches a reloading server
  type        Generates CRUD actions for type
  version     Version will output the current build information

Flags:
  -h, --help     help for starport
  -t, --toggle   Help message for toggle

Use "starport [command] --help" for more information about a command.
```

This command list was generated for Starport version v0.14.0.
