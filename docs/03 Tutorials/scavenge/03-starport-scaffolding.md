---
order: 3
---

# Starport

We'll be using a tool called [starport](https://github.com/tendermint/starport) to help us spin up a boilerplate app quickly. 

You can install `starport` via `npm`, `brew`, or building it from [source](https://github.com/tendermint/starport).

#### npm install
```bash
npm i -g @tendermint/starport
```

#### homebrew installation
```bash
brew install tendermint/tap/starport
```

#### building from source
```bash
git clone https://github.com/tendermint/starport && cd starport && make
```


Afterwards, you can enter in `starport` in your terminal, and should see the following help text displayed:
```sh
$ starport
A tool for scaffolding out Cosmos applications

Usage:
  starport [command]

Available Commands:
  add         Adds a feature to a project.
  app         Generates an empty application
  help        Help about any command
  serve       Launches a reloading server
  type        Generates CRUD actions for type
  version     Version will output the current build information

Flags:
  -h, --help     help for starport
  -t, --toggle   Help message for toggle

Use "starport [command] --help" for more information about a command.
```

Now that the `starport` command is available, you can scaffold an application by using the `starport app` command:

```bash
$ starport app --help
Generates an empty application

Usage:
  starport app [github.com/org/repo] [flags]

Flags:
      --address-prefix string   Address prefix (default "cosmos")
  -h, --help                    help for app
```

Let's start by scaffolding our `scavenge` application with `starport app`. This should generate a directory of folders called `scavenge` inside your current working directory, as well as scaffold our `scavenge` module. 

```bash
$ starport app github.com/github-username/scavenge

⭐️ Successfully created a Cosmos app 'scavenge'.
👉 Get started with the following commands:

 % cd scavenge
 % starport serve

NOTE: add --verbose flag for verbose (detailed) output.
```

You've successfully scaffolded a Cosmos SDK application using `starport`! In the next step, we're going to run the application using the instructions provided. 