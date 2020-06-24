# Starport CLI

Code scaffolding tool `starport` for Cosmos SDK applications.

## Installation

### macOS

```
brew install tendermint/tap/starport
```

### GNU/Linux

Will be released as a Snap soon. Use `git clone` until then.

### Build locally

```
git clone https://github.com/tendermint/starport && cd starport && make
```

## Creating an application

```
starport app [modulePath]
```

This command creates an empty template for a Cosmos SDK application. By default it also includes a module with the same name as the package.

To create a new application called `blog`, run:

```
starport app github.com/example/blog
```

## Running an application

```
starport serve
```

To start the server, go into you application's directory and run `starport serve`. This commands installs dependencies, builds and initializes the app and runs both Tendermint RPC server (by default on `localhost:26657`) as well as LCD (by default on `localhost:1317`) with hot reloading enabled.

Note: depending on your OS and firewall settings, you may have to accept a prompt asking if your application's binary (`blogd` in this case) can accept external connections.

## Creating types

```
starport type [typeName] [field1] [field2:bool] ...
```

This command generates messages, handlers, keepers, CLI and REST clients and type definition for `typeName` type. A type can have any number of `field` arguments. By default fields are strings, but `bool` and `int` are supported.

For example,

```
starport type post title body
```

This command generates a type `Post` with two fields: `title` and `body`.

To add a post run `blogcli tx blog create-post "My title" "This is a blog" --from=me`.

## Front-end application

By default the generator creates a front-end application in `ui` directory. If you have Node.js installed, run `cd ui && npm i && npm run serve` to launch the app. With this app you can generate accounts, request tokens from a faucet, list and create objects generated with `starport type`.
